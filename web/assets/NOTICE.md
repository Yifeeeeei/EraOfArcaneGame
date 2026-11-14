# Frontend asset sources

## Arcane emblem

`brand/arcane-emblem.png` was created with the built-in imagegen tool on
2026-11-14 from the user's supplied black-and-white arcane line drawing.
It preserves the supplied interwoven ribbon design, with antique-gold shading
and a transparent background. The source drawing retains its original ownership.
The exact generation prompt and visual check are kept with the ignored local
review record under `docs/project/payment-and-emblem-20261114/`.

## Original concept backgrounds

- `backgrounds/royal-hall.png`
- `backgrounds/battlefield.png`

These backgrounds were newly generated with the built-in imagegen tool for the
Era of Arcane frontend proposal. They are copied unchanged from the approved
local proposal under `docs/project/frontend-visual-proposal-v1/assets/generated/`.
The original prompts and generation record remain in that local proposal's
`GENERATED-PROMPTS.md` and `ASSETS.md`.

No Hearthstone, MTG Arena, Yu-Gi-Oh! Master Duel, or third-party commercial UI
artwork was used as an input or copied into these assets.

## Existing Era of Arcane artwork

- `card-back.png`: unchanged project card back from
  [ArcaneImages/card_back/crad_back0.png](https://github.com/Yifeeeeei/ArcaneImages/blob/main/card_back/crad_back0.png).
- `elements/{air,dark,earth,fire,light,none,water}.png`: unchanged copies of the
  approved proposal's small element glyphs, originally sourced from the project
  deck composer's `resources/placeholders/ele_*.png` assets.
- `elements/{attack,duration,life,power}.png`: unchanged stat glyphs from the
  same project deck composer.
- `stories/{EmptyThrone,Knight,Moon}.md`: current story text copied unchanged
  from the project's EraOfArcane website source. These replaceable content
  files are independent of game logic and interface structure.

These existing project resources retain their original ownership and applicable
project permissions. This notice does not assign them a new third-party licence.
The large elemental illustrations from the homepage are not used as UI glyphs.

## Original card artwork

`card-art-index.js` contains only original-art URLs for the currently supported
cards. Generate it with `python scripts/sync-card-art-index.py` after running
ArcaneImages' `upload.py`. The image repository's `drawings/` files
are unchanged copies of `arcane_go/drawings/`; no new ownership or licence is
assigned to these project images. Artwork pixels are not duplicated here.

Players load artwork from the public ArcaneImages GitHub Pages site, with the
repository's raw files as a fallback. The game server does not read or serve
an adjacent image checkout, and players need no local image repository.

## Fonts

The proposal's font sources are Google Fonts
[Cinzel](https://github.com/google/fonts/tree/main/ofl/cinzel) and
[Noto Serif SC](https://github.com/google/fonts/tree/main/ofl/notoserifsc).
Their original SIL Open Font License 1.1 notices are retained unchanged in:

- `fonts/cinzel-OFL.txt`
- `fonts/noto-serif-sc-OFL.txt`

Any bundled webfont or subset derived from these sources remains covered by
its corresponding OFL notice. Font software must not be sold by itself.

The bundled variable WOFF2 subsets are generated with
`python scripts/subset-ui-fonts.py`. The script reads all HTML, JavaScript and
CSS under `web/`, `data/supported_card_infos.json`, and the existing stories
`EmptyThrone.md`, `Moon.md`, and `Knight.md` bundled in `web/assets/stories/`,
then adds ASCII and common Chinese punctuation. The original
TTF files are not included in the web assets.

<!-- font-subset-table:start -->
| Font | WOFF2 bytes | Retained Unicode codepoints | Preserved weight axis |
| --- | ---: | ---: | --- |
| `fonts/noto-serif-sc.woff2` | 783,900 | 2,132 | 200–900 |
| `fonts/cinzel.woff2` | 26,448 | 108 | 400–900 |
<!-- font-subset-table:end -->

All requested codepoints supported by each original font are retained in the
subset's cmap; no supported characters were lost. The original Noto Serif SC
covers every requested CJK character. Symbols outside an original font's cmap
continue to use the CSS system fallback. Regenerate subsets after introducing
new UI copy, card metadata or story text. Exact sizes, hashes and coverage
counts are recorded in `fonts/subset-manifest.json`.
