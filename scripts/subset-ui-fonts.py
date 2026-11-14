#!/usr/bin/env python3
"""Build local OFL WOFF2 subsets without changing card data or source text.

Run from any directory:
  python scripts/subset-ui-fonts.py
  python scripts/subset-ui-fonts.py --source-dir /path/to/original-fonts

Requires fonttools[woff] and brotli. An ignored tmp/font-tools installation is
also supported. Original fonts are not copied into the web distribution.
"""

from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path
import re
import sys
import unicodedata


REPO = Path(__file__).resolve().parents[1]
LOCAL_TOOLS = REPO / "tmp" / "font-tools"
if LOCAL_TOOLS.is_dir():
    sys.path.insert(0, str(LOCAL_TOOLS))

try:
    import fontTools
    from fontTools import subset
    from fontTools.ttLib import TTFont
except ImportError as error:
    raise SystemExit(
        "Install fonttools[woff] and brotli, or place them in tmp/font-tools."
    ) from error


STORIES = ("EmptyThrone.md", "Moon.md", "Knight.md")
PUNCTUATION = "\u3000，。！？、；：‘’“”「」『』（）《》〈〉【】〔〕［］｛｝…—–·￥％＋－×÷＝＜＞～｜／＼"
FONT_NAMES = ("noto-serif-sc", "cinzel")


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def collect_text(story_root: Path) -> tuple[set[int], list[dict[str, str]]]:
    paths = sorted(
        path
        for path in (REPO / "web").rglob("*")
        if path.is_file() and path.suffix.lower() in {".html", ".js", ".css"}
    )
    paths.append(REPO / "data" / "supported_card_infos.json")
    paths.extend(story_root / name for name in STORIES)
    codepoints = set(range(0x20, 0x7F))
    codepoints.update(map(ord, PUNCTUATION))
    sources = []
    for path in paths:
        if not path.is_file():
            raise FileNotFoundError(f"Required font text source is missing: {path}")
        text = path.read_text(encoding="utf-8-sig")
        codepoints.update(
            ord(char)
            for char in text
            if ord(char) >= 0x20
            and not 0xD800 <= ord(char) <= 0xDFFF
            and unicodedata.category(char) != "Cc"
        )
        try:
            display = path.relative_to(REPO).as_posix()
        except ValueError:
            display = "../EraOfArcane/mds/" + path.name
        sources.append({"path": display, "sha256": sha256(path)})
    return codepoints, sources


def axes(font: TTFont) -> list[dict[str, int | float | str]]:
    if "fvar" not in font:
        return []
    return [
        {
            "tag": axis.axisTag,
            "min": axis.minValue,
            "default": axis.defaultValue,
            "max": axis.maxValue,
        }
        for axis in font["fvar"].axes
    ]


def build_font(name: str, source_dir: Path, output_dir: Path, requested: set[int]) -> dict:
    source = source_dir / f"{name}.ttf"
    target = output_dir / f"{name}.woff2"
    font = TTFont(source, recalcTimestamp=False)
    original_cmap = set(font.getBestCmap())
    wanted = original_cmap & requested
    original_axes = axes(font)
    original_glyphs = len(font.getGlyphOrder())
    options = subset.Options()
    options.layout_features = ["*"]
    options.name_IDs = ["*"]
    options.name_legacy = True
    options.name_languages = ["*"]
    options.notdef_glyph = True
    options.notdef_outline = True
    options.recalc_timestamp = False
    subsetter = subset.Subsetter(options=options)
    subsetter.populate(unicodes=wanted)
    subsetter.subset(font)
    font.flavor = "woff2"
    font.save(target)
    font.close()

    with TTFont(target, recalcTimestamp=False) as result:
        result_cmap = set(result.getBestCmap())
        missing = wanted - result_cmap
        if missing:
            raise RuntimeError(f"{name} lost {len(missing)} required source glyphs")
        if axes(result) != original_axes:
            raise RuntimeError(f"{name} variable axes changed during subsetting")
        glyphs = len(result.getGlyphOrder())

    unsupported = requested - original_cmap
    missing_cjk = sorted(cp for cp in unsupported if 0x3400 <= cp <= 0x9FFF)
    if name == "noto-serif-sc" and missing_cjk:
        raise RuntimeError(
            "Source Noto Serif SC does not cover required CJK characters: "
            + ", ".join(f"U+{cp:04X}" for cp in missing_cjk)
        )
    return {
        "file": f"{name}.woff2",
        "source_file": source.name,
        "source_sha256": sha256(source),
        "source_bytes": source.stat().st_size,
        "bytes": target.stat().st_size,
        "sha256": sha256(target),
        "requested_codepoints": len(requested),
        "source_supported_requested_codepoints": len(wanted),
        "output_cmap_codepoints": len(result_cmap),
        "missing_supported_codepoints": len(missing),
        "source_glyphs": original_glyphs,
        "subset_glyphs": glyphs,
        "axes": original_axes,
        "source_unsupported_codepoints": [f"U+{cp:04X}" for cp in sorted(unsupported)],
    }


def update_notice(output_dir: Path, reports: list[dict]) -> None:
    notice = output_dir.parent / "NOTICE.md"
    if not notice.is_file():
        return
    text = notice.read_text(encoding="utf-8")
    start = "<!-- font-subset-table:start -->"
    end = "<!-- font-subset-table:end -->"
    if start not in text or end not in text:
        return
    rows = [
        start,
        "| Font | WOFF2 bytes | Retained Unicode codepoints | Preserved weight axis |",
        "| --- | ---: | ---: | --- |",
    ]
    for report in reports:
        axis = next((item for item in report["axes"] if item["tag"] == "wght"), None)
        weights = f"{axis['min']:g}–{axis['max']:g}" if axis else "static"
        rows.append(
            f"| `fonts/{report['file']}` | {report['bytes']:,} | "
            f"{report['output_cmap_codepoints']:,} | {weights} |"
        )
    rows.append(end)
    notice.write_text(
        re.sub(re.escape(start) + r".*?" + re.escape(end), "\n".join(rows), text, flags=re.S),
        encoding="utf-8",
    )


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--source-dir",
        type=Path,
        default=REPO / "docs/project/frontend-visual-proposal-v1/assets/fonts",
        help="Directory containing the original noto-serif-sc.ttf and cinzel.ttf",
    )
    parser.add_argument("--output-dir", type=Path, default=REPO / "web/assets/fonts")
    parser.add_argument("--story-root", type=Path, default=REPO / "web/assets/stories")
    parser.add_argument(
        "--report",
        type=Path,
        default=REPO / "docs/project/frontend-implementation-20261114/font-subset-report.json",
        help="Detailed local report, including unsupported source characters",
    )
    args = parser.parse_args()
    requested, sources = collect_text(args.story_root.resolve())
    args.output_dir.mkdir(parents=True, exist_ok=True)
    reports = []
    for name in FONT_NAMES:
        report = build_font(name, args.source_dir, args.output_dir, requested)
        reports.append(report)
        print(
            f"{name}: {report['source_bytes']:,} -> {report['bytes']:,} bytes; "
            f"cmap {report['output_cmap_codepoints']}; "
            f"missing supported {report['missing_supported_codepoints']}; "
            f"axes {json.dumps(report['axes'])}",
            flush=True,
        )
    details = {
        "fonttools_version": fontTools.__version__,
        "requested_codepoints": len(requested),
        "text_sources": sources,
        "fonts": reports,
    }
    args.report.parent.mkdir(parents=True, exist_ok=True)
    args.report.write_text(json.dumps(details, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    compact_fonts = [
        {key: value for key, value in report.items() if key != "source_unsupported_codepoints"}
        for report in reports
    ]
    manifest = {
        "generator": "scripts/subset-ui-fonts.py",
        "fonttools_version": fontTools.__version__,
        "text_source_count": len(sources),
        "requested_codepoints": len(requested),
        "cmap_coverage": "All requested codepoints supported by each original font are retained.",
        "fallback": "Characters outside an original font's cmap require the CSS system fallback.",
        "fonts": compact_fonts,
    }
    (args.output_dir / "subset-manifest.json").write_text(
        json.dumps(manifest, ensure_ascii=False, indent=2) + "\n", encoding="utf-8"
    )
    update_notice(args.output_dir, reports)


if __name__ == "__main__":
    main()
