#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const root = path.resolve(__dirname, '..');
const cardsPath = path.join(root, 'data', 'supported_card_infos.json');
const catalogPath = path.join(root, 'server', 'game', 'card_effects_catalog.go');
const outputPath = path.join(root, 'EFFECT_COVERAGE.md');
const gameDir = path.join(root, 'server', 'game');
const tmpDir = path.join(root, 'tmp');

const cards = JSON.parse(fs.readFileSync(cardsPath, 'utf8'));
const catalog = fs.readFileSync(catalogPath, 'utf8');
const implemented = new Set(
  [...catalog.matchAll(/"(\d+)":\s*func\(\) CardBehavior/g)].map((match) => match[1]),
);
const testText = fs.readdirSync(gameDir)
  .filter((file) => file.endsWith('_test.go'))
  .map((file) => fs.readFileSync(path.join(gameDir, file), 'utf8'))
  .join('\n');

const frontendReports = fs.existsSync(tmpDir)
  ? fs.readdirSync(tmpDir)
    .filter((file) => /^frontend-card-operation-report.*\.json$/.test(file))
    .map((file) => {
      try {
        return { file, data: JSON.parse(fs.readFileSync(path.join(tmpDir, file), 'utf8')) };
      } catch {
        return null;
      }
    })
    .filter(Boolean)
  : [];
const frontendByCard = new Map();
for (const report of frontendReports) {
  for (const result of report.data.results || []) {
    if (!result || !result.number) continue;
    const previous = frontendByCard.get(result.number);
    const current = {
      status: result.status || 'unknown',
      file: report.file,
      generatedAt: report.data.generatedAt || '',
      reason: result.reason || '',
    };
    if (!previous || current.status === 'pass' || current.generatedAt > previous.generatedAt) {
      frontendByCard.set(result.number, current);
    }
  }
}

const genericRules = [
  ['速攻', /速攻/],
  ['穿透', /穿透/],
  ['冷却', /冷却\d*/],
  ['防御限制', /防御|不可用于防御/],
  ['咒术立即结算', /咒术/],
  ['隐蔽', /隐蔽/],
  ['引魔', /引魔/],
  ['护盾', /护盾/],
  ['临时', /临时/],
  ['状态容器', /冻结|晕眩|眩晕|石化|点燃|虚弱/],
  ['命中状态', /点燃\d+|冻结\d+|晕眩\d+|眩晕\d+|石化\d+|虚弱\d+/],
  ['获得元素', /获得\d+\\[火水地气光暗无]/],
  ['基础范围', /范围:方阵|范围:纵列|范围:前排|范围:溅射|范围:全场/],
  ['法术攻威被动', /你的(所有|火焰|水纹|地脉|大气|光辉|暗影|奥术|水纹和大气)?法术(?:在攻击时|攻击和强化攻击时)?\+\d+\\(威|攻)/],
];

const customEffectPattern =
  /入场|遗言|绝技|回合技|消耗:|吞噬|绑定|异能|反制|命中|防御失败|获得|抽|检索|展示|标记|费用|减少|增加|\+|摧毁|消灭|丢弃|选择|献祭|召唤|学习|装备|法力范围|范围:方阵|范围:纵列|范围:溅射|无法|不可|每当|如果|当|游戏开始前|起始手牌|换牌/;

const mechanicGroups = [
  ['选择/待选动作', /选择|任意数量/],
  ['检索/看牌/洗牌/牌堆顶', /检索|查看|牌堆|卡组顶|洗回|翻开|重洗/],
  ['绑定技能/衍生牌', /绑定技能|衍生|置于.*手牌|加入手牌/],
  ['费用修改', /费用|花费.*减少|减少\d+\\|减|\+1\\威|\+2\\威|\+1\\攻|\+2\\攻|法术\+|技能.*获得/],
  ['状态/标记/计数器', /点燃|冻结|晕眩|眩晕|石化|虚弱|护盾|隐蔽|标记|充能/],
  ['吞噬/献祭/死亡触发', /吞噬|献祭|死亡|遗言|消灭|摧毁/],
  ['范围/目标规则', /范围:方阵|范围:纵列|范围:溅射|范围:前排|法力范围|全场|穿透/],
  ['主动能力', /绝技|回合技|消耗:/],
  ['反制/触发监听', /反制|每当|当敌方|当对方|如果|在你/],
];

function genericTags(description) {
  return genericRules
    .filter(([, pattern]) => pattern.test(description || ''))
    .map(([name]) => name);
}

function coverageStatus(card) {
  const description = (card.description || '').trim();
  if (implemented.has(card.number)) return '已实现';
  if (!description) return '白板/无效果';
  if (isGenericOnlyDescription(description)) return '通用机制';
  if (customEffectPattern.test(description)) return '未实现';
  if (genericTags(description).length) return '通用机制';
  return '需复核';
}

function isGenericOnlyDescription(description) {
  const normalized = description
    .replace(/咒术/g, '')
    .replace(/法术/g, '')
    .replace(/创造|驱动|幻变|聚能|灵媒/g, '')
    .trim();
  const parts = normalized.split(/[.。,，]/).map((part) => part.trim()).filter(Boolean);
  if (!parts.length) return true;
  return parts.every((part) =>
    /^速攻$/.test(part) ||
    /^穿透$/.test(part) ||
    /^冷却\d+$/.test(part) ||
    /^防御$/.test(part) ||
    /^范围:(方阵|纵列|前排|溅射|全场)$/.test(part) ||
    /^(点燃|冻结|晕眩|眩晕|石化|虚弱)\d+$/.test(part) ||
    /^命中:使目标伙伴(点燃|冻结|晕眩|眩晕|石化|虚弱)\d+$/.test(part) ||
    /^获得\d+\\[火水地气光暗无]$/.test(part) ||
    /^你的(所有|火焰|水纹|地脉|大气|光辉|暗影|奥术|水纹和大气)?法术(?:在攻击时|攻击和强化攻击时)?\+\d+\\(威|攻)$/.test(part) ||
    /^无法强化或被强化$/.test(part) ||
    /^不可用于防御$/.test(part) ||
    /^无法用于强化$/.test(part) ||
    /^不可用于强化$/.test(part)
  );
}

function statusKey(status) {
  return {
    已实现: 'implemented',
    未实现: 'missing',
    通用机制: 'generic',
    '白板/无效果': 'vanilla',
    需复核: 'review',
  }[status] || status;
}

function countBy(items, keyFn) {
  return items.reduce((counts, item) => {
    const key = keyFn(item);
    counts[key] = (counts[key] || 0) + 1;
    return counts;
  }, {});
}

function escapeCell(value) {
  return String(value || '').replace(/\|/g, '\\|').replace(/\n/g, '<br>');
}

function semanticTestStatus(card) {
  const mentions = (testText.match(new RegExp(card.number, 'g')) || []).length;
  if (mentions > 0) return { status: '已点名', mentions };
  if (!customEffectPattern.test(card.description || '') && isGenericOnlyDescription(card.description || '')) {
    return { status: '通用/白板', mentions };
  }
  return { status: '未点名', mentions };
}

function frontendStatus(card) {
  const result = frontendByCard.get(card.number);
  if (!result) return '未记录';
  if (result.status === 'pass') return `已操作(${result.file})`;
  return `${result.status}(${result.file})`;
}

function riskLevel(row) {
  if (row.runtimeStatus === '未实现' || row.runtimeStatus === '需复核') return '高';
  if (row.semanticStatus === '未点名' && row.runtimeStatus === '已实现') return '高';
  if (row.semanticStatus === '未点名') return '中';
  if (row.frontendStatus === '未记录' && row.runtimeStatus !== '白板/无效果') return '中';
  return '低';
}

const rows = cards.map((card) => {
  const status = coverageStatus(card);
  const semantic = semanticTestStatus(card);
  return {
    number: card.number,
    name: card.name,
    type: card.type,
    element: card.category,
    runtimeStatus: status,
    semanticStatus: semantic.status,
    semanticMentions: semantic.mentions,
    frontendStatus: frontendStatus(card),
    generic: genericTags(card.description).join('、'),
    description: card.description || '',
  };
});
for (const row of rows) {
  row.risk = riskLevel(row);
}

const counts = countBy(rows, (row) => row.runtimeStatus);
const semanticCounts = countBy(rows, (row) => row.semanticStatus);
const frontendCounts = countBy(rows, (row) => row.frontendStatus.startsWith('已操作') ? '已操作' : row.frontendStatus === '未记录' ? '未记录' : '未通过/旧失败');
const riskCounts = countBy(rows, (row) => row.risk);
const byType = {};
for (const row of rows) {
  byType[row.type] ||= {};
  byType[row.type][row.runtimeStatus] = (byType[row.type][row.runtimeStatus] || 0) + 1;
}

const missingRows = rows.filter((row) => row.runtimeStatus === '未实现');
const mechanicCounts = {};
for (const [name, pattern] of mechanicGroups) {
  mechanicCounts[name] = missingRows.filter((row) => pattern.test(row.description)).length;
}
const nextRows = rows
  .filter((row) => row.risk === '高')
  .sort((a, b) => {
    const runtimeRank = { 未实现: 0, 需复核: 1, 已实现: 2, 通用机制: 3, '白板/无效果': 4 };
    return (runtimeRank[a.runtimeStatus] ?? 9) - (runtimeRank[b.runtimeStatus] ?? 9) || a.number.localeCompare(b.number);
  });

const lines = [];
lines.push('# 基础包卡牌效果覆盖表');
lines.push('');
lines.push('> 这个文件由 `node tools/effect-coverage-report.js --write` 生成，用来追踪基础包卡牌效果是否有运行时代码、后端语义测试、以及前端操作记录。');
lines.push('');
lines.push('## 总览');
lines.push('');
lines.push(`- 基础包卡牌总数：${rows.length}`);
lines.push(`- 已实现专属行为：${counts['已实现'] || 0}`);
lines.push(`- 只依赖通用机制：${counts['通用机制'] || 0}`);
lines.push(`- 白板/无效果：${counts['白板/无效果'] || 0}`);
lines.push(`- 需复核：${counts['需复核'] || 0}`);
lines.push(`- 未实现：${counts['未实现'] || 0}`);
lines.push(`- 后端语义测试已点名：${semanticCounts['已点名'] || 0}`);
lines.push(`- 后端语义测试未点名：${semanticCounts['未点名'] || 0}`);
lines.push(`- 前端已有操作记录：${frontendCounts['已操作'] || 0}`);
lines.push(`- 高风险待验证：${riskCounts['高'] || 0}`);
lines.push('');
lines.push('## 按类型统计');
lines.push('');
lines.push('| 类型 | 已实现 | 通用机制 | 白板/无效果 | 需复核 | 未实现 |');
lines.push('|---|---:|---:|---:|---:|---:|');
for (const type of ['人物', '伙伴', '技能', '道具']) {
  const bucket = byType[type] || {};
  lines.push(`| ${type} | ${bucket['已实现'] || 0} | ${bucket['通用机制'] || 0} | ${bucket['白板/无效果'] || 0} | ${bucket['需复核'] || 0} | ${bucket['未实现'] || 0} |`);
}
lines.push('');
lines.push('## 未实现机制分布');
lines.push('');
lines.push('| 机制族 | 未实现卡数 |');
lines.push('|---|---:|');
for (const [name, count] of Object.entries(mechanicCounts)) {
  lines.push(`| ${name} | ${count} |`);
}
lines.push('');
lines.push('## 状态说明');
lines.push('');
lines.push('- `已实现`：已在 `server/game/card_effects_catalog.go` 注册专属 Go 行为。');
lines.push('- `通用机制`：目前没有专属行为，但描述主要落在引擎已有关键词机制上。仍需要前端实测。');
lines.push('- `白板/无效果`：描述为空，可以按普通卡运行。');
lines.push('- `需复核`：描述不为空，但脚本无法判断是否需要代码。');
lines.push('- `未实现`：描述明显要求行为代码或尚不存在的通用机制。');
lines.push('- `后端语义测试` 的 `已点名` 表示测试文件里直接出现过该卡编号；它不是完整正确性的证明，但能区分“只被全量 smoke 扫过”和“有专门语义断言”。');
lines.push('- `前端操作` 来自 `tmp/frontend-card-operation-report*.json`，是历史浏览器操作记录；旧失败不一定代表当前仍失败，但必须重新验证。');
lines.push('- `风险` 主要用于排队：已实现但没有语义测试、未实现、需复核的卡优先处理。');
lines.push('');
lines.push('## 下一批高风险清单');
lines.push('');
lines.push('| 编号 | 名称 | 类型 | 运行时代码 | 后端语义测试 | 前端操作 | 描述 |');
lines.push('|---|---|---|---|---|---|---|');
for (const row of nextRows.slice(0, 60)) {
  lines.push(`| ${row.number} | ${escapeCell(row.name)} | ${row.type} | ${row.runtimeStatus} | ${row.semanticStatus} | ${escapeCell(row.frontendStatus)} | ${escapeCell(row.description)} |`);
}
lines.push('');
lines.push('## 全量清单');
lines.push('');
lines.push('| 编号 | 名称 | 类型 | 属性 | 运行时代码 | 后端语义测试 | 前端操作 | 风险 | 通用机制 | 描述 |');
lines.push('|---|---|---|---|---|---|---|---|---|---|');
for (const row of rows) {
  lines.push(`| ${row.number} | ${escapeCell(row.name)} | ${row.type} | ${row.element} | ${row.runtimeStatus} | ${row.semanticStatus}${row.semanticMentions ? `(${row.semanticMentions})` : ''} | ${escapeCell(row.frontendStatus)} | ${row.risk} | ${escapeCell(row.generic)} | ${escapeCell(row.description)} |`);
}
lines.push('');

const report = `${lines.join('\n')}`;

if (process.argv.includes('--write')) {
  fs.writeFileSync(outputPath, report);
  console.log(`Wrote ${path.relative(root, outputPath)} (${rows.length} cards)`);
} else {
  console.log(JSON.stringify({
    total: rows.length,
    counts: Object.fromEntries(Object.entries(counts).map(([key, value]) => [statusKey(key), value])),
    semanticCounts,
    frontendCounts,
    riskCounts,
    byType,
    mechanicCounts,
    highRisk: nextRows.slice(0, 60).map((row) => ({
      number: row.number,
      name: row.name,
      type: row.type,
      runtimeStatus: row.runtimeStatus,
      semanticStatus: row.semanticStatus,
      frontendStatus: row.frontendStatus,
      description: row.description,
    })),
  }, null, 2));
}
