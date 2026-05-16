#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const root = path.resolve(__dirname, '..');
const cardsPath = path.join(root, 'data', 'supported_card_infos.json');
const catalogPath = path.join(root, 'server', 'game', 'card_effects_catalog.go');
const outputPath = path.join(root, 'EFFECT_COVERAGE.md');

const cards = JSON.parse(fs.readFileSync(cardsPath, 'utf8'));
const catalog = fs.readFileSync(catalogPath, 'utf8');
const implemented = new Set(
  [...catalog.matchAll(/"(\d+)":\s*func\(\) CardBehavior/g)].map((match) => match[1]),
);

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

const rows = cards.map((card) => {
  const status = coverageStatus(card);
  return {
    number: card.number,
    name: card.name,
    type: card.type,
    element: card.category,
    status,
    generic: genericTags(card.description).join('、'),
    description: card.description || '',
  };
});

const counts = countBy(rows, (row) => row.status);
const byType = {};
for (const row of rows) {
  byType[row.type] ||= {};
  byType[row.type][row.status] = (byType[row.type][row.status] || 0) + 1;
}

const missingRows = rows.filter((row) => row.status === '未实现');
const mechanicCounts = {};
for (const [name, pattern] of mechanicGroups) {
  mechanicCounts[name] = missingRows.filter((row) => pattern.test(row.description)).length;
}

const lines = [];
lines.push('# 基础包卡牌效果覆盖表');
lines.push('');
lines.push('> 这个文件由 `node tools/effect-coverage-report.js --write` 生成，用来追踪基础包卡牌效果是否真的有运行时代码。');
lines.push('');
lines.push('## 总览');
lines.push('');
lines.push(`- 基础包卡牌总数：${rows.length}`);
lines.push(`- 已实现专属行为：${counts['已实现'] || 0}`);
lines.push(`- 只依赖通用机制：${counts['通用机制'] || 0}`);
lines.push(`- 白板/无效果：${counts['白板/无效果'] || 0}`);
lines.push(`- 需复核：${counts['需复核'] || 0}`);
lines.push(`- 未实现：${counts['未实现'] || 0}`);
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
lines.push('');
lines.push('## 全量清单');
lines.push('');
lines.push('| 编号 | 名称 | 类型 | 属性 | 状态 | 通用机制 | 描述 |');
lines.push('|---|---|---|---|---|---|---|');
for (const row of rows) {
  lines.push(`| ${row.number} | ${escapeCell(row.name)} | ${row.type} | ${row.element} | ${row.status} | ${escapeCell(row.generic)} | ${escapeCell(row.description)} |`);
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
    byType,
    mechanicCounts,
  }, null, 2));
}
