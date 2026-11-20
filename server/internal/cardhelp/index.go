// Package cardhelp builds display-only keyword references offline. It must not
// be used by the rules engine to determine a card's abilities or legal actions.
package cardhelp

import (
	"eraofarcane/model"
	"regexp"
	"strings"
)

var defenseOnly = regexp.MustCompile(`(?:^|[.。;；\n])\s*防御(?:$|[.。;；\n,，])`)
var enterHeader = regexp.MustCompile(`(?:^|[.。;；\n])\s*(?:入场(?:[,，、]遗言)?|遗言[,，、]入场)\s*[:：]`)
var areaHeader = regexp.MustCompile(`范围\s*[:：]`)

// References describes only the printed text, not current instance state. A
// reference can explain a term mentioned by a card without granting that term.
func References(card *model.Card) []string {
	d := card.Description
	refs := []string{}
	add := func(id string, matched bool) {
		if matched {
			refs = append(refs, id)
		}
	}
	add("K01", enterHeader.MatchString(d) || strings.Contains(d, "入场效果"))
	for _, rule := range []struct{ id, word string }{
		{"K02", "遗言"}, {"K03", "主动"}, {"K04", "诱发"}, {"K05", "光环"},
		{"K06", "回合技"}, {"K07", "绝技"}, {"K08", "祈咒"}, {"K09", "反制"},
		{"K10", "速攻"}, {"K11", "临时"}, {"K12", "异能"},
	} {
		add(rule.id, strings.Contains(d, rule.word))
	}
	limited := defenseOnly.MatchString(d)
	add("K13", limited)
	for _, rule := range []struct{ id, word string }{
		{"K14", "命中"}, {"K15", "精通"}, {"K16", "检索"}, {"K17", "翻取"},
		{"K18", "献祭"}, {"K19", "吞噬"}, {"K20", "绑定"},
	} {
		add(rule.id, strings.Contains(d, rule.word))
	}
	add("K21", strings.Contains(card.Tag, "衍生") || strings.Contains(d, "衍生"))
	add("K22", strings.Contains(d, "引魔"))
	add("K23", strings.Contains(d, "穿透"))
	add("K24", areaHeader.MatchString(strings.ReplaceAll(d, "法力范围", "")))
	for _, rule := range []struct{ id, word string }{
		{"K25", "前排"}, {"K26", "纵列"}, {"K27", "方阵"}, {"K28", "溅射"},
		{"K29", "全场"}, {"K30", "点燃"}, {"K31", "冻结"},
	} {
		add(rule.id, strings.Contains(d, rule.word))
	}
	add("K32", strings.Contains(d, "眩晕") || strings.Contains(d, "晕眩"))
	for _, rule := range []struct{ id, word string }{
		{"K33", "石化"}, {"K34", "虚弱"}, {"K35", "冷却"}, {"K36", "隐蔽"}, {"K37", "护盾"},
		{"A01", "负载"}, {"A02", "消耗"}, {"A03", "透支"}, {"A04", "强化"},
		{"A05", "重置"}, {"A06", "法力范围"},
	} {
		add(rule.id, strings.Contains(d, rule.word))
	}
	add("A07", !limited && strings.Contains(d, "防御"))
	for _, rule := range []struct{ id, word string }{
		{"A08", "标记物"}, {"A09", "负面状态"}, {"S01", "红月"},
		{"S02", "移出游戏"}, {"S03", "严格"},
	} {
		add(rule.id, strings.Contains(d, rule.word))
	}
	return refs
}
