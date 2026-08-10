// Code generated from data/supported_card_infos.json; DO NOT EDIT BY HAND.
package cards

import "eraofarcane/model"

// CardDefinition is the compiled, code-owned definition of a playable card.
// Runtime card instances point at the model.Card returned by these definitions.
type CardDefinition interface {
	ID() string
	Name() string
	Kind() string
	Element() string
	Card() model.Card
}

type CardDef1001101 struct{}

func (CardDef1001101) ID() string      { return "1001101" }
func (CardDef1001101) Name() string    { return "弃子" }
func (CardDef1001101) Kind() string    { return "伙伴" }
func (CardDef1001101) Element() string { return "无" }

func (CardDef1001101) Card() model.Card {
	return model.Card{
		Number:          "1001101",
		Type:            "伙伴",
		Name:            "弃子",
		Category:        "无",
		Tag:             "衍生-造物",
		Description:     "遗言:对周围单位造成1点伤害,如果有单位因此死亡,在该位置召唤1个弃子",
		Quote:           "弃子非懦,舍得非怯",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1001101.jpg",
	}
}

type CardDef1011001 struct{}

func (CardDef1011001) ID() string      { return "1011001" }
func (CardDef1011001) Name() string    { return "魔龙 奥瑞" }
func (CardDef1011001) Kind() string    { return "伙伴" }
func (CardDef1011001) Element() string { return "无" }

func (CardDef1011001) Card() model.Card {
	return model.Card{
		Number:          "1011001",
		Type:            "伙伴",
		Name:            "魔龙 奥瑞",
		Category:        "无",
		Tag:             "传奇-龙",
		Description:     "引魔.绑定技能:破灭魔光",
		Quote:           "吾即是始源,吾即是终焉,不生不灭,万法归一",
		ElementsCost:    map[string]int{"地": 1, "无": 5, "气": 1, "水": 1, "火": 1},
		ElementsGain:    map[string]int{"地": 1, "无": 2, "气": 1, "水": 1, "火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            5,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"3001001"},
		OutputPath:      "output\\基础包\\伙伴\\无\\1011001.jpg",
	}
}

type CardDef1011002 struct{}

func (CardDef1011002) ID() string      { return "1011002" }
func (CardDef1011002) Name() string    { return "巫师之塔 通天阁" }
func (CardDef1011002) Kind() string    { return "伙伴" }
func (CardDef1011002) Element() string { return "无" }

func (CardDef1011002) Card() model.Card {
	return model.Card{
		Number:          "1011002",
		Type:            "伙伴",
		Name:            "巫师之塔 通天阁",
		Category:        "无",
		Tag:             "传奇-造物",
		Description:     "引魔.入场:你的技能槽每有1个技能,就获得1点其属性对应的元素.光环:你的法力范围变为全场",
		Quote:           "俯瞰着整片大陆,难道你不想将世界收入囊中吗?",
		ElementsCost:    map[string]int{"无": 8},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1011002.jpg",
	}
}

type CardDef1011003 struct{}

func (CardDef1011003) ID() string      { return "1011003" }
func (CardDef1011003) Name() string    { return "盟主 法罗兰克" }
func (CardDef1011003) Kind() string    { return "伙伴" }
func (CardDef1011003) Element() string { return "无" }

func (CardDef1011003) Card() model.Card {
	return model.Card{
		Number:          "1011003",
		Type:            "伙伴",
		Name:            "盟主 法罗兰克",
		Category:        "无",
		Tag:             "传奇-巫师",
		Description:     "入场:获得等同于所有相邻伙伴负载的负载.绑定技能:纯净奥术",
		Quote:           "我们继承先贤的智慧,我们维护人间的秩序",
		ElementsCost:    map[string]int{"无": 9},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"3001002"},
		OutputPath:      "output\\基础包\\伙伴\\无\\1011003.jpg",
	}
}

type CardDef1011101 struct{}

func (CardDef1011101) ID() string      { return "1011101" }
func (CardDef1011101) Name() string    { return "收藏家 珊瑚 芬洛" }
func (CardDef1011101) Kind() string    { return "伙伴" }
func (CardDef1011101) Element() string { return "无" }

func (CardDef1011101) Card() model.Card {
	return model.Card{
		Number:          "1011101",
		Type:            "伙伴",
		Name:            "收藏家 珊瑚 芬洛",
		Category:        "无",
		Tag:             "传奇-人类",
		Description:     "诱发回合技:每当你装备1个道具,抽1张牌.诱发回合技:每当你打出1个消耗品,获得1\\无",
		Quote:           "\"没错啊,小岛上全是我的收藏品,不过这些东西加一起都比不上我的宝贝女儿\"",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1011101.jpg",
	}
}

type CardDef1011102 struct{}

func (CardDef1011102) ID() string      { return "1011102" }
func (CardDef1011102) Name() string    { return "\"指挥家\" 洛斯" }
func (CardDef1011102) Kind() string    { return "伙伴" }
func (CardDef1011102) Element() string { return "无" }

func (CardDef1011102) Card() model.Card {
	return model.Card{
		Number:          "1011102",
		Type:            "伙伴",
		Name:            "\"指挥家\" 洛斯",
		Category:        "无",
		Tag:             "传奇-人类",
		Description:     "引魔.诱发:每当场上4张卡牌被消耗,可以为你装备一个落幕提琴.主动:消耗此卡才能发动,重置你的所有落幕提琴",
		Quote:           "在曾经辉煌的大厅,一人继续排练着无人能听到的乐章",
		ElementsCost:    map[string]int{"无": 8},
		ElementsGain:    map[string]int{"无": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2001101"},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1011102.jpg",
	}
}

type CardDef1011103 struct{}

func (CardDef1011103) ID() string      { return "1011103" }
func (CardDef1011103) Name() string    { return "\"弈者\"" }
func (CardDef1011103) Kind() string    { return "伙伴" }
func (CardDef1011103) Element() string { return "无" }

func (CardDef1011103) Card() model.Card {
	return model.Card{
		Number:          "1011103",
		Type:            "伙伴",
		Name:            "\"弈者\"",
		Category:        "无",
		Tag:             "传奇-人类",
		Description:     "引魔.绑定技能:入局.主动绝技:对场上所有弃子造成1点伤害",
		Quote:           "世事如棋,步步谋局;权争似弈,胜负难明",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"1001101", "3001101"},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1011103.jpg",
	}
}

type CardDef1021001 struct{}

func (CardDef1021001) ID() string      { return "1021001" }
func (CardDef1021001) Name() string    { return "巫师的学徒" }
func (CardDef1021001) Kind() string    { return "伙伴" }
func (CardDef1021001) Element() string { return "无" }

func (CardDef1021001) Card() model.Card {
	return model.Card{
		Number:          "1021001",
		Type:            "伙伴",
		Name:            "巫师的学徒",
		Category:        "无",
		Tag:             "巫师",
		Description:     "",
		Quote:           "人总是在年轻时才有无限的潜力",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021001.jpg",
	}
}

type CardDef1021002 struct{}

func (CardDef1021002) ID() string      { return "1021002" }
func (CardDef1021002) Name() string    { return "学院导师" }
func (CardDef1021002) Kind() string    { return "伙伴" }
func (CardDef1021002) Element() string { return "无" }

func (CardDef1021002) Card() model.Card {
	return model.Card{
		Number:          "1021002",
		Type:            "伙伴",
		Name:            "学院导师",
		Category:        "无",
		Tag:             "巫师",
		Description:     "",
		Quote:           "别看他们似乎因循守旧,这里才是真正远离尘嚣的好去处",
		ElementsCost:    map[string]int{"无": 4},
		ElementsGain:    map[string]int{"地": 1, "气": 1, "水": 1, "火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021002.jpg",
	}
}

type CardDef1021003 struct{}

func (CardDef1021003) ID() string      { return "1021003" }
func (CardDef1021003) Name() string    { return "誓约巫师" }
func (CardDef1021003) Kind() string    { return "伙伴" }
func (CardDef1021003) Element() string { return "无" }

func (CardDef1021003) Card() model.Card {
	return model.Card{
		Number:          "1021003",
		Type:            "伙伴",
		Name:            "誓约巫师",
		Category:        "无",
		Tag:             "巫师",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 6},
		ElementsGain:    map[string]int{"无": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021003.jpg",
	}
}

type CardDef1021004 struct{}

func (CardDef1021004) ID() string      { return "1021004" }
func (CardDef1021004) Name() string    { return "守护骑士" }
func (CardDef1021004) Kind() string    { return "伙伴" }
func (CardDef1021004) Element() string { return "无" }

func (CardDef1021004) Card() model.Card {
	return model.Card{
		Number:          "1021004",
		Type:            "伙伴",
		Name:            "守护骑士",
		Category:        "无",
		Tag:             "人类",
		Description:     "",
		Quote:           "守护誓约的巫师将得到贴身的守护骑士和孤岛的漫漫长夜",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021004.jpg",
	}
}

type CardDef1021005 struct{}

func (CardDef1021005) ID() string      { return "1021005" }
func (CardDef1021005) Name() string    { return "内阁巫师" }
func (CardDef1021005) Kind() string    { return "伙伴" }
func (CardDef1021005) Element() string { return "无" }

func (CardDef1021005) Card() model.Card {
	return model.Card{
		Number:          "1021005",
		Type:            "伙伴",
		Name:            "内阁巫师",
		Category:        "无",
		Tag:             "巫师",
		Description:     "",
		Quote:           "想在内阁里站稳脚跟,权术和心计往往比法力更重要",
		ElementsCost:    map[string]int{"无": 4},
		ElementsGain:    map[string]int{"光": 1, "无": 1, "暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021005.jpg",
	}
}

type CardDef1021006 struct{}

func (CardDef1021006) ID() string      { return "1021006" }
func (CardDef1021006) Name() string    { return "杂货商贩" }
func (CardDef1021006) Kind() string    { return "伙伴" }
func (CardDef1021006) Element() string { return "无" }

func (CardDef1021006) Card() model.Card {
	return model.Card{
		Number:          "1021006",
		Type:            "伙伴",
		Name:            "杂货商贩",
		Category:        "无",
		Tag:             "人类",
		Description:     "入场:抽2张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021006.jpg",
	}
}

type CardDef1021007 struct{}

func (CardDef1021007) ID() string      { return "1021007" }
func (CardDef1021007) Name() string    { return "回收小精灵" }
func (CardDef1021007) Kind() string    { return "伙伴" }
func (CardDef1021007) Element() string { return "无" }

func (CardDef1021007) Card() model.Card {
	return model.Card{
		Number:          "1021007",
		Type:            "伙伴",
		Name:            "回收小精灵",
		Category:        "无",
		Tag:             "精灵",
		Description:     "入场:将你弃牌堆的1张牌放到卡组顶",
		Quote:           "学院喜闻乐见的小帮手,任何东西丢了都可以找它",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021007.jpg",
	}
}

type CardDef1021008 struct{}

func (CardDef1021008) ID() string      { return "1021008" }
func (CardDef1021008) Name() string    { return "预见先知" }
func (CardDef1021008) Kind() string    { return "伙伴" }
func (CardDef1021008) Element() string { return "无" }

func (CardDef1021008) Card() model.Card {
	return model.Card{
		Number:          "1021008",
		Type:            "伙伴",
		Name:            "预见先知",
		Category:        "无",
		Tag:             "精灵",
		Description:     "诱发:回合开始抽牌前,你可以查看牌堆顶的1张牌,将其放回牌堆顶或牌堆底",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021008.jpg",
	}
}

type CardDef1021009 struct{}

func (CardDef1021009) ID() string      { return "1021009" }
func (CardDef1021009) Name() string    { return "竞技场虚像" }
func (CardDef1021009) Kind() string    { return "伙伴" }
func (CardDef1021009) Element() string { return "无" }

func (CardDef1021009) Card() model.Card {
	return model.Card{
		Number:          "1021009",
		Type:            "伙伴",
		Name:            "竞技场虚像",
		Category:        "无",
		Tag:             "造物",
		Description:     "光环:不会受到法术攻击以外的伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021009.jpg",
	}
}

type CardDef1021010 struct{}

func (CardDef1021010) ID() string      { return "1021010" }
func (CardDef1021010) Name() string    { return "专精法师" }
func (CardDef1021010) Kind() string    { return "伙伴" }
func (CardDef1021010) Element() string { return "无" }

func (CardDef1021010) Card() model.Card {
	return model.Card{
		Number:          "1021010",
		Type:            "伙伴",
		Name:            "专精法师",
		Category:        "无",
		Tag:             "巫师",
		Description:     "入场:选择任意1个属性,此卡的负载变为该种属性",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 4},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021010.jpg",
	}
}

type CardDef1021011 struct{}

func (CardDef1021011) ID() string      { return "1021011" }
func (CardDef1021011) Name() string    { return "屠魔者杀手" }
func (CardDef1021011) Kind() string    { return "伙伴" }
func (CardDef1021011) Element() string { return "无" }

func (CardDef1021011) Card() model.Card {
	return model.Card{
		Number:          "1021011",
		Type:            "伙伴",
		Name:            "屠魔者杀手",
		Category:        "无",
		Tag:             "人类",
		Description:     "速攻.",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021011.jpg",
	}
}

type CardDef1021012 struct{}

func (CardDef1021012) ID() string      { return "1021012" }
func (CardDef1021012) Name() string    { return "黑市商贩" }
func (CardDef1021012) Kind() string    { return "伙伴" }
func (CardDef1021012) Element() string { return "无" }

func (CardDef1021012) Card() model.Card {
	return model.Card{
		Number:          "1021012",
		Type:            "伙伴",
		Name:            "黑市商贩",
		Category:        "无",
		Tag:             "人类",
		Description:     "主动绝技:从你的手牌或者装备区弃置1张道具牌才能发动,抽2张牌",
		Quote:           "\"巫师老爷,您手里那些东西在我们这可都是宝\"",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021012.jpg",
	}
}

type CardDef1021013 struct{}

func (CardDef1021013) ID() string      { return "1021013" }
func (CardDef1021013) Name() string    { return "屠魔者武士" }
func (CardDef1021013) Kind() string    { return "伙伴" }
func (CardDef1021013) Element() string { return "无" }

func (CardDef1021013) Card() model.Card {
	return model.Card{
		Number:          "1021013",
		Type:            "伙伴",
		Name:            "屠魔者武士",
		Category:        "无",
		Tag:             "人类",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021013.jpg",
	}
}

type CardDef1021014 struct{}

func (CardDef1021014) ID() string      { return "1021014" }
func (CardDef1021014) Name() string    { return "急不可耐的小师弟" }
func (CardDef1021014) Kind() string    { return "伙伴" }
func (CardDef1021014) Element() string { return "无" }

func (CardDef1021014) Card() model.Card {
	return model.Card{
		Number:          "1021014",
		Type:            "伙伴",
		Name:            "急不可耐的小师弟",
		Category:        "无",
		Tag:             "人类",
		Description:     "入场:本回合你学习的下一个技能获得\"速攻\"",
		Quote:           "来不及了,就这个吧!",
		ElementsCost:    map[string]int{"无": 1, "暗": 1},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021014.jpg",
	}
}

type CardDef1021015 struct{}

func (CardDef1021015) ID() string      { return "1021015" }
func (CardDef1021015) Name() string    { return "精力充沛的大师兄" }
func (CardDef1021015) Kind() string    { return "伙伴" }
func (CardDef1021015) Element() string { return "无" }

func (CardDef1021015) Card() model.Card {
	return model.Card{
		Number:          "1021015",
		Type:            "伙伴",
		Name:            "精力充沛的大师兄",
		Category:        "无",
		Tag:             "巫师",
		Description:     "入场:本回合你施放的下一个技能不需要冷却",
		Quote:           "没想到吧,再来一次!",
		ElementsCost:    map[string]int{"光": 1, "无": 2},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021015.jpg",
	}
}

type CardDef1021016 struct{}

func (CardDef1021016) ID() string      { return "1021016" }
func (CardDef1021016) Name() string    { return "奥术盔甲匠" }
func (CardDef1021016) Kind() string    { return "伙伴" }
func (CardDef1021016) Element() string { return "无" }

func (CardDef1021016) Card() model.Card {
	return model.Card{
		Number:          "1021016",
		Type:            "伙伴",
		Name:            "奥术盔甲匠",
		Category:        "无",
		Tag:             "人类",
		Description:     "入场:如果你没有装备,检索1个装备道具",
		Quote:           "\"不不不,您这套行头出去可不像样啊\"",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021016.jpg",
	}
}

type CardDef1021017 struct{}

func (CardDef1021017) ID() string      { return "1021017" }
func (CardDef1021017) Name() string    { return "符文师" }
func (CardDef1021017) Kind() string    { return "伙伴" }
func (CardDef1021017) Element() string { return "无" }

func (CardDef1021017) Card() model.Card {
	return model.Card{
		Number:          "1021017",
		Type:            "伙伴",
		Name:            "符文师",
		Category:        "无",
		Tag:             "巫师",
		Description:     "入场:丢弃1张手牌才能发动,检索1个符文或卷轴",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021017.jpg",
	}
}

type CardDef1021018 struct{}

func (CardDef1021018) ID() string      { return "1021018" }
func (CardDef1021018) Name() string    { return "奥术壁垒" }
func (CardDef1021018) Kind() string    { return "伙伴" }
func (CardDef1021018) Element() string { return "无" }

func (CardDef1021018) Card() model.Card {
	return model.Card{
		Number:          "1021018",
		Type:            "伙伴",
		Name:            "奥术壁垒",
		Category:        "无",
		Tag:             "造物",
		Description:     "遗言:对方获得2\\无",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            5,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\无\\1021018.jpg",
	}
}

type CardDef1021101 struct{}

func (CardDef1021101) ID() string      { return "1021101" }
func (CardDef1021101) Name() string    { return "私家教师" }
func (CardDef1021101) Kind() string    { return "伙伴" }
func (CardDef1021101) Element() string { return "无" }

func (CardDef1021101) Card() model.Card {
	return model.Card{
		Number:          "1021101",
		Type:            "伙伴",
		Name:            "私家教师",
		Category:        "无",
		Tag:             "巫师",
		Description:     "入场:学习1个学习花费小于4的法术,无需花费",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3, "水": 1},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021101.jpg",
	}
}

type CardDef1021102 struct{}

func (CardDef1021102) ID() string      { return "1021102" }
func (CardDef1021102) Name() string    { return "剑术师傅" }
func (CardDef1021102) Kind() string    { return "伙伴" }
func (CardDef1021102) Element() string { return "无" }

func (CardDef1021102) Card() model.Card {
	return model.Card{
		Number:          "1021102",
		Type:            "伙伴",
		Name:            "剑术师傅",
		Category:        "无",
		Tag:             "人类",
		Description:     "入场:使1个相邻友方伙伴获得+1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2, "火": 1},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021102.jpg",
	}
}

type CardDef1021103 struct{}

func (CardDef1021103) ID() string      { return "1021103" }
func (CardDef1021103) Name() string    { return "皇城结界兽" }
func (CardDef1021103) Kind() string    { return "伙伴" }
func (CardDef1021103) Element() string { return "无" }

func (CardDef1021103) Card() model.Card {
	return model.Card{
		Number:          "1021103",
		Type:            "伙伴",
		Name:            "皇城结界兽",
		Category:        "无",
		Tag:             "异兽",
		Description:     "入场:获得护盾2",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021103.jpg",
	}
}

type CardDef1021104 struct{}

func (CardDef1021104) ID() string      { return "1021104" }
func (CardDef1021104) Name() string    { return "次元撕裂兽" }
func (CardDef1021104) Kind() string    { return "伙伴" }
func (CardDef1021104) Element() string { return "无" }

func (CardDef1021104) Card() model.Card {
	return model.Card{
		Number:          "1021104",
		Type:            "伙伴",
		Name:            "次元撕裂兽",
		Category:        "无",
		Tag:             "异兽",
		Description:     "入场:将法力范围内的1个敌方伙伴移出游戏",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "无": 5, "暗": 1},
		ElementsGain:    map[string]int{"光": 1, "暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021104.jpg",
	}
}

type CardDef1021105 struct{}

func (CardDef1021105) ID() string      { return "1021105" }
func (CardDef1021105) Name() string    { return "皇城征税员" }
func (CardDef1021105) Kind() string    { return "伙伴" }
func (CardDef1021105) Element() string { return "无" }

func (CardDef1021105) Card() model.Card {
	return model.Card{
		Number:          "1021105",
		Type:            "伙伴",
		Name:            "皇城征税员",
		Category:        "无",
		Tag:             "人类",
		Description:     "入场:直到对手的下个回合结束,对手每抽1张牌,你获得1\\无",
		Quote:           "内政大臣们正在考虑取消税负政策,具体做法是:将税金更名为\"义务奉献\"",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021105.jpg",
	}
}

type CardDef1021106 struct{}

func (CardDef1021106) ID() string      { return "1021106" }
func (CardDef1021106) Name() string    { return "云霄城富豪" }
func (CardDef1021106) Kind() string    { return "伙伴" }
func (CardDef1021106) Element() string { return "无" }

func (CardDef1021106) Card() model.Card {
	return model.Card{
		Number:          "1021106",
		Type:            "伙伴",
		Name:            "云霄城富豪",
		Category:        "无",
		Tag:             "人类",
		Description:     "主动:消耗此卡才能发动,双方各抽1张牌(次序由你选择)",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "气": 1},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021106.jpg",
	}
}

type CardDef1021107 struct{}

func (CardDef1021107) ID() string      { return "1021107" }
func (CardDef1021107) Name() string    { return "天才少年" }
func (CardDef1021107) Kind() string    { return "伙伴" }
func (CardDef1021107) Element() string { return "无" }

func (CardDef1021107) Card() model.Card {
	return model.Card{
		Number:          "1021107",
		Type:            "伙伴",
		Name:            "天才少年",
		Category:        "无",
		Tag:             "人类",
		Description:     "精通2:此卡获得任意1点奥术元素以外的负载",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021107.jpg",
	}
}

type CardDef1021108 struct{}

func (CardDef1021108) ID() string      { return "1021108" }
func (CardDef1021108) Name() string    { return "炼金术学徒" }
func (CardDef1021108) Kind() string    { return "伙伴" }
func (CardDef1021108) Element() string { return "无" }

func (CardDef1021108) Card() model.Card {
	return model.Card{
		Number:          "1021108",
		Type:            "伙伴",
		Name:            "炼金术学徒",
		Category:        "无",
		Tag:             "巫师",
		Description:     "主动:消耗此卡才能发动,将1\\无转化为任意2点奥术元素以外的元素",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021108.jpg",
	}
}

type CardDef1021109 struct{}

func (CardDef1021109) ID() string      { return "1021109" }
func (CardDef1021109) Name() string    { return "教廷特使" }
func (CardDef1021109) Kind() string    { return "伙伴" }
func (CardDef1021109) Element() string { return "无" }

func (CardDef1021109) Card() model.Card {
	return model.Card{
		Number:          "1021109",
		Type:            "伙伴",
		Name:            "教廷特使",
		Category:        "无",
		Tag:             "人类",
		Description:     "主动绝技:移除1张友方卡牌全部负面效果",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "无": 1},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021109.jpg",
	}
}

type CardDef1021110 struct{}

func (CardDef1021110) ID() string      { return "1021110" }
func (CardDef1021110) Name() string    { return "岩壁护卫军" }
func (CardDef1021110) Kind() string    { return "伙伴" }
func (CardDef1021110) Element() string { return "无" }

func (CardDef1021110) Card() model.Card {
	return model.Card{
		Number:          "1021110",
		Type:            "伙伴",
		Name:            "岩壁护卫军",
		Category:        "无",
		Tag:             "造物",
		Description:     "引魔.诱发绝技:当敌方法术命中且你没有护盾时,获得护盾2",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021110.jpg",
	}
}

type CardDef1021111 struct{}

func (CardDef1021111) ID() string      { return "1021111" }
func (CardDef1021111) Name() string    { return "孤星勇者" }
func (CardDef1021111) Kind() string    { return "伙伴" }
func (CardDef1021111) Element() string { return "无" }

func (CardDef1021111) Card() model.Card {
	return model.Card{
		Number:          "1021111",
		Type:            "伙伴",
		Name:            "孤星勇者",
		Category:        "无",
		Tag:             "人类",
		Description:     "你每有1张其他手牌,此卡的入场花费+1\\无",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021111.jpg",
	}
}

type CardDef1021112 struct{}

func (CardDef1021112) ID() string      { return "1021112" }
func (CardDef1021112) Name() string    { return "奥术纯净体" }
func (CardDef1021112) Kind() string    { return "伙伴" }
func (CardDef1021112) Element() string { return "无" }

func (CardDef1021112) Card() model.Card {
	return model.Card{
		Number:          "1021112",
		Type:            "伙伴",
		Name:            "奥术纯净体",
		Category:        "无",
		Tag:             "造物",
		Description:     "此卡的入场花费必须严格为奥术元素",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 4},
		ElementsGain:    map[string]int{"无": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021112.jpg",
	}
}

type CardDef1021113 struct{}

func (CardDef1021113) ID() string      { return "1021113" }
func (CardDef1021113) Name() string    { return "魔法飞蛾" }
func (CardDef1021113) Kind() string    { return "伙伴" }
func (CardDef1021113) Element() string { return "无" }

func (CardDef1021113) Card() model.Card {
	return model.Card{
		Number:          "1021113",
		Type:            "伙伴",
		Name:            "魔法飞蛾",
		Category:        "无",
		Tag:             "野兽",
		Description:     "诱发:在你施放一个聚能技能后,可以从卡组抽取本卡",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "地": 1, "气": 1, "水": 1, "火": 1},
		ElementsGain:    map[string]int{"无": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021113.jpg",
	}
}

type CardDef1021114 struct{}

func (CardDef1021114) ID() string      { return "1021114" }
func (CardDef1021114) Name() string    { return "蛰蛙" }
func (CardDef1021114) Kind() string    { return "伙伴" }
func (CardDef1021114) Element() string { return "无" }

func (CardDef1021114) Card() model.Card {
	return model.Card{
		Number:          "1021114",
		Type:            "伙伴",
		Name:            "蛰蛙",
		Category:        "无",
		Tag:             "野兽",
		Description:     "诱发:每当你使用一个驱动或聚能技能后,本回合你的法术+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 4, "气": 1},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021114.jpg",
	}
}

type CardDef1021115 struct{}

func (CardDef1021115) ID() string      { return "1021115" }
func (CardDef1021115) Name() string    { return "九霄刺客" }
func (CardDef1021115) Kind() string    { return "伙伴" }
func (CardDef1021115) Element() string { return "无" }

func (CardDef1021115) Card() model.Card {
	return model.Card{
		Number:          "1021115",
		Type:            "伙伴",
		Name:            "九霄刺客",
		Category:        "无",
		Tag:             "人类",
		Description:     "入场:将1张九霄印记加入对手手牌.遗言:将4张九霄印记洗入对手卡组",
		Quote:           "\"想看看云霄城的夜色吗,可能是最后一次了哦\"",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2001102"},
		OutputPath:      "output\\王权纷争\\伙伴\\无\\1021115.jpg",
	}
}

type CardDef1111001 struct{}

func (CardDef1111001) ID() string      { return "1111001" }
func (CardDef1111001) Name() string    { return "火龙 \"辉煌\"" }
func (CardDef1111001) Kind() string    { return "伙伴" }
func (CardDef1111001) Element() string { return "火" }

func (CardDef1111001) Card() model.Card {
	return model.Card{
		Number:          "1111001",
		Type:            "伙伴",
		Name:            "火龙 \"辉煌\"",
		Category:        "火",
		Tag:             "传奇-龙",
		Description:     "吞噬:3\\火.引魔.绑定技能:火焰吐息",
		Quote:           "不可一世的巨龙,为了残存的火焰元素,向巫师们尽忠",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{"火": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"3101001"},
		OutputPath:      "output\\基础包\\伙伴\\火\\1111001.jpg",
	}
}

type CardDef1111002 struct{}

func (CardDef1111002) ID() string      { return "1111002" }
func (CardDef1111002) Name() string    { return "炎狱大将军 狄斯托德" }
func (CardDef1111002) Kind() string    { return "伙伴" }
func (CardDef1111002) Element() string { return "火" }

func (CardDef1111002) Card() model.Card {
	return model.Card{
		Number:          "1111002",
		Type:            "伙伴",
		Name:            "炎狱大将军 狄斯托德",
		Category:        "火",
		Tag:             "传奇-巫师",
		Description:     "诱发:每当对方召唤1个伙伴时,使其获得点燃2和石化2",
		Quote:           "\"我的使命,就是将地狱带到人间\"",
		ElementsCost:    map[string]int{"地": 2, "火": 6},
		ElementsGain:    map[string]int{"地": 1, "火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1111002.jpg",
	}
}

type CardDef1111003 struct{}

func (CardDef1111003) ID() string      { return "1111003" }
func (CardDef1111003) Name() string    { return "毕方" }
func (CardDef1111003) Kind() string    { return "伙伴" }
func (CardDef1111003) Element() string { return "火" }

func (CardDef1111003) Card() model.Card {
	return model.Card{
		Number:          "1111003",
		Type:            "伙伴",
		Name:            "毕方",
		Category:        "火",
		Tag:             "传奇-异兽",
		Description:     "引魔.光环:敌方单位受到的点燃伤害+1",
		Quote:           "始于灰烬,终于灰烬",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1111003.jpg",
	}
}

type CardDef1111101 struct{}

func (CardDef1111101) ID() string      { return "1111101" }
func (CardDef1111101) Name() string    { return "无上女王 黛琳 凯尔特" }
func (CardDef1111101) Kind() string    { return "伙伴" }
func (CardDef1111101) Element() string { return "火" }

func (CardDef1111101) Card() model.Card {
	return model.Card{
		Number:          "1111101",
		Type:            "伙伴",
		Name:            "无上女王 黛琳 凯尔特",
		Category:        "火",
		Tag:             "传奇-巫师",
		Description:     "入场:从手牌召唤任意数量火焰伙伴到此卡相邻位置,无需花费,直到下个回合结束此卡和此卡召唤的卡牌免疫所有伤害和负面效果",
		Quote:           "我心即是无上意志,我行即是万物准则",
		ElementsCost:    map[string]int{"无": 3, "火": 9},
		ElementsGain:    map[string]int{"无": 1, "火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1111101.jpg",
	}
}

type CardDef1111102 struct{}

func (CardDef1111102) ID() string      { return "1111102" }
func (CardDef1111102) Name() string    { return "\"流放者\" 索拓尔" }
func (CardDef1111102) Kind() string    { return "伙伴" }
func (CardDef1111102) Element() string { return "火" }

func (CardDef1111102) Card() model.Card {
	return model.Card{
		Number:          "1111102",
		Type:            "伙伴",
		Name:            "\"流放者\" 索拓尔",
		Category:        "火",
		Tag:             "传奇-巫师",
		Description:     "光环:你的所有法术额外选择与原本目标相邻的所有单位为目标",
		Quote:           "传说凯尔特草原曾是繁茂的森林,直到索拓尔用漫天的火焰将其焚烧殆尽,尔后万物复苏,草原遂成",
		ElementsCost:    map[string]int{"无": 1, "气": 2, "火": 4},
		ElementsGain:    map[string]int{"气": 1, "火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1111102.jpg",
	}
}

type CardDef1111103 struct{}

func (CardDef1111103) ID() string      { return "1111103" }
func (CardDef1111103) Name() string    { return "贪婪暴君 卡姆 弗卡莱诺" }
func (CardDef1111103) Kind() string    { return "伙伴" }
func (CardDef1111103) Element() string { return "火" }

func (CardDef1111103) Card() model.Card {
	return model.Card{
		Number:          "1111103",
		Type:            "伙伴",
		Name:            "贪婪暴君 卡姆 弗卡莱诺",
		Category:        "火",
		Tag:             "传奇-人类",
		Description:     "引魔.光环:双方手牌里的卡牌入场花费+1\\无",
		Quote:           "至少,他境内的岩浆比谁都多",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1111103.jpg",
	}
}

type CardDef1121001 struct{}

func (CardDef1121001) ID() string      { return "1121001" }
func (CardDef1121001) Name() string    { return "火焰精灵" }
func (CardDef1121001) Kind() string    { return "伙伴" }
func (CardDef1121001) Element() string { return "火" }

func (CardDef1121001) Card() model.Card {
	return model.Card{
		Number:          "1121001",
		Type:            "伙伴",
		Name:            "火焰精灵",
		Category:        "火",
		Tag:             "精灵",
		Description:     "诱发:每当此卡被消耗时,获得点燃1",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121001.jpg",
	}
}

type CardDef1121002 struct{}

func (CardDef1121002) ID() string      { return "1121002" }
func (CardDef1121002) Name() string    { return "活泼的炉火" }
func (CardDef1121002) Kind() string    { return "伙伴" }
func (CardDef1121002) Element() string { return "火" }

func (CardDef1121002) Card() model.Card {
	return model.Card{
		Number:          "1121002",
		Type:            "伙伴",
		Name:            "活泼的炉火",
		Category:        "火",
		Tag:             "造物",
		Description:     "入场:抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "火": 1},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121002.jpg",
	}
}

type CardDef1121003 struct{}

func (CardDef1121003) ID() string      { return "1121003" }
func (CardDef1121003) Name() string    { return "锻石工匠" }
func (CardDef1121003) Kind() string    { return "伙伴" }
func (CardDef1121003) Element() string { return "火" }

func (CardDef1121003) Card() model.Card {
	return model.Card{
		Number:          "1121003",
		Type:            "伙伴",
		Name:            "锻石工匠",
		Category:        "火",
		Tag:             "人类",
		Description:     "主动:消耗此卡才能发动,使你的1个法术在本回合+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121003.jpg",
	}
}

type CardDef1121004 struct{}

func (CardDef1121004) ID() string      { return "1121004" }
func (CardDef1121004) Name() string    { return "凯尔特雄狮" }
func (CardDef1121004) Kind() string    { return "伙伴" }
func (CardDef1121004) Element() string { return "火" }

func (CardDef1121004) Card() model.Card {
	return model.Card{
		Number:          "1121004",
		Type:            "伙伴",
		Name:            "凯尔特雄狮",
		Category:        "火",
		Tag:             "野兽",
		Description:     "光环:你的所有法术+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "火": 4},
		ElementsGain:    map[string]int{"地": 1, "火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121004.jpg",
	}
}

type CardDef1121005 struct{}

func (CardDef1121005) ID() string      { return "1121005" }
func (CardDef1121005) Name() string    { return "熔岩傀儡" }
func (CardDef1121005) Kind() string    { return "伙伴" }
func (CardDef1121005) Element() string { return "火" }

func (CardDef1121005) Card() model.Card {
	return model.Card{
		Number:          "1121005",
		Type:            "伙伴",
		Name:            "熔岩傀儡",
		Category:        "火",
		Tag:             "恶魔",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{"地": 1, "火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121005.jpg",
	}
}

type CardDef1121006 struct{}

func (CardDef1121006) ID() string      { return "1121006" }
func (CardDef1121006) Name() string    { return "熔岩烽蛇" }
func (CardDef1121006) Kind() string    { return "伙伴" }
func (CardDef1121006) Element() string { return "火" }

func (CardDef1121006) Card() model.Card {
	return model.Card{
		Number:          "1121006",
		Type:            "伙伴",
		Name:            "熔岩烽蛇",
		Category:        "火",
		Tag:             "野兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{"火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121006.jpg",
	}
}

type CardDef1121007 struct{}

func (CardDef1121007) ID() string      { return "1121007" }
func (CardDef1121007) Name() string    { return "至纯之火" }
func (CardDef1121007) Kind() string    { return "伙伴" }
func (CardDef1121007) Element() string { return "火" }

func (CardDef1121007) Card() model.Card {
	return model.Card{
		Number:          "1121007",
		Type:            "伙伴",
		Name:            "至纯之火",
		Category:        "火",
		Tag:             "造物",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{"火": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121007.jpg",
	}
}

type CardDef1121008 struct{}

func (CardDef1121008) ID() string      { return "1121008" }
func (CardDef1121008) Name() string    { return "炎狱使者" }
func (CardDef1121008) Kind() string    { return "伙伴" }
func (CardDef1121008) Element() string { return "火" }

func (CardDef1121008) Card() model.Card {
	return model.Card{
		Number:          "1121008",
		Type:            "伙伴",
		Name:            "炎狱使者",
		Category:        "火",
		Tag:             "恶魔",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "火": 6},
		ElementsGain:    map[string]int{"地": 1, "火": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121008.jpg",
	}
}

type CardDef1121009 struct{}

func (CardDef1121009) ID() string      { return "1121009" }
func (CardDef1121009) Name() string    { return "赤鹰" }
func (CardDef1121009) Kind() string    { return "伙伴" }
func (CardDef1121009) Element() string { return "火" }

func (CardDef1121009) Card() model.Card {
	return model.Card{
		Number:          "1121009",
		Type:            "伙伴",
		Name:            "赤鹰",
		Category:        "火",
		Tag:             "野兽",
		Description:     "入场:检索1个入场花费大于等于4的火焰伙伴",
		Quote:           "每只赤鹰都想成为凤凰,前提是它们身上的羽毛还没有被猎人拔光",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121009.jpg",
	}
}

type CardDef1121010 struct{}

func (CardDef1121010) ID() string      { return "1121010" }
func (CardDef1121010) Name() string    { return "火焰艺人" }
func (CardDef1121010) Kind() string    { return "伙伴" }
func (CardDef1121010) Element() string { return "火" }

func (CardDef1121010) Card() model.Card {
	return model.Card{
		Number:          "1121010",
		Type:            "伙伴",
		Name:            "火焰艺人",
		Category:        "火",
		Tag:             "巫师",
		Description:     "主动绝技:消耗此卡才能发动,重置你的另1张人物牌以外的火焰牌",
		Quote:           "\"红色的火焰?那也太低级了!\"",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121010.jpg",
	}
}

type CardDef1121011 struct{}

func (CardDef1121011) ID() string      { return "1121011" }
func (CardDef1121011) Name() string    { return "火山飞龙" }
func (CardDef1121011) Kind() string    { return "伙伴" }
func (CardDef1121011) Element() string { return "火" }

func (CardDef1121011) Card() model.Card {
	return model.Card{
		Number:          "1121011",
		Type:            "伙伴",
		Name:            "火山飞龙",
		Category:        "火",
		Tag:             "龙",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "火": 6},
		ElementsGain:    map[string]int{"地": 1, "气": 1, "火": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121011.jpg",
	}
}

type CardDef1121012 struct{}

func (CardDef1121012) ID() string      { return "1121012" }
func (CardDef1121012) Name() string    { return "火焰洞察者" }
func (CardDef1121012) Kind() string    { return "伙伴" }
func (CardDef1121012) Element() string { return "火" }

func (CardDef1121012) Card() model.Card {
	return model.Card{
		Number:          "1121012",
		Type:            "伙伴",
		Name:            "火焰洞察者",
		Category:        "火",
		Tag:             "巫师",
		Description:     "诱发回合技:若有单位受到火焰伤害,抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121012.jpg",
	}
}

type CardDef1121013 struct{}

func (CardDef1121013) ID() string      { return "1121013" }
func (CardDef1121013) Name() string    { return "纵火者" }
func (CardDef1121013) Kind() string    { return "伙伴" }
func (CardDef1121013) Element() string { return "火" }

func (CardDef1121013) Card() model.Card {
	return model.Card{
		Number:          "1121013",
		Type:            "伙伴",
		Name:            "纵火者",
		Category:        "火",
		Tag:             "巫师",
		Description:     "诱发回合技:在你使用1个火焰法术后,可以使法力范围内的1个单位点燃1",
		Quote:           "谁能拒绝将一切燃烧殆尽的快乐呢",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{"火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121013.jpg",
	}
}

type CardDef1121014 struct{}

func (CardDef1121014) ID() string      { return "1121014" }
func (CardDef1121014) Name() string    { return "火荆" }
func (CardDef1121014) Kind() string    { return "伙伴" }
func (CardDef1121014) Element() string { return "火" }

func (CardDef1121014) Card() model.Card {
	return model.Card{
		Number:          "1121014",
		Type:            "伙伴",
		Name:            "火荆",
		Category:        "火",
		Tag:             "植物",
		Description:     "遗言:使法力范围内的1个敌人点燃1",
		Quote:           "穿过如火荆棘,方见白洁花簇",
		ElementsCost:    map[string]int{"火": 1},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121014.jpg",
	}
}

type CardDef1121015 struct{}

func (CardDef1121015) ID() string      { return "1121015" }
func (CardDef1121015) Name() string    { return "火云法师" }
func (CardDef1121015) Kind() string    { return "伙伴" }
func (CardDef1121015) Element() string { return "火" }

func (CardDef1121015) Card() model.Card {
	return model.Card{
		Number:          "1121015",
		Type:            "伙伴",
		Name:            "火云法师",
		Category:        "火",
		Tag:             "巫师",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{"气": 1, "火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121015.jpg",
	}
}

type CardDef1121016 struct{}

func (CardDef1121016) ID() string      { return "1121016" }
func (CardDef1121016) Name() string    { return "舞火者" }
func (CardDef1121016) Kind() string    { return "伙伴" }
func (CardDef1121016) Element() string { return "火" }

func (CardDef1121016) Card() model.Card {
	return model.Card{
		Number:          "1121016",
		Type:            "伙伴",
		Name:            "舞火者",
		Category:        "火",
		Tag:             "人类",
		Description:     "入场,遗言:使你场上的所有火焰卡牌免疫负面状态(仍可处于)直到你的下一次回合结束",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\火\\1121016.jpg",
	}
}

type CardDef1121101 struct{}

func (CardDef1121101) ID() string      { return "1121101" }
func (CardDef1121101) Name() string    { return "火山蝾螈" }
func (CardDef1121101) Kind() string    { return "伙伴" }
func (CardDef1121101) Element() string { return "火" }

func (CardDef1121101) Card() model.Card {
	return model.Card{
		Number:          "1121101",
		Type:            "伙伴",
		Name:            "火山蝾螈",
		Category:        "火",
		Tag:             "野兽",
		Description:     "精通2:献祭此卡并从手牌召唤1个入场费用小于8的火焰伙伴,无需花费.",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121101.jpg",
	}
}

type CardDef1121102 struct{}

func (CardDef1121102) ID() string      { return "1121102" }
func (CardDef1121102) Name() string    { return "火山谷底巨兽" }
func (CardDef1121102) Kind() string    { return "伙伴" }
func (CardDef1121102) Element() string { return "火" }

func (CardDef1121102) Card() model.Card {
	return model.Card{
		Number:          "1121102",
		Type:            "伙伴",
		Name:            "火山谷底巨兽",
		Category:        "火",
		Tag:             "异兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 6},
		ElementsGain:    map[string]int{"地": 1, "火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            5,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121102.jpg",
	}
}

type CardDef1121103 struct{}

func (CardDef1121103) ID() string      { return "1121103" }
func (CardDef1121103) Name() string    { return "烽火台守卫" }
func (CardDef1121103) Kind() string    { return "伙伴" }
func (CardDef1121103) Element() string { return "火" }

func (CardDef1121103) Card() model.Card {
	return model.Card{
		Number:          "1121103",
		Type:            "伙伴",
		Name:            "烽火台守卫",
		Category:        "火",
		Tag:             "人类",
		Description:     "入场:如果你场上的单位数量少于敌方,获得护盾3",
		Quote:           "时刻保持警戒!",
		ElementsCost:    map[string]int{"无": 1, "火": 1},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121103.jpg",
	}
}

type CardDef1121104 struct{}

func (CardDef1121104) ID() string      { return "1121104" }
func (CardDef1121104) Name() string    { return "熔岩堡战车" }
func (CardDef1121104) Kind() string    { return "伙伴" }
func (CardDef1121104) Element() string { return "火" }

func (CardDef1121104) Card() model.Card {
	return model.Card{
		Number:          "1121104",
		Type:            "伙伴",
		Name:            "熔岩堡战车",
		Category:        "火",
		Tag:             "机械",
		Description:     "此卡攻击附加点燃1",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121104.jpg",
	}
}

type CardDef1121105 struct{}

func (CardDef1121105) ID() string      { return "1121105" }
func (CardDef1121105) Name() string    { return "弗卡莱诺近卫" }
func (CardDef1121105) Kind() string    { return "伙伴" }
func (CardDef1121105) Element() string { return "火" }

func (CardDef1121105) Card() model.Card {
	return model.Card{
		Number:          "1121105",
		Type:            "伙伴",
		Name:            "弗卡莱诺近卫",
		Category:        "火",
		Tag:             "人类",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{"火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121105.jpg",
	}
}

type CardDef1121106 struct{}

func (CardDef1121106) ID() string      { return "1121106" }
func (CardDef1121106) Name() string    { return "弗卡莱诺皇家驯兽师" }
func (CardDef1121106) Kind() string    { return "伙伴" }
func (CardDef1121106) Element() string { return "火" }

func (CardDef1121106) Card() model.Card {
	return model.Card{
		Number:          "1121106",
		Type:            "伙伴",
		Name:            "弗卡莱诺皇家驯兽师",
		Category:        "火",
		Tag:             "人类",
		Description:     "入场:你的下一个野兽或异兽火焰伙伴花费-2",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121106.jpg",
	}
}

type CardDef1121107 struct{}

func (CardDef1121107) ID() string      { return "1121107" }
func (CardDef1121107) Name() string    { return "失控的神火兽" }
func (CardDef1121107) Kind() string    { return "伙伴" }
func (CardDef1121107) Element() string { return "火" }

func (CardDef1121107) Card() model.Card {
	return model.Card{
		Number:          "1121107",
		Type:            "伙伴",
		Name:            "失控的神火兽",
		Category:        "火",
		Tag:             "异兽",
		Description:     "光环:双方法术在攻击时+2\\威",
		Quote:           "通常神火军团不会雨天出动,不是因为害怕雨水,而是害怕闪电让他们的坐骑彻底失控",
		ElementsCost:    map[string]int{"气": 1, "火": 4},
		ElementsGain:    map[string]int{"气": 1, "火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121107.jpg",
	}
}

type CardDef1121108 struct{}

func (CardDef1121108) ID() string      { return "1121108" }
func (CardDef1121108) Name() string    { return "火蝴蝶" }
func (CardDef1121108) Kind() string    { return "伙伴" }
func (CardDef1121108) Element() string { return "火" }

func (CardDef1121108) Card() model.Card {
	return model.Card{
		Number:          "1121108",
		Type:            "伙伴",
		Name:            "火蝴蝶",
		Category:        "火",
		Tag:             "野兽",
		Description:     "主动回合技:负载临时改为1\\气",
		Quote:           "想在熔岩堡这么恶劣的环境下生存,各种生灵也得有些自己的绝活",
		ElementsCost:    map[string]int{"火": 1},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121108.jpg",
	}
}

type CardDef1121109 struct{}

func (CardDef1121109) ID() string      { return "1121109" }
func (CardDef1121109) Name() string    { return "神火兽骑手" }
func (CardDef1121109) Kind() string    { return "伙伴" }
func (CardDef1121109) Element() string { return "火" }

func (CardDef1121109) Card() model.Card {
	return model.Card{
		Number:          "1121109",
		Type:            "伙伴",
		Name:            "神火兽骑手",
		Category:        "火",
		Tag:             "人类",
		Description:     "主动绝技:消耗1个友方其他火焰伙伴才能发动,下一次你的火焰法术\\威上升其入场花费元素的数值",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 5},
		ElementsGain:    map[string]int{"气": 1, "火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121109.jpg",
	}
}

type CardDef1121110 struct{}

func (CardDef1121110) ID() string      { return "1121110" }
func (CardDef1121110) Name() string    { return "熔岩堡档案员" }
func (CardDef1121110) Kind() string    { return "伙伴" }
func (CardDef1121110) Element() string { return "火" }

func (CardDef1121110) Card() model.Card {
	return model.Card{
		Number:          "1121110",
		Type:            "伙伴",
		Name:            "熔岩堡档案员",
		Category:        "火",
		Tag:             "巫师",
		Description:     "诱发绝技:在你使用1个创造法术后,翻取1个卷轴或符文",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121110.jpg",
	}
}

type CardDef1121111 struct{}

func (CardDef1121111) ID() string      { return "1121111" }
func (CardDef1121111) Name() string    { return "孤星火种" }
func (CardDef1121111) Kind() string    { return "伙伴" }
func (CardDef1121111) Element() string { return "火" }

func (CardDef1121111) Card() model.Card {
	return model.Card{
		Number:          "1121111",
		Type:            "伙伴",
		Name:            "孤星火种",
		Category:        "火",
		Tag:             "精灵",
		Description:     "诱发绝技:在场上其他伙伴牌受到火焰伤害后,此卡获得负载+1\\火",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121111.jpg",
	}
}

type CardDef1121112 struct{}

func (CardDef1121112) ID() string      { return "1121112" }
func (CardDef1121112) Name() string    { return "火花飞蛾" }
func (CardDef1121112) Kind() string    { return "伙伴" }
func (CardDef1121112) Element() string { return "火" }

func (CardDef1121112) Card() model.Card {
	return model.Card{
		Number:          "1121112",
		Type:            "伙伴",
		Name:            "火花飞蛾",
		Category:        "火",
		Tag:             "野兽",
		Description:     "诱发:在一个火焰法术命中后,可以展示手牌中的此卡并-1入场花费",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121112.jpg",
	}
}

type CardDef1121113 struct{}

func (CardDef1121113) ID() string      { return "1121113" }
func (CardDef1121113) Name() string    { return "熔岩堡地狱犬" }
func (CardDef1121113) Kind() string    { return "伙伴" }
func (CardDef1121113) Element() string { return "火" }

func (CardDef1121113) Card() model.Card {
	return model.Card{
		Number:          "1121113",
		Type:            "伙伴",
		Name:            "熔岩堡地狱犬",
		Category:        "火",
		Tag:             "异兽",
		Description:     "诱发回合技:此卡被其他卡牌效果消耗后,必须选择法力范围内2个不同单位造成1点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 5},
		ElementsGain:    map[string]int{"火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121113.jpg",
	}
}

type CardDef1121114 struct{}

func (CardDef1121114) ID() string      { return "1121114" }
func (CardDef1121114) Name() string    { return "军团将星" }
func (CardDef1121114) Kind() string    { return "伙伴" }
func (CardDef1121114) Element() string { return "火" }

func (CardDef1121114) Card() model.Card {
	return model.Card{
		Number:          "1121114",
		Type:            "伙伴",
		Name:            "军团将星",
		Category:        "火",
		Tag:             "人类",
		Description:     "引魔.祈咒:直到下个回合结束,你的所有火焰法术获得+2\\威或+1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "气": 1, "火": 6},
		ElementsGain:    map[string]int{"无": 1, "火": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            5,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121114.jpg",
	}
}

type CardDef1121115 struct{}

func (CardDef1121115) ID() string      { return "1121115" }
func (CardDef1121115) Name() string    { return "军团参谋" }
func (CardDef1121115) Kind() string    { return "伙伴" }
func (CardDef1121115) Element() string { return "火" }

func (CardDef1121115) Card() model.Card {
	return model.Card{
		Number:          "1121115",
		Type:            "伙伴",
		Name:            "军团参谋",
		Category:        "火",
		Tag:             "巫师",
		Description:     "诱发回合技:在你使用1个创造种类的技能后,可以翻取1个火焰属性消耗品道具但在回合结束时丢弃",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{"火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\火\\1121115.jpg",
	}
}

type CardDef1201101 struct{}

func (CardDef1201101) ID() string      { return "1201101" }
func (CardDef1201101) Name() string    { return "凛冰之龙" }
func (CardDef1201101) Kind() string    { return "伙伴" }
func (CardDef1201101) Element() string { return "水" }

func (CardDef1201101) Card() model.Card {
	return model.Card{
		Number:          "1201101",
		Type:            "伙伴",
		Name:            "凛冰之龙",
		Category:        "水",
		Tag:             "衍生-龙",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"水": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1201101.jpg",
	}
}

type CardDef1211001 struct{}

func (CardDef1211001) ID() string      { return "1211001" }
func (CardDef1211001) Name() string    { return "人鱼 菲尔" }
func (CardDef1211001) Kind() string    { return "伙伴" }
func (CardDef1211001) Element() string { return "水" }

func (CardDef1211001) Card() model.Card {
	return model.Card{
		Number:          "1211001",
		Type:            "伙伴",
		Name:            "人鱼 菲尔",
		Category:        "水",
		Tag:             "传奇-异兽",
		Description:     "祈咒:如果此卡相邻没有伙伴,检索1张水纹伙伴",
		Quote:           "合上眼睛,她回忆起曾经的海浪,礁石...以及她信任的那位王子",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1211001.jpg",
	}
}

type CardDef1211002 struct{}

func (CardDef1211002) ID() string      { return "1211002" }
func (CardDef1211002) Name() string    { return "深渊巨口 利维坦" }
func (CardDef1211002) Kind() string    { return "伙伴" }
func (CardDef1211002) Element() string { return "水" }

func (CardDef1211002) Card() model.Card {
	return model.Card{
		Number:          "1211002",
		Type:            "伙伴",
		Name:            "深渊巨口 利维坦",
		Category:        "水",
		Tag:             "传奇-异兽-深渊",
		Description:     "主动:消耗此卡才能发动,消灭法力范围内1个伙伴,下个你的回合不能使用此效果",
		Quote:           "\"只有一人能从巨口中存活,他便是那怪物命中注定的斩杀者\"",
		ElementsCost:    map[string]int{"暗": 2, "水": 4},
		ElementsGain:    map[string]int{"暗": 1, "水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1211002.jpg",
	}
}

type CardDef1211003 struct{}

func (CardDef1211003) ID() string      { return "1211003" }
func (CardDef1211003) Name() string    { return "\"雪女\" 天户凌" }
func (CardDef1211003) Kind() string    { return "伙伴" }
func (CardDef1211003) Element() string { return "水" }

func (CardDef1211003) Card() model.Card {
	return model.Card{
		Number:          "1211003",
		Type:            "伙伴",
		Name:            "\"雪女\" 天户凌",
		Category:        "水",
		Tag:             "传奇-巫师",
		Description:     "引魔.诱发回合技3:在你检索1张水纹卡牌后,选择1个法力范围内的敌人,使其冻结1",
		Quote:           "我梦到北方纷飞的大雪,那里便是我的归宿",
		ElementsCost:    map[string]int{"水": 5},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1211003.jpg",
	}
}

type CardDef1211101 struct{}

func (CardDef1211101) ID() string      { return "1211101" }
func (CardDef1211101) Name() string    { return "雾之国主 那兰提" }
func (CardDef1211101) Kind() string    { return "伙伴" }
func (CardDef1211101) Element() string { return "水" }

func (CardDef1211101) Card() model.Card {
	return model.Card{
		Number:          "1211101",
		Type:            "伙伴",
		Name:            "雾之国主 那兰提",
		Category:        "水",
		Tag:             "传奇-巫师",
		Description:     "入场:所有友方没有隐蔽的单位获得隐蔽2",
		Quote:           "如果你妄想进犯我的国度,我的子民将会让你看到他们在迷雾中的另一面",
		ElementsCost:    map[string]int{"气": 2, "水": 4},
		ElementsGain:    map[string]int{"气": 1, "水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1211101.jpg",
	}
}

type CardDef1211102 struct{}

func (CardDef1211102) ID() string      { return "1211102" }
func (CardDef1211102) Name() string    { return "花海梦鲸" }
func (CardDef1211102) Kind() string    { return "伙伴" }
func (CardDef1211102) Element() string { return "水" }

func (CardDef1211102) Card() model.Card {
	return model.Card{
		Number:          "1211102",
		Type:            "伙伴",
		Name:            "花海梦鲸",
		Category:        "水",
		Tag:             "传奇-野兽",
		Description:     "入场:将3种不同的衍生卡牌幻创之梦洗入你的卡组.诱发:当你累计使用2次创造法术后,检索1张幻创之梦",
		Quote:           "亦真亦幻,如痴如梦",
		ElementsCost:    map[string]int{"水": 5},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2201101", "2201102", "2201103"},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1211102.jpg",
	}
}

type CardDef1211103 struct{}

func (CardDef1211103) ID() string      { return "1211103" }
func (CardDef1211103) Name() string    { return "海上巾帼 珊瑚 雯迪" }
func (CardDef1211103) Kind() string    { return "伙伴" }
func (CardDef1211103) Element() string { return "水" }

func (CardDef1211103) Card() model.Card {
	return model.Card{
		Number:          "1211103",
		Type:            "伙伴",
		Name:            "海上巾帼 珊瑚 雯迪",
		Category:        "水",
		Tag:             "传奇-人类",
		Description:     "回合技2:在你使用1个使用花费小于3的法术后,可以花费2\\水将其重置",
		Quote:           "\"忘了那个老头吧,现在是年轻人的时代\"",
		ElementsCost:    map[string]int{"水": 6},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1211103.jpg",
	}
}

type CardDef1221001 struct{}

func (CardDef1221001) ID() string      { return "1221001" }
func (CardDef1221001) Name() string    { return "海豚伙伴" }
func (CardDef1221001) Kind() string    { return "伙伴" }
func (CardDef1221001) Element() string { return "水" }

func (CardDef1221001) Card() model.Card {
	return model.Card{
		Number:          "1221001",
		Type:            "伙伴",
		Name:            "海豚伙伴",
		Category:        "水",
		Tag:             "野兽",
		Description:     "诱发:当1个其他友方单位将要受到致命伤害时,将此卡献祭才能发动,防止该伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221001.jpg",
	}
}

type CardDef1221002 struct{}

func (CardDef1221002) ID() string      { return "1221002" }
func (CardDef1221002) Name() string    { return "冰原法师" }
func (CardDef1221002) Kind() string    { return "伙伴" }
func (CardDef1221002) Element() string { return "水" }

func (CardDef1221002) Card() model.Card {
	return model.Card{
		Number:          "1221002",
		Type:            "伙伴",
		Name:            "冰原法师",
		Category:        "水",
		Tag:             "巫师",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 4},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221002.jpg",
	}
}

type CardDef1221003 struct{}

func (CardDef1221003) ID() string      { return "1221003" }
func (CardDef1221003) Name() string    { return "掠夺者海盗" }
func (CardDef1221003) Kind() string    { return "伙伴" }
func (CardDef1221003) Element() string { return "水" }

func (CardDef1221003) Card() model.Card {
	return model.Card{
		Number:          "1221003",
		Type:            "伙伴",
		Name:            "掠夺者海盗",
		Category:        "水",
		Tag:             "人类",
		Description:     "",
		Quote:           "四大洋境内最臭名昭著的群体",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{"暗": 1, "水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221003.jpg",
	}
}

type CardDef1221004 struct{}

func (CardDef1221004) ID() string      { return "1221004" }
func (CardDef1221004) Name() string    { return "寒霜傀儡" }
func (CardDef1221004) Kind() string    { return "伙伴" }
func (CardDef1221004) Element() string { return "水" }

func (CardDef1221004) Card() model.Card {
	return model.Card{
		Number:          "1221004",
		Type:            "伙伴",
		Name:            "寒霜傀儡",
		Category:        "水",
		Tag:             "造物",
		Description:     "入场:对法力范围内1个敌方伙伴冻结1",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221004.jpg",
	}
}

type CardDef1221005 struct{}

func (CardDef1221005) ID() string      { return "1221005" }
func (CardDef1221005) Name() string    { return "西境海妖" }
func (CardDef1221005) Kind() string    { return "伙伴" }
func (CardDef1221005) Element() string { return "水" }

func (CardDef1221005) Card() model.Card {
	return model.Card{
		Number:          "1221005",
		Type:            "伙伴",
		Name:            "西境海妖",
		Category:        "水",
		Tag:             "异兽",
		Description:     "祈咒:选择法力范围内的1个敌方伙伴,将其横置",
		Quote:           "很难想象什么样的水手没能禁住诱惑",
		ElementsCost:    map[string]int{"无": 1, "水": 4},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221005.jpg",
	}
}

type CardDef1221006 struct{}

func (CardDef1221006) ID() string      { return "1221006" }
func (CardDef1221006) Name() string    { return "水栖狸猫" }
func (CardDef1221006) Kind() string    { return "伙伴" }
func (CardDef1221006) Element() string { return "水" }

func (CardDef1221006) Card() model.Card {
	return model.Card{
		Number:          "1221006",
		Type:            "伙伴",
		Name:            "水栖狸猫",
		Category:        "水",
		Tag:             "野兽",
		Description:     "光环:在本卡相邻有2个及以上水纹伙伴时,本卡负载+1\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221006.jpg",
	}
}

type CardDef1221007 struct{}

func (CardDef1221007) ID() string      { return "1221007" }
func (CardDef1221007) Name() string    { return "冰原狼" }
func (CardDef1221007) Kind() string    { return "伙伴" }
func (CardDef1221007) Element() string { return "水" }

func (CardDef1221007) Card() model.Card {
	return model.Card{
		Number:          "1221007",
		Type:            "伙伴",
		Name:            "冰原狼",
		Category:        "水",
		Tag:             "野兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221007.jpg",
	}
}

type CardDef1221008 struct{}

func (CardDef1221008) ID() string      { return "1221008" }
func (CardDef1221008) Name() string    { return "冰域恶魔" }
func (CardDef1221008) Kind() string    { return "伙伴" }
func (CardDef1221008) Element() string { return "水" }

func (CardDef1221008) Card() model.Card {
	return model.Card{
		Number:          "1221008",
		Type:            "伙伴",
		Name:            "冰域恶魔",
		Category:        "水",
		Tag:             "恶魔",
		Description:     "入场:对法力范围内的所有敌人冻结1",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 5},
		ElementsGain:    map[string]int{"暗": 1, "水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221008.jpg",
	}
}

type CardDef1221009 struct{}

func (CardDef1221009) ID() string      { return "1221009" }
func (CardDef1221009) Name() string    { return "南海海怪" }
func (CardDef1221009) Kind() string    { return "伙伴" }
func (CardDef1221009) Element() string { return "水" }

func (CardDef1221009) Card() model.Card {
	return model.Card{
		Number:          "1221009",
		Type:            "伙伴",
		Name:            "南海海怪",
		Category:        "水",
		Tag:             "异兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 5},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221009.jpg",
	}
}

type CardDef1221010 struct{}

func (CardDef1221010) ID() string      { return "1221010" }
func (CardDef1221010) Name() string    { return "护壁者" }
func (CardDef1221010) Kind() string    { return "伙伴" }
func (CardDef1221010) Element() string { return "水" }

func (CardDef1221010) Card() model.Card {
	return model.Card{
		Number:          "1221010",
		Type:            "伙伴",
		Name:            "护壁者",
		Category:        "水",
		Tag:             "巫师",
		Description:     "入场:直到下个回合结束所有法术\\攻变为0",
		Quote:           "在极北冰原之上,护壁者们或许是在守护,或许是在封印,人们唯一知道的是他们不欢迎任何冒险家",
		ElementsCost:    map[string]int{"水": 7},
		ElementsGain:    map[string]int{"水": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221010.jpg",
	}
}

type CardDef1221011 struct{}

func (CardDef1221011) ID() string      { return "1221011" }
func (CardDef1221011) Name() string    { return "凛冬城术士" }
func (CardDef1221011) Kind() string    { return "伙伴" }
func (CardDef1221011) Element() string { return "水" }

func (CardDef1221011) Card() model.Card {
	return model.Card{
		Number:          "1221011",
		Type:            "伙伴",
		Name:            "凛冬城术士",
		Category:        "水",
		Tag:             "巫师",
		Description:     "主动绝技:本回合你的下一次法术获得冻结1",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221011.jpg",
	}
}

type CardDef1221012 struct{}

func (CardDef1221012) ID() string      { return "1221012" }
func (CardDef1221012) Name() string    { return "龙王子裔" }
func (CardDef1221012) Kind() string    { return "伙伴" }
func (CardDef1221012) Element() string { return "水" }

func (CardDef1221012) Card() model.Card {
	return model.Card{
		Number:          "1221012",
		Type:            "伙伴",
		Name:            "龙王子裔",
		Category:        "水",
		Tag:             "龙",
		Description:     "精通2:检索1个水纹伙伴并使其入场花费减少1\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221012.jpg",
	}
}

type CardDef1221013 struct{}

func (CardDef1221013) ID() string      { return "1221013" }
func (CardDef1221013) Name() string    { return "唤雨师" }
func (CardDef1221013) Kind() string    { return "伙伴" }
func (CardDef1221013) Element() string { return "水" }

func (CardDef1221013) Card() model.Card {
	return model.Card{
		Number:          "1221013",
		Type:            "伙伴",
		Name:            "唤雨师",
		Category:        "水",
		Tag:             "巫师",
		Description:     "光环:你的水纹和大气法术+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "水": 3},
		ElementsGain:    map[string]int{"气": 1, "水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221013.jpg",
	}
}

type CardDef1221014 struct{}

func (CardDef1221014) ID() string      { return "1221014" }
func (CardDef1221014) Name() string    { return "北海飞鱼" }
func (CardDef1221014) Kind() string    { return "伙伴" }
func (CardDef1221014) Element() string { return "水" }

func (CardDef1221014) Card() model.Card {
	return model.Card{
		Number:          "1221014",
		Type:            "伙伴",
		Name:            "北海飞鱼",
		Category:        "水",
		Tag:             "野兽",
		Description:     "主动回合技:负载临时改为1\\气",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221014.jpg",
	}
}

type CardDef1221015 struct{}

func (CardDef1221015) ID() string      { return "1221015" }
func (CardDef1221015) Name() string    { return "眺望者商舰" }
func (CardDef1221015) Kind() string    { return "伙伴" }
func (CardDef1221015) Element() string { return "水" }

func (CardDef1221015) Card() model.Card {
	return model.Card{
		Number:          "1221015",
		Type:            "伙伴",
		Name:            "眺望者商舰",
		Category:        "水",
		Tag:             "机械",
		Description:     "祈咒:检索1个水纹卡牌,然后选择1张手牌洗回卡组",
		Quote:           "群屿大陆绝不会停止的两件事:战争和贸易",
		ElementsCost:    map[string]int{"气": 2, "水": 4},
		ElementsGain:    map[string]int{"气": 1, "水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221015.jpg",
	}
}

type CardDef1221016 struct{}

func (CardDef1221016) ID() string      { return "1221016" }
func (CardDef1221016) Name() string    { return "冰刺堡垒" }
func (CardDef1221016) Kind() string    { return "伙伴" }
func (CardDef1221016) Element() string { return "水" }

func (CardDef1221016) Card() model.Card {
	return model.Card{
		Number:          "1221016",
		Type:            "伙伴",
		Name:            "冰刺堡垒",
		Category:        "水",
		Tag:             "造物",
		Description:     "诱发:每当此卡受到敌方伤害,选择法力范围内1个敌人冻结1,如果已冻结,则改为造成1点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\水\\1221016.jpg",
	}
}

type CardDef1221101 struct{}

func (CardDef1221101) ID() string      { return "1221101" }
func (CardDef1221101) Name() string    { return "掠夺者海盗船" }
func (CardDef1221101) Kind() string    { return "伙伴" }
func (CardDef1221101) Element() string { return "水" }

func (CardDef1221101) Card() model.Card {
	return model.Card{
		Number:          "1221101",
		Type:            "伙伴",
		Name:            "掠夺者海盗船",
		Category:        "水",
		Tag:             "机械",
		Description:     "",
		Quote:           "追寻宝藏的海盗,永远无法见到海洋最珍贵的宝物",
		ElementsCost:    map[string]int{"暗": 1, "水": 4},
		ElementsGain:    map[string]int{"暗": 1, "水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221101.jpg",
	}
}

type CardDef1221102 struct{}

func (CardDef1221102) ID() string      { return "1221102" }
func (CardDef1221102) Name() string    { return "雾之国法师" }
func (CardDef1221102) Kind() string    { return "伙伴" }
func (CardDef1221102) Element() string { return "水" }

func (CardDef1221102) Card() model.Card {
	return model.Card{
		Number:          "1221102",
		Type:            "伙伴",
		Name:            "雾之国法师",
		Category:        "水",
		Tag:             "巫师",
		Description:     "回合技:使1个没有隐蔽的其他友方单位隐蔽2",
		Quote:           "\"没错,我之前只是在保存实力\"",
		ElementsCost:    map[string]int{"气": 1, "水": 4},
		ElementsGain:    map[string]int{"气": 1, "水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221102.jpg",
	}
}

type CardDef1221103 struct{}

func (CardDef1221103) ID() string      { return "1221103" }
func (CardDef1221103) Name() string    { return "凛冬城射手" }
func (CardDef1221103) Kind() string    { return "伙伴" }
func (CardDef1221103) Element() string { return "水" }

func (CardDef1221103) Card() model.Card {
	return model.Card{
		Number:          "1221103",
		Type:            "伙伴",
		Name:            "凛冬城射手",
		Category:        "水",
		Tag:             "人类",
		Description:     "引魔.此卡不位于前排也能进行攻击",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221103.jpg",
	}
}

type CardDef1221104 struct{}

func (CardDef1221104) ID() string      { return "1221104" }
func (CardDef1221104) Name() string    { return "冰原猛犸" }
func (CardDef1221104) Kind() string    { return "伙伴" }
func (CardDef1221104) Element() string { return "水" }

func (CardDef1221104) Card() model.Card {
	return model.Card{
		Number:          "1221104",
		Type:            "伙伴",
		Name:            "冰原猛犸",
		Category:        "水",
		Tag:             "野兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 5},
		ElementsGain:    map[string]int{"水": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221104.jpg",
	}
}

type CardDef1221105 struct{}

func (CardDef1221105) ID() string      { return "1221105" }
func (CardDef1221105) Name() string    { return "雾之国舞女" }
func (CardDef1221105) Kind() string    { return "伙伴" }
func (CardDef1221105) Element() string { return "水" }

func (CardDef1221105) Card() model.Card {
	return model.Card{
		Number:          "1221105",
		Type:            "伙伴",
		Name:            "雾之国舞女",
		Category:        "水",
		Tag:             "精灵",
		Description:     "入场:使法力范围内1个伙伴隐蔽2",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "水": 2},
		ElementsGain:    map[string]int{"气": 1, "水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221105.jpg",
	}
}

type CardDef1221106 struct{}

func (CardDef1221106) ID() string      { return "1221106" }
func (CardDef1221106) Name() string    { return "镜花海之莲" }
func (CardDef1221106) Kind() string    { return "伙伴" }
func (CardDef1221106) Element() string { return "水" }

func (CardDef1221106) Card() model.Card {
	return model.Card{
		Number:          "1221106",
		Type:            "伙伴",
		Name:            "镜花海之莲",
		Category:        "水",
		Tag:             "植物",
		Description:     "引魔.祈咒:获得负载+1\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221106.jpg",
	}
}

type CardDef1221107 struct{}

func (CardDef1221107) ID() string      { return "1221107" }
func (CardDef1221107) Name() string    { return "心莲守护者" }
func (CardDef1221107) Kind() string    { return "伙伴" }
func (CardDef1221107) Element() string { return "水" }

func (CardDef1221107) Card() model.Card {
	return model.Card{
		Number:          "1221107",
		Type:            "伙伴",
		Name:            "心莲守护者",
		Category:        "水",
		Tag:             "巫师",
		Description:     "入场:获得护盾2",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 4},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221107.jpg",
	}
}

type CardDef1221108 struct{}

func (CardDef1221108) ID() string      { return "1221108" }
func (CardDef1221108) Name() string    { return "心莲镜魔师" }
func (CardDef1221108) Kind() string    { return "伙伴" }
func (CardDef1221108) Element() string { return "水" }

func (CardDef1221108) Card() model.Card {
	return model.Card{
		Number:          "1221108",
		Type:            "伙伴",
		Name:            "心莲镜魔师",
		Category:        "水",
		Tag:             "巫师",
		Description:     "引魔.回合技:在你使用1个创造种类的技能后,可以翻取1张水纹属性道具牌,如果是反制卡牌可以立刻盖放在场上并且入场费用变为0",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 5},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221108.jpg",
	}
}

type CardDef1221109 struct{}

func (CardDef1221109) ID() string      { return "1221109" }
func (CardDef1221109) Name() string    { return "雾霭幽魂" }
func (CardDef1221109) Kind() string    { return "伙伴" }
func (CardDef1221109) Element() string { return "水" }

func (CardDef1221109) Card() model.Card {
	return model.Card{
		Number:          "1221109",
		Type:            "伙伴",
		Name:            "雾霭幽魂",
		Category:        "水",
		Tag:             "造物",
		Description:     "入场:获得隐蔽3.光环:此卡具有隐蔽时获得负载+2\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221109.jpg",
	}
}

type CardDef1221110 struct{}

func (CardDef1221110) ID() string      { return "1221110" }
func (CardDef1221110) Name() string    { return "掠夺者幽灵船长" }
func (CardDef1221110) Kind() string    { return "伙伴" }
func (CardDef1221110) Element() string { return "水" }

func (CardDef1221110) Card() model.Card {
	return model.Card{
		Number:          "1221110",
		Type:            "伙伴",
		Name:            "掠夺者幽灵船长",
		Category:        "水",
		Tag:             "人类",
		Description:     "光环:你的其他名字带有\"掠夺者\"的伙伴牌获得负载+1\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1, "水": 3},
		ElementsGain:    map[string]int{"暗": 1, "水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221110.jpg",
	}
}

type CardDef1221111 struct{}

func (CardDef1221111) ID() string      { return "1221111" }
func (CardDef1221111) Name() string    { return "掠夺者炮手" }
func (CardDef1221111) Kind() string    { return "伙伴" }
func (CardDef1221111) Element() string { return "水" }

func (CardDef1221111) Card() model.Card {
	return model.Card{
		Number:          "1221111",
		Type:            "伙伴",
		Name:            "掠夺者炮手",
		Category:        "水",
		Tag:             "人类",
		Description:     "诱发绝技:在你的法术命中后,随机弃置敌方1张手牌",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1, "水": 1},
		ElementsGain:    map[string]int{"暗": 1, "水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221111.jpg",
	}
}

type CardDef1221112 struct{}

func (CardDef1221112) ID() string      { return "1221112" }
func (CardDef1221112) Name() string    { return "水魔导师" }
func (CardDef1221112) Kind() string    { return "伙伴" }
func (CardDef1221112) Element() string { return "水" }

func (CardDef1221112) Card() model.Card {
	return model.Card{
		Number:          "1221112",
		Type:            "伙伴",
		Name:            "水魔导师",
		Category:        "水",
		Tag:             "巫师",
		Description:     "主动绝技:重置你的1个使用花费小于3的水纹法术",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221112.jpg",
	}
}

type CardDef1221113 struct{}

func (CardDef1221113) ID() string      { return "1221113" }
func (CardDef1221113) Name() string    { return "凛冬城象骑兵" }
func (CardDef1221113) Kind() string    { return "伙伴" }
func (CardDef1221113) Element() string { return "水" }

func (CardDef1221113) Card() model.Card {
	return model.Card{
		Number:          "1221113",
		Type:            "伙伴",
		Name:            "凛冬城象骑兵",
		Category:        "水",
		Tag:             "人类",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "水": 5},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            5,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221113.jpg",
	}
}

type CardDef1221114 struct{}

func (CardDef1221114) ID() string      { return "1221114" }
func (CardDef1221114) Name() string    { return "玉面雪狐" }
func (CardDef1221114) Kind() string    { return "伙伴" }
func (CardDef1221114) Element() string { return "水" }

func (CardDef1221114) Card() model.Card {
	return model.Card{
		Number:          "1221114",
		Type:            "伙伴",
		Name:            "玉面雪狐",
		Category:        "水",
		Tag:             "野兽",
		Description:     "诱发绝技:当敌方使用法术攻击时,你可以立刻移动此卡并获得2\\水,敌方需要重新选择目标",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221114.jpg",
	}
}

type CardDef1221115 struct{}

func (CardDef1221115) ID() string      { return "1221115" }
func (CardDef1221115) Name() string    { return "凛冬城御魔师" }
func (CardDef1221115) Kind() string    { return "伙伴" }
func (CardDef1221115) Element() string { return "水" }

func (CardDef1221115) Card() model.Card {
	return model.Card{
		Number:          "1221115",
		Type:            "伙伴",
		Name:            "凛冬城御魔师",
		Category:        "水",
		Tag:             "巫师",
		Description:     "祈咒:使你已学习的所有技能下一次使用花费-1\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 7},
		ElementsGain:    map[string]int{"水": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\水\\1221115.jpg",
	}
}

type CardDef1311001 struct{}

func (CardDef1311001) ID() string      { return "1311001" }
func (CardDef1311001) Name() string    { return "大鹏" }
func (CardDef1311001) Kind() string    { return "伙伴" }
func (CardDef1311001) Element() string { return "气" }

func (CardDef1311001) Card() model.Card {
	return model.Card{
		Number:          "1311001",
		Type:            "伙伴",
		Name:            "大鹏",
		Category:        "气",
		Tag:             "传奇-野兽",
		Description:     "入场:翻开卡组顶8张牌,抽取其中入场花费小于3的卡牌,重洗你的卡组,在本回合结束时必须丢弃这些这些被抽取的卡牌",
		Quote:           "跨四境,御九霄",
		ElementsCost:    map[string]int{"光": 1, "气": 5},
		ElementsGain:    map[string]int{"光": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1311001.jpg",
	}
}

type CardDef1311002 struct{}

func (CardDef1311002) ID() string      { return "1311002" }
func (CardDef1311002) Name() string    { return "\"风暴之女\" 艾拉雅" }
func (CardDef1311002) Kind() string    { return "伙伴" }
func (CardDef1311002) Element() string { return "气" }

func (CardDef1311002) Card() model.Card {
	return model.Card{
		Number:          "1311002",
		Type:            "伙伴",
		Name:            "\"风暴之女\" 艾拉雅",
		Category:        "气",
		Tag:             "传奇-巫师",
		Description:     "绑定技能:风暴之怒.光环:在你的手牌数大于等于手牌上限时,风暴之怒视为已经生效",
		Quote:           "我并不操纵风暴,我就是风暴",
		ElementsCost:    map[string]int{"气": 6},
		ElementsGain:    map[string]int{"气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"3301001"},
		OutputPath:      "output\\基础包\\伙伴\\气\\1311002.jpg",
	}
}

type CardDef1311003 struct{}

func (CardDef1311003) ID() string      { return "1311003" }
func (CardDef1311003) Name() string    { return "\"风刃\" 卡琳娜" }
func (CardDef1311003) Kind() string    { return "伙伴" }
func (CardDef1311003) Element() string { return "气" }

func (CardDef1311003) Card() model.Card {
	return model.Card{
		Number:          "1311003",
		Type:            "伙伴",
		Name:            "\"风刃\" 卡琳娜",
		Category:        "气",
		Tag:             "传奇-人类",
		Description:     "引魔.光环:你的没有穿透的大气技能获得穿透和使用花费+1\\气(不需要选择目标的技能不受影响)",
		Quote:           "匕首的价值,取决于它架在谁的脖子上",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1311003.jpg",
	}
}

type CardDef1311101 struct{}

func (CardDef1311101) ID() string      { return "1311101" }
func (CardDef1311101) Name() string    { return "斯帕罗 银叶" }
func (CardDef1311101) Kind() string    { return "伙伴" }
func (CardDef1311101) Element() string { return "气" }

func (CardDef1311101) Card() model.Card {
	return model.Card{
		Number:          "1311101",
		Type:            "伙伴",
		Name:            "斯帕罗 银叶",
		Category:        "气",
		Tag:             "传奇-巫师",
		Description:     "入场:对法力范围内1名敌人造成等同于你本回合丢弃手牌数量的伤害,最高3点",
		Quote:           "\"好久不见,我的哥哥.上次还是...还是你把我送进这间牢房的时候\"",
		ElementsCost:    map[string]int{"无": 1, "气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1311101.jpg",
	}
}

type CardDef1311102 struct{}

func (CardDef1311102) ID() string      { return "1311102" }
func (CardDef1311102) Name() string    { return "云顶商行 克罗斯" }
func (CardDef1311102) Kind() string    { return "伙伴" }
func (CardDef1311102) Element() string { return "气" }

func (CardDef1311102) Card() model.Card {
	return model.Card{
		Number:          "1311102",
		Type:            "伙伴",
		Name:            "云顶商行 克罗斯",
		Category:        "气",
		Tag:             "传奇-人类",
		Description:     "光环:如果你的卡组为空,你可以从对手的卡组抽牌,那些卡牌的入场花费和负载全部变为等量的\\无",
		Quote:           "\"我有的赚,大家都有的赚\"",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1311102.jpg",
	}
}

type CardDef1311103 struct{}

func (CardDef1311103) ID() string      { return "1311103" }
func (CardDef1311103) Name() string    { return "九霄议庭言主 麦阿提" }
func (CardDef1311103) Kind() string    { return "伙伴" }
func (CardDef1311103) Element() string { return "气" }

func (CardDef1311103) Card() model.Card {
	return model.Card{
		Number:          "1311103",
		Type:            "伙伴",
		Name:            "九霄议庭言主 麦阿提",
		Category:        "气",
		Tag:             "传奇-人类",
		Description:     "光环:对方手牌上限-1,每当对方手牌超过上限必须立刻弃牌至手牌上限",
		Quote:           "\"我们会照顾所有人的需求,在有序的前提下\"",
		ElementsCost:    map[string]int{"气": 5},
		ElementsGain:    map[string]int{"光": 1, "气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1311103.jpg",
	}
}

type CardDef1321001 struct{}

func (CardDef1321001) ID() string      { return "1321001" }
func (CardDef1321001) Name() string    { return "渡鸦信使" }
func (CardDef1321001) Kind() string    { return "伙伴" }
func (CardDef1321001) Element() string { return "气" }

func (CardDef1321001) Card() model.Card {
	return model.Card{
		Number:          "1321001",
		Type:            "伙伴",
		Name:            "渡鸦信使",
		Category:        "气",
		Tag:             "野兽",
		Description:     "主动:消耗此卡才能发动,抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "气": 1},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321001.jpg",
	}
}

type CardDef1321002 struct{}

func (CardDef1321002) ID() string      { return "1321002" }
func (CardDef1321002) Name() string    { return "随风旅行者" }
func (CardDef1321002) Kind() string    { return "伙伴" }
func (CardDef1321002) Element() string { return "气" }

func (CardDef1321002) Card() model.Card {
	return model.Card{
		Number:          "1321002",
		Type:            "伙伴",
		Name:            "随风旅行者",
		Category:        "气",
		Tag:             "精灵",
		Description:     "入场:获得1\\气.遗言:抽1张牌",
		Quote:           "记得避开艾拉雅",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321002.jpg",
	}
}

type CardDef1321003 struct{}

func (CardDef1321003) ID() string      { return "1321003" }
func (CardDef1321003) Name() string    { return "魔法蒲公英" }
func (CardDef1321003) Kind() string    { return "伙伴" }
func (CardDef1321003) Element() string { return "气" }

func (CardDef1321003) Card() model.Card {
	return model.Card{
		Number:          "1321003",
		Type:            "伙伴",
		Name:            "魔法蒲公英",
		Category:        "气",
		Tag:             "植物",
		Description:     "诱发:当你抽到此卡时,将其展示.入场:如果你在本回合抽到此卡,抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321003.jpg",
	}
}

type CardDef1321004 struct{}

func (CardDef1321004) ID() string      { return "1321004" }
func (CardDef1321004) Name() string    { return "雷电元素" }
func (CardDef1321004) Kind() string    { return "伙伴" }
func (CardDef1321004) Element() string { return "气" }

func (CardDef1321004) Card() model.Card {
	return model.Card{
		Number:          "1321004",
		Type:            "伙伴",
		Name:            "雷电元素",
		Category:        "气",
		Tag:             "造物",
		Description:     "入场:使法力范围内1个伙伴晕眩1",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{"光": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321004.jpg",
	}
}

type CardDef1321005 struct{}

func (CardDef1321005) ID() string      { return "1321005" }
func (CardDef1321005) Name() string    { return "驭风师" }
func (CardDef1321005) Kind() string    { return "伙伴" }
func (CardDef1321005) Element() string { return "气" }

func (CardDef1321005) Card() model.Card {
	return model.Card{
		Number:          "1321005",
		Type:            "伙伴",
		Name:            "驭风师",
		Category:        "气",
		Tag:             "巫师",
		Description:     "主动绝技:丢弃任意数量的手牌,每张牌使你获得1\\气",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 4},
		ElementsGain:    map[string]int{"气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321005.jpg",
	}
}

type CardDef1321006 struct{}

func (CardDef1321006) ID() string      { return "1321006" }
func (CardDef1321006) Name() string    { return "雷霆兽" }
func (CardDef1321006) Kind() string    { return "伙伴" }
func (CardDef1321006) Element() string { return "气" }

func (CardDef1321006) Card() model.Card {
	return model.Card{
		Number:          "1321006",
		Type:            "伙伴",
		Name:            "雷霆兽",
		Category:        "气",
		Tag:             "异兽",
		Description:     "光环:你的大气法术+1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "气": 4},
		ElementsGain:    map[string]int{"光": 1, "气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321006.jpg",
	}
}

type CardDef1321007 struct{}

func (CardDef1321007) ID() string      { return "1321007" }
func (CardDef1321007) Name() string    { return "工蜂骑士" }
func (CardDef1321007) Kind() string    { return "伙伴" }
func (CardDef1321007) Element() string { return "气" }

func (CardDef1321007) Card() model.Card {
	return model.Card{
		Number:          "1321007",
		Type:            "伙伴",
		Name:            "工蜂骑士",
		Category:        "气",
		Tag:             "人类",
		Description:     "",
		Quote:           "工蜂可能是一次性的,但是骑士不是",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321007.jpg",
	}
}

type CardDef1321008 struct{}

func (CardDef1321008) ID() string      { return "1321008" }
func (CardDef1321008) Name() string    { return "风息奔马" }
func (CardDef1321008) Kind() string    { return "伙伴" }
func (CardDef1321008) Element() string { return "气" }

func (CardDef1321008) Card() model.Card {
	return model.Card{
		Number:          "1321008",
		Type:            "伙伴",
		Name:            "风息奔马",
		Category:        "气",
		Tag:             "野兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{"气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321008.jpg",
	}
}

type CardDef1321009 struct{}

func (CardDef1321009) ID() string      { return "1321009" }
func (CardDef1321009) Name() string    { return "风魔" }
func (CardDef1321009) Kind() string    { return "伙伴" }
func (CardDef1321009) Element() string { return "气" }

func (CardDef1321009) Card() model.Card {
	return model.Card{
		Number:          "1321009",
		Type:            "伙伴",
		Name:            "风魔",
		Category:        "气",
		Tag:             "恶魔",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2, "气": 5},
		ElementsGain:    map[string]int{"气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            5,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321009.jpg",
	}
}

type CardDef1321010 struct{}

func (CardDef1321010) ID() string      { return "1321010" }
func (CardDef1321010) Name() string    { return "风暴奇美拉" }
func (CardDef1321010) Kind() string    { return "伙伴" }
func (CardDef1321010) Element() string { return "气" }

func (CardDef1321010) Card() model.Card {
	return model.Card{
		Number:          "1321010",
		Type:            "伙伴",
		Name:            "风暴奇美拉",
		Category:        "气",
		Tag:             "异兽",
		Description:     "引魔.吞噬:3\\气.光环:你的大气法术使用花费减少1\\气",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{"气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321010.jpg",
	}
}

type CardDef1321011 struct{}

func (CardDef1321011) ID() string      { return "1321011" }
func (CardDef1321011) Name() string    { return "雷精灵" }
func (CardDef1321011) Kind() string    { return "伙伴" }
func (CardDef1321011) Element() string { return "气" }

func (CardDef1321011) Card() model.Card {
	return model.Card{
		Number:          "1321011",
		Type:            "伙伴",
		Name:            "雷精灵",
		Category:        "气",
		Tag:             "精灵",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{"光": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321011.jpg",
	}
}

type CardDef1321012 struct{}

func (CardDef1321012) ID() string      { return "1321012" }
func (CardDef1321012) Name() string    { return "风灵媒师" }
func (CardDef1321012) Kind() string    { return "伙伴" }
func (CardDef1321012) Element() string { return "气" }

func (CardDef1321012) Card() model.Card {
	return model.Card{
		Number:          "1321012",
		Type:            "伙伴",
		Name:            "风灵媒师",
		Category:        "气",
		Tag:             "巫师",
		Description:     "诱发回合技:在你使用1个大气技能后,抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 4},
		ElementsGain:    map[string]int{"气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321012.jpg",
	}
}

type CardDef1321013 struct{}

func (CardDef1321013) ID() string      { return "1321013" }
func (CardDef1321013) Name() string    { return "传送法师" }
func (CardDef1321013) Kind() string    { return "伙伴" }
func (CardDef1321013) Element() string { return "气" }

func (CardDef1321013) Card() model.Card {
	return model.Card{
		Number:          "1321013",
		Type:            "伙伴",
		Name:            "传送法师",
		Category:        "气",
		Tag:             "巫师",
		Description:     "主动回合技:移动1个友方单位",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "气": 1},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321013.jpg",
	}
}

type CardDef1321014 struct{}

func (CardDef1321014) ID() string      { return "1321014" }
func (CardDef1321014) Name() string    { return "风息谷雷鸟" }
func (CardDef1321014) Kind() string    { return "伙伴" }
func (CardDef1321014) Element() string { return "气" }

func (CardDef1321014) Card() model.Card {
	return model.Card{
		Number:          "1321014",
		Type:            "伙伴",
		Name:            "风息谷雷鸟",
		Category:        "气",
		Tag:             "异兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 5},
		ElementsGain:    map[string]int{"气": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321014.jpg",
	}
}

type CardDef1321015 struct{}

func (CardDef1321015) ID() string      { return "1321015" }
func (CardDef1321015) Name() string    { return "风语者" }
func (CardDef1321015) Kind() string    { return "伙伴" }
func (CardDef1321015) Element() string { return "气" }

func (CardDef1321015) Card() model.Card {
	return model.Card{
		Number:          "1321015",
		Type:            "伙伴",
		Name:            "风语者",
		Category:        "气",
		Tag:             "精灵",
		Description:     "诱发回合技:当你丢弃手牌时,获得1\\气",
		Quote:           "那些逝去的终会回来,我在风中听到了它们的低语",
		ElementsCost:    map[string]int{"气": 4},
		ElementsGain:    map[string]int{"气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321015.jpg",
	}
}

type CardDef1321016 struct{}

func (CardDef1321016) ID() string      { return "1321016" }
func (CardDef1321016) Name() string    { return "雷傀儡" }
func (CardDef1321016) Kind() string    { return "伙伴" }
func (CardDef1321016) Element() string { return "气" }

func (CardDef1321016) Card() model.Card {
	return model.Card{
		Number:          "1321016",
		Type:            "伙伴",
		Name:            "雷傀儡",
		Category:        "气",
		Tag:             "造物",
		Description:     "遗言:对手丢弃1张手牌",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\气\\1321016.jpg",
	}
}

type CardDef1321101 struct{}

func (CardDef1321101) ID() string      { return "1321101" }
func (CardDef1321101) Name() string    { return "翱翔者E2型运输舰" }
func (CardDef1321101) Kind() string    { return "伙伴" }
func (CardDef1321101) Element() string { return "气" }

func (CardDef1321101) Card() model.Card {
	return model.Card{
		Number:          "1321101",
		Type:            "伙伴",
		Name:            "翱翔者E2型运输舰",
		Category:        "气",
		Tag:             "机械",
		Description:     "祈咒:抽2张牌,或挑选弃牌堆的2张大气卡牌洗回牌组",
		Quote:           "西部重工负责制造以及保修期内的维修,商品使用期间造成的连带损害由使用者承担.",
		ElementsCost:    map[string]int{"气": 6},
		ElementsGain:    map[string]int{"气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321101.jpg",
	}
}

type CardDef1321102 struct{}

func (CardDef1321102) ID() string      { return "1321102" }
func (CardDef1321102) Name() string    { return "花斑麻雀" }
func (CardDef1321102) Kind() string    { return "伙伴" }
func (CardDef1321102) Element() string { return "气" }

func (CardDef1321102) Card() model.Card {
	return model.Card{
		Number:          "1321102",
		Type:            "伙伴",
		Name:            "花斑麻雀",
		Category:        "气",
		Tag:             "野兽",
		Description:     "诱发:如果本卡被从手牌弃掉,可以花费1\\气召唤本卡",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321102.jpg",
	}
}

type CardDef1321103 struct{}

func (CardDef1321103) ID() string      { return "1321103" }
func (CardDef1321103) Name() string    { return "孤星塔守望者" }
func (CardDef1321103) Kind() string    { return "伙伴" }
func (CardDef1321103) Element() string { return "气" }

func (CardDef1321103) Card() model.Card {
	return model.Card{
		Number:          "1321103",
		Type:            "伙伴",
		Name:            "孤星塔守望者",
		Category:        "气",
		Tag:             "人类",
		Description:     "主动绝技:丢弃至多3张手牌并获得等量护盾",
		Quote:           "放下尘世因缘,方能成为守望者",
		ElementsCost:    map[string]int{"无": 1, "气": 1},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321103.jpg",
	}
}

type CardDef1321104 struct{}

func (CardDef1321104) ID() string      { return "1321104" }
func (CardDef1321104) Name() string    { return "织雾者" }
func (CardDef1321104) Kind() string    { return "伙伴" }
func (CardDef1321104) Element() string { return "气" }

func (CardDef1321104) Card() model.Card {
	return model.Card{
		Number:          "1321104",
		Type:            "伙伴",
		Name:            "织雾者",
		Category:        "气",
		Tag:             "巫师",
		Description:     "入场:使任意1个敌方单位隐蔽2",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{"气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321104.jpg",
	}
}

type CardDef1321105 struct{}

func (CardDef1321105) ID() string      { return "1321105" }
func (CardDef1321105) Name() string    { return "幻术师" }
func (CardDef1321105) Kind() string    { return "伙伴" }
func (CardDef1321105) Element() string { return "气" }

func (CardDef1321105) Card() model.Card {
	return model.Card{
		Number:          "1321105",
		Type:            "伙伴",
		Name:            "幻术师",
		Category:        "气",
		Tag:             "巫师",
		Description:     "主动绝技:将你的1个入场花费小于6的伙伴移回手牌,并获得等同于其负载的元素费用",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "气": 2},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321105.jpg",
	}
}

type CardDef1321106 struct{}

func (CardDef1321106) ID() string      { return "1321106" }
func (CardDef1321106) Name() string    { return "银叶游侠" }
func (CardDef1321106) Kind() string    { return "伙伴" }
func (CardDef1321106) Element() string { return "气" }

func (CardDef1321106) Card() model.Card {
	return model.Card{
		Number:          "1321106",
		Type:            "伙伴",
		Name:            "银叶游侠",
		Category:        "气",
		Tag:             "人类",
		Description:     "主动:消耗此卡才能发动,你的下一次法术+1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321106.jpg",
	}
}

type CardDef1321107 struct{}

func (CardDef1321107) ID() string      { return "1321107" }
func (CardDef1321107) Name() string    { return "云霄城大盗" }
func (CardDef1321107) Kind() string    { return "伙伴" }
func (CardDef1321107) Element() string { return "气" }

func (CardDef1321107) Card() model.Card {
	return model.Card{
		Number:          "1321107",
		Type:            "伙伴",
		Name:            "云霄城大盗",
		Category:        "气",
		Tag:             "人类",
		Description:     "入场:双方各自随机丢弃1张手牌",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321107.jpg",
	}
}

type CardDef1321108 struct{}

func (CardDef1321108) ID() string      { return "1321108" }
func (CardDef1321108) Name() string    { return "翡翠蜂鸟" }
func (CardDef1321108) Kind() string    { return "伙伴" }
func (CardDef1321108) Element() string { return "气" }

func (CardDef1321108) Card() model.Card {
	return model.Card{
		Number:          "1321108",
		Type:            "伙伴",
		Name:            "翡翠蜂鸟",
		Category:        "气",
		Tag:             "野兽",
		Description:     "入场:如果你的手牌数量小于2张,抽2张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321108.jpg",
	}
}

type CardDef1321109 struct{}

func (CardDef1321109) ID() string      { return "1321109" }
func (CardDef1321109) Name() string    { return "风暴之角" }
func (CardDef1321109) Kind() string    { return "伙伴" }
func (CardDef1321109) Element() string { return "气" }

func (CardDef1321109) Card() model.Card {
	return model.Card{
		Number:          "1321109",
		Type:            "伙伴",
		Name:            "风暴之角",
		Category:        "气",
		Tag:             "异兽",
		Description:     "绝技:丢弃1张手牌,翻取1张大气装备",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 4},
		ElementsGain:    map[string]int{"气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321109.jpg",
	}
}

type CardDef1321110 struct{}

func (CardDef1321110) ID() string      { return "1321110" }
func (CardDef1321110) Name() string    { return "银叶信使" }
func (CardDef1321110) Kind() string    { return "伙伴" }
func (CardDef1321110) Element() string { return "气" }

func (CardDef1321110) Card() model.Card {
	return model.Card{
		Number:          "1321110",
		Type:            "伙伴",
		Name:            "银叶信使",
		Category:        "气",
		Tag:             "人类",
		Description:     "入场:检索1张失落的银叶花",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321110.jpg",
	}
}

type CardDef1321111 struct{}

func (CardDef1321111) ID() string      { return "1321111" }
func (CardDef1321111) Name() string    { return "雷光战士" }
func (CardDef1321111) Kind() string    { return "伙伴" }
func (CardDef1321111) Element() string { return "气" }

func (CardDef1321111) Card() model.Card {
	return model.Card{
		Number:          "1321111",
		Type:            "伙伴",
		Name:            "雷光战士",
		Category:        "气",
		Tag:             "人类",
		Description:     "入场:你每装备有1件雷光道具(同时具有\\气和\\光负载),此卡可以获得以下1项:+2\\血,+1\\攻,负载+1\\气,负载+1\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "气": 4},
		ElementsGain:    map[string]int{"光": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321111.jpg",
	}
}

type CardDef1321112 struct{}

func (CardDef1321112) ID() string      { return "1321112" }
func (CardDef1321112) Name() string    { return "九霄接头人" }
func (CardDef1321112) Kind() string    { return "伙伴" }
func (CardDef1321112) Element() string { return "气" }

func (CardDef1321112) Card() model.Card {
	return model.Card{
		Number:          "1321112",
		Type:            "伙伴",
		Name:            "九霄接头人",
		Category:        "气",
		Tag:             "人类",
		Description:     "引魔.祈咒:如果对方手牌未达上限,将1张九霄印记加入对手的手牌",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2001102"},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321112.jpg",
	}
}

type CardDef1321113 struct{}

func (CardDef1321113) ID() string      { return "1321113" }
func (CardDef1321113) Name() string    { return "议庭传信鸽" }
func (CardDef1321113) Kind() string    { return "伙伴" }
func (CardDef1321113) Element() string { return "气" }

func (CardDef1321113) Card() model.Card {
	return model.Card{
		Number:          "1321113",
		Type:            "伙伴",
		Name:            "议庭传信鸽",
		Category:        "气",
		Tag:             "野兽",
		Description:     "入场:将1张九霄印记加入对手手牌",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2001102"},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321113.jpg",
	}
}

type CardDef1321114 struct{}

func (CardDef1321114) ID() string      { return "1321114" }
func (CardDef1321114) Name() string    { return "议庭执行者" }
func (CardDef1321114) Kind() string    { return "伙伴" }
func (CardDef1321114) Element() string { return "气" }

func (CardDef1321114) Card() model.Card {
	return model.Card{
		Number:          "1321114",
		Type:            "伙伴",
		Name:            "议庭执行者",
		Category:        "气",
		Tag:             "巫师",
		Description:     "入场:随机丢弃敌方1张手牌,如果是九霄印记再随机丢弃1张",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 5},
		ElementsGain:    map[string]int{"气": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321114.jpg",
	}
}

type CardDef1321115 struct{}

func (CardDef1321115) ID() string      { return "1321115" }
func (CardDef1321115) Name() string    { return "苍穹描摹者" }
func (CardDef1321115) Kind() string    { return "伙伴" }
func (CardDef1321115) Element() string { return "气" }

func (CardDef1321115) Card() model.Card {
	return model.Card{
		Number:          "1321115",
		Type:            "伙伴",
		Name:            "苍穹描摹者",
		Category:        "气",
		Tag:             "巫师",
		Description:     "入场:复制你场上另一个使用花费小于6的大气卡牌的入场效果",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 5},
		ElementsGain:    map[string]int{"气": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\气\\1321115.jpg",
	}
}

type CardDef1401001 struct{}

func (CardDef1401001) ID() string      { return "1401001" }
func (CardDef1401001) Name() string    { return "生命种子" }
func (CardDef1401001) Kind() string    { return "伙伴" }
func (CardDef1401001) Element() string { return "地" }

func (CardDef1401001) Card() model.Card {
	return model.Card{
		Number:          "1401001",
		Type:            "伙伴",
		Name:            "生命种子",
		Category:        "地",
		Tag:             "衍生-植物",
		Description:     "精通2:可以献祭此卡并从你的手牌中召唤1个地脉伙伴(无需花费),它会继承此卡的生命和负载加成",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1401001.jpg",
	}
}

type CardDef1401002 struct{}

func (CardDef1401002) ID() string      { return "1401002" }
func (CardDef1401002) Name() string    { return "灵兽 辛柯" }
func (CardDef1401002) Kind() string    { return "伙伴" }
func (CardDef1401002) Element() string { return "地" }

func (CardDef1401002) Card() model.Card {
	return model.Card{
		Number:          "1401002",
		Type:            "伙伴",
		Name:            "灵兽 辛柯",
		Category:        "地",
		Tag:             "衍生-野兽",
		Description:     "诱发:当友方单位受到敌方伤害后,可以从卡组或手牌召唤此卡,无需入场花费",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{"地": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1401002.jpg",
	}
}

type CardDef1401101 struct{}

func (CardDef1401101) ID() string      { return "1401101" }
func (CardDef1401101) Name() string    { return "普通蜥蜴" }
func (CardDef1401101) Kind() string    { return "伙伴" }
func (CardDef1401101) Element() string { return "地" }

func (CardDef1401101) Card() model.Card {
	return model.Card{
		Number:          "1401101",
		Type:            "伙伴",
		Name:            "普通蜥蜴",
		Category:        "地",
		Tag:             "衍生-野兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1401101.jpg",
	}
}

type CardDef1411001 struct{}

func (CardDef1411001) ID() string      { return "1411001" }
func (CardDef1411001) Name() string    { return "\"轮回不息\" 大德鲁伊 烟尘" }
func (CardDef1411001) Kind() string    { return "伙伴" }
func (CardDef1411001) Element() string { return "地" }

func (CardDef1411001) Card() model.Card {
	return model.Card{
		Number:          "1411001",
		Type:            "伙伴",
		Name:            "\"轮回不息\" 大德鲁伊 烟尘",
		Category:        "地",
		Tag:             "传奇-巫师",
		Description:     "诱发绝技:当1个友方伙伴死亡时,可以召唤1个生命种子,它会继承该伙伴的所有生命和负载加成",
		Quote:           "四百岁的大德鲁伊最讨厌的三件事情:休息,水晶蜘蛛,以及被一棵树称作孩子",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"1401001"},
		OutputPath:      "output\\基础包\\伙伴\\地\\1411001.jpg",
	}
}

type CardDef1411002 struct{}

func (CardDef1411002) ID() string      { return "1411002" }
func (CardDef1411002) Name() string    { return "\"知识古树\" 深耕" }
func (CardDef1411002) Kind() string    { return "伙伴" }
func (CardDef1411002) Element() string { return "地" }

func (CardDef1411002) Card() model.Card {
	return model.Card{
		Number:          "1411002",
		Type:            "伙伴",
		Name:            "\"知识古树\" 深耕",
		Category:        "地",
		Tag:             "传奇-植物",
		Description:     "入场:你的精通立刻达到最高",
		Quote:           "\"我亲爱的孩子,只要问我,我便会答复,有关生命,平衡,以及大地\"",
		ElementsCost:    map[string]int{"地": 5},
		ElementsGain:    map[string]int{"地": 1, "无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1411002.jpg",
	}
}

type CardDef1411003 struct{}

func (CardDef1411003) ID() string      { return "1411003" }
func (CardDef1411003) Name() string    { return "沙之魔巫 梭默" }
func (CardDef1411003) Kind() string    { return "伙伴" }
func (CardDef1411003) Element() string { return "地" }

func (CardDef1411003) Card() model.Card {
	return model.Card{
		Number:          "1411003",
		Type:            "伙伴",
		Name:            "沙之魔巫 梭默",
		Category:        "地",
		Tag:             "传奇-巫师",
		Description:     "光环:你的没有范围效果的地脉法术获得范围:方阵",
		Quote:           "沙瓦尔大陆的守护者,旅行者的指路人,屠魔者的灾星",
		ElementsCost:    map[string]int{"地": 4, "气": 2},
		ElementsGain:    map[string]int{"地": 2, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1411003.jpg",
	}
}

type CardDef1411101 struct{}

func (CardDef1411101) ID() string      { return "1411101" }
func (CardDef1411101) Name() string    { return "苍老者 弗兰肯 拜利兰" }
func (CardDef1411101) Kind() string    { return "伙伴" }
func (CardDef1411101) Element() string { return "地" }

func (CardDef1411101) Card() model.Card {
	return model.Card{
		Number:          "1411101",
		Type:            "伙伴",
		Name:            "苍老者 弗兰肯 拜利兰",
		Category:        "地",
		Tag:             "传奇-巫师",
		Description:     "祈咒:此卡失去1点负载.诱发:每当此卡失去1点负载,选择1个敌方法术永久减少2\\威",
		Quote:           "日月轮换,朽与不朽皆在一念间",
		ElementsCost:    map[string]int{"地": 7, "无": 2},
		ElementsGain:    map[string]int{"地": 3, "无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1411101.jpg",
	}
}

type CardDef1411102 struct{}

func (CardDef1411102) ID() string      { return "1411102" }
func (CardDef1411102) Name() string    { return "谧语精灵王 辛达瑞尔" }
func (CardDef1411102) Kind() string    { return "伙伴" }
func (CardDef1411102) Element() string { return "地" }

func (CardDef1411102) Card() model.Card {
	return model.Card{
		Number:          "1411102",
		Type:            "伙伴",
		Name:            "谧语精灵王 辛达瑞尔",
		Category:        "地",
		Tag:             "传奇-精灵",
		Description:     "入场:在上个回合里,敌方法术每满足以下一项(命中3次及以上,命中3个目标及以上,造成3点伤害及以上)可以选择任意1个不同的敌人,对那些敌人造成2点伤害",
		Quote:           "\"凡事不过三,否则我们必将以牙还牙,以眼还眼\"",
		ElementsCost:    map[string]int{"地": 4, "无": 1},
		ElementsGain:    map[string]int{"地": 1, "暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1411102.jpg",
	}
}

type CardDef1411103 struct{}

func (CardDef1411103) ID() string      { return "1411103" }
func (CardDef1411103) Name() string    { return "百兽之王 莱恩克塞斯" }
func (CardDef1411103) Kind() string    { return "伙伴" }
func (CardDef1411103) Element() string { return "地" }

func (CardDef1411103) Card() model.Card {
	return model.Card{
		Number:          "1411103",
		Type:            "伙伴",
		Name:            "百兽之王 莱恩克塞斯",
		Category:        "地",
		Tag:             "传奇-野兽",
		Description:     "入场:翻取1个伙伴,如果是地脉伙伴则将其召唤,无需花费",
		Quote:           "\"该死的杂种...不好,怎么还有这么多\"",
		ElementsCost:    map[string]int{"地": 8},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1411103.jpg",
	}
}

type CardDef1421001 struct{}

func (CardDef1421001) ID() string      { return "1421001" }
func (CardDef1421001) Name() string    { return "流沙法师" }
func (CardDef1421001) Kind() string    { return "伙伴" }
func (CardDef1421001) Element() string { return "地" }

func (CardDef1421001) Card() model.Card {
	return model.Card{
		Number:          "1421001",
		Type:            "伙伴",
		Name:            "流沙法师",
		Category:        "地",
		Tag:             "巫师",
		Description:     "入场:使1个无视范围的敌人石化1",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421001.jpg",
	}
}

type CardDef1421002 struct{}

func (CardDef1421002) ID() string      { return "1421002" }
func (CardDef1421002) Name() string    { return "祝祷祭师" }
func (CardDef1421002) Kind() string    { return "伙伴" }
func (CardDef1421002) Element() string { return "地" }

func (CardDef1421002) Card() model.Card {
	return model.Card{
		Number:          "1421002",
		Type:            "伙伴",
		Name:            "祝祷祭师",
		Category:        "地",
		Tag:             "巫师",
		Description:     "光环:此卡和相邻单位不受负面状态影响(仍可处于负面状态)",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "无": 1},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421002.jpg",
	}
}

type CardDef1421003 struct{}

func (CardDef1421003) ID() string      { return "1421003" }
func (CardDef1421003) Name() string    { return "成长的树人" }
func (CardDef1421003) Kind() string    { return "伙伴" }
func (CardDef1421003) Element() string { return "地" }

func (CardDef1421003) Card() model.Card {
	return model.Card{
		Number:          "1421003",
		Type:            "伙伴",
		Name:            "成长的树人",
		Category:        "地",
		Tag:             "植物",
		Description:     "精通2,4:此卡负载+1\\地或者+1\\血",
		Quote:           "每一棵小树,都是未来树林的支柱",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421003.jpg",
	}
}

type CardDef1421004 struct{}

func (CardDef1421004) ID() string      { return "1421004" }
func (CardDef1421004) Name() string    { return "森林守卫" }
func (CardDef1421004) Kind() string    { return "伙伴" }
func (CardDef1421004) Element() string { return "地" }

func (CardDef1421004) Card() model.Card {
	return model.Card{
		Number:          "1421004",
		Type:            "伙伴",
		Name:            "森林守卫",
		Category:        "地",
		Tag:             "造物",
		Description:     "精通1:+1\\血.精通3:负载+1\\地.精通5:+2\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421004.jpg",
	}
}

type CardDef1421005 struct{}

func (CardDef1421005) ID() string      { return "1421005" }
func (CardDef1421005) Name() string    { return "磐石元素" }
func (CardDef1421005) Kind() string    { return "伙伴" }
func (CardDef1421005) Element() string { return "地" }

func (CardDef1421005) Card() model.Card {
	return model.Card{
		Number:          "1421005",
		Type:            "伙伴",
		Name:            "磐石元素",
		Category:        "地",
		Tag:             "造物",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 6},
		ElementsGain:    map[string]int{"地": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421005.jpg",
	}
}

type CardDef1421006 struct{}

func (CardDef1421006) ID() string      { return "1421006" }
func (CardDef1421006) Name() string    { return "林地变形者" }
func (CardDef1421006) Kind() string    { return "伙伴" }
func (CardDef1421006) Element() string { return "地" }

func (CardDef1421006) Card() model.Card {
	return model.Card{
		Number:          "1421006",
		Type:            "伙伴",
		Name:            "林地变形者",
		Category:        "地",
		Tag:             "精灵",
		Description:     "",
		Quote:           "了解自然,亲近自然,融入自然",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421006.jpg",
	}
}

type CardDef1421007 struct{}

func (CardDef1421007) ID() string      { return "1421007" }
func (CardDef1421007) Name() string    { return "高地泰坦" }
func (CardDef1421007) Kind() string    { return "伙伴" }
func (CardDef1421007) Element() string { return "地" }

func (CardDef1421007) Card() model.Card {
	return model.Card{
		Number:          "1421007",
		Type:            "伙伴",
		Name:            "高地泰坦",
		Category:        "地",
		Tag:             "精灵",
		Description:     "光环:未被强化的法术对本卡造成的伤害+1",
		Quote:           "\"身高小于180米的泰坦都是残疾\"——斯卡尔蒂 罗佳",
		ElementsCost:    map[string]int{"地": 6, "无": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          3,
		Life:            7,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421007.jpg",
	}
}

type CardDef1421008 struct{}

func (CardDef1421008) ID() string      { return "1421008" }
func (CardDef1421008) Name() string    { return "岩山翼龙" }
func (CardDef1421008) Kind() string    { return "伙伴" }
func (CardDef1421008) Element() string { return "地" }

func (CardDef1421008) Card() model.Card {
	return model.Card{
		Number:          "1421008",
		Type:            "伙伴",
		Name:            "岩山翼龙",
		Category:        "地",
		Tag:             "野兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{"地": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421008.jpg",
	}
}

type CardDef1421009 struct{}

func (CardDef1421009) ID() string      { return "1421009" }
func (CardDef1421009) Name() string    { return "被祝福的少女" }
func (CardDef1421009) Kind() string    { return "伙伴" }
func (CardDef1421009) Element() string { return "地" }

func (CardDef1421009) Card() model.Card {
	return model.Card{
		Number:          "1421009",
		Type:            "伙伴",
		Name:            "被祝福的少女",
		Category:        "地",
		Tag:             "精灵",
		Description:     "祈咒:使1个相邻地脉伙伴获得负载+1\\地",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421009.jpg",
	}
}

type CardDef1421010 struct{}

func (CardDef1421010) ID() string      { return "1421010" }
func (CardDef1421010) Name() string    { return "种植园丁" }
func (CardDef1421010) Kind() string    { return "伙伴" }
func (CardDef1421010) Element() string { return "地" }

func (CardDef1421010) Card() model.Card {
	return model.Card{
		Number:          "1421010",
		Type:            "伙伴",
		Name:            "种植园丁",
		Category:        "地",
		Tag:             "精灵",
		Description:     "诱发:你的卡牌每次获得负载,在此卡上放置1个标记.主动回合技:取除2个标记才能发动,抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421010.jpg",
	}
}

type CardDef1421011 struct{}

func (CardDef1421011) ID() string      { return "1421011" }
func (CardDef1421011) Name() string    { return "大长老" }
func (CardDef1421011) Kind() string    { return "伙伴" }
func (CardDef1421011) Element() string { return "地" }

func (CardDef1421011) Card() model.Card {
	return model.Card{
		Number:          "1421011",
		Type:            "伙伴",
		Name:            "大长老",
		Category:        "地",
		Tag:             "精灵",
		Description:     "精通1,3:下一次学习地脉技能的花费-2",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421011.jpg",
	}
}

type CardDef1421012 struct{}

func (CardDef1421012) ID() string      { return "1421012" }
func (CardDef1421012) Name() string    { return "林地飞鼠" }
func (CardDef1421012) Kind() string    { return "伙伴" }
func (CardDef1421012) Element() string { return "地" }

func (CardDef1421012) Card() model.Card {
	return model.Card{
		Number:          "1421012",
		Type:            "伙伴",
		Name:            "林地飞鼠",
		Category:        "地",
		Tag:             "野兽",
		Description:     "主动回合技:负载临时改为1\\气",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421012.jpg",
	}
}

type CardDef1421013 struct{}

func (CardDef1421013) ID() string      { return "1421013" }
func (CardDef1421013) Name() string    { return "岩山恐兽" }
func (CardDef1421013) Kind() string    { return "伙伴" }
func (CardDef1421013) Element() string { return "地" }

func (CardDef1421013) Card() model.Card {
	return model.Card{
		Number:          "1421013",
		Type:            "伙伴",
		Name:            "岩山恐兽",
		Category:        "地",
		Tag:             "野兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 6, "无": 1},
		ElementsGain:    map[string]int{"地": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421013.jpg",
	}
}

type CardDef1421014 struct{}

func (CardDef1421014) ID() string      { return "1421014" }
func (CardDef1421014) Name() string    { return "风息谷旅商" }
func (CardDef1421014) Kind() string    { return "伙伴" }
func (CardDef1421014) Element() string { return "地" }

func (CardDef1421014) Card() model.Card {
	return model.Card{
		Number:          "1421014",
		Type:            "伙伴",
		Name:            "风息谷旅商",
		Category:        "地",
		Tag:             "人类",
		Description:     "入场:你的场上每有1个野兽,精灵或植物,抽1张牌(最多3张)",
		Quote:           "风息谷大概是商人最受欢迎的地方了,缺点是没什么客人",
		ElementsCost:    map[string]int{"地": 2, "无": 1},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421014.jpg",
	}
}

type CardDef1421015 struct{}

func (CardDef1421015) ID() string      { return "1421015" }
func (CardDef1421015) Name() string    { return "苍绿之龙" }
func (CardDef1421015) Kind() string    { return "伙伴" }
func (CardDef1421015) Element() string { return "地" }

func (CardDef1421015) Card() model.Card {
	return model.Card{
		Number:          "1421015",
		Type:            "伙伴",
		Name:            "苍绿之龙",
		Category:        "地",
		Tag:             "龙",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 7},
		ElementsGain:    map[string]int{"地": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421015.jpg",
	}
}

type CardDef1421016 struct{}

func (CardDef1421016) ID() string      { return "1421016" }
func (CardDef1421016) Name() string    { return "食腐者" }
func (CardDef1421016) Kind() string    { return "伙伴" }
func (CardDef1421016) Element() string { return "地" }

func (CardDef1421016) Card() model.Card {
	return model.Card{
		Number:          "1421016",
		Type:            "伙伴",
		Name:            "食腐者",
		Category:        "地",
		Tag:             "野兽",
		Description:     "引魔.诱发:每当其他友方单位受到对方伤害,你获得2\\地",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\地\\1421016.jpg",
	}
}

type CardDef1421101 struct{}

func (CardDef1421101) ID() string      { return "1421101" }
func (CardDef1421101) Name() string    { return "岩壁刺球" }
func (CardDef1421101) Kind() string    { return "伙伴" }
func (CardDef1421101) Element() string { return "地" }

func (CardDef1421101) Card() model.Card {
	return model.Card{
		Number:          "1421101",
		Type:            "伙伴",
		Name:            "岩壁刺球",
		Category:        "地",
		Tag:             "异兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421101.jpg",
	}
}

type CardDef1421102 struct{}

func (CardDef1421102) ID() string      { return "1421102" }
func (CardDef1421102) Name() string    { return "翡翠守卫" }
func (CardDef1421102) Kind() string    { return "伙伴" }
func (CardDef1421102) Element() string { return "地" }

func (CardDef1421102) Card() model.Card {
	return model.Card{
		Number:          "1421102",
		Type:            "伙伴",
		Name:            "翡翠守卫",
		Category:        "地",
		Tag:             "造物",
		Description:     "入场:如果你没有护盾,获得护盾2",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2, "无": 1},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421102.jpg",
	}
}

type CardDef1421103 struct{}

func (CardDef1421103) ID() string      { return "1421103" }
func (CardDef1421103) Name() string    { return "寄生虫" }
func (CardDef1421103) Kind() string    { return "伙伴" }
func (CardDef1421103) Element() string { return "地" }

func (CardDef1421103) Card() model.Card {
	return model.Card{
		Number:          "1421103",
		Type:            "伙伴",
		Name:            "寄生虫",
		Category:        "地",
		Tag:             "野兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{"地": 1, "暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421103.jpg",
	}
}

type CardDef1421104 struct{}

func (CardDef1421104) ID() string      { return "1421104" }
func (CardDef1421104) Name() string    { return "拜利兰森林熊" }
func (CardDef1421104) Kind() string    { return "伙伴" }
func (CardDef1421104) Element() string { return "地" }

func (CardDef1421104) Card() model.Card {
	return model.Card{
		Number:          "1421104",
		Type:            "伙伴",
		Name:            "拜利兰森林熊",
		Category:        "地",
		Tag:             "野兽",
		Description:     "引魔.入场:获得护盾3",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 6},
		ElementsGain:    map[string]int{"地": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421104.jpg",
	}
}

type CardDef1421105 struct{}

func (CardDef1421105) ID() string      { return "1421105" }
func (CardDef1421105) Name() string    { return "失活的根须" }
func (CardDef1421105) Kind() string    { return "伙伴" }
func (CardDef1421105) Element() string { return "地" }

func (CardDef1421105) Card() model.Card {
	return model.Card{
		Number:          "1421105",
		Type:            "伙伴",
		Name:            "失活的根须",
		Category:        "地",
		Tag:             "植物",
		Description:     "祈咒:若此卡没有负载,获得负载1\\地",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421105.jpg",
	}
}

type CardDef1421106 struct{}

func (CardDef1421106) ID() string      { return "1421106" }
func (CardDef1421106) Name() string    { return "幻影蜥蜴" }
func (CardDef1421106) Kind() string    { return "伙伴" }
func (CardDef1421106) Element() string { return "地" }

func (CardDef1421106) Card() model.Card {
	return model.Card{
		Number:          "1421106",
		Type:            "伙伴",
		Name:            "幻影蜥蜴",
		Category:        "地",
		Tag:             "野兽",
		Description:     "诱发绝技:在你使用一个灵媒技能后消耗此卡才能发动,此卡分裂为两个普通蜥蜴",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"1401101"},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421106.jpg",
	}
}

type CardDef1421107 struct{}

func (CardDef1421107) ID() string      { return "1421107" }
func (CardDef1421107) Name() string    { return "龙血树精" }
func (CardDef1421107) Kind() string    { return "伙伴" }
func (CardDef1421107) Element() string { return "地" }

func (CardDef1421107) Card() model.Card {
	return model.Card{
		Number:          "1421107",
		Type:            "伙伴",
		Name:            "龙血树精",
		Category:        "地",
		Tag:             "精灵",
		Description:     "入场:使1个友方卡牌失去1点负载,此卡获得负载+1\\暗",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421107.jpg",
	}
}

type CardDef1421108 struct{}

func (CardDef1421108) ID() string      { return "1421108" }
func (CardDef1421108) Name() string    { return "凯尔特灵鹿" }
func (CardDef1421108) Kind() string    { return "伙伴" }
func (CardDef1421108) Element() string { return "地" }

func (CardDef1421108) Card() model.Card {
	return model.Card{
		Number:          "1421108",
		Type:            "伙伴",
		Name:            "凯尔特灵鹿",
		Category:        "地",
		Tag:             "野兽",
		Description:     "诱发回合技:在一个灵媒技能被使用后,重置本卡",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 6},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421108.jpg",
	}
}

type CardDef1421109 struct{}

func (CardDef1421109) ID() string      { return "1421109" }
func (CardDef1421109) Name() string    { return "地穴巨蝠" }
func (CardDef1421109) Kind() string    { return "伙伴" }
func (CardDef1421109) Element() string { return "地" }

func (CardDef1421109) Card() model.Card {
	return model.Card{
		Number:          "1421109",
		Type:            "伙伴",
		Name:            "地穴巨蝠",
		Category:        "地",
		Tag:             "野兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{"地": 2, "暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421109.jpg",
	}
}

type CardDef1421110 struct{}

func (CardDef1421110) ID() string      { return "1421110" }
func (CardDef1421110) Name() string    { return "岩壁巨像" }
func (CardDef1421110) Kind() string    { return "伙伴" }
func (CardDef1421110) Element() string { return "地" }

func (CardDef1421110) Card() model.Card {
	return model.Card{
		Number:          "1421110",
		Type:            "伙伴",
		Name:            "岩壁巨像",
		Category:        "地",
		Tag:             "造物",
		Description:     "光环:如果你没有学习任何法术,你召唤的每个地脉伙伴+1\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421110.jpg",
	}
}

type CardDef1421111 struct{}

func (CardDef1421111) ID() string      { return "1421111" }
func (CardDef1421111) Name() string    { return "岩壁魔怪" }
func (CardDef1421111) Kind() string    { return "伙伴" }
func (CardDef1421111) Element() string { return "地" }

func (CardDef1421111) Card() model.Card {
	return model.Card{
		Number:          "1421111",
		Type:            "伙伴",
		Name:            "岩壁魔怪",
		Category:        "地",
		Tag:             "造物",
		Description:     "光环:如果你没有学习任何法术,此卡每次最多受到1点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421111.jpg",
	}
}

type CardDef1421112 struct{}

func (CardDef1421112) ID() string      { return "1421112" }
func (CardDef1421112) Name() string    { return "沙尘恶魔" }
func (CardDef1421112) Kind() string    { return "伙伴" }
func (CardDef1421112) Element() string { return "地" }

func (CardDef1421112) Card() model.Card {
	return model.Card{
		Number:          "1421112",
		Type:            "伙伴",
		Name:            "沙尘恶魔",
		Category:        "地",
		Tag:             "恶魔",
		Description:     "引魔.祈咒:使前排敌人石化1",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421112.jpg",
	}
}

type CardDef1421113 struct{}

func (CardDef1421113) ID() string      { return "1421113" }
func (CardDef1421113) Name() string    { return "岩壁修道士" }
func (CardDef1421113) Kind() string    { return "伙伴" }
func (CardDef1421113) Element() string { return "地" }

func (CardDef1421113) Card() model.Card {
	return model.Card{
		Number:          "1421113",
		Type:            "伙伴",
		Name:            "岩壁修道士",
		Category:        "地",
		Tag:             "巫师",
		Description:     "光环:如果你没有学习任何法术,每回合命中的第一个敌方法术\\攻变为0",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421113.jpg",
	}
}

type CardDef1421114 struct{}

func (CardDef1421114) ID() string      { return "1421114" }
func (CardDef1421114) Name() string    { return "巨型沙虫" }
func (CardDef1421114) Kind() string    { return "伙伴" }
func (CardDef1421114) Element() string { return "地" }

func (CardDef1421114) Card() model.Card {
	return model.Card{
		Number:          "1421114",
		Type:            "伙伴",
		Name:            "巨型沙虫",
		Category:        "地",
		Tag:             "野兽",
		Description:     "诱发:每当此卡受到伤害,获得隐蔽1",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 6},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421114.jpg",
	}
}

type CardDef1421115 struct{}

func (CardDef1421115) ID() string      { return "1421115" }
func (CardDef1421115) Name() string    { return "地卜行者" }
func (CardDef1421115) Kind() string    { return "伙伴" }
func (CardDef1421115) Element() string { return "地" }

func (CardDef1421115) Card() model.Card {
	return model.Card{
		Number:          "1421115",
		Type:            "伙伴",
		Name:            "地卜行者",
		Category:        "地",
		Tag:             "巫师",
		Description:     "入场:抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\地\\1421115.jpg",
	}
}

type CardDef1501001 struct{}

func (CardDef1501001) ID() string      { return "1501001" }
func (CardDef1501001) Name() string    { return "孪生天使" }
func (CardDef1501001) Kind() string    { return "伙伴" }
func (CardDef1501001) Element() string { return "光" }

func (CardDef1501001) Card() model.Card {
	return model.Card{
		Number:          "1501001",
		Type:            "伙伴",
		Name:            "孪生天使",
		Category:        "光",
		Tag:             "衍生-精灵",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1501001.jpg",
	}
}

type CardDef1511001 struct{}

func (CardDef1511001) ID() string      { return "1511001" }
func (CardDef1511001) Name() string    { return "白袍大贤者 掌号使" }
func (CardDef1511001) Kind() string    { return "伙伴" }
func (CardDef1511001) Element() string { return "光" }

func (CardDef1511001) Card() model.Card {
	return model.Card{
		Number:          "1511001",
		Type:            "伙伴",
		Name:            "白袍大贤者 掌号使",
		Category:        "光",
		Tag:             "传奇-巫师",
		Description:     "主动绝技:选择法力范围内的1个敌方伙伴,支付其入场花费才能发动,获得其控制权",
		Quote:           "既为巫师,承天命,领人事,广布恩泽,守一方之序,岂有为祸纷乱之理",
		ElementsCost:    map[string]int{"光": 5, "无": 3},
		ElementsGain:    map[string]int{"光": 2, "无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1511001.jpg",
	}
}

type CardDef1511002 struct{}

func (CardDef1511002) ID() string      { return "1511002" }
func (CardDef1511002) Name() string    { return "大法师 伦德萨尔" }
func (CardDef1511002) Kind() string    { return "伙伴" }
func (CardDef1511002) Element() string { return "光" }

func (CardDef1511002) Card() model.Card {
	return model.Card{
		Number:          "1511002",
		Type:            "伙伴",
		Name:            "大法师 伦德萨尔",
		Category:        "光",
		Tag:             "传奇-巫师",
		Description:     "入场,遗言:使你的一个法术永久获得+3\\威或+1\\攻",
		Quote:           "\"真是没完没了,幸好我还剩下一招,一个真正的绝活!\"",
		ElementsCost:    map[string]int{"光": 7},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1511002.jpg",
	}
}

type CardDef1511003 struct{}

func (CardDef1511003) ID() string      { return "1511003" }
func (CardDef1511003) Name() string    { return "天枢圣兽 珀伽索斯" }
func (CardDef1511003) Kind() string    { return "伙伴" }
func (CardDef1511003) Element() string { return "光" }

func (CardDef1511003) Card() model.Card {
	return model.Card{
		Number:          "1511003",
		Type:            "伙伴",
		Name:            "天枢圣兽 珀伽索斯",
		Category:        "光",
		Tag:             "传奇-异兽",
		Description:     "引魔.光环:敌方法术对天枢圣兽 珀伽索斯以外的友方单位造成伤害变为0",
		Quote:           "希望之名,慈悲之怀",
		ElementsCost:    map[string]int{"光": 6},
		ElementsGain:    map[string]int{"光": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1511003.jpg",
	}
}

type CardDef1511101 struct{}

func (CardDef1511101) ID() string      { return "1511101" }
func (CardDef1511101) Name() string    { return "末路的王子 灰烬 凯尔特" }
func (CardDef1511101) Kind() string    { return "伙伴" }
func (CardDef1511101) Element() string { return "光" }

func (CardDef1511101) Card() model.Card {
	return model.Card{
		Number:          "1511101",
		Type:            "伙伴",
		Name:            "末路的王子 灰烬 凯尔特",
		Category:        "光",
		Tag:             "传奇-人类",
		Description:     "诱发:当你的护盾被对方打破时,抽2张牌并获得2\\光.入场:获得护盾2",
		Quote:           "\"既然神无动于衷,那拯救只好由我代劳\"",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1511101.jpg",
	}
}

type CardDef1511102 struct{}

func (CardDef1511102) ID() string      { return "1511102" }
func (CardDef1511102) Name() string    { return "孤星之魂 凯拉莫将军" }
func (CardDef1511102) Kind() string    { return "伙伴" }
func (CardDef1511102) Element() string { return "光" }

func (CardDef1511102) Card() model.Card {
	return model.Card{
		Number:          "1511102",
		Type:            "伙伴",
		Name:            "孤星之魂 凯拉莫将军",
		Category:        "光",
		Tag:             "传奇-人类",
		Description:     "诱发:每当此卡受到1次敌方伤害,如果此卡没有任何相邻伙伴,获得护盾1和+1\\攻",
		Quote:           "\"我从不介意与鼠辈为敌\"",
		ElementsCost:    map[string]int{"光": 6, "无": 1},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1511102.jpg",
	}
}

type CardDef1511103 struct{}

func (CardDef1511103) ID() string      { return "1511103" }
func (CardDef1511103) Name() string    { return "\"玫瑰先知\" 洛莉" }
func (CardDef1511103) Kind() string    { return "伙伴" }
func (CardDef1511103) Element() string { return "光" }

func (CardDef1511103) Card() model.Card {
	return model.Card{
		Number:          "1511103",
		Type:            "伙伴",
		Name:            "\"玫瑰先知\" 洛莉",
		Category:        "光",
		Tag:             "传奇-巫师",
		Description:     "诱发:每当对方洗切卡组,观看敌方卡组顶3张牌,将其以任意顺序放回卡组顶或卡组底",
		Quote:           "我看到了玫瑰被鲜血侵染,却无法将这罪孽斩断",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1511103.jpg",
	}
}

type CardDef1521001 struct{}

func (CardDef1521001) ID() string      { return "1521001" }
func (CardDef1521001) Name() string    { return "治疗术士" }
func (CardDef1521001) Kind() string    { return "伙伴" }
func (CardDef1521001) Element() string { return "光" }

func (CardDef1521001) Card() model.Card {
	return model.Card{
		Number:          "1521001",
		Type:            "伙伴",
		Name:            "治疗术士",
		Category:        "光",
		Tag:             "巫师",
		Description:     "祈咒:使1个友方单位回复1\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521001.jpg",
	}
}

type CardDef1521002 struct{}

func (CardDef1521002) ID() string      { return "1521002" }
func (CardDef1521002) Name() string    { return "光铸泰坦" }
func (CardDef1521002) Kind() string    { return "伙伴" }
func (CardDef1521002) Element() string { return "光" }

func (CardDef1521002) Card() model.Card {
	return model.Card{
		Number:          "1521002",
		Type:            "伙伴",
		Name:            "光铸泰坦",
		Category:        "光",
		Tag:             "精灵",
		Description:     "入场:抽2张牌.光环:驱动、神秘和聚能法术对本卡造成的伤害+1",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 5},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521002.jpg",
	}
}

type CardDef1521003 struct{}

func (CardDef1521003) ID() string      { return "1521003" }
func (CardDef1521003) Name() string    { return "七神侍从" }
func (CardDef1521003) Kind() string    { return "伙伴" }
func (CardDef1521003) Element() string { return "光" }

func (CardDef1521003) Card() model.Card {
	return model.Card{
		Number:          "1521003",
		Type:            "伙伴",
		Name:            "七神侍从",
		Category:        "光",
		Tag:             "人类",
		Description:     "",
		Quote:           "具体待遇还得看你侍奉哪个神",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521003.jpg",
	}
}

type CardDef1521004 struct{}

func (CardDef1521004) ID() string      { return "1521004" }
func (CardDef1521004) Name() string    { return "誓约之泉的守卫" }
func (CardDef1521004) Kind() string    { return "伙伴" }
func (CardDef1521004) Element() string { return "光" }

func (CardDef1521004) Card() model.Card {
	return model.Card{
		Number:          "1521004",
		Type:            "伙伴",
		Name:            "誓约之泉的守卫",
		Category:        "光",
		Tag:             "造物",
		Description:     "",
		Quote:           "被巫师学徒们戏称为\"澡堂门卫\"",
		ElementsCost:    map[string]int{"光": 4},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521004.jpg",
	}
}

type CardDef1521005 struct{}

func (CardDef1521005) ID() string      { return "1521005" }
func (CardDef1521005) Name() string    { return "双生天使" }
func (CardDef1521005) Kind() string    { return "伙伴" }
func (CardDef1521005) Element() string { return "光" }

func (CardDef1521005) Card() model.Card {
	return model.Card{
		Number:          "1521005",
		Type:            "伙伴",
		Name:            "双生天使",
		Category:        "光",
		Tag:             "精灵",
		Description:     "入场:将1张衍生卡牌孪生天使置于你的手牌",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"1501001"},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521005.jpg",
	}
}

type CardDef1521006 struct{}

func (CardDef1521006) ID() string      { return "1521006" }
func (CardDef1521006) Name() string    { return "生命之花" }
func (CardDef1521006) Kind() string    { return "伙伴" }
func (CardDef1521006) Element() string { return "光" }

func (CardDef1521006) Card() model.Card {
	return model.Card{
		Number:          "1521006",
		Type:            "伙伴",
		Name:            "生命之花",
		Category:        "光",
		Tag:             "植物",
		Description:     "入场:使1个其他友方单位+1\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521006.jpg",
	}
}

type CardDef1521007 struct{}

func (CardDef1521007) ID() string      { return "1521007" }
func (CardDef1521007) Name() string    { return "虹之天使" }
func (CardDef1521007) Kind() string    { return "伙伴" }
func (CardDef1521007) Element() string { return "光" }

func (CardDef1521007) Card() model.Card {
	return model.Card{
		Number:          "1521007",
		Type:            "伙伴",
		Name:            "虹之天使",
		Category:        "光",
		Tag:             "精灵",
		Description:     "光环:你的光辉元素可以当做任意元素使用",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521007.jpg",
	}
}

type CardDef1521008 struct{}

func (CardDef1521008) ID() string      { return "1521008" }
func (CardDef1521008) Name() string    { return "御座的圣翼" }
func (CardDef1521008) Kind() string    { return "伙伴" }
func (CardDef1521008) Element() string { return "光" }

func (CardDef1521008) Card() model.Card {
	return model.Card{
		Number:          "1521008",
		Type:            "伙伴",
		Name:            "御座的圣翼",
		Category:        "光",
		Tag:             "精灵",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 4},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521008.jpg",
	}
}

type CardDef1521009 struct{}

func (CardDef1521009) ID() string      { return "1521009" }
func (CardDef1521009) Name() string    { return "天马骑士" }
func (CardDef1521009) Kind() string    { return "伙伴" }
func (CardDef1521009) Element() string { return "光" }

func (CardDef1521009) Card() model.Card {
	return model.Card{
		Number:          "1521009",
		Type:            "伙伴",
		Name:            "天马骑士",
		Category:        "光",
		Tag:             "人类",
		Description:     "入场:检索1张独角天马",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521009.jpg",
	}
}

type CardDef1521010 struct{}

func (CardDef1521010) ID() string      { return "1521010" }
func (CardDef1521010) Name() string    { return "神护者" }
func (CardDef1521010) Kind() string    { return "伙伴" }
func (CardDef1521010) Element() string { return "光" }

func (CardDef1521010) Card() model.Card {
	return model.Card{
		Number:          "1521010",
		Type:            "伙伴",
		Name:            "神护者",
		Category:        "光",
		Tag:             "巫师",
		Description:     "光环:此卡免疫负面状态",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 5},
		ElementsGain:    map[string]int{"光": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521010.jpg",
	}
}

type CardDef1521011 struct{}

func (CardDef1521011) ID() string      { return "1521011" }
func (CardDef1521011) Name() string    { return "日轮法师" }
func (CardDef1521011) Kind() string    { return "伙伴" }
func (CardDef1521011) Element() string { return "光" }

func (CardDef1521011) Card() model.Card {
	return model.Card{
		Number:          "1521011",
		Type:            "伙伴",
		Name:            "日轮法师",
		Category:        "光",
		Tag:             "巫师",
		Description:     "主动绝技:重置你的1个光辉法术",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 5},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521011.jpg",
	}
}

type CardDef1521012 struct{}

func (CardDef1521012) ID() string      { return "1521012" }
func (CardDef1521012) Name() string    { return "独角天马" }
func (CardDef1521012) Kind() string    { return "伙伴" }
func (CardDef1521012) Element() string { return "光" }

func (CardDef1521012) Card() model.Card {
	return model.Card{
		Number:          "1521012",
		Type:            "伙伴",
		Name:            "独角天马",
		Category:        "光",
		Tag:             "异兽",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521012.jpg",
	}
}

type CardDef1521013 struct{}

func (CardDef1521013) ID() string      { return "1521013" }
func (CardDef1521013) Name() string    { return "神火兽" }
func (CardDef1521013) Kind() string    { return "伙伴" }
func (CardDef1521013) Element() string { return "光" }

func (CardDef1521013) Card() model.Card {
	return model.Card{
		Number:          "1521013",
		Type:            "伙伴",
		Name:            "神火兽",
		Category:        "光",
		Tag:             "异兽",
		Description:     "光环:你的法术在攻击时+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3, "火": 1},
		ElementsGain:    map[string]int{"光": 1, "火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521013.jpg",
	}
}

type CardDef1521014 struct{}

func (CardDef1521014) ID() string      { return "1521014" }
func (CardDef1521014) Name() string    { return "炬之女巫" }
func (CardDef1521014) Kind() string    { return "伙伴" }
func (CardDef1521014) Element() string { return "光" }

func (CardDef1521014) Card() model.Card {
	return model.Card{
		Number:          "1521014",
		Type:            "伙伴",
		Name:            "炬之女巫",
		Category:        "光",
		Tag:             "巫师",
		Description:     "入场:本卡获得点燃2.祈咒:使1个相邻伙伴获得负载+1\\光",
		Quote:           "\"不必惧怕黑暗,我会为你带来光明\"",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521014.jpg",
	}
}

type CardDef1521015 struct{}

func (CardDef1521015) ID() string      { return "1521015" }
func (CardDef1521015) Name() string    { return "烬之女巫" }
func (CardDef1521015) Kind() string    { return "伙伴" }
func (CardDef1521015) Element() string { return "光" }

func (CardDef1521015) Card() model.Card {
	return model.Card{
		Number:          "1521015",
		Type:            "伙伴",
		Name:            "烬之女巫",
		Category:        "光",
		Tag:             "巫师",
		Description:     "入场:本卡获得点燃3.遗言:使你的1个法术永久+2\\威",
		Quote:           "\"不必担心坎坷,我会为你开启坦途\"",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{"光": 1, "火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521015.jpg",
	}
}

type CardDef1521016 struct{}

func (CardDef1521016) ID() string      { return "1521016" }
func (CardDef1521016) Name() string    { return "索洛城的坚守者" }
func (CardDef1521016) Kind() string    { return "伙伴" }
func (CardDef1521016) Element() string { return "光" }

func (CardDef1521016) Card() model.Card {
	return model.Card{
		Number:          "1521016",
		Type:            "伙伴",
		Name:            "索洛城的坚守者",
		Category:        "光",
		Tag:             "人类",
		Description:     "诱发:此卡在满血受到治疗效果时,获得+1\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\光\\1521016.jpg",
	}
}

type CardDef1521101 struct{}

func (CardDef1521101) ID() string      { return "1521101" }
func (CardDef1521101) Name() string    { return "月霞之灵" }
func (CardDef1521101) Kind() string    { return "伙伴" }
func (CardDef1521101) Element() string { return "光" }

func (CardDef1521101) Card() model.Card {
	return model.Card{
		Number:          "1521101",
		Type:            "伙伴",
		Name:            "月霞之灵",
		Category:        "光",
		Tag:             "精灵",
		Description:     "引魔.光环:你的法术+2\\威,在你使用法术攻击后失去此光环效果",
		Quote:           "冷月孤光,夜灵幽守",
		ElementsCost:    map[string]int{"光": 5},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521101.jpg",
	}
}

type CardDef1521102 struct{}

func (CardDef1521102) ID() string      { return "1521102" }
func (CardDef1521102) Name() string    { return "神圣之子" }
func (CardDef1521102) Kind() string    { return "伙伴" }
func (CardDef1521102) Element() string { return "光" }

func (CardDef1521102) Card() model.Card {
	return model.Card{
		Number:          "1521102",
		Type:            "伙伴",
		Name:            "神圣之子",
		Category:        "光",
		Tag:             "造物",
		Description:     "诱发绝技:此卡获得生命或负载时,额外获得负载+1\\光或+1\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521102.jpg",
	}
}

type CardDef1521103 struct{}

func (CardDef1521103) ID() string      { return "1521103" }
func (CardDef1521103) Name() string    { return "孤星城的守护灵" }
func (CardDef1521103) Kind() string    { return "伙伴" }
func (CardDef1521103) Element() string { return "光" }

func (CardDef1521103) Card() model.Card {
	return model.Card{
		Number:          "1521103",
		Type:            "伙伴",
		Name:            "孤星城的守护灵",
		Category:        "光",
		Tag:             "造物",
		Description:     "入场:使1个友方伙伴+1\\血.遗言:使1个友方伙伴负载+1\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521103.jpg",
	}
}

type CardDef1521104 struct{}

func (CardDef1521104) ID() string      { return "1521104" }
func (CardDef1521104) Name() string    { return "旭日之龙" }
func (CardDef1521104) Kind() string    { return "伙伴" }
func (CardDef1521104) Element() string { return "光" }

func (CardDef1521104) Card() model.Card {
	return model.Card{
		Number:          "1521104",
		Type:            "伙伴",
		Name:            "旭日之龙",
		Category:        "光",
		Tag:             "龙",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 5},
		ElementsGain:    map[string]int{"光": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521104.jpg",
	}
}

type CardDef1521105 struct{}

func (CardDef1521105) ID() string      { return "1521105" }
func (CardDef1521105) Name() string    { return "辉之都祭司" }
func (CardDef1521105) Kind() string    { return "伙伴" }
func (CardDef1521105) Element() string { return "光" }

func (CardDef1521105) Card() model.Card {
	return model.Card{
		Number:          "1521105",
		Type:            "伙伴",
		Name:            "辉之都祭司",
		Category:        "光",
		Tag:             "巫师",
		Description:     "诱发绝技:敌方法术命中时可以发动,将该法术即将造成的伤害变为目标获得等量的点燃",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3, "火": 1},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521105.jpg",
	}
}

type CardDef1521106 struct{}

func (CardDef1521106) ID() string      { return "1521106" }
func (CardDef1521106) Name() string    { return "教廷驱魔师" }
func (CardDef1521106) Kind() string    { return "伙伴" }
func (CardDef1521106) Element() string { return "光" }

func (CardDef1521106) Card() model.Card {
	return model.Card{
		Number:          "1521106",
		Type:            "伙伴",
		Name:            "教廷驱魔师",
		Category:        "光",
		Tag:             "巫师",
		Description:     "入场:移除1张友方卡牌上的所有负面状态,每移除1层获得1\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521106.jpg",
	}
}

type CardDef1521107 struct{}

func (CardDef1521107) ID() string      { return "1521107" }
func (CardDef1521107) Name() string    { return "辉之圣防军" }
func (CardDef1521107) Kind() string    { return "伙伴" }
func (CardDef1521107) Element() string { return "光" }

func (CardDef1521107) Card() model.Card {
	return model.Card{
		Number:          "1521107",
		Type:            "伙伴",
		Name:            "辉之圣防军",
		Category:        "光",
		Tag:             "人类",
		Description:     "如果上个回合友方单位受到过伤害,此卡无需入场花费",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521107.jpg",
	}
}

type CardDef1521108 struct{}

func (CardDef1521108) ID() string      { return "1521108" }
func (CardDef1521108) Name() string    { return "矛盾的骑士" }
func (CardDef1521108) Kind() string    { return "伙伴" }
func (CardDef1521108) Element() string { return "光" }

func (CardDef1521108) Card() model.Card {
	return model.Card{
		Number:          "1521108",
		Type:            "伙伴",
		Name:            "矛盾的骑士",
		Category:        "光",
		Tag:             "人类",
		Description:     "遗言:为对手召唤本卡,本卡生命上限永久-1(由对手决定位置)",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521108.jpg",
	}
}

type CardDef1521109 struct{}

func (CardDef1521109) ID() string      { return "1521109" }
func (CardDef1521109) Name() string    { return "辉之天使" }
func (CardDef1521109) Kind() string    { return "伙伴" }
func (CardDef1521109) Element() string { return "光" }

func (CardDef1521109) Card() model.Card {
	return model.Card{
		Number:          "1521109",
		Type:            "伙伴",
		Name:            "辉之天使",
		Category:        "光",
		Tag:             "精灵",
		Description:     "你可以把其他任意元素当做光辉元素使用",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521109.jpg",
	}
}

type CardDef1521110 struct{}

func (CardDef1521110) ID() string      { return "1521110" }
func (CardDef1521110) Name() string    { return "议庭言客" }
func (CardDef1521110) Kind() string    { return "伙伴" }
func (CardDef1521110) Element() string { return "光" }

func (CardDef1521110) Card() model.Card {
	return model.Card{
		Number:          "1521110",
		Type:            "伙伴",
		Name:            "议庭言客",
		Category:        "光",
		Tag:             "人类",
		Description:     "入场:将4张九霄印记洗入对手的卡组.遗言:对方将卡组里1张九霄印记洗放在卡组顶端",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2001102"},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521110.jpg",
	}
}

type CardDef1521111 struct{}

func (CardDef1521111) ID() string      { return "1521111" }
func (CardDef1521111) Name() string    { return "议庭执政官" }
func (CardDef1521111) Kind() string    { return "伙伴" }
func (CardDef1521111) Element() string { return "光" }

func (CardDef1521111) Card() model.Card {
	return model.Card{
		Number:          "1521111",
		Type:            "伙伴",
		Name:            "议庭执政官",
		Category:        "光",
		Tag:             "巫师",
		Description:     "诱发:每当敌方召唤1个伙伴,将3张九霄印记洗入对手卡组",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 5},
		ElementsGain:    map[string]int{"光": 2, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2001102"},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521111.jpg",
	}
}

type CardDef1521112 struct{}

func (CardDef1521112) ID() string      { return "1521112" }
func (CardDef1521112) Name() string    { return "议庭护法" }
func (CardDef1521112) Kind() string    { return "伙伴" }
func (CardDef1521112) Element() string { return "光" }

func (CardDef1521112) Card() model.Card {
	return model.Card{
		Number:          "1521112",
		Type:            "伙伴",
		Name:            "议庭护法",
		Category:        "光",
		Tag:             "巫师",
		Description:     "主动回合技:翻开对手卡组顶5张牌,如果其中有九霄印记则你可以调整这5张牌的顺序并放回卡组顶,否则对方洗混卡组",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3, "气": 1},
		ElementsGain:    map[string]int{"光": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2001102"},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521112.jpg",
	}
}

type CardDef1521113 struct{}

func (CardDef1521113) ID() string      { return "1521113" }
func (CardDef1521113) Name() string    { return "辉之都戒卫犬" }
func (CardDef1521113) Kind() string    { return "伙伴" }
func (CardDef1521113) Element() string { return "光" }

func (CardDef1521113) Card() model.Card {
	return model.Card{
		Number:          "1521113",
		Type:            "伙伴",
		Name:            "辉之都戒卫犬",
		Category:        "光",
		Tag:             "野兽",
		Description:     "遗言:如果此卡是被敌方杀死,翻取1个伙伴牌并使其入场花费-1\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521113.jpg",
	}
}

type CardDef1521114 struct{}

func (CardDef1521114) ID() string      { return "1521114" }
func (CardDef1521114) Name() string    { return "辉之都祈祷者" }
func (CardDef1521114) Kind() string    { return "伙伴" }
func (CardDef1521114) Element() string { return "光" }

func (CardDef1521114) Card() model.Card {
	return model.Card{
		Number:          "1521114",
		Type:            "伙伴",
		Name:            "辉之都祈祷者",
		Category:        "光",
		Tag:             "人类",
		Description:     "入场:每有1个受伤的友方单位,获得1\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521114.jpg",
	}
}

type CardDef1521115 struct{}

func (CardDef1521115) ID() string      { return "1521115" }
func (CardDef1521115) Name() string    { return "孤星铁骑士" }
func (CardDef1521115) Kind() string    { return "伙伴" }
func (CardDef1521115) Element() string { return "光" }

func (CardDef1521115) Card() model.Card {
	return model.Card{
		Number:          "1521115",
		Type:            "伙伴",
		Name:            "孤星铁骑士",
		Category:        "光",
		Tag:             "人类",
		Description:     "入场:如果此卡处于前排且没有任何相邻伙伴,此卡获得负载+1\\光+1\\血",
		Quote:           "\"铁骑士誓为孤星战至生命的最后一刻!\"",
		ElementsCost:    map[string]int{"光": 5},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\光\\1521115.jpg",
	}
}

type CardDef1601101 struct{}

func (CardDef1601101) ID() string      { return "1601101" }
func (CardDef1601101) Name() string    { return "血影之躯" }
func (CardDef1601101) Kind() string    { return "伙伴" }
func (CardDef1601101) Element() string { return "暗" }

func (CardDef1601101) Card() model.Card {
	return model.Card{
		Number:          "1601101",
		Type:            "伙伴",
		Name:            "血影之躯",
		Category:        "暗",
		Tag:             "衍生-恶魔",
		Description:     "主动:取除1个红月标记物才能发动,使你的下一次法术可以额外选择1个无视范围的目标.红月结束后将此卡替换回红月魔巫 瑟薇安娜并将其重置",
		Quote:           "\"终于...开始吧,埋葬这孱弱的世界\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1601101.jpg",
	}
}

type CardDef1611001 struct{}

func (CardDef1611001) ID() string      { return "1611001" }
func (CardDef1611001) Name() string    { return "\"观察者\" 欧柯茹" }
func (CardDef1611001) Kind() string    { return "伙伴" }
func (CardDef1611001) Element() string { return "暗" }

func (CardDef1611001) Card() model.Card {
	return model.Card{
		Number:          "1611001",
		Type:            "伙伴",
		Name:            "\"观察者\" 欧柯茹",
		Category:        "暗",
		Tag:             "传奇-恶魔",
		Description:     "入场:查看卡组顶5张牌,你可以将其抽取或以任意顺序放回卡组顶部、底部,每抽取1张,对你的人物造成1点伤害",
		Quote:           "\"接受真相,比真相本身更加残酷...\"",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1611001.jpg",
	}
}

type CardDef1611002 struct{}

func (CardDef1611002) ID() string      { return "1611002" }
func (CardDef1611002) Name() string    { return "黑袍执行官 无心" }
func (CardDef1611002) Kind() string    { return "伙伴" }
func (CardDef1611002) Element() string { return "暗" }

func (CardDef1611002) Card() model.Card {
	return model.Card{
		Number:          "1611002",
		Type:            "伙伴",
		Name:            "黑袍执行官 无心",
		Category:        "暗",
		Tag:             "传奇-巫师",
		Description:     "诱发:每当你献祭或吞噬1个伙伴,根据其生命值在此卡上放置暗影标记物.主动绝技:选择法力范围内的1个伙伴,取除其生命值数量的暗影标记物并将其消灭",
		Quote:           "\"你厌恶的事实,不过是万物轮回的流程\"",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1611002.jpg",
	}
}

type CardDef1611003 struct{}

func (CardDef1611003) ID() string      { return "1611003" }
func (CardDef1611003) Name() string    { return "\"穿心人\"" }
func (CardDef1611003) Kind() string    { return "伙伴" }
func (CardDef1611003) Element() string { return "暗" }

func (CardDef1611003) Card() model.Card {
	return model.Card{
		Number:          "1611003",
		Type:            "伙伴",
		Name:            "\"穿心人\"",
		Category:        "暗",
		Tag:             "传奇-巫师",
		Description:     "入场:将1张衍生道具幻痛加入手牌.幻痛在触发时可以额外选择1个敌方法术",
		Quote:           "\"你亲手为我刻下的伤痕,现在如数奉还\"",
		ElementsCost:    map[string]int{"暗": 6},
		ElementsGain:    map[string]int{"暗": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2601001"},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1611003.jpg",
	}
}

type CardDef1611101 struct{}

func (CardDef1611101) ID() string      { return "1611101" }
func (CardDef1611101) Name() string    { return "红月魔巫 瑟薇安娜" }
func (CardDef1611101) Kind() string    { return "伙伴" }
func (CardDef1611101) Element() string { return "暗" }

func (CardDef1611101) Card() model.Card {
	return model.Card{
		Number:          "1611101",
		Type:            "伙伴",
		Name:            "红月魔巫 瑟薇安娜",
		Category:        "暗",
		Tag:             "传奇-巫师",
		Description:     "入场,祈咒:在你的红月上放置一个红月标记物,这些标记物使红月生效时额外给其他法术+1\\威.红月生效期间,将此卡替换为血影之躯",
		Quote:           "\"如此执着,如此狡诈,如此野心...你让我回想起了,我选中你的时候\"——安迪斯",
		ElementsCost:    map[string]int{"暗": 7},
		ElementsGain:    map[string]int{"暗": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"1601101"},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1611101.jpg",
	}
}

type CardDef1611102 struct{}

func (CardDef1611102) ID() string      { return "1611102" }
func (CardDef1611102) Name() string    { return "蔷薇花园的血荆棘" }
func (CardDef1611102) Kind() string    { return "伙伴" }
func (CardDef1611102) Element() string { return "暗" }

func (CardDef1611102) Card() model.Card {
	return model.Card{
		Number:          "1611102",
		Type:            "伙伴",
		Name:            "蔷薇花园的血荆棘",
		Category:        "暗",
		Tag:             "传奇-植物",
		Description:     "遗言:因友方卡牌攻击或效果而死亡的场合,可以花费1\\暗重新召唤",
		Quote:           "\"伯爵,如果我没记错的话,好像它比昨天繁茂了不少\"",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1611102.jpg",
	}
}

type CardDef1611103 struct{}

func (CardDef1611103) ID() string      { return "1611103" }
func (CardDef1611103) Name() string    { return "鲜血贵公子 罗伯特 黑松" }
func (CardDef1611103) Kind() string    { return "伙伴" }
func (CardDef1611103) Element() string { return "暗" }

func (CardDef1611103) Card() model.Card {
	return model.Card{
		Number:          "1611103",
		Type:            "伙伴",
		Name:            "鲜血贵公子 罗伯特 黑松",
		Category:        "暗",
		Tag:             "传奇-巫师",
		Description:     "诱发:每当1个友方单位受到友方伤害,放置1个标记物,每当友方单位因友方伤害或效果死亡,放置2个标记物.主动回合技:移除3个标记物,此卡获得+1\\血或负载+1\\暗或+1\\攻",
		Quote:           "\"这并不是我的爱好,这是我的生活\"",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1611103.jpg",
	}
}

type CardDef1621001 struct{}

func (CardDef1621001) ID() string      { return "1621001" }
func (CardDef1621001) Name() string    { return "冥界信鸽" }
func (CardDef1621001) Kind() string    { return "伙伴" }
func (CardDef1621001) Element() string { return "暗" }

func (CardDef1621001) Card() model.Card {
	return model.Card{
		Number:          "1621001",
		Type:            "伙伴",
		Name:            "冥界信鸽",
		Category:        "暗",
		Tag:             "野兽",
		Description:     "遗言:抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "暗": 1},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621001.jpg",
	}
}

type CardDef1621002 struct{}

func (CardDef1621002) ID() string      { return "1621002" }
func (CardDef1621002) Name() string    { return "元素躯壳" }
func (CardDef1621002) Kind() string    { return "伙伴" }
func (CardDef1621002) Element() string { return "暗" }

func (CardDef1621002) Card() model.Card {
	return model.Card{
		Number:          "1621002",
		Type:            "伙伴",
		Name:            "元素躯壳",
		Category:        "暗",
		Tag:             "造物",
		Description:     "遗言:获得1\\无",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621002.jpg",
	}
}

type CardDef1621003 struct{}

func (CardDef1621003) ID() string      { return "1621003" }
func (CardDef1621003) Name() string    { return "恐惧魔" }
func (CardDef1621003) Kind() string    { return "伙伴" }
func (CardDef1621003) Element() string { return "暗" }

func (CardDef1621003) Card() model.Card {
	return model.Card{
		Number:          "1621003",
		Type:            "伙伴",
		Name:            "恐惧魔",
		Category:        "暗",
		Tag:             "恶魔",
		Description:     "吞噬:3\\血",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621003.jpg",
	}
}

type CardDef1621004 struct{}

func (CardDef1621004) ID() string      { return "1621004" }
func (CardDef1621004) Name() string    { return "巫术祭司" }
func (CardDef1621004) Kind() string    { return "伙伴" }
func (CardDef1621004) Element() string { return "暗" }

func (CardDef1621004) Card() model.Card {
	return model.Card{
		Number:          "1621004",
		Type:            "伙伴",
		Name:            "巫术祭司",
		Category:        "暗",
		Tag:             "巫师",
		Description:     "主动绝技:献祭你的1个伙伴,使另一个角色获得其生命值",
		Quote:           "\"又到了献祭的时刻......\"",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621004.jpg",
	}
}

type CardDef1621005 struct{}

func (CardDef1621005) ID() string      { return "1621005" }
func (CardDef1621005) Name() string    { return "诅咒魔像" }
func (CardDef1621005) Kind() string    { return "伙伴" }
func (CardDef1621005) Element() string { return "暗" }

func (CardDef1621005) Card() model.Card {
	return model.Card{
		Number:          "1621005",
		Type:            "伙伴",
		Name:            "诅咒魔像",
		Category:        "暗",
		Tag:             "造物",
		Description:     "入场:使1个敌方法术获得虚弱2",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621005.jpg",
	}
}

type CardDef1621006 struct{}

func (CardDef1621006) ID() string      { return "1621006" }
func (CardDef1621006) Name() string    { return "梦魇" }
func (CardDef1621006) Kind() string    { return "伙伴" }
func (CardDef1621006) Element() string { return "暗" }

func (CardDef1621006) Card() model.Card {
	return model.Card{
		Number:          "1621006",
		Type:            "伙伴",
		Name:            "梦魇",
		Category:        "暗",
		Tag:             "恶魔",
		Description:     "诱发:每当其他友方单位死亡后,此卡获得+1\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621006.jpg",
	}
}

type CardDef1621007 struct{}

func (CardDef1621007) ID() string      { return "1621007" }
func (CardDef1621007) Name() string    { return "巫师的人偶" }
func (CardDef1621007) Kind() string    { return "伙伴" }
func (CardDef1621007) Element() string { return "暗" }

func (CardDef1621007) Card() model.Card {
	return model.Card{
		Number:          "1621007",
		Type:            "伙伴",
		Name:            "巫师的人偶",
		Category:        "暗",
		Tag:             "造物",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"光": 1, "暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621007.jpg",
	}
}

type CardDef1621008 struct{}

func (CardDef1621008) ID() string      { return "1621008" }
func (CardDef1621008) Name() string    { return "南境奴隶" }
func (CardDef1621008) Kind() string    { return "伙伴" }
func (CardDef1621008) Element() string { return "暗" }

func (CardDef1621008) Card() model.Card {
	return model.Card{
		Number:          "1621008",
		Type:            "伙伴",
		Name:            "南境奴隶",
		Category:        "暗",
		Tag:             "人类",
		Description:     "",
		Quote:           "不是不报,时候未到",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621008.jpg",
	}
}

type CardDef1621009 struct{}

func (CardDef1621009) ID() string      { return "1621009" }
func (CardDef1621009) Name() string    { return "唤魔邪术士" }
func (CardDef1621009) Kind() string    { return "伙伴" }
func (CardDef1621009) Element() string { return "暗" }

func (CardDef1621009) Card() model.Card {
	return model.Card{
		Number:          "1621009",
		Type:            "伙伴",
		Name:            "唤魔邪术士",
		Category:        "暗",
		Tag:             "巫师",
		Description:     "诱发回合技:在你的1个伙伴死亡后,检索1个暗影造物或恶魔",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 5},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621009.jpg",
	}
}

type CardDef1621010 struct{}

func (CardDef1621010) ID() string      { return "1621010" }
func (CardDef1621010) Name() string    { return "恶魔尊主" }
func (CardDef1621010) Kind() string    { return "伙伴" }
func (CardDef1621010) Element() string { return "暗" }

func (CardDef1621010) Card() model.Card {
	return model.Card{
		Number:          "1621010",
		Type:            "伙伴",
		Name:            "恶魔尊主",
		Category:        "暗",
		Tag:             "恶魔",
		Description:     "吞噬:4\\暗",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"暗": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621010.jpg",
	}
}

type CardDef1621011 struct{}

func (CardDef1621011) ID() string      { return "1621011" }
func (CardDef1621011) Name() string    { return "白骨骑士" }
func (CardDef1621011) Kind() string    { return "伙伴" }
func (CardDef1621011) Element() string { return "暗" }

func (CardDef1621011) Card() model.Card {
	return model.Card{
		Number:          "1621011",
		Type:            "伙伴",
		Name:            "白骨骑士",
		Category:        "暗",
		Tag:             "造物",
		Description:     "遗言:重新召唤此伙伴,并失去此遗言",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621011.jpg",
	}
}

type CardDef1621012 struct{}

func (CardDef1621012) ID() string      { return "1621012" }
func (CardDef1621012) Name() string    { return "灵魂祭司" }
func (CardDef1621012) Kind() string    { return "伙伴" }
func (CardDef1621012) Element() string { return "暗" }

func (CardDef1621012) Card() model.Card {
	return model.Card{
		Number:          "1621012",
		Type:            "伙伴",
		Name:            "灵魂祭司",
		Category:        "暗",
		Tag:             "巫师",
		Description:     "主动绝技:献祭1个友方伙伴,抽2张牌",
		Quote:           "\"人人都在卖命,至少我这的价格是公平的\"",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621012.jpg",
	}
}

type CardDef1621013 struct{}

func (CardDef1621013) ID() string      { return "1621013" }
func (CardDef1621013) Name() string    { return "言灵" }
func (CardDef1621013) Kind() string    { return "伙伴" }
func (CardDef1621013) Element() string { return "暗" }

func (CardDef1621013) Card() model.Card {
	return model.Card{
		Number:          "1621013",
		Type:            "伙伴",
		Name:            "言灵",
		Category:        "暗",
		Tag:             "造物",
		Description:     "诱发回合技:对方使用技能后,可以使敌方所有横置的法术虚弱1",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 6},
		ElementsGain:    map[string]int{"无": 1, "暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621013.jpg",
	}
}

type CardDef1621014 struct{}

func (CardDef1621014) ID() string      { return "1621014" }
func (CardDef1621014) Name() string    { return "恶魔仆从" }
func (CardDef1621014) Kind() string    { return "伙伴" }
func (CardDef1621014) Element() string { return "暗" }

func (CardDef1621014) Card() model.Card {
	return model.Card{
		Number:          "1621014",
		Type:            "伙伴",
		Name:            "恶魔仆从",
		Category:        "暗",
		Tag:             "恶魔",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621014.jpg",
	}
}

type CardDef1621015 struct{}

func (CardDef1621015) ID() string      { return "1621015" }
func (CardDef1621015) Name() string    { return "人面枭" }
func (CardDef1621015) Kind() string    { return "伙伴" }
func (CardDef1621015) Element() string { return "暗" }

func (CardDef1621015) Card() model.Card {
	return model.Card{
		Number:          "1621015",
		Type:            "伙伴",
		Name:            "人面枭",
		Category:        "暗",
		Tag:             "异兽",
		Description:     "",
		Quote:           "不要去问枭鸟为谁鸣",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621015.jpg",
	}
}

type CardDef1621016 struct{}

func (CardDef1621016) ID() string      { return "1621016" }
func (CardDef1621016) Name() string    { return "复仇死者" }
func (CardDef1621016) Kind() string    { return "伙伴" }
func (CardDef1621016) Element() string { return "暗" }

func (CardDef1621016) Card() model.Card {
	return model.Card{
		Number:          "1621016",
		Type:            "伙伴",
		Name:            "复仇死者",
		Category:        "暗",
		Tag:             "造物",
		Description:     "遗言:对此卡造成致命伤害来源一方的人物牌受到2点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\伙伴\\暗\\1621016.jpg",
	}
}

type CardDef1621101 struct{}

func (CardDef1621101) ID() string      { return "1621101" }
func (CardDef1621101) Name() string    { return "苦痛之魂" }
func (CardDef1621101) Kind() string    { return "伙伴" }
func (CardDef1621101) Element() string { return "暗" }

func (CardDef1621101) Card() model.Card {
	return model.Card{
		Number:          "1621101",
		Type:            "伙伴",
		Name:            "苦痛之魂",
		Category:        "暗",
		Tag:             "造物",
		Description:     "诱发回合技:在此卡受到伤害后,获得负载+1\\暗",
		Quote:           "天知道这些贵族在城堡里饲养了些什么",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621101.jpg",
	}
}

type CardDef1621102 struct{}

func (CardDef1621102) ID() string      { return "1621102" }
func (CardDef1621102) Name() string    { return "苦痛复仇者" }
func (CardDef1621102) Kind() string    { return "伙伴" }
func (CardDef1621102) Element() string { return "暗" }

func (CardDef1621102) Card() model.Card {
	return model.Card{
		Number:          "1621102",
		Type:            "伙伴",
		Name:            "苦痛复仇者",
		Category:        "暗",
		Tag:             "造物",
		Description:     "诱发回合技:在此卡受到伤害后,获得+1\\攻",
		Quote:           "等你知道了,你就没命了",
		ElementsCost:    map[string]int{"暗": 6},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621102.jpg",
	}
}

type CardDef1621103 struct{}

func (CardDef1621103) ID() string      { return "1621103" }
func (CardDef1621103) Name() string    { return "鲜血傀儡" }
func (CardDef1621103) Kind() string    { return "伙伴" }
func (CardDef1621103) Element() string { return "暗" }

func (CardDef1621103) Card() model.Card {
	return model.Card{
		Number:          "1621103",
		Type:            "伙伴",
		Name:            "鲜血傀儡",
		Category:        "暗",
		Tag:             "造物",
		Description:     "入场:对你的人物造成2点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621103.jpg",
	}
}

type CardDef1621104 struct{}

func (CardDef1621104) ID() string      { return "1621104" }
func (CardDef1621104) Name() string    { return "蔷薇花园园丁" }
func (CardDef1621104) Kind() string    { return "伙伴" }
func (CardDef1621104) Element() string { return "暗" }

func (CardDef1621104) Card() model.Card {
	return model.Card{
		Number:          "1621104",
		Type:            "伙伴",
		Name:            "蔷薇花园园丁",
		Category:        "暗",
		Tag:             "人类",
		Description:     "诱发回合技:当1个单位死亡,使1个友方单位回复2\\血",
		Quote:           "黑魔法不仅能杀人,也能滋养人",
		ElementsCost:    map[string]int{"光": 1, "暗": 4},
		ElementsGain:    map[string]int{"光": 1, "暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621104.jpg",
	}
}

type CardDef1621105 struct{}

func (CardDef1621105) ID() string      { return "1621105" }
func (CardDef1621105) Name() string    { return "混沌胚胎" }
func (CardDef1621105) Kind() string    { return "伙伴" }
func (CardDef1621105) Element() string { return "暗" }

func (CardDef1621105) Card() model.Card {
	return model.Card{
		Number:          "1621105",
		Type:            "伙伴",
		Name:            "混沌胚胎",
		Category:        "暗",
		Tag:             "造物",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2, "暗": 2},
		ElementsGain:    map[string]int{"光": 1, "暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621105.jpg",
	}
}

type CardDef1621106 struct{}

func (CardDef1621106) ID() string      { return "1621106" }
func (CardDef1621106) Name() string    { return "猎魂者" }
func (CardDef1621106) Kind() string    { return "伙伴" }
func (CardDef1621106) Element() string { return "暗" }

func (CardDef1621106) Card() model.Card {
	return model.Card{
		Number:          "1621106",
		Type:            "伙伴",
		Name:            "猎魂者",
		Category:        "暗",
		Tag:             "恶魔",
		Description:     "诱发回合技:当你的法术命中后,在其上放置1个灵魂标记物,每个灵魂标记物使其永久+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 6},
		ElementsGain:    map[string]int{"暗": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621106.jpg",
	}
}

type CardDef1621107 struct{}

func (CardDef1621107) ID() string      { return "1621107" }
func (CardDef1621107) Name() string    { return "蔷薇死神" }
func (CardDef1621107) Kind() string    { return "伙伴" }
func (CardDef1621107) Element() string { return "暗" }

func (CardDef1621107) Card() model.Card {
	return model.Card{
		Number:          "1621107",
		Type:            "伙伴",
		Name:            "蔷薇死神",
		Category:        "暗",
		Tag:             "恶魔",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 6},
		ElementsGain:    map[string]int{"暗": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621107.jpg",
	}
}

type CardDef1621108 struct{}

func (CardDef1621108) ID() string      { return "1621108" }
func (CardDef1621108) Name() string    { return "恶魔之子" }
func (CardDef1621108) Kind() string    { return "伙伴" }
func (CardDef1621108) Element() string { return "暗" }

func (CardDef1621108) Card() model.Card {
	return model.Card{
		Number:          "1621108",
		Type:            "伙伴",
		Name:            "恶魔之子",
		Category:        "暗",
		Tag:             "人类",
		Description:     "吞噬:1个暗影伙伴",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621108.jpg",
	}
}

type CardDef1621109 struct{}

func (CardDef1621109) ID() string      { return "1621109" }
func (CardDef1621109) Name() string    { return "猩红之翼" }
func (CardDef1621109) Kind() string    { return "伙伴" }
func (CardDef1621109) Element() string { return "暗" }

func (CardDef1621109) Card() model.Card {
	return model.Card{
		Number:          "1621109",
		Type:            "伙伴",
		Name:            "猩红之翼",
		Category:        "暗",
		Tag:             "恶魔",
		Description:     "诱发:每次红月生效后,对法力范围内1个单位造成1点伤害并使此卡+1\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621109.jpg",
	}
}

type CardDef1621110 struct{}

func (CardDef1621110) ID() string      { return "1621110" }
func (CardDef1621110) Name() string    { return "猩红魔兽" }
func (CardDef1621110) Kind() string    { return "伙伴" }
func (CardDef1621110) Element() string { return "暗" }

func (CardDef1621110) Card() model.Card {
	return model.Card{
		Number:          "1621110",
		Type:            "伙伴",
		Name:            "猩红魔兽",
		Category:        "暗",
		Tag:             "野兽",
		Description:     "引魔.光环:红月生效期间,你的暗影法术+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621110.jpg",
	}
}

type CardDef1621111 struct{}

func (CardDef1621111) ID() string      { return "1621111" }
func (CardDef1621111) Name() string    { return "红月先知" }
func (CardDef1621111) Kind() string    { return "伙伴" }
func (CardDef1621111) Element() string { return "暗" }

func (CardDef1621111) Card() model.Card {
	return model.Card{
		Number:          "1621111",
		Type:            "伙伴",
		Name:            "红月先知",
		Category:        "暗",
		Tag:             "巫师",
		Description:     "入场,遗言:使你当前(如果红月已生效)或下一次红月的冷却层数-1",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621111.jpg",
	}
}

type CardDef1621112 struct{}

func (CardDef1621112) ID() string      { return "1621112" }
func (CardDef1621112) Name() string    { return "谧语精灵猎手" }
func (CardDef1621112) Kind() string    { return "伙伴" }
func (CardDef1621112) Element() string { return "暗" }

func (CardDef1621112) Card() model.Card {
	return model.Card{
		Number:          "1621112",
		Type:            "伙伴",
		Name:            "谧语精灵猎手",
		Category:        "暗",
		Tag:             "精灵",
		Description:     "遗言:对任意1个敌人造成1点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621112.jpg",
	}
}

type CardDef1621113 struct{}

func (CardDef1621113) ID() string      { return "1621113" }
func (CardDef1621113) Name() string    { return "谧语精灵祭司" }
func (CardDef1621113) Kind() string    { return "伙伴" }
func (CardDef1621113) Element() string { return "暗" }

func (CardDef1621113) Card() model.Card {
	return model.Card{
		Number:          "1621113",
		Type:            "伙伴",
		Name:            "谧语精灵祭司",
		Category:        "暗",
		Tag:             "精灵",
		Description:     "遗言:使1个友方伙伴负载+1\\暗",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "暗": 1},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            2,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621113.jpg",
	}
}

type CardDef1621114 struct{}

func (CardDef1621114) ID() string      { return "1621114" }
func (CardDef1621114) Name() string    { return "灵魂共生体" }
func (CardDef1621114) Kind() string    { return "伙伴" }
func (CardDef1621114) Element() string { return "暗" }

func (CardDef1621114) Card() model.Card {
	return model.Card{
		Number:          "1621114",
		Type:            "伙伴",
		Name:            "灵魂共生体",
		Category:        "暗",
		Tag:             "造物",
		Description:     "遗言:给你的最多2个法术放置1个灵魂标记物,每个灵魂标记物使该法术+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 5},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            3,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621114.jpg",
	}
}

type CardDef1621115 struct{}

func (CardDef1621115) ID() string      { return "1621115" }
func (CardDef1621115) Name() string    { return "灵魂吸食者" }
func (CardDef1621115) Kind() string    { return "伙伴" }
func (CardDef1621115) Element() string { return "暗" }

func (CardDef1621115) Card() model.Card {
	return model.Card{
		Number:          "1621115",
		Type:            "伙伴",
		Name:            "灵魂吸食者",
		Category:        "暗",
		Tag:             "恶魔",
		Description:     "引魔.主动回合技:移除你场上的1个灵魂标记物,抽2张牌并获得2\\暗",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 6},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            4,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\伙伴\\暗\\1621115.jpg",
	}
}

type CardDef2001101 struct{}

func (CardDef2001101) ID() string      { return "2001101" }
func (CardDef2001101) Name() string    { return "落幕提琴" }
func (CardDef2001101) Kind() string    { return "道具" }
func (CardDef2001101) Element() string { return "无" }

func (CardDef2001101) Card() model.Card {
	return model.Card{
		Number:          "2001101",
		Type:            "道具",
		Name:            "落幕提琴",
		Category:        "无",
		Tag:             "衍生-装备",
		Description:     "花费2\\无才能进行攻击",
		Quote:           "\"乐曲终会完结,你总要走出这里\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2001101.jpg",
	}
}

type CardDef2001102 struct{}

func (CardDef2001102) ID() string      { return "2001102" }
func (CardDef2001102) Name() string    { return "九霄印记" }
func (CardDef2001102) Kind() string    { return "道具" }
func (CardDef2001102) Element() string { return "无" }

func (CardDef2001102) Card() model.Card {
	return model.Card{
		Number:          "2001102",
		Type:            "道具",
		Name:            "九霄印记",
		Category:        "无",
		Tag:             "衍生-消耗品",
		Description:     "当你丢弃本卡时,对你的人物造成2点伤害",
		Quote:           "\"是的,我们是自愿退出原则\"",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2001102.jpg",
	}
}

type CardDef2011001 struct{}

func (CardDef2011001) ID() string      { return "2011001" }
func (CardDef2011001) Name() string    { return "大法师之杖" }
func (CardDef2011001) Kind() string    { return "道具" }
func (CardDef2011001) Element() string { return "无" }

func (CardDef2011001) Card() model.Card {
	return model.Card{
		Number:          "2011001",
		Type:            "道具",
		Name:            "大法师之杖",
		Category:        "无",
		Tag:             "传奇-装备-武器",
		Description:     "入场:从你的技能池将1个法术置于此卡上.主动绝技:花费元素使用此卡上的1个技能,然后将该卡牌从游戏中移除",
		Quote:           "等你成为了大法师,你的法杖就是大法师之杖",
		ElementsCost:    map[string]int{"无": 7},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2011001.jpg",
	}
}

type CardDef2011002 struct{}

func (CardDef2011002) ID() string      { return "2011002" }
func (CardDef2011002) Name() string    { return "统御者之冠" }
func (CardDef2011002) Kind() string    { return "道具" }
func (CardDef2011002) Element() string { return "无" }

func (CardDef2011002) Card() model.Card {
	return model.Card{
		Number:          "2011002",
		Type:            "道具",
		Name:            "统御者之冠",
		Category:        "无",
		Tag:             "传奇-装备-饰物",
		Description:     "入场:此后本局游戏你召唤的所有伙伴负载变为0",
		Quote:           "巫师盟主年年换,今年花又落谁家",
		ElementsCost:    map[string]int{"无": 8},
		ElementsGain:    map[string]int{"无": 6},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2011002.jpg",
	}
}

type CardDef2011003 struct{}

func (CardDef2011003) ID() string      { return "2011003" }
func (CardDef2011003) Name() string    { return "君王法袍 至贤" }
func (CardDef2011003) Kind() string    { return "道具" }
func (CardDef2011003) Element() string { return "无" }

func (CardDef2011003) Card() model.Card {
	return model.Card{
		Number:          "2011003",
		Type:            "道具",
		Name:            "君王法袍 至贤",
		Category:        "无",
		Tag:             "传奇-装备-防具",
		Description:     "诱发:当敌方法术命中时,你可以将1张技能牌从技能池移出游戏来发动,该敌方攻击法术在本回合-2\\攻",
		Quote:           "它实在是太重了,大部分时间都被挂在架子上",
		ElementsCost:    map[string]int{"无": 5},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2011003.jpg",
	}
}

type CardDef2011101 struct{}

func (CardDef2011101) ID() string      { return "2011101" }
func (CardDef2011101) Name() string    { return "奥术铠甲 天穹" }
func (CardDef2011101) Kind() string    { return "道具" }
func (CardDef2011101) Element() string { return "无" }

func (CardDef2011101) Card() model.Card {
	return model.Card{
		Number:          "2011101",
		Type:            "道具",
		Name:            "奥术铠甲 天穹",
		Category:        "无",
		Tag:             "传奇-装备-防具",
		Description:     "光环:你的护盾在回合结束时不会减少.入场:获得护盾2,本局游戏中你不能再获得护盾",
		Quote:           "它最大的弱点是需要他人帮忙穿上,至于为什么它的主人死后总是侍从接手,你就别问了",
		ElementsCost:    map[string]int{"无": 6},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2011101.jpg",
	}
}

type CardDef2011102 struct{}

func (CardDef2011102) ID() string      { return "2011102" }
func (CardDef2011102) Name() string    { return "预知宝珠" }
func (CardDef2011102) Kind() string    { return "道具" }
func (CardDef2011102) Element() string { return "无" }

func (CardDef2011102) Card() model.Card {
	return model.Card{
		Number:          "2011102",
		Type:            "道具",
		Name:            "预知宝珠",
		Category:        "无",
		Tag:             "传奇-装备-神器",
		Description:     "光环:你可以随时查看你的卡组顶3张牌.主动:消耗此卡才能发动,将你的卡组顶3张牌以任意顺序放回卡组顶或卡组底",
		Quote:           "预知不远处的未来,或者尝试改变它",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2011102.jpg",
	}
}

type CardDef2021001 struct{}

func (CardDef2021001) ID() string      { return "2021001" }
func (CardDef2021001) Name() string    { return "秘法宝典" }
func (CardDef2021001) Kind() string    { return "道具" }
func (CardDef2021001) Element() string { return "无" }

func (CardDef2021001) Card() model.Card {
	return model.Card{
		Number:          "2021001",
		Type:            "道具",
		Name:            "秘法宝典",
		Category:        "无",
		Tag:             "装备",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021001.jpg",
	}
}

type CardDef2021002 struct{}

func (CardDef2021002) ID() string      { return "2021002" }
func (CardDef2021002) Name() string    { return "记忆项链" }
func (CardDef2021002) Kind() string    { return "道具" }
func (CardDef2021002) Element() string { return "无" }

func (CardDef2021002) Card() model.Card {
	return model.Card{
		Number:          "2021002",
		Type:            "道具",
		Name:            "记忆项链",
		Category:        "无",
		Tag:             "装备-饰物",
		Description:     "光环:你的技能槽位+1",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021002.jpg",
	}
}

type CardDef2021003 struct{}

func (CardDef2021003) ID() string      { return "2021003" }
func (CardDef2021003) Name() string    { return "随心魔杖" }
func (CardDef2021003) Kind() string    { return "道具" }
func (CardDef2021003) Element() string { return "无" }

func (CardDef2021003) Card() model.Card {
	return model.Card{
		Number:          "2021003",
		Type:            "道具",
		Name:            "随心魔杖",
		Category:        "无",
		Tag:             "装备-武器",
		Description:     "主动:消耗此卡才能发动,将你的1个使用花费小于3的法术重置",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 4},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021003.jpg",
	}
}

type CardDef2021004 struct{}

func (CardDef2021004) ID() string      { return "2021004" }
func (CardDef2021004) Name() string    { return "巫师权杖" }
func (CardDef2021004) Kind() string    { return "道具" }
func (CardDef2021004) Element() string { return "无" }

func (CardDef2021004) Card() model.Card {
	return model.Card{
		Number:          "2021004",
		Type:            "道具",
		Name:            "巫师权杖",
		Category:        "无",
		Tag:             "装备-武器",
		Description:     "光环:你的法术+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 5},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021004.jpg",
	}
}

type CardDef2021005 struct{}

func (CardDef2021005) ID() string      { return "2021005" }
func (CardDef2021005) Name() string    { return "瓶装元素" }
func (CardDef2021005) Kind() string    { return "道具" }
func (CardDef2021005) Element() string { return "无" }

func (CardDef2021005) Card() model.Card {
	return model.Card{
		Number:          "2021005",
		Type:            "道具",
		Name:            "瓶装元素",
		Category:        "无",
		Tag:             "消耗品-药剂",
		Description:     "获得1\\无",
		Quote:           "如果你的包还没装满,塞一瓶这个总归能派上用场",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021005.jpg",
	}
}

type CardDef2021006 struct{}

func (CardDef2021006) ID() string      { return "2021006" }
func (CardDef2021006) Name() string    { return "百宝锦囊" }
func (CardDef2021006) Kind() string    { return "道具" }
func (CardDef2021006) Element() string { return "无" }

func (CardDef2021006) Card() model.Card {
	return model.Card{
		Number:          "2021006",
		Type:            "道具",
		Name:            "百宝锦囊",
		Category:        "无",
		Tag:             "装备",
		Description:     "主动:消耗并献祭此卡才能发动,从卡组检索1张消耗品道具牌",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021006.jpg",
	}
}

type CardDef2021007 struct{}

func (CardDef2021007) ID() string      { return "2021007" }
func (CardDef2021007) Name() string    { return "巫师齐射线列" }
func (CardDef2021007) Kind() string    { return "道具" }
func (CardDef2021007) Element() string { return "无" }

func (CardDef2021007) Card() model.Card {
	return model.Card{
		Number:          "2021007",
		Type:            "道具",
		Name:            "巫师齐射线列",
		Category:        "无",
		Tag:             "消耗品-卷轴",
		Description:     "如果你的场上至少有7个伙伴,重置你的一个法术,下一次它的范围变成AOE:前排",
		Quote:           "\"准备,1,2…我还没数3呢!\"",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021007.jpg",
	}
}

type CardDef2021008 struct{}

func (CardDef2021008) ID() string      { return "2021008" }
func (CardDef2021008) Name() string    { return "魔法石" }
func (CardDef2021008) Kind() string    { return "道具" }
func (CardDef2021008) Element() string { return "无" }

func (CardDef2021008) Card() model.Card {
	return model.Card{
		Number:          "2021008",
		Type:            "道具",
		Name:            "魔法石",
		Category:        "无",
		Tag:             "装备",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 4},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021008.jpg",
	}
}

type CardDef2021009 struct{}

func (CardDef2021009) ID() string      { return "2021009" }
func (CardDef2021009) Name() string    { return "誓约之戒" }
func (CardDef2021009) Kind() string    { return "道具" }
func (CardDef2021009) Element() string { return "无" }

func (CardDef2021009) Card() model.Card {
	return model.Card{
		Number:          "2021009",
		Type:            "道具",
		Name:            "誓约之戒",
		Category:        "无",
		Tag:             "装备-饰物",
		Description:     "光环:你的法术在攻击和强化攻击时-2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021009.jpg",
	}
}

type CardDef2021010 struct{}

func (CardDef2021010) ID() string      { return "2021010" }
func (CardDef2021010) Name() string    { return "封印卷轴" }
func (CardDef2021010) Kind() string    { return "道具" }
func (CardDef2021010) Element() string { return "无" }

func (CardDef2021010) Card() model.Card {
	return model.Card{
		Number:          "2021010",
		Type:            "道具",
		Name:            "封印卷轴",
		Category:        "无",
		Tag:             "消耗品-卷轴",
		Description:     "如果敌方有4个及以上的技能,选择其中1个,使其直到下个回合结束不能使用",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021010.jpg",
	}
}

type CardDef2021011 struct{}

func (CardDef2021011) ID() string      { return "2021011" }
func (CardDef2021011) Name() string    { return "生命护符" }
func (CardDef2021011) Kind() string    { return "道具" }
func (CardDef2021011) Element() string { return "无" }

func (CardDef2021011) Card() model.Card {
	return model.Card{
		Number:          "2021011",
		Type:            "道具",
		Name:            "生命护符",
		Category:        "无",
		Tag:             "装备-饰物",
		Description:     "入场:使1个友方角色+1\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021011.jpg",
	}
}

type CardDef2021012 struct{}

func (CardDef2021012) ID() string      { return "2021012" }
func (CardDef2021012) Name() string    { return "速写卷轴" }
func (CardDef2021012) Kind() string    { return "道具" }
func (CardDef2021012) Element() string { return "无" }

func (CardDef2021012) Card() model.Card {
	return model.Card{
		Number:          "2021012",
		Type:            "道具",
		Name:            "速写卷轴",
		Category:        "无",
		Tag:             "消耗品-卷轴",
		Description:     "释放1个你已经学习的法术并支付其使用花费,无需消耗该法术",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021012.jpg",
	}
}

type CardDef2021013 struct{}

func (CardDef2021013) ID() string      { return "2021013" }
func (CardDef2021013) Name() string    { return "断绝之刃" }
func (CardDef2021013) Kind() string    { return "道具" }
func (CardDef2021013) Element() string { return "无" }

func (CardDef2021013) Card() model.Card {
	return model.Card{
		Number:          "2021013",
		Type:            "道具",
		Name:            "断绝之刃",
		Category:        "无",
		Tag:             "装备-武器",
		Description:     "光环:你的法术攻击和强化攻击时+2\\威,你的法术无法用于防御",
		Quote:           "\"一帮盛气凌人的王子王孙...我真是受够了\"",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021013.jpg",
	}
}

type CardDef2021014 struct{}

func (CardDef2021014) ID() string      { return "2021014" }
func (CardDef2021014) Name() string    { return "法力增强剂A型" }
func (CardDef2021014) Kind() string    { return "道具" }
func (CardDef2021014) Element() string { return "无" }

func (CardDef2021014) Card() model.Card {
	return model.Card{
		Number:          "2021014",
		Type:            "道具",
		Name:            "法力增强剂A型",
		Category:        "无",
		Tag:             "消耗品-药剂",
		Description:     "本回合你的下1次技能使用花费为0",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021014.jpg",
	}
}

type CardDef2021015 struct{}

func (CardDef2021015) ID() string      { return "2021015" }
func (CardDef2021015) Name() string    { return "法力增强剂C型" }
func (CardDef2021015) Kind() string    { return "道具" }
func (CardDef2021015) Element() string { return "无" }

func (CardDef2021015) Card() model.Card {
	return model.Card{
		Number:          "2021015",
		Type:            "道具",
		Name:            "法力增强剂C型",
		Category:        "无",
		Tag:             "消耗品-药剂",
		Description:     "本回合你的法术使用花费为0,但在使用后获得冷却2",
		Quote:           "有相关研究表明,其记忆衰退的副作用是由于过量的法力涌入,以致对大脑造成不可逆的伤害.",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021015.jpg",
	}
}

type CardDef2021016 struct{}

func (CardDef2021016) ID() string      { return "2021016" }
func (CardDef2021016) Name() string    { return "纹饰佩剑" }
func (CardDef2021016) Kind() string    { return "道具" }
func (CardDef2021016) Element() string { return "无" }

func (CardDef2021016) Card() model.Card {
	return model.Card{
		Number:          "2021016",
		Type:            "道具",
		Name:            "纹饰佩剑",
		Category:        "无",
		Tag:             "装备-武器",
		Description:     "",
		Quote:           "北加雷亚大陆的巫师通常都带有佩剑以示其贵族身份,当然砍人也不错",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021016.jpg",
	}
}

type CardDef2021017 struct{}

func (CardDef2021017) ID() string      { return "2021017" }
func (CardDef2021017) Name() string    { return "旅行行囊" }
func (CardDef2021017) Kind() string    { return "道具" }
func (CardDef2021017) Element() string { return "无" }

func (CardDef2021017) Card() model.Card {
	return model.Card{
		Number:          "2021017",
		Type:            "道具",
		Name:            "旅行行囊",
		Category:        "无",
		Tag:             "装备",
		Description:     "光环:你的道具槽位+3",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021017.jpg",
	}
}

type CardDef2021018 struct{}

func (CardDef2021018) ID() string      { return "2021018" }
func (CardDef2021018) Name() string    { return "奥术符文" }
func (CardDef2021018) Kind() string    { return "道具" }
func (CardDef2021018) Element() string { return "无" }

func (CardDef2021018) Card() model.Card {
	return model.Card{
		Number:          "2021018",
		Type:            "道具",
		Name:            "奥术符文",
		Category:        "无",
		Tag:             "消耗品-符文",
		Description:     "反制:当敌方使用法术时,使你的1个法术在本回合+3\\威(敌方可以继续进行强化)",
		Quote:           "广泛镶嵌于屠魔者的剑上",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021018.jpg",
	}
}

type CardDef2021019 struct{}

func (CardDef2021019) ID() string      { return "2021019" }
func (CardDef2021019) Name() string    { return "诅咒卷轴" }
func (CardDef2021019) Kind() string    { return "道具" }
func (CardDef2021019) Element() string { return "无" }

func (CardDef2021019) Card() model.Card {
	return model.Card{
		Number:          "2021019",
		Type:            "道具",
		Name:            "诅咒卷轴",
		Category:        "无",
		Tag:             "消耗品-卷轴",
		Description:     "抽2张牌,但在本回合结束时将那些牌丢弃",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021019.jpg",
	}
}

type CardDef2021020 struct{}

func (CardDef2021020) ID() string      { return "2021020" }
func (CardDef2021020) Name() string    { return "假面" }
func (CardDef2021020) Kind() string    { return "道具" }
func (CardDef2021020) Element() string { return "无" }

func (CardDef2021020) Card() model.Card {
	return model.Card{
		Number:          "2021020",
		Type:            "道具",
		Name:            "假面",
		Category:        "无",
		Tag:             "装备-饰物",
		Description:     "光环:你的人物的负载变为等量的奥术元素\\无",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021020.jpg",
	}
}

type CardDef2021021 struct{}

func (CardDef2021021) ID() string      { return "2021021" }
func (CardDef2021021) Name() string    { return "聚能卷轴" }
func (CardDef2021021) Kind() string    { return "道具" }
func (CardDef2021021) Element() string { return "无" }

func (CardDef2021021) Card() model.Card {
	return model.Card{
		Number:          "2021021",
		Type:            "道具",
		Name:            "聚能卷轴",
		Category:        "无",
		Tag:             "消耗品-卷轴",
		Description:     "在你的下个回合开始时获得5\\无",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021021.jpg",
	}
}

type CardDef2021022 struct{}

func (CardDef2021022) ID() string      { return "2021022" }
func (CardDef2021022) Name() string    { return "反制符文" }
func (CardDef2021022) Kind() string    { return "道具" }
func (CardDef2021022) Element() string { return "无" }

func (CardDef2021022) Card() model.Card {
	return model.Card{
		Number:          "2021022",
		Type:            "道具",
		Name:            "反制符文",
		Category:        "无",
		Tag:             "消耗品-符文",
		Description:     "反制:当敌方使用卷轴或符文时,将其无效",
		Quote:           "聪明反被聪明误",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021022.jpg",
	}
}

type CardDef2021023 struct{}

func (CardDef2021023) ID() string      { return "2021023" }
func (CardDef2021023) Name() string    { return "奥术魔法筒" }
func (CardDef2021023) Kind() string    { return "道具" }
func (CardDef2021023) Element() string { return "无" }

func (CardDef2021023) Card() model.Card {
	return model.Card{
		Number:          "2021023",
		Type:            "道具",
		Name:            "奥术魔法筒",
		Category:        "无",
		Tag:             "装备",
		Description:     "入场:放置3个标记物.主动回合技:消耗此卡并取除1个标记物才能发动,获得2\\无",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\无\\2021023.jpg",
	}
}

type CardDef2021101 struct{}

func (CardDef2021101) ID() string      { return "2021101" }
func (CardDef2021101) Name() string    { return "失落的银叶花" }
func (CardDef2021101) Kind() string    { return "道具" }
func (CardDef2021101) Element() string { return "无" }

func (CardDef2021101) Card() model.Card {
	return model.Card{
		Number:          "2021101",
		Type:            "道具",
		Name:            "失落的银叶花",
		Category:        "无",
		Tag:             "消耗品",
		Description:     "抽2张牌,弃1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021101.jpg",
	}
}

type CardDef2021102 struct{}

func (CardDef2021102) ID() string      { return "2021102" }
func (CardDef2021102) Name() string    { return "破魔之刃" }
func (CardDef2021102) Kind() string    { return "道具" }
func (CardDef2021102) Element() string { return "无" }

func (CardDef2021102) Card() model.Card {
	return model.Card{
		Number:          "2021102",
		Type:            "道具",
		Name:            "破魔之刃",
		Category:        "无",
		Tag:             "装备-武器",
		Description:     "入场:使对方失去护盾3",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021102.jpg",
	}
}

type CardDef2021103 struct{}

func (CardDef2021103) ID() string      { return "2021103" }
func (CardDef2021103) Name() string    { return "幻雾药剂" }
func (CardDef2021103) Kind() string    { return "道具" }
func (CardDef2021103) Element() string { return "无" }

func (CardDef2021103) Card() model.Card {
	return model.Card{
		Number:          "2021103",
		Type:            "道具",
		Name:            "幻雾药剂",
		Category:        "无",
		Tag:             "消耗品-药剂",
		Description:     "使法力范围内的1个伙伴隐蔽2",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021103.jpg",
	}
}

type CardDef2021104 struct{}

func (CardDef2021104) ID() string      { return "2021104" }
func (CardDef2021104) Name() string    { return "五色珊瑚" }
func (CardDef2021104) Kind() string    { return "道具" }
func (CardDef2021104) Element() string { return "无" }

func (CardDef2021104) Card() model.Card {
	return model.Card{
		Number:          "2021104",
		Type:            "道具",
		Name:            "五色珊瑚",
		Category:        "无",
		Tag:             "装备",
		Description:     "入场:此卡获得2点不同的奥术元素以外的负载",
		Quote:           "之所以人们称之为五色珊瑚,是因为人们从来没有发现过黑色的品种",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021104.jpg",
	}
}

type CardDef2021105 struct{}

func (CardDef2021105) ID() string      { return "2021105" }
func (CardDef2021105) Name() string    { return "珍宝柜" }
func (CardDef2021105) Kind() string    { return "道具" }
func (CardDef2021105) Element() string { return "无" }

func (CardDef2021105) Card() model.Card {
	return model.Card{
		Number:          "2021105",
		Type:            "道具",
		Name:            "珍宝柜",
		Category:        "无",
		Tag:             "装备",
		Description:     "光环:你的道具槽位+1,且可以同时装备任意数量的同种装备",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021105.jpg",
	}
}

type CardDef2021106 struct{}

func (CardDef2021106) ID() string      { return "2021106" }
func (CardDef2021106) Name() string    { return "奥术容器" }
func (CardDef2021106) Kind() string    { return "道具" }
func (CardDef2021106) Element() string { return "无" }

func (CardDef2021106) Card() model.Card {
	return model.Card{
		Number:          "2021106",
		Type:            "道具",
		Name:            "奥术容器",
		Category:        "无",
		Tag:             "装备",
		Description:     "此卡提供的奥术元素不能当作其他元素使用",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021106.jpg",
	}
}

type CardDef2021107 struct{}

func (CardDef2021107) ID() string      { return "2021107" }
func (CardDef2021107) Name() string    { return "重塑" }
func (CardDef2021107) Kind() string    { return "道具" }
func (CardDef2021107) Element() string { return "无" }

func (CardDef2021107) Card() model.Card {
	return model.Card{
		Number:          "2021107",
		Type:            "道具",
		Name:            "重塑",
		Category:        "无",
		Tag:             "消耗品-卷轴",
		Description:     "弃掉你所有手牌,抽2张牌",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021107.jpg",
	}
}

type CardDef2021108 struct{}

func (CardDef2021108) ID() string      { return "2021108" }
func (CardDef2021108) Name() string    { return "通灵盘" }
func (CardDef2021108) Kind() string    { return "道具" }
func (CardDef2021108) Element() string { return "无" }

func (CardDef2021108) Card() model.Card {
	return model.Card{
		Number:          "2021108",
		Type:            "道具",
		Name:            "通灵盘",
		Category:        "无",
		Tag:             "装备",
		Description:     "光环:你的灵媒技能使用花费-1",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021108.jpg",
	}
}

type CardDef2021109 struct{}

func (CardDef2021109) ID() string      { return "2021109" }
func (CardDef2021109) Name() string    { return "氏族战锤" }
func (CardDef2021109) Kind() string    { return "道具" }
func (CardDef2021109) Element() string { return "无" }

func (CardDef2021109) Card() model.Card {
	return model.Card{
		Number:          "2021109",
		Type:            "道具",
		Name:            "氏族战锤",
		Category:        "无",
		Tag:             "装备-武器",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021109.jpg",
	}
}

type CardDef2021110 struct{}

func (CardDef2021110) ID() string      { return "2021110" }
func (CardDef2021110) Name() string    { return "灵守护符" }
func (CardDef2021110) Kind() string    { return "道具" }
func (CardDef2021110) Element() string { return "无" }

func (CardDef2021110) Card() model.Card {
	return model.Card{
		Number:          "2021110",
		Type:            "道具",
		Name:            "灵守护符",
		Category:        "无",
		Tag:             "装备",
		Description:     "光环:如果这是你唯一的装备,此卡负载+1\\无",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{"无": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021110.jpg",
	}
}

type CardDef2021111 struct{}

func (CardDef2021111) ID() string      { return "2021111" }
func (CardDef2021111) Name() string    { return "注能符文A型" }
func (CardDef2021111) Kind() string    { return "道具" }
func (CardDef2021111) Element() string { return "无" }

func (CardDef2021111) Card() model.Card {
	return model.Card{
		Number:          "2021111",
		Type:            "道具",
		Name:            "注能符文A型",
		Category:        "无",
		Tag:             "消耗品-符文",
		Description:     "使你的1个法术永久+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021111.jpg",
	}
}

type CardDef2021112 struct{}

func (CardDef2021112) ID() string      { return "2021112" }
func (CardDef2021112) Name() string    { return "注能符文B型" }
func (CardDef2021112) Kind() string    { return "道具" }
func (CardDef2021112) Element() string { return "无" }

func (CardDef2021112) Card() model.Card {
	return model.Card{
		Number:          "2021112",
		Type:            "道具",
		Name:            "注能符文B型",
		Category:        "无",
		Tag:             "消耗品-符文",
		Description:     "使你的1个法术永久+1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021112.jpg",
	}
}

type CardDef2021113 struct{}

func (CardDef2021113) ID() string      { return "2021113" }
func (CardDef2021113) Name() string    { return "奥术结界卷轴" }
func (CardDef2021113) Kind() string    { return "道具" }
func (CardDef2021113) Element() string { return "无" }

func (CardDef2021113) Card() model.Card {
	return model.Card{
		Number:          "2021113",
		Type:            "道具",
		Name:            "奥术结界卷轴",
		Category:        "无",
		Tag:             "消耗品-卷轴",
		Description:     "反制:当敌方法术命中时,获得护盾2",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021113.jpg",
	}
}

type CardDef2021114 struct{}

func (CardDef2021114) ID() string      { return "2021114" }
func (CardDef2021114) Name() string    { return "神护符文" }
func (CardDef2021114) Kind() string    { return "道具" }
func (CardDef2021114) Element() string { return "无" }

func (CardDef2021114) Card() model.Card {
	return model.Card{
		Number:          "2021114",
		Type:            "道具",
		Name:            "神护符文",
		Category:        "无",
		Tag:             "消耗品-符文",
		Description:     "反制:当1个友方单位受到致命伤害时,防止该伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021114.jpg",
	}
}

type CardDef2021115 struct{}

func (CardDef2021115) ID() string      { return "2021115" }
func (CardDef2021115) Name() string    { return "注能符文E型" }
func (CardDef2021115) Kind() string    { return "道具" }
func (CardDef2021115) Element() string { return "无" }

func (CardDef2021115) Card() model.Card {
	return model.Card{
		Number:          "2021115",
		Type:            "道具",
		Name:            "注能符文E型",
		Category:        "无",
		Tag:             "消耗品-符文",
		Description:     "反制:在你使用法术进行防御或强化防御后,将那些法术重置",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021115.jpg",
	}
}

type CardDef2021116 struct{}

func (CardDef2021116) ID() string      { return "2021116" }
func (CardDef2021116) Name() string    { return "奥能炸弹" }
func (CardDef2021116) Kind() string    { return "道具" }
func (CardDef2021116) Element() string { return "无" }

func (CardDef2021116) Card() model.Card {
	return model.Card{
		Number:          "2021116",
		Type:            "道具",
		Name:            "奥能炸弹",
		Category:        "无",
		Tag:             "消耗品",
		Description:     "对法力范围内的1个伙伴造成2点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\无\\2021116.jpg",
	}
}

type CardDef2111001 struct{}

func (CardDef2111001) ID() string      { return "2111001" }
func (CardDef2111001) Name() string    { return "火龙之心" }
func (CardDef2111001) Kind() string    { return "道具" }
func (CardDef2111001) Element() string { return "火" }

func (CardDef2111001) Card() model.Card {
	return model.Card{
		Number:          "2111001",
		Type:            "道具",
		Name:            "火龙之心",
		Category:        "火",
		Tag:             "传奇-装备-神器",
		Description:     "主动回合技:献祭包含最多3点\\火负载的卡牌,每1点\\火使下一次火焰法术获得+1\\攻或者+3\\威",
		Quote:           "辉煌死后,索拓尔使用它的心脏炼制了这件神器.只是没人看见火龙是怎么死的",
		ElementsCost:    map[string]int{"火": 6},
		ElementsGain:    map[string]int{"火": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2111001.jpg",
	}
}

type CardDef2111002 struct{}

func (CardDef2111002) ID() string      { return "2111002" }
func (CardDef2111002) Name() string    { return "努尔之眼" }
func (CardDef2111002) Kind() string    { return "道具" }
func (CardDef2111002) Element() string { return "火" }

func (CardDef2111002) Card() model.Card {
	return model.Card{
		Number:          "2111002",
		Type:            "道具",
		Name:            "努尔之眼",
		Category:        "火",
		Tag:             "传奇-装备-神器",
		Description:     "诱发:每当1个单位受到1点火焰伤害时,在此卡上放置1个火焰标记物.祈咒:移除此卡所有火焰标记物,根据数量执行以下效果.0个:摧毁此卡;1个:获得2\\火;2个:本回合你的火焰法术+2\\威;3个:本回合你的火焰法术+1\\攻;4个及以上:造成2点火焰伤害(不放置标记物)",
		Quote:           "目光透过烈焰,女巫记下了每一位观众的面容",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2111002.jpg",
	}
}

type CardDef2111101 struct{}

func (CardDef2111101) ID() string      { return "2111101" }
func (CardDef2111101) Name() string    { return "神火杖 赤空" }
func (CardDef2111101) Kind() string    { return "道具" }
func (CardDef2111101) Element() string { return "火" }

func (CardDef2111101) Card() model.Card {
	return model.Card{
		Number:          "2111101",
		Type:            "道具",
		Name:            "神火杖 赤空",
		Category:        "火",
		Tag:             "传奇-装备-武器",
		Description:     "诱发:当你的火焰法术命中后可以消耗此卡来发动,使该技能永久获得穿透和+1\\威",
		Quote:           "刹那间赤焰划过长空,热浪奔涌",
		ElementsCost:    map[string]int{"气": 1, "火": 3},
		ElementsGain:    map[string]int{"火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2111101.jpg",
	}
}

type CardDef2111102 struct{}

func (CardDef2111102) ID() string      { return "2111102" }
func (CardDef2111102) Name() string    { return "熔岩魔甲 业炎" }
func (CardDef2111102) Kind() string    { return "道具" }
func (CardDef2111102) Element() string { return "火" }

func (CardDef2111102) Card() model.Card {
	return model.Card{
		Number:          "2111102",
		Type:            "道具",
		Name:            "熔岩魔甲 业炎",
		Category:        "火",
		Tag:             "传奇-装备-防具",
		Description:     "诱发:当敌方法术命中时,可以献祭此卡发动,获得护盾2.遗言:本回合如果你的护盾被打破,从手牌、卡组装备1个熔火战铠,无需花费",
		Quote:           "注意别穿反了",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{"地": 1, "火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2111102.jpg",
	}
}

type CardDef2121001 struct{}

func (CardDef2121001) ID() string      { return "2121001" }
func (CardDef2121001) Name() string    { return "凤凰之羽" }
func (CardDef2121001) Kind() string    { return "道具" }
func (CardDef2121001) Element() string { return "火" }

func (CardDef2121001) Card() model.Card {
	return model.Card{
		Number:          "2121001",
		Type:            "道具",
		Name:            "凤凰之羽",
		Category:        "火",
		Tag:             "装备-神器",
		Description:     "入场:放置4个标记物.主动回合技:取除1个标记物才能发动,获得1\\火",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121001.jpg",
	}
}

type CardDef2121002 struct{}

func (CardDef2121002) ID() string      { return "2121002" }
func (CardDef2121002) Name() string    { return "火焰符文" }
func (CardDef2121002) Kind() string    { return "道具" }
func (CardDef2121002) Element() string { return "火" }

func (CardDef2121002) Card() model.Card {
	return model.Card{
		Number:          "2121002",
		Type:            "道具",
		Name:            "火焰符文",
		Category:        "火",
		Tag:             "消耗品-符文",
		Description:     "反制:当有单位被消耗时,使其获得点燃1",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121002.jpg",
	}
}

type CardDef2121003 struct{}

func (CardDef2121003) ID() string      { return "2121003" }
func (CardDef2121003) Name() string    { return "灼烧卷轴" }
func (CardDef2121003) Kind() string    { return "道具" }
func (CardDef2121003) Element() string { return "火" }

func (CardDef2121003) Card() model.Card {
	return model.Card{
		Number:          "2121003",
		Type:            "道具",
		Name:            "灼烧卷轴",
		Category:        "火",
		Tag:             "消耗品-法术卷轴-聚能",
		Description:     "点燃1",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121003.jpg",
	}
}

type CardDef2121004 struct{}

func (CardDef2121004) ID() string      { return "2121004" }
func (CardDef2121004) Name() string    { return "火焰箭" }
func (CardDef2121004) Kind() string    { return "道具" }
func (CardDef2121004) Element() string { return "火" }

func (CardDef2121004) Card() model.Card {
	return model.Card{
		Number:          "2121004",
		Type:            "道具",
		Name:            "火焰箭",
		Category:        "火",
		Tag:             "装备-武器",
		Description:     "主动:消耗并献祭此卡才能发动,对任意1个敌人造成1点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121004.jpg",
	}
}

type CardDef2121005 struct{}

func (CardDef2121005) ID() string      { return "2121005" }
func (CardDef2121005) Name() string    { return "神炎魔咒药剂" }
func (CardDef2121005) Kind() string    { return "道具" }
func (CardDef2121005) Element() string { return "火" }

func (CardDef2121005) Card() model.Card {
	return model.Card{
		Number:          "2121005",
		Type:            "道具",
		Name:            "神炎魔咒药剂",
		Category:        "火",
		Tag:             "消耗品-药剂",
		Description:     "直到下个回合结束你的法术+2\\威,你的人物获得点燃1",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121005.jpg",
	}
}

type CardDef2121006 struct{}

func (CardDef2121006) ID() string      { return "2121006" }
func (CardDef2121006) Name() string    { return "火焰面甲" }
func (CardDef2121006) Kind() string    { return "道具" }
func (CardDef2121006) Element() string { return "火" }

func (CardDef2121006) Card() model.Card {
	return model.Card{
		Number:          "2121006",
		Type:            "道具",
		Name:            "火焰面甲",
		Category:        "火",
		Tag:             "装备-饰物",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 1},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121006.jpg",
	}
}

type CardDef2121007 struct{}

func (CardDef2121007) ID() string      { return "2121007" }
func (CardDef2121007) Name() string    { return "舞火战裙" }
func (CardDef2121007) Kind() string    { return "道具" }
func (CardDef2121007) Element() string { return "火" }

func (CardDef2121007) Card() model.Card {
	return model.Card{
		Number:          "2121007",
		Type:            "道具",
		Name:            "舞火战裙",
		Category:        "火",
		Tag:             "装备-防具",
		Description:     "主动绝技:移除1个友方火焰单位所有负面状态",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "火": 1},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121007.jpg",
	}
}

type CardDef2121008 struct{}

func (CardDef2121008) ID() string      { return "2121008" }
func (CardDef2121008) Name() string    { return "烈焰风暴卷轴" }
func (CardDef2121008) Kind() string    { return "道具" }
func (CardDef2121008) Element() string { return "火" }

func (CardDef2121008) Card() model.Card {
	return model.Card{
		Number:          "2121008",
		Type:            "道具",
		Name:            "烈焰风暴卷轴",
		Category:        "火",
		Tag:             "消耗品-法术卷轴-创造",
		Description:     "范围:方阵.点燃1",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121008.jpg",
	}
}

type CardDef2121009 struct{}

func (CardDef2121009) ID() string      { return "2121009" }
func (CardDef2121009) Name() string    { return "烈焰障壁卷轴" }
func (CardDef2121009) Kind() string    { return "道具" }
func (CardDef2121009) Element() string { return "火" }

func (CardDef2121009) Card() model.Card {
	return model.Card{
		Number:          "2121009",
		Type:            "道具",
		Name:            "烈焰障壁卷轴",
		Category:        "火",
		Tag:             "消耗品-法术卷轴-创造",
		Description:     "防御",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121009.jpg",
	}
}

type CardDef2121010 struct{}

func (CardDef2121010) ID() string      { return "2121010" }
func (CardDef2121010) Name() string    { return "炽火链鞭" }
func (CardDef2121010) Kind() string    { return "道具" }
func (CardDef2121010) Element() string { return "火" }

func (CardDef2121010) Card() model.Card {
	return model.Card{
		Number:          "2121010",
		Type:            "道具",
		Name:            "炽火链鞭",
		Category:        "火",
		Tag:             "装备-武器",
		Description:     "",
		Quote:           "想不到你还好这口",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121010.jpg",
	}
}

type CardDef2121011 struct{}

func (CardDef2121011) ID() string      { return "2121011" }
func (CardDef2121011) Name() string    { return "火流星卷轴" }
func (CardDef2121011) Kind() string    { return "道具" }
func (CardDef2121011) Element() string { return "火" }

func (CardDef2121011) Card() model.Card {
	return model.Card{
		Number:          "2121011",
		Type:            "道具",
		Name:            "火流星卷轴",
		Category:        "火",
		Tag:             "消耗品-法术卷轴-创造",
		Description:     "穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121011.jpg",
	}
}

type CardDef2121012 struct{}

func (CardDef2121012) ID() string      { return "2121012" }
func (CardDef2121012) Name() string    { return "狱火符文" }
func (CardDef2121012) Kind() string    { return "道具" }
func (CardDef2121012) Element() string { return "火" }

func (CardDef2121012) Card() model.Card {
	return model.Card{
		Number:          "2121012",
		Type:            "道具",
		Name:            "狱火符文",
		Category:        "火",
		Tag:             "消耗品-符文",
		Description:     "反制:当敌方召唤1个伙伴时,使其获得晕眩2,石化2,点燃2",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121012.jpg",
	}
}

type CardDef2121013 struct{}

func (CardDef2121013) ID() string      { return "2121013" }
func (CardDef2121013) Name() string    { return "熔火战铠" }
func (CardDef2121013) Kind() string    { return "道具" }
func (CardDef2121013) Element() string { return "火" }

func (CardDef2121013) Card() model.Card {
	return model.Card{
		Number:          "2121013",
		Type:            "道具",
		Name:            "熔火战铠",
		Category:        "火",
		Tag:             "装备-防具",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 5},
		ElementsGain:    map[string]int{"火": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121013.jpg",
	}
}

type CardDef2121014 struct{}

func (CardDef2121014) ID() string      { return "2121014" }
func (CardDef2121014) Name() string    { return "火匣子" }
func (CardDef2121014) Kind() string    { return "道具" }
func (CardDef2121014) Element() string { return "火" }

func (CardDef2121014) Card() model.Card {
	return model.Card{
		Number:          "2121014",
		Type:            "道具",
		Name:            "火匣子",
		Category:        "火",
		Tag:             "装备",
		Description:     "入场:放置3个标记物.主动回合技:消耗此卡并取除1个标记物才能发动,获得2\\火",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\火\\2121014.jpg",
	}
}

type CardDef2121101 struct{}

func (CardDef2121101) ID() string      { return "2121101" }
func (CardDef2121101) Name() string    { return "熔岩堡的灰烬" }
func (CardDef2121101) Kind() string    { return "道具" }
func (CardDef2121101) Element() string { return "火" }

func (CardDef2121101) Card() model.Card {
	return model.Card{
		Number:          "2121101",
		Type:            "道具",
		Name:            "熔岩堡的灰烬",
		Category:        "火",
		Tag:             "消耗品",
		Description:     "将你场上或技能池1个火焰技能移出游戏,翻取1个入场花费数量更高的火焰卡牌并使其入场花费-1\\火",
		Quote:           "一点点法力就能让它复燃",
		ElementsCost:    map[string]int{"火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121101.jpg",
	}
}

type CardDef2121102 struct{}

func (CardDef2121102) ID() string      { return "2121102" }
func (CardDef2121102) Name() string    { return "火云扇" }
func (CardDef2121102) Kind() string    { return "道具" }
func (CardDef2121102) Element() string { return "火" }

func (CardDef2121102) Card() model.Card {
	return model.Card{
		Number:          "2121102",
		Type:            "道具",
		Name:            "火云扇",
		Category:        "火",
		Tag:             "装备-武器",
		Description:     "光环:你的火焰和大气法术可以选择正前方没有单位的非前排敌人",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "气": 1, "火": 2},
		ElementsGain:    map[string]int{"气": 1, "火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121102.jpg",
	}
}

type CardDef2121103 struct{}

func (CardDef2121103) ID() string      { return "2121103" }
func (CardDef2121103) Name() string    { return "浴火之翼" }
func (CardDef2121103) Kind() string    { return "道具" }
func (CardDef2121103) Element() string { return "火" }

func (CardDef2121103) Card() model.Card {
	return model.Card{
		Number:          "2121103",
		Type:            "道具",
		Name:            "浴火之翼",
		Category:        "火",
		Tag:             "装备-神器",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{"火": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121103.jpg",
	}
}

type CardDef2121104 struct{}

func (CardDef2121104) ID() string      { return "2121104" }
func (CardDef2121104) Name() string    { return "浴火重生卷轴" }
func (CardDef2121104) Kind() string    { return "道具" }
func (CardDef2121104) Element() string { return "火" }

func (CardDef2121104) Card() model.Card {
	return model.Card{
		Number:          "2121104",
		Type:            "道具",
		Name:            "浴火重生卷轴",
		Category:        "火",
		Tag:             "消耗品-卷轴",
		Description:     "反制:敌方回合结束时,从弃牌堆复活1个当回合死亡的友方火焰伙伴并将其重置,无需花费",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121104.jpg",
	}
}

type CardDef2121105 struct{}

func (CardDef2121105) ID() string      { return "2121105" }
func (CardDef2121105) Name() string    { return "魔焰阵卷轴" }
func (CardDef2121105) Kind() string    { return "道具" }
func (CardDef2121105) Element() string { return "火" }

func (CardDef2121105) Card() model.Card {
	return model.Card{
		Number:          "2121105",
		Type:            "道具",
		Name:            "魔焰阵卷轴",
		Category:        "火",
		Tag:             "消耗品-法术卷轴-幻变",
		Description:     "范围:溅射.使用时必须消耗1个入场花费大于4的火焰伙伴,并使此卡\\威上升其入场花费的数值",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           2,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121105.jpg",
	}
}

type CardDef2121106 struct{}

func (CardDef2121106) ID() string      { return "2121106" }
func (CardDef2121106) Name() string    { return "驯兽项圈" }
func (CardDef2121106) Kind() string    { return "道具" }
func (CardDef2121106) Element() string { return "火" }

func (CardDef2121106) Card() model.Card {
	return model.Card{
		Number:          "2121106",
		Type:            "道具",
		Name:            "驯兽项圈",
		Category:        "火",
		Tag:             "装备",
		Description:     "入场:选择你的1个巫师以外的火焰伙伴.主动回合技:消耗该火焰伙伴并获得其入场花费的元素",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121106.jpg",
	}
}

type CardDef2121107 struct{}

func (CardDef2121107) ID() string      { return "2121107" }
func (CardDef2121107) Name() string    { return "神火集结号" }
func (CardDef2121107) Kind() string    { return "道具" }
func (CardDef2121107) Element() string { return "火" }

func (CardDef2121107) Card() model.Card {
	return model.Card{
		Number:          "2121107",
		Type:            "道具",
		Name:            "神火集结号",
		Category:        "火",
		Tag:             "装备",
		Description:     "入场:翻取2张火焰伙伴牌",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{"火": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121107.jpg",
	}
}

type CardDef2121108 struct{}

func (CardDef2121108) ID() string      { return "2121108" }
func (CardDef2121108) Name() string    { return "燃烬卷轴" }
func (CardDef2121108) Kind() string    { return "道具" }
func (CardDef2121108) Element() string { return "火" }

func (CardDef2121108) Card() model.Card {
	return model.Card{
		Number:          "2121108",
		Type:            "道具",
		Name:            "燃烬卷轴",
		Category:        "火",
		Tag:             "消耗品-卷轴",
		Description:     "消耗1个友方火焰伙伴,获得其入场花费的元素",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121108.jpg",
	}
}

type CardDef2121109 struct{}

func (CardDef2121109) ID() string      { return "2121109" }
func (CardDef2121109) Name() string    { return "唤灵术卷轴 烈焰犬" }
func (CardDef2121109) Kind() string    { return "道具" }
func (CardDef2121109) Element() string { return "火" }

func (CardDef2121109) Card() model.Card {
	return model.Card{
		Number:          "2121109",
		Type:            "道具",
		Name:            "唤灵术卷轴 烈焰犬",
		Category:        "火",
		Tag:             "消耗品-法术卷轴-创造",
		Description:     "若此卡防御成功,对法力范围内1个敌人造成2点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121109.jpg",
	}
}

type CardDef2121110 struct{}

func (CardDef2121110) ID() string      { return "2121110" }
func (CardDef2121110) Name() string    { return "供奉之炬" }
func (CardDef2121110) Kind() string    { return "道具" }
func (CardDef2121110) Element() string { return "火" }

func (CardDef2121110) Card() model.Card {
	return model.Card{
		Number:          "2121110",
		Type:            "道具",
		Name:            "供奉之炬",
		Category:        "火",
		Tag:             "消耗品",
		Description:     "将你场上的1个火焰法术移出游戏,使你的另一个火焰法术永久增加该法术的\\威和\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121110.jpg",
	}
}

type CardDef2121111 struct{}

func (CardDef2121111) ID() string      { return "2121111" }
func (CardDef2121111) Name() string    { return "洛普修斯之怒" }
func (CardDef2121111) Kind() string    { return "道具" }
func (CardDef2121111) Element() string { return "火" }

func (CardDef2121111) Card() model.Card {
	return model.Card{
		Number:          "2121111",
		Type:            "道具",
		Name:            "洛普修斯之怒",
		Category:        "火",
		Tag:             "消耗品-法术卷轴-神秘",
		Description:     "此卡攻击时,直到完全结算完毕对方不能使用反制或诱发效果",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          4,
		Life:            -1,
		Duration:        -1,
		Power:           9,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121111.jpg",
	}
}

type CardDef2121112 struct{}

func (CardDef2121112) ID() string      { return "2121112" }
func (CardDef2121112) Name() string    { return "炎流卷轴" }
func (CardDef2121112) Kind() string    { return "道具" }
func (CardDef2121112) Element() string { return "火" }

func (CardDef2121112) Card() model.Card {
	return model.Card{
		Number:          "2121112",
		Type:            "道具",
		Name:            "炎流卷轴",
		Category:        "火",
		Tag:             "消耗品-法术卷轴-驱动",
		Description:     "范围:纵列.点燃1",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\火\\2121112.jpg",
	}
}

type CardDef2201101 struct{}

func (CardDef2201101) ID() string      { return "2201101" }
func (CardDef2201101) Name() string    { return "幻创之梦-绽放" }
func (CardDef2201101) Kind() string    { return "道具" }
func (CardDef2201101) Element() string { return "水" }

func (CardDef2201101) Card() model.Card {
	return model.Card{
		Number:          "2201101",
		Type:            "道具",
		Name:            "幻创之梦-绽放",
		Category:        "水",
		Tag:             "衍生-消耗品",
		Description:     "抽3张牌",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2201101.jpg",
	}
}

type CardDef2201102 struct{}

func (CardDef2201102) ID() string      { return "2201102" }
func (CardDef2201102) Name() string    { return "幻创之梦-幻能" }
func (CardDef2201102) Kind() string    { return "道具" }
func (CardDef2201102) Element() string { return "水" }

func (CardDef2201102) Card() model.Card {
	return model.Card{
		Number:          "2201102",
		Type:            "道具",
		Name:            "幻创之梦-幻能",
		Category:        "水",
		Tag:             "衍生-消耗品",
		Description:     "获得3\\无",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2201102.jpg",
	}
}

type CardDef2201103 struct{}

func (CardDef2201103) ID() string      { return "2201103" }
func (CardDef2201103) Name() string    { return "幻创之梦-波纹" }
func (CardDef2201103) Kind() string    { return "道具" }
func (CardDef2201103) Element() string { return "水" }

func (CardDef2201103) Card() model.Card {
	return model.Card{
		Number:          "2201103",
		Type:            "道具",
		Name:            "幻创之梦-波纹",
		Category:        "水",
		Tag:             "衍生-消耗品",
		Description:     "对前排任意敌人造成共计3点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2201103.jpg",
	}
}

type CardDef2211001 struct{}

func (CardDef2211001) ID() string      { return "2211001" }
func (CardDef2211001) Name() string    { return "人鱼之泪" }
func (CardDef2211001) Kind() string    { return "道具" }
func (CardDef2211001) Element() string { return "水" }

func (CardDef2211001) Card() model.Card {
	return model.Card{
		Number:          "2211001",
		Type:            "道具",
		Name:            "人鱼之泪",
		Category:        "水",
		Tag:             "传奇-装备-神器",
		Description:     "主动:将此卡从游戏中移除才能发动,复活你的1个死亡伙伴但只有1\\血",
		Quote:           "巴特尔从沙滩上惊醒,身旁的泡沫顷刻间被浪花融化",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2211001.jpg",
	}
}

type CardDef2211002 struct{}

func (CardDef2211002) ID() string      { return "2211002" }
func (CardDef2211002) Name() string    { return "嗜魔弓 凛冬" }
func (CardDef2211002) Kind() string    { return "道具" }
func (CardDef2211002) Element() string { return "水" }

func (CardDef2211002) Card() model.Card {
	return model.Card{
		Number:          "2211002",
		Type:            "道具",
		Name:            "嗜魔弓 凛冬",
		Category:        "水",
		Tag:             "传奇-装备-武器",
		Description:     "诱发:每当有玩家使用法术时,在此卡上放置1个水纹标记物.绑定技能:凛冬将至",
		Quote:           "一箭霜降,两箭严寒",
		ElementsCost:    map[string]int{"水": 7},
		ElementsGain:    map[string]int{"水": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"3201002"},
		OutputPath:      "output\\基础包\\道具\\水\\2211002.jpg",
	}
}

type CardDef2211101 struct{}

func (CardDef2211101) ID() string      { return "2211101" }
func (CardDef2211101) Name() string    { return "珊瑚秘宝 深邃之剑" }
func (CardDef2211101) Kind() string    { return "道具" }
func (CardDef2211101) Element() string { return "水" }

func (CardDef2211101) Card() model.Card {
	return model.Card{
		Number:          "2211101",
		Type:            "道具",
		Name:            "珊瑚秘宝 深邃之剑",
		Category:        "水",
		Tag:             "传奇-装备-武器",
		Description:     "此卡无法被任何卡牌检索或翻取.诱发:当你从卡组抽到这张牌时,如果当前敌方法术总威力比友方高,可以展示此卡并对法力范围内所有敌人造成2点伤害",
		Quote:           "\"我的孩子,我希望你永远都不会有使用它的一天,但假如那一天真的到来,它会给你足够的勇气去斩断一切\"",
		ElementsCost:    map[string]int{"水": 6},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2211101.jpg",
	}
}

type CardDef2211102 struct{}

func (CardDef2211102) ID() string      { return "2211102" }
func (CardDef2211102) Name() string    { return "玛涅斯之予夺" }
func (CardDef2211102) Kind() string    { return "道具" }
func (CardDef2211102) Element() string { return "水" }

func (CardDef2211102) Card() model.Card {
	return model.Card{
		Number:          "2211102",
		Type:            "道具",
		Name:            "玛涅斯之予夺",
		Category:        "水",
		Tag:             "传奇-装备-武器",
		Description:     "入场:从以下两张效果中选择1个直到本局游戏结束:此后你学习的所有水纹法术+2\\威;你的1个水纹法术+3\\威+1\\攻,你不能再学习任何水纹法术",
		Quote:           "\"看来这次,我必须亲力亲为\"",
		ElementsCost:    map[string]int{"水": 4},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2211102.jpg",
	}
}

type CardDef2221001 struct{}

func (CardDef2221001) ID() string      { return "2221001" }
func (CardDef2221001) Name() string    { return "冰霜之心" }
func (CardDef2221001) Kind() string    { return "道具" }
func (CardDef2221001) Element() string { return "水" }

func (CardDef2221001) Card() model.Card {
	return model.Card{
		Number:          "2221001",
		Type:            "道具",
		Name:            "冰霜之心",
		Category:        "水",
		Tag:             "装备-防具",
		Description:     "诱发:敌方法术命中时献祭此卡可以发动,该法术伤害变为0",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221001.jpg",
	}
}

type CardDef2221002 struct{}

func (CardDef2221002) ID() string      { return "2221002" }
func (CardDef2221002) Name() string    { return "冰霜符文" }
func (CardDef2221002) Kind() string    { return "道具" }
func (CardDef2221002) Element() string { return "水" }

func (CardDef2221002) Card() model.Card {
	return model.Card{
		Number:          "2221002",
		Type:            "道具",
		Name:            "冰霜符文",
		Category:        "水",
		Tag:             "消耗品-符文",
		Description:     "反制:当有敌方伙伴被消耗时,使其冻结1",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221002.jpg",
	}
}

type CardDef2221003 struct{}

func (CardDef2221003) ID() string      { return "2221003" }
func (CardDef2221003) Name() string    { return "冰封卷轴" }
func (CardDef2221003) Kind() string    { return "道具" }
func (CardDef2221003) Element() string { return "水" }

func (CardDef2221003) Card() model.Card {
	return model.Card{
		Number:          "2221003",
		Type:            "道具",
		Name:            "冰封卷轴",
		Category:        "水",
		Tag:             "消耗品-卷轴",
		Description:     "使所有前排敌人冻结1",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221003.jpg",
	}
}

type CardDef2221004 struct{}

func (CardDef2221004) ID() string      { return "2221004" }
func (CardDef2221004) Name() string    { return "玛涅斯之杖" }
func (CardDef2221004) Kind() string    { return "道具" }
func (CardDef2221004) Element() string { return "水" }

func (CardDef2221004) Card() model.Card {
	return model.Card{
		Number:          "2221004",
		Type:            "道具",
		Name:            "玛涅斯之杖",
		Category:        "水",
		Tag:             "装备-武器",
		Description:     "光环:你的水纹法术+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221004.jpg",
	}
}

type CardDef2221005 struct{}

func (CardDef2221005) ID() string      { return "2221005" }
func (CardDef2221005) Name() string    { return "精力药剂" }
func (CardDef2221005) Kind() string    { return "道具" }
func (CardDef2221005) Element() string { return "水" }

func (CardDef2221005) Card() model.Card {
	return model.Card{
		Number:          "2221005",
		Type:            "道具",
		Name:            "精力药剂",
		Category:        "水",
		Tag:             "消耗品-药剂",
		Description:     "反制:对方回合结束时,将你的全部法术重置",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221005.jpg",
	}
}

type CardDef2221006 struct{}

func (CardDef2221006) ID() string      { return "2221006" }
func (CardDef2221006) Name() string    { return "海之眷顾" }
func (CardDef2221006) Kind() string    { return "道具" }
func (CardDef2221006) Element() string { return "水" }

func (CardDef2221006) Card() model.Card {
	return model.Card{
		Number:          "2221006",
		Type:            "道具",
		Name:            "海之眷顾",
		Category:        "水",
		Tag:             "装备-饰物",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 5},
		ElementsGain:    map[string]int{"水": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221006.jpg",
	}
}

type CardDef2221007 struct{}

func (CardDef2221007) ID() string      { return "2221007" }
func (CardDef2221007) Name() string    { return "凝霜手镯" }
func (CardDef2221007) Kind() string    { return "道具" }
func (CardDef2221007) Element() string { return "水" }

func (CardDef2221007) Card() model.Card {
	return model.Card{
		Number:          "2221007",
		Type:            "道具",
		Name:            "凝霜手镯",
		Category:        "水",
		Tag:             "装备-饰物",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221007.jpg",
	}
}

type CardDef2221008 struct{}

func (CardDef2221008) ID() string      { return "2221008" }
func (CardDef2221008) Name() string    { return "水形之束卷轴" }
func (CardDef2221008) Kind() string    { return "道具" }
func (CardDef2221008) Element() string { return "水" }

func (CardDef2221008) Card() model.Card {
	return model.Card{
		Number:          "2221008",
		Type:            "道具",
		Name:            "水形之束卷轴",
		Category:        "水",
		Tag:             "消耗品-法术卷轴-驱动",
		Description:     "命中:若目标为伙伴牌,消耗该伙伴",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221008.jpg",
	}
}

type CardDef2221009 struct{}

func (CardDef2221009) ID() string      { return "2221009" }
func (CardDef2221009) Name() string    { return "寒冰爆裂卷轴" }
func (CardDef2221009) Kind() string    { return "道具" }
func (CardDef2221009) Element() string { return "水" }

func (CardDef2221009) Card() model.Card {
	return model.Card{
		Number:          "2221009",
		Type:            "道具",
		Name:            "寒冰爆裂卷轴",
		Category:        "水",
		Tag:             "消耗品-法术卷轴-聚能",
		Description:     "范围:溅射.冻结1",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221009.jpg",
	}
}

type CardDef2221010 struct{}

func (CardDef2221010) ID() string      { return "2221010" }
func (CardDef2221010) Name() string    { return "潮涌符文" }
func (CardDef2221010) Kind() string    { return "道具" }
func (CardDef2221010) Element() string { return "水" }

func (CardDef2221010) Card() model.Card {
	return model.Card{
		Number:          "2221010",
		Type:            "道具",
		Name:            "潮涌符文",
		Category:        "水",
		Tag:             "消耗品-符文",
		Description:     "反制:当敌方在一个回合内抽第三张牌时,使你的1个水纹伙伴获得负载+2\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221010.jpg",
	}
}

type CardDef2221011 struct{}

func (CardDef2221011) ID() string      { return "2221011" }
func (CardDef2221011) Name() string    { return "恩惠之雨" }
func (CardDef2221011) Kind() string    { return "道具" }
func (CardDef2221011) Element() string { return "水" }

func (CardDef2221011) Card() model.Card {
	return model.Card{
		Number:          "2221011",
		Type:            "道具",
		Name:            "恩惠之雨",
		Category:        "水",
		Tag:             "消耗品-卷轴",
		Description:     "反制:当1个友方单位受伤后,使所有友方单位回复2\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221011.jpg",
	}
}

type CardDef2221012 struct{}

func (CardDef2221012) ID() string      { return "2221012" }
func (CardDef2221012) Name() string    { return "水行之靴" }
func (CardDef2221012) Kind() string    { return "道具" }
func (CardDef2221012) Element() string { return "水" }

func (CardDef2221012) Card() model.Card {
	return model.Card{
		Number:          "2221012",
		Type:            "道具",
		Name:            "水行之靴",
		Category:        "水",
		Tag:             "装备-防具",
		Description:     "光环:在你的人物与至少3个水纹伙伴相邻时,此卡负载+1\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221012.jpg",
	}
}

type CardDef2221013 struct{}

func (CardDef2221013) ID() string      { return "2221013" }
func (CardDef2221013) Name() string    { return "深寒诅咒卷轴" }
func (CardDef2221013) Kind() string    { return "道具" }
func (CardDef2221013) Element() string { return "水" }

func (CardDef2221013) Card() model.Card {
	return model.Card{
		Number:          "2221013",
		Type:            "道具",
		Name:            "深寒诅咒卷轴",
		Category:        "水",
		Tag:             "消耗品-卷轴",
		Description:     "使法力范围内的1个敌方伙伴永久冻结",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221013.jpg",
	}
}

type CardDef2221014 struct{}

func (CardDef2221014) ID() string      { return "2221014" }
func (CardDef2221014) Name() string    { return "水之咏叹" }
func (CardDef2221014) Kind() string    { return "道具" }
func (CardDef2221014) Element() string { return "水" }

func (CardDef2221014) Card() model.Card {
	return model.Card{
		Number:          "2221014",
		Type:            "道具",
		Name:            "水之咏叹",
		Category:        "水",
		Tag:             "装备",
		Description:     "入场:放置4个标记物.主动回合技:消耗此卡并取除1个标记物才能发动,获得3\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\水\\2221014.jpg",
	}
}

type CardDef2221101 struct{}

func (CardDef2221101) ID() string      { return "2221101" }
func (CardDef2221101) Name() string    { return "镜花海的天泉" }
func (CardDef2221101) Kind() string    { return "道具" }
func (CardDef2221101) Element() string { return "水" }

func (CardDef2221101) Card() model.Card {
	return model.Card{
		Number:          "2221101",
		Type:            "道具",
		Name:            "镜花海的天泉",
		Category:        "水",
		Tag:             "消耗品-药剂",
		Description:     "使你的1个法术使用花费永久减少1\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221101.jpg",
	}
}

type CardDef2221102 struct{}

func (CardDef2221102) ID() string      { return "2221102" }
func (CardDef2221102) Name() string    { return "海洋之盾卷轴" }
func (CardDef2221102) Kind() string    { return "道具" }
func (CardDef2221102) Element() string { return "水" }

func (CardDef2221102) Card() model.Card {
	return model.Card{
		Number:          "2221102",
		Type:            "道具",
		Name:            "海洋之盾卷轴",
		Category:        "水",
		Tag:             "消耗品-卷轴",
		Description:     "获得护盾2",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221102.jpg",
	}
}

type CardDef2221103 struct{}

func (CardDef2221103) ID() string      { return "2221103" }
func (CardDef2221103) Name() string    { return "冰锁符文" }
func (CardDef2221103) Kind() string    { return "道具" }
func (CardDef2221103) Element() string { return "水" }

func (CardDef2221103) Card() model.Card {
	return model.Card{
		Number:          "2221103",
		Type:            "道具",
		Name:            "冰锁符文",
		Category:        "水",
		Tag:             "消耗品-符文",
		Description:     "反制:当对方学习1个技能时,使其直到下个回合结束不能使用",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221103.jpg",
	}
}

type CardDef2221104 struct{}

func (CardDef2221104) ID() string      { return "2221104" }
func (CardDef2221104) Name() string    { return "水镜卷轴" }
func (CardDef2221104) Kind() string    { return "道具" }
func (CardDef2221104) Element() string { return "水" }

func (CardDef2221104) Card() model.Card {
	return model.Card{
		Number:          "2221104",
		Type:            "道具",
		Name:            "水镜卷轴",
		Category:        "水",
		Tag:             "消耗品-法术卷轴-幻变",
		Description:     "复制你使用的上一个花费小于3的水纹法术(存在强化时,只复制主法术)",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221104.jpg",
	}
}

type CardDef2221105 struct{}

func (CardDef2221105) ID() string      { return "2221105" }
func (CardDef2221105) Name() string    { return "掠夺者黑帆" }
func (CardDef2221105) Kind() string    { return "道具" }
func (CardDef2221105) Element() string { return "水" }

func (CardDef2221105) Card() model.Card {
	return model.Card{
		Number:          "2221105",
		Type:            "道具",
		Name:            "掠夺者黑帆",
		Category:        "水",
		Tag:             "消耗品",
		Description:     "检索1个掠夺者伙伴,如果你的场上已有掠夺者伙伴则使检索的卡牌入场花费减少1\\水或1\\暗",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221105.jpg",
	}
}

type CardDef2221106 struct{}

func (CardDef2221106) ID() string      { return "2221106" }
func (CardDef2221106) Name() string    { return "凛冰法袍" }
func (CardDef2221106) Kind() string    { return "道具" }
func (CardDef2221106) Element() string { return "水" }

func (CardDef2221106) Card() model.Card {
	return model.Card{
		Number:          "2221106",
		Type:            "道具",
		Name:            "凛冰法袍",
		Category:        "水",
		Tag:             "装备-防具",
		Description:     "诱发绝技:当友方水纹单位受到敌方法术伤害后,使法力范围内所有敌人冻结1",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 4},
		ElementsGain:    map[string]int{"水": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221106.jpg",
	}
}

type CardDef2221107 struct{}

func (CardDef2221107) ID() string      { return "2221107" }
func (CardDef2221107) Name() string    { return "水纹之镜" }
func (CardDef2221107) Kind() string    { return "道具" }
func (CardDef2221107) Element() string { return "水" }

func (CardDef2221107) Card() model.Card {
	return model.Card{
		Number:          "2221107",
		Type:            "道具",
		Name:            "水纹之镜",
		Category:        "水",
		Tag:             "装备-神器",
		Description:     "入场:使你的1个水纹法术使用花费-1\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221107.jpg",
	}
}

type CardDef2221108 struct{}

func (CardDef2221108) ID() string      { return "2221108" }
func (CardDef2221108) Name() string    { return "西境航海图" }
func (CardDef2221108) Kind() string    { return "道具" }
func (CardDef2221108) Element() string { return "水" }

func (CardDef2221108) Card() model.Card {
	return model.Card{
		Number:          "2221108",
		Type:            "道具",
		Name:            "西境航海图",
		Category:        "水",
		Tag:             "装备",
		Description:     "入场:选择你的1个水纹法术.光环:使该法术获得穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{"水": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221108.jpg",
	}
}

type CardDef2221109 struct{}

func (CardDef2221109) ID() string      { return "2221109" }
func (CardDef2221109) Name() string    { return "速射冰弹" }
func (CardDef2221109) Kind() string    { return "道具" }
func (CardDef2221109) Element() string { return "水" }

func (CardDef2221109) Card() model.Card {
	return model.Card{
		Number:          "2221109",
		Type:            "道具",
		Name:            "速射冰弹",
		Category:        "水",
		Tag:             "消耗品",
		Description:     "本回合你的下一张消耗品道具或法术的打出或使用花费-3\\水",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221109.jpg",
	}
}

type CardDef2221110 struct{}

func (CardDef2221110) ID() string      { return "2221110" }
func (CardDef2221110) Name() string    { return "残霜飞雪卷轴" }
func (CardDef2221110) Kind() string    { return "道具" }
func (CardDef2221110) Element() string { return "水" }

func (CardDef2221110) Card() model.Card {
	return model.Card{
		Number:          "2221110",
		Type:            "道具",
		Name:            "残霜飞雪卷轴",
		Category:        "水",
		Tag:             "消耗品-法术卷轴-幻变",
		Description:     "范围:全场.本回合你每使用过1个水纹法术,此卡获得+3\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           0,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221110.jpg",
	}
}

type CardDef2221111 struct{}

func (CardDef2221111) ID() string      { return "2221111" }
func (CardDef2221111) Name() string    { return "冰魄印 淬" }
func (CardDef2221111) Kind() string    { return "道具" }
func (CardDef2221111) Element() string { return "水" }

func (CardDef2221111) Card() model.Card {
	return model.Card{
		Number:          "2221111",
		Type:            "道具",
		Name:            "冰魄印 淬",
		Category:        "水",
		Tag:             "消耗品-符文",
		Description:     "反制:敌方使用的法术总威力大于10时才能发动,将其总威力减半(向上取整),敌方可以继续强化",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221111.jpg",
	}
}

type CardDef2221112 struct{}

func (CardDef2221112) ID() string      { return "2221112" }
func (CardDef2221112) Name() string    { return "冰魄印 消" }
func (CardDef2221112) Kind() string    { return "道具" }
func (CardDef2221112) Element() string { return "水" }

func (CardDef2221112) Card() model.Card {
	return model.Card{
		Number:          "2221112",
		Type:            "道具",
		Name:            "冰魄印 消",
		Category:        "水",
		Tag:             "消耗品-符文",
		Description:     "反制:敌方使用1个法术强化时,若威力小于5,将其无效,敌方可以继续强化",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\水\\2221112.jpg",
	}
}

type CardDef2311001 struct{}

func (CardDef2311001) ID() string      { return "2311001" }
func (CardDef2311001) Name() string    { return "雷之源" }
func (CardDef2311001) Kind() string    { return "道具" }
func (CardDef2311001) Element() string { return "气" }

func (CardDef2311001) Card() model.Card {
	return model.Card{
		Number:          "2311001",
		Type:            "道具",
		Name:            "雷之源",
		Category:        "气",
		Tag:             "传奇-装备-神器",
		Description:     "光环:你的卡牌入场花费,学习花费和使用花费减少1\\气",
		Quote:           "四境雷动,不过弹指之间",
		ElementsCost:    map[string]int{"气": 8},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2311001.jpg",
	}
}

type CardDef2311002 struct{}

func (CardDef2311002) ID() string      { return "2311002" }
func (CardDef2311002) Name() string    { return "唤雷震鼓" }
func (CardDef2311002) Kind() string    { return "道具" }
func (CardDef2311002) Element() string { return "气" }

func (CardDef2311002) Card() model.Card {
	return model.Card{
		Number:          "2311002",
		Type:            "道具",
		Name:            "唤雷震鼓",
		Category:        "气",
		Tag:             "传奇-装备",
		Description:     "诱发:每当你抽1张牌时,可以将其展示并在此卡上放置1个标记.主动回合技:移除3个标记才能发动,本回合你的大气法术获得+1\\攻或者晕眩1",
		Quote:           "沉重的闷响传向远方,然后从远方传来",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2311002.jpg",
	}
}

type CardDef2311101 struct{}

func (CardDef2311101) ID() string      { return "2311101" }
func (CardDef2311101) Name() string    { return "云霄城的天顶石" }
func (CardDef2311101) Kind() string    { return "道具" }
func (CardDef2311101) Element() string { return "气" }

func (CardDef2311101) Card() model.Card {
	return model.Card{
		Number:          "2311101",
		Type:            "道具",
		Name:            "云霄城的天顶石",
		Category:        "气",
		Tag:             "传奇-装备-神器",
		Description:     "诱发:每当有玩家抽1张牌,在此卡上放置1个标记物.诱发:当此卡标记物达到5时必须移除所有标记物,此时抽卡的玩家前排单位受到1点伤害和晕眩1",
		Quote:           "这份杰作从第一天开始就驱动着云霄城,唯一的制造者早已埋没于历史",
		ElementsCost:    map[string]int{"光": 1, "气": 4},
		ElementsGain:    map[string]int{"无": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2311101.jpg",
	}
}

type CardDef2311102 struct{}

func (CardDef2311102) ID() string      { return "2311102" }
func (CardDef2311102) Name() string    { return "兰普斯之剑" }
func (CardDef2311102) Kind() string    { return "道具" }
func (CardDef2311102) Element() string { return "气" }

func (CardDef2311102) Card() model.Card {
	return model.Card{
		Number:          "2311102",
		Type:            "道具",
		Name:            "兰普斯之剑",
		Category:        "气",
		Tag:             "传奇-装备-武器",
		Description:     "主动:献祭此卡才能发动,丢弃任意数量大气手牌,下次对方回合结束时将那个数值的伤害分配给法力范围内的敌方伙伴",
		Quote:           "久远的记忆,浮现在被称作麻雀之人的眼前",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2311102.jpg",
	}
}

type CardDef2321001 struct{}

func (CardDef2321001) ID() string      { return "2321001" }
func (CardDef2321001) Name() string    { return "风息罗盘" }
func (CardDef2321001) Kind() string    { return "道具" }
func (CardDef2321001) Element() string { return "气" }

func (CardDef2321001) Card() model.Card {
	return model.Card{
		Number:          "2321001",
		Type:            "道具",
		Name:            "风息罗盘",
		Category:        "气",
		Tag:             "装备",
		Description:     "诱发回合技3:当你抽到1张大气卡牌时,你可以将其展示然后此卡临时获得负载1点\\气",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321001.jpg",
	}
}

type CardDef2321002 struct{}

func (CardDef2321002) ID() string      { return "2321002" }
func (CardDef2321002) Name() string    { return "闪电符文" }
func (CardDef2321002) Kind() string    { return "道具" }
func (CardDef2321002) Element() string { return "气" }

func (CardDef2321002) Card() model.Card {
	return model.Card{
		Number:          "2321002",
		Type:            "道具",
		Name:            "闪电符文",
		Category:        "气",
		Tag:             "消耗品-符文",
		Description:     "反制:当有敌人被消耗时,使其与1个相邻单位晕眩1",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321002.jpg",
	}
}

type CardDef2321003 struct{}

func (CardDef2321003) ID() string      { return "2321003" }
func (CardDef2321003) Name() string    { return "雷暴卷轴" }
func (CardDef2321003) Kind() string    { return "道具" }
func (CardDef2321003) Element() string { return "气" }

func (CardDef2321003) Card() model.Card {
	return model.Card{
		Number:          "2321003",
		Type:            "道具",
		Name:            "雷暴卷轴",
		Category:        "气",
		Tag:             "消耗品-法术卷轴-聚能",
		Description:     "范围:方阵.命中:使所有命中伙伴晕眩1",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321003.jpg",
	}
}

type CardDef2321004 struct{}

func (CardDef2321004) ID() string      { return "2321004" }
func (CardDef2321004) Name() string    { return "雷霆魔杖" }
func (CardDef2321004) Kind() string    { return "道具" }
func (CardDef2321004) Element() string { return "气" }

func (CardDef2321004) Card() model.Card {
	return model.Card{
		Number:          "2321004",
		Type:            "道具",
		Name:            "雷霆魔杖",
		Category:        "气",
		Tag:             "装备-武器",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321004.jpg",
	}
}

type CardDef2321005 struct{}

func (CardDef2321005) ID() string      { return "2321005" }
func (CardDef2321005) Name() string    { return "唤风卷轴" }
func (CardDef2321005) Kind() string    { return "道具" }
func (CardDef2321005) Element() string { return "气" }

func (CardDef2321005) Card() model.Card {
	return model.Card{
		Number:          "2321005",
		Type:            "道具",
		Name:            "唤风卷轴",
		Category:        "气",
		Tag:             "消耗品-卷轴",
		Description:     "抽2张牌,但在下个你的回合开始不抽牌",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321005.jpg",
	}
}

type CardDef2321006 struct{}

func (CardDef2321006) ID() string      { return "2321006" }
func (CardDef2321006) Name() string    { return "瓶中闪电" }
func (CardDef2321006) Kind() string    { return "道具" }
func (CardDef2321006) Element() string { return "气" }

func (CardDef2321006) Card() model.Card {
	return model.Card{
		Number:          "2321006",
		Type:            "道具",
		Name:            "瓶中闪电",
		Category:        "气",
		Tag:             "消耗品-药剂",
		Description:     "获得2\\气,使1个友方大气单位获得晕眩2",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321006.jpg",
	}
}

type CardDef2321007 struct{}

func (CardDef2321007) ID() string      { return "2321007" }
func (CardDef2321007) Name() string    { return "风语之戒" }
func (CardDef2321007) Kind() string    { return "道具" }
func (CardDef2321007) Element() string { return "气" }

func (CardDef2321007) Card() model.Card {
	return model.Card{
		Number:          "2321007",
		Type:            "道具",
		Name:            "风语之戒",
		Category:        "气",
		Tag:             "装备-饰物",
		Description:     "入场:抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321007.jpg",
	}
}

type CardDef2321008 struct{}

func (CardDef2321008) ID() string      { return "2321008" }
func (CardDef2321008) Name() string    { return "旋风卷轴" }
func (CardDef2321008) Kind() string    { return "道具" }
func (CardDef2321008) Element() string { return "气" }

func (CardDef2321008) Card() model.Card {
	return model.Card{
		Number:          "2321008",
		Type:            "道具",
		Name:            "旋风卷轴",
		Category:        "气",
		Tag:             "消耗品-卷轴",
		Description:     "摧毁敌方任意1个入场花费小于5的装备道具",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321008.jpg",
	}
}

type CardDef2321009 struct{}

func (CardDef2321009) ID() string      { return "2321009" }
func (CardDef2321009) Name() string    { return "连锁闪电卷轴" }
func (CardDef2321009) Kind() string    { return "道具" }
func (CardDef2321009) Element() string { return "气" }

func (CardDef2321009) Card() model.Card {
	return model.Card{
		Number:          "2321009",
		Type:            "道具",
		Name:            "连锁闪电卷轴",
		Category:        "气",
		Tag:             "消耗品-法术卷轴-创造",
		Description:     "命中:抽1张牌或者检索1张连锁闪电卷轴",
		Quote:           "当局已经控制了利普兹学院体育馆400人触电事件的肇事学生霍尔顿·弗雷,警方认为证据确凿:他是唯一一个上体育课换上橡胶手套和雨靴的学生.",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321009.jpg",
	}
}

type CardDef2321010 struct{}

func (CardDef2321010) ID() string      { return "2321010" }
func (CardDef2321010) Name() string    { return "幻术卷轴" }
func (CardDef2321010) Kind() string    { return "道具" }
func (CardDef2321010) Element() string { return "气" }

func (CardDef2321010) Card() model.Card {
	return model.Card{
		Number:          "2321010",
		Type:            "道具",
		Name:            "幻术卷轴",
		Category:        "气",
		Tag:             "消耗品-卷轴",
		Description:     "反制:当敌方使用法术攻击时,重新排列你场上的所有单位,对方需要重新选择目标",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321010.jpg",
	}
}

type CardDef2321011 struct{}

func (CardDef2321011) ID() string      { return "2321011" }
func (CardDef2321011) Name() string    { return "传送符文" }
func (CardDef2321011) Kind() string    { return "道具" }
func (CardDef2321011) Element() string { return "气" }

func (CardDef2321011) Card() model.Card {
	return model.Card{
		Number:          "2321011",
		Type:            "道具",
		Name:            "传送符文",
		Category:        "气",
		Tag:             "消耗品-符文",
		Description:     "反制:当1个伙伴被召唤或消耗后,将其移动到另一位置",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321011.jpg",
	}
}

type CardDef2321012 struct{}

func (CardDef2321012) ID() string      { return "2321012" }
func (CardDef2321012) Name() string    { return "随风斗篷" }
func (CardDef2321012) Kind() string    { return "道具" }
func (CardDef2321012) Element() string { return "气" }

func (CardDef2321012) Card() model.Card {
	return model.Card{
		Number:          "2321012",
		Type:            "道具",
		Name:            "随风斗篷",
		Category:        "气",
		Tag:             "装备-防具",
		Description:     "主动绝技:移动你的人物至另一位置",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321012.jpg",
	}
}

type CardDef2321013 struct{}

func (CardDef2321013) ID() string      { return "2321013" }
func (CardDef2321013) Name() string    { return "驭风杖" }
func (CardDef2321013) Kind() string    { return "道具" }
func (CardDef2321013) Element() string { return "气" }

func (CardDef2321013) Card() model.Card {
	return model.Card{
		Number:          "2321013",
		Type:            "道具",
		Name:            "驭风杖",
		Category:        "气",
		Tag:             "装备-武器",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 5},
		ElementsGain:    map[string]int{"气": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321013.jpg",
	}
}

type CardDef2321014 struct{}

func (CardDef2321014) ID() string      { return "2321014" }
func (CardDef2321014) Name() string    { return "聆风羽毛笔" }
func (CardDef2321014) Kind() string    { return "道具" }
func (CardDef2321014) Element() string { return "气" }

func (CardDef2321014) Card() model.Card {
	return model.Card{
		Number:          "2321014",
		Type:            "道具",
		Name:            "聆风羽毛笔",
		Category:        "气",
		Tag:             "装备",
		Description:     "入场:放置3个标记物.主动回合技:消耗此卡并取除1个标记物才能发动,抽1张牌,获得1\\气",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\气\\2321014.jpg",
	}
}

type CardDef2321101 struct{}

func (CardDef2321101) ID() string      { return "2321101" }
func (CardDef2321101) Name() string    { return "雷之链" }
func (CardDef2321101) Kind() string    { return "道具" }
func (CardDef2321101) Element() string { return "气" }

func (CardDef2321101) Card() model.Card {
	return model.Card{
		Number:          "2321101",
		Type:            "道具",
		Name:            "雷之链",
		Category:        "气",
		Tag:             "装备-武器",
		Description:     "主动:消耗此卡才能发动,你的下一次驱动法术可以额外选择1个无视范围的目标",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "气": 2},
		ElementsGain:    map[string]int{"光": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321101.jpg",
	}
}

type CardDef2321102 struct{}

func (CardDef2321102) ID() string      { return "2321102" }
func (CardDef2321102) Name() string    { return "风之轮回" }
func (CardDef2321102) Kind() string    { return "道具" }
func (CardDef2321102) Element() string { return "气" }

func (CardDef2321102) Card() model.Card {
	return model.Card{
		Number:          "2321102",
		Type:            "道具",
		Name:            "风之轮回",
		Category:        "气",
		Tag:             "装备-神器",
		Description:     "主动:消耗并献祭此卡才能发动,从弃牌堆将任意数量的大气卡牌洗回卡组",
		Quote:           "那些逝去的已经归来,我在风中看到了它们的身影",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321102.jpg",
	}
}

type CardDef2321103 struct{}

func (CardDef2321103) ID() string      { return "2321103" }
func (CardDef2321103) Name() string    { return "雷鸣之息" }
func (CardDef2321103) Kind() string    { return "道具" }
func (CardDef2321103) Element() string { return "气" }

func (CardDef2321103) Card() model.Card {
	return model.Card{
		Number:          "2321103",
		Type:            "道具",
		Name:            "雷鸣之息",
		Category:        "气",
		Tag:             "消耗品-药剂",
		Description:     "此卡被使用或从手牌丢弃时获得1\\气",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321103.jpg",
	}
}

type CardDef2321104 struct{}

func (CardDef2321104) ID() string      { return "2321104" }
func (CardDef2321104) Name() string    { return "雷光头冠" }
func (CardDef2321104) Kind() string    { return "道具" }
func (CardDef2321104) Element() string { return "气" }

func (CardDef2321104) Card() model.Card {
	return model.Card{
		Number:          "2321104",
		Type:            "道具",
		Name:            "雷光头冠",
		Category:        "气",
		Tag:             "装备-饰物",
		Description:     "祈咒:你的下一次聚能法术+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "气": 2},
		ElementsGain:    map[string]int{"光": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321104.jpg",
	}
}

type CardDef2321105 struct{}

func (CardDef2321105) ID() string      { return "2321105" }
func (CardDef2321105) Name() string    { return "雷光战铠" }
func (CardDef2321105) Kind() string    { return "道具" }
func (CardDef2321105) Element() string { return "气" }

func (CardDef2321105) Card() model.Card {
	return model.Card{
		Number:          "2321105",
		Type:            "道具",
		Name:            "雷光战铠",
		Category:        "气",
		Tag:             "装备-防具",
		Description:     "光环:如果你同时装备至少3张雷光道具(同时具有\\气和\\光负载),你的聚能和驱动法术+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "气": 4},
		ElementsGain:    map[string]int{"光": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321105.jpg",
	}
}

type CardDef2321106 struct{}

func (CardDef2321106) ID() string      { return "2321106" }
func (CardDef2321106) Name() string    { return "无尽风潮" }
func (CardDef2321106) Kind() string    { return "道具" }
func (CardDef2321106) Element() string { return "气" }

func (CardDef2321106) Card() model.Card {
	return model.Card{
		Number:          "2321106",
		Type:            "道具",
		Name:            "无尽风潮",
		Category:        "气",
		Tag:             "消耗品-法术卷轴-驱动",
		Description:     "命中:本卡不送去弃牌堆改为回到你的手牌,本卡永久获得+2\\威,入场费用+1\\气,但本回合不能再使用",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           2,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321106.jpg",
	}
}

type CardDef2321107 struct{}

func (CardDef2321107) ID() string      { return "2321107" }
func (CardDef2321107) Name() string    { return "飞鸽拘捕令" }
func (CardDef2321107) Kind() string    { return "道具" }
func (CardDef2321107) Element() string { return "气" }

func (CardDef2321107) Card() model.Card {
	return model.Card{
		Number:          "2321107",
		Type:            "道具",
		Name:            "飞鸽拘捕令",
		Category:        "气",
		Tag:             "装备",
		Description:     "诱发回合技:当你的法术命中后,将1张九霄印记加入对手手牌",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2001102"},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321107.jpg",
	}
}

type CardDef2321108 struct{}

func (CardDef2321108) ID() string      { return "2321108" }
func (CardDef2321108) Name() string    { return "散去" }
func (CardDef2321108) Kind() string    { return "道具" }
func (CardDef2321108) Element() string { return "气" }

func (CardDef2321108) Card() model.Card {
	return model.Card{
		Number:          "2321108",
		Type:            "道具",
		Name:            "散去",
		Category:        "气",
		Tag:             "消耗品-卷轴",
		Description:     "反制:在友方大气单位受到1次伤害后,直到你的回合开始,使其免疫所有伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321108.jpg",
	}
}

type CardDef2321109 struct{}

func (CardDef2321109) ID() string      { return "2321109" }
func (CardDef2321109) Name() string    { return "幻雾面罩" }
func (CardDef2321109) Kind() string    { return "道具" }
func (CardDef2321109) Element() string { return "气" }

func (CardDef2321109) Card() model.Card {
	return model.Card{
		Number:          "2321109",
		Type:            "道具",
		Name:            "幻雾面罩",
		Category:        "气",
		Tag:             "装备-防具",
		Description:     "诱发绝技:敌方法术命中时,丢弃最多3张手牌,使该法术\\攻下降丢弃的数值",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{"气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321109.jpg",
	}
}

type CardDef2321110 struct{}

func (CardDef2321110) ID() string      { return "2321110" }
func (CardDef2321110) Name() string    { return "飞鸽急袭令" }
func (CardDef2321110) Kind() string    { return "道具" }
func (CardDef2321110) Element() string { return "气" }

func (CardDef2321110) Card() model.Card {
	return model.Card{
		Number:          "2321110",
		Type:            "道具",
		Name:            "飞鸽急袭令",
		Category:        "气",
		Tag:             "消耗品-卷轴",
		Description:     "选择你在本回合学习的速攻法术,使其下一次使用时+1\\攻+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321110.jpg",
	}
}

type CardDef2321111 struct{}

func (CardDef2321111) ID() string      { return "2321111" }
func (CardDef2321111) Name() string    { return "反击风洞卷轴" }
func (CardDef2321111) Kind() string    { return "道具" }
func (CardDef2321111) Element() string { return "气" }

func (CardDef2321111) Card() model.Card {
	return model.Card{
		Number:          "2321111",
		Type:            "道具",
		Name:            "反击风洞卷轴",
		Category:        "气",
		Tag:             "消耗品-卷轴",
		Description:     "反制:在敌方1个范围效果以外的法术未命中或被无效后,视为对任意1个敌方单位使用该法术,包含所有强化效果.",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321111.jpg",
	}
}

type CardDef2321112 struct{}

func (CardDef2321112) ID() string      { return "2321112" }
func (CardDef2321112) Name() string    { return "撕裂冲击卷轴" }
func (CardDef2321112) Kind() string    { return "道具" }
func (CardDef2321112) Element() string { return "气" }

func (CardDef2321112) Card() model.Card {
	return model.Card{
		Number:          "2321112",
		Type:            "道具",
		Name:            "撕裂冲击卷轴",
		Category:        "气",
		Tag:             "消耗品-法术卷轴-聚能",
		Description:     "范围:纵列.命中:将共计3点伤害分配给目标范围内的单位.",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\气\\2321112.jpg",
	}
}

type CardDef2411001 struct{}

func (CardDef2411001) ID() string      { return "2411001" }
func (CardDef2411001) Name() string    { return "古树之心" }
func (CardDef2411001) Kind() string    { return "道具" }
func (CardDef2411001) Element() string { return "地" }

func (CardDef2411001) Card() model.Card {
	return model.Card{
		Number:          "2411001",
		Type:            "道具",
		Name:            "古树之心",
		Category:        "地",
		Tag:             "传奇-装备-神器",
		Description:     "诱发回合技:友方单位获得负载时可以使其+1\\血,或获得生命时可以使其负载+1\\地",
		Quote:           "在这参天的身躯之下潜藏的意志究竟是什么或许只是知识本身",
		ElementsCost:    map[string]int{"地": 5},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2411001.jpg",
	}
}

type CardDef2411002 struct{}

func (CardDef2411002) ID() string      { return "2411002" }
func (CardDef2411002) Name() string    { return "裂地巨剑 阿托比斯" }
func (CardDef2411002) Kind() string    { return "道具" }
func (CardDef2411002) Element() string { return "地" }

func (CardDef2411002) Card() model.Card {
	return model.Card{
		Number:          "2411002",
		Type:            "道具",
		Name:            "裂地巨剑 阿托比斯",
		Category:        "地",
		Tag:             "传奇-装备-武器",
		Description:     "主动:消耗此卡才能发动,本回合下一次法术获得+4\\威且范围变为前排,或者+2\\攻且范围变为纵列",
		Quote:           "\"太无聊了,让我们一起闹出点大动静吧\"",
		ElementsCost:    map[string]int{"地": 6},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2411002.jpg",
	}
}

type CardDef2411101 struct{}

func (CardDef2411101) ID() string      { return "2411101" }
func (CardDef2411101) Name() string    { return "翡翠永生" }
func (CardDef2411101) Kind() string    { return "道具" }
func (CardDef2411101) Element() string { return "地" }

func (CardDef2411101) Card() model.Card {
	return model.Card{
		Number:          "2411101",
		Type:            "道具",
		Name:            "翡翠永生",
		Category:        "地",
		Tag:             "传奇-装备-饰物",
		Description:     "光环:如果你拥有护盾,防止所有友方单位受到的伤害,也不会受负面状态影响(仍可处于负面状态).入场:获得护盾2",
		Quote:           "永远的守护,永生的束缚",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2411101.jpg",
	}
}

type CardDef2411102 struct{}

func (CardDef2411102) ID() string      { return "2411102" }
func (CardDef2411102) Name() string    { return "腐朽的古树之心" }
func (CardDef2411102) Kind() string    { return "道具" }
func (CardDef2411102) Element() string { return "地" }

func (CardDef2411102) Card() model.Card {
	return model.Card{
		Number:          "2411102",
		Type:            "道具",
		Name:            "腐朽的古树之心",
		Category:        "地",
		Tag:             "传奇-装备-神器",
		Description:     "诱发:双方玩家每使用2个法术,该玩家必须移除自己场上1点负载.诱发:如果此卡失去所有负载,摧毁此卡",
		Quote:           "这还算是原来的那个古树之心吗?",
		ElementsCost:    map[string]int{"地": 5},
		ElementsGain:    map[string]int{"地": 2, "暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2411102.jpg",
	}
}

type CardDef2421001 struct{}

func (CardDef2421001) ID() string      { return "2421001" }
func (CardDef2421001) Name() string    { return "知识古树的关怀" }
func (CardDef2421001) Kind() string    { return "道具" }
func (CardDef2421001) Element() string { return "地" }

func (CardDef2421001) Card() model.Card {
	return model.Card{
		Number:          "2421001",
		Type:            "道具",
		Name:            "知识古树的关怀",
		Category:        "地",
		Tag:             "装备-武器",
		Description:     "诱发:当你的卡牌达到精通时,可以消耗此卡来发动,抽1张牌并获得1\\地",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421001.jpg",
	}
}

type CardDef2421002 struct{}

func (CardDef2421002) ID() string      { return "2421002" }
func (CardDef2421002) Name() string    { return "生长药水" }
func (CardDef2421002) Kind() string    { return "道具" }
func (CardDef2421002) Element() string { return "地" }

func (CardDef2421002) Card() model.Card {
	return model.Card{
		Number:          "2421002",
		Type:            "道具",
		Name:            "生长药水",
		Category:        "地",
		Tag:             "消耗品-药剂",
		Description:     "重置你的1个地脉伙伴",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421002.jpg",
	}
}

type CardDef2421003 struct{}

func (CardDef2421003) ID() string      { return "2421003" }
func (CardDef2421003) Name() string    { return "坚固卷轴" }
func (CardDef2421003) Kind() string    { return "道具" }
func (CardDef2421003) Element() string { return "地" }

func (CardDef2421003) Card() model.Card {
	return model.Card{
		Number:          "2421003",
		Type:            "道具",
		Name:            "坚固卷轴",
		Category:        "地",
		Tag:             "消耗品-卷轴",
		Description:     "直到下个回合结束,使1个友方单位免疫最多3点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421003.jpg",
	}
}

type CardDef2421004 struct{}

func (CardDef2421004) ID() string      { return "2421004" }
func (CardDef2421004) Name() string    { return "德鲁伊水平测试" }
func (CardDef2421004) Kind() string    { return "道具" }
func (CardDef2421004) Element() string { return "地" }

func (CardDef2421004) Card() model.Card {
	return model.Card{
		Number:          "2421004",
		Type:            "道具",
		Name:            "德鲁伊水平测试",
		Category:        "地",
		Tag:             "消耗品-卷轴",
		Description:     "你所有负载大于2的伙伴获得负载+1\\地",
		Quote:           "\"我重申一遍,考试期间禁止向那棵老树提问!\"——大德鲁伊烟尘",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421004.jpg",
	}
}

type CardDef2421005 struct{}

func (CardDef2421005) ID() string      { return "2421005" }
func (CardDef2421005) Name() string    { return "石化卷轴" }
func (CardDef2421005) Kind() string    { return "道具" }
func (CardDef2421005) Element() string { return "地" }

func (CardDef2421005) Card() model.Card {
	return model.Card{
		Number:          "2421005",
		Type:            "道具",
		Name:            "石化卷轴",
		Category:        "地",
		Tag:             "消耗品-卷轴",
		Description:     "使1个无视范围的单位石化2",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421005.jpg",
	}
}

type CardDef2421006 struct{}

func (CardDef2421006) ID() string      { return "2421006" }
func (CardDef2421006) Name() string    { return "磐藤胸甲" }
func (CardDef2421006) Kind() string    { return "道具" }
func (CardDef2421006) Element() string { return "地" }

func (CardDef2421006) Card() model.Card {
	return model.Card{
		Number:          "2421006",
		Type:            "道具",
		Name:            "磐藤胸甲",
		Category:        "地",
		Tag:             "装备-防具",
		Description:     "入场:你的人物获得+2\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421006.jpg",
	}
}

type CardDef2421007 struct{}

func (CardDef2421007) ID() string      { return "2421007" }
func (CardDef2421007) Name() string    { return "寄生之触" }
func (CardDef2421007) Kind() string    { return "道具" }
func (CardDef2421007) Element() string { return "地" }

func (CardDef2421007) Card() model.Card {
	return model.Card{
		Number:          "2421007",
		Type:            "道具",
		Name:            "寄生之触",
		Category:        "地",
		Tag:             "装备",
		Description:     "精通1:负载+1\\地",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421007.jpg",
	}
}

type CardDef2421008 struct{}

func (CardDef2421008) ID() string      { return "2421008" }
func (CardDef2421008) Name() string    { return "巨石阵卷轴" }
func (CardDef2421008) Kind() string    { return "道具" }
func (CardDef2421008) Element() string { return "地" }

func (CardDef2421008) Card() model.Card {
	return model.Card{
		Number:          "2421008",
		Type:            "道具",
		Name:            "巨石阵卷轴",
		Category:        "地",
		Tag:             "消耗品-法术卷轴-创造",
		Description:     "范围:方阵",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2, "无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421008.jpg",
	}
}

type CardDef2421009 struct{}

func (CardDef2421009) ID() string      { return "2421009" }
func (CardDef2421009) Name() string    { return "森林之矢卷轴" }
func (CardDef2421009) Kind() string    { return "道具" }
func (CardDef2421009) Element() string { return "地" }

func (CardDef2421009) Card() model.Card {
	return model.Card{
		Number:          "2421009",
		Type:            "道具",
		Name:            "森林之矢卷轴",
		Category:        "地",
		Tag:             "消耗品-法术卷轴-驱动",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421009.jpg",
	}
}

type CardDef2421010 struct{}

func (CardDef2421010) ID() string      { return "2421010" }
func (CardDef2421010) Name() string    { return "自然封印卷轴" }
func (CardDef2421010) Kind() string    { return "道具" }
func (CardDef2421010) Element() string { return "地" }

func (CardDef2421010) Card() model.Card {
	return model.Card{
		Number:          "2421010",
		Type:            "道具",
		Name:            "自然封印卷轴",
		Category:        "地",
		Tag:             "消耗品-卷轴",
		Description:     "直到你下个回合的回合结束,所有法术\\攻变为0",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421010.jpg",
	}
}

type CardDef2421011 struct{}

func (CardDef2421011) ID() string      { return "2421011" }
func (CardDef2421011) Name() string    { return "精灵铠" }
func (CardDef2421011) Kind() string    { return "道具" }
func (CardDef2421011) Element() string { return "地" }

func (CardDef2421011) Card() model.Card {
	return model.Card{
		Number:          "2421011",
		Type:            "道具",
		Name:            "精灵铠",
		Category:        "地",
		Tag:             "装备-防具",
		Description:     "祈咒:为你的人物回复1\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3, "无": 1},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421011.jpg",
	}
}

type CardDef2421012 struct{}

func (CardDef2421012) ID() string      { return "2421012" }
func (CardDef2421012) Name() string    { return "地脉灵石" }
func (CardDef2421012) Kind() string    { return "道具" }
func (CardDef2421012) Element() string { return "地" }

func (CardDef2421012) Card() model.Card {
	return model.Card{
		Number:          "2421012",
		Type:            "道具",
		Name:            "地脉灵石",
		Category:        "地",
		Tag:             "装备-神器",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 6, "无": 1},
		ElementsGain:    map[string]int{"地": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421012.jpg",
	}
}

type CardDef2421013 struct{}

func (CardDef2421013) ID() string      { return "2421013" }
func (CardDef2421013) Name() string    { return "《地理学入门》" }
func (CardDef2421013) Kind() string    { return "道具" }
func (CardDef2421013) Element() string { return "地" }

func (CardDef2421013) Card() model.Card {
	return model.Card{
		Number:          "2421013",
		Type:            "道具",
		Name:            "《地理学入门》",
		Category:        "地",
		Tag:             "装备",
		Description:     "光环:你原始入场花费大于5的卡牌入场费用减少2\\地",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421013.jpg",
	}
}

type CardDef2421014 struct{}

func (CardDef2421014) ID() string      { return "2421014" }
func (CardDef2421014) Name() string    { return "森之贮藏" }
func (CardDef2421014) Kind() string    { return "道具" }
func (CardDef2421014) Element() string { return "地" }

func (CardDef2421014) Card() model.Card {
	return model.Card{
		Number:          "2421014",
		Type:            "道具",
		Name:            "森之贮藏",
		Category:        "地",
		Tag:             "装备",
		Description:     "入场:放置4个标记物.主动:消耗此卡并取除1个标记物才能发动,获得4\\地",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\地\\2421014.jpg",
	}
}

type CardDef2421101 struct{}

func (CardDef2421101) ID() string      { return "2421101" }
func (CardDef2421101) Name() string    { return "秋暮耳环" }
func (CardDef2421101) Kind() string    { return "道具" }
func (CardDef2421101) Element() string { return "地" }

func (CardDef2421101) Card() model.Card {
	return model.Card{
		Number:          "2421101",
		Type:            "道具",
		Name:            "秋暮耳环",
		Category:        "地",
		Tag:             "装备-饰物",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421101.jpg",
	}
}

type CardDef2421102 struct{}

func (CardDef2421102) ID() string      { return "2421102" }
func (CardDef2421102) Name() string    { return "蔷薇之鞭" }
func (CardDef2421102) Kind() string    { return "道具" }
func (CardDef2421102) Element() string { return "地" }

func (CardDef2421102) Card() model.Card {
	return model.Card{
		Number:          "2421102",
		Type:            "道具",
		Name:            "蔷薇之鞭",
		Category:        "地",
		Tag:             "装备-武器",
		Description:     "回合技:每当友方卡牌负载数量减少后,此卡获得负载+1\\暗,最多2点",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2, "无": 1},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421102.jpg",
	}
}

type CardDef2421103 struct{}

func (CardDef2421103) ID() string      { return "2421103" }
func (CardDef2421103) Name() string    { return "捕梦网" }
func (CardDef2421103) Kind() string    { return "道具" }
func (CardDef2421103) Element() string { return "地" }

func (CardDef2421103) Card() model.Card {
	return model.Card{
		Number:          "2421103",
		Type:            "道具",
		Name:            "捕梦网",
		Category:        "地",
		Tag:             "装备-武器",
		Description:     "入场:你已学习的灵媒法术永久+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{"地": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421103.jpg",
	}
}

type CardDef2421104 struct{}

func (CardDef2421104) ID() string      { return "2421104" }
func (CardDef2421104) Name() string    { return "血蔷薇契约" }
func (CardDef2421104) Kind() string    { return "道具" }
func (CardDef2421104) Element() string { return "地" }

func (CardDef2421104) Card() model.Card {
	return model.Card{
		Number:          "2421104",
		Type:            "道具",
		Name:            "血蔷薇契约",
		Category:        "地",
		Tag:             "消耗品-卷轴",
		Description:     "将你的一个法术变为你的一个地脉或暗影伙伴的绑定技能,并且上升该伙伴负载数量的\\威(在其死亡后将该绑定技能从游戏中移除)",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2, "无": 1, "暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421104.jpg",
	}
}

type CardDef2421105 struct{}

func (CardDef2421105) ID() string      { return "2421105" }
func (CardDef2421105) Name() string    { return "自然交感" }
func (CardDef2421105) Kind() string    { return "道具" }
func (CardDef2421105) Element() string { return "地" }

func (CardDef2421105) Card() model.Card {
	return model.Card{
		Number:          "2421105",
		Type:            "道具",
		Name:            "自然交感",
		Category:        "地",
		Tag:             "消耗品-卷轴",
		Description:     "选择你场上2个地脉伙伴,重新分配他们的负载",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421105.jpg",
	}
}

type CardDef2421106 struct{}

func (CardDef2421106) ID() string      { return "2421106" }
func (CardDef2421106) Name() string    { return "苍老药剂" }
func (CardDef2421106) Kind() string    { return "道具" }
func (CardDef2421106) Element() string { return "地" }

func (CardDef2421106) Card() model.Card {
	return model.Card{
		Number:          "2421106",
		Type:            "道具",
		Name:            "苍老药剂",
		Category:        "地",
		Tag:             "消耗品-药剂",
		Description:     "移除友方卡牌负载的1\\地,使其立刻达到下一次精通",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421106.jpg",
	}
}

type CardDef2421107 struct{}

func (CardDef2421107) ID() string      { return "2421107" }
func (CardDef2421107) Name() string    { return "翡翠结界卷轴" }
func (CardDef2421107) Kind() string    { return "道具" }
func (CardDef2421107) Element() string { return "地" }

func (CardDef2421107) Card() model.Card {
	return model.Card{
		Number:          "2421107",
		Type:            "道具",
		Name:            "翡翠结界卷轴",
		Category:        "地",
		Tag:             "消耗品-卷轴",
		Description:     "当敌方法术数量比我方多时,每多1个获得护盾1",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421107.jpg",
	}
}

type CardDef2421108 struct{}

func (CardDef2421108) ID() string      { return "2421108" }
func (CardDef2421108) Name() string    { return "翡翠果" }
func (CardDef2421108) Kind() string    { return "道具" }
func (CardDef2421108) Element() string { return "地" }

func (CardDef2421108) Card() model.Card {
	return model.Card{
		Number:          "2421108",
		Type:            "道具",
		Name:            "翡翠果",
		Category:        "地",
		Tag:             "装备",
		Description:     "入场:使一个友方伙伴获得除\\地与\\无外的任意1点负载",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "无": 2},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421108.jpg",
	}
}

type CardDef2421109 struct{}

func (CardDef2421109) ID() string      { return "2421109" }
func (CardDef2421109) Name() string    { return "地穴精灵矿镐" }
func (CardDef2421109) Kind() string    { return "道具" }
func (CardDef2421109) Element() string { return "地" }

func (CardDef2421109) Card() model.Card {
	return model.Card{
		Number:          "2421109",
		Type:            "道具",
		Name:            "地穴精灵矿镐",
		Category:        "地",
		Tag:             "装备-武器",
		Description:     "消耗:选择1个种类(伙伴或道具),在5张牌之内翻取1张所选择种类的卡牌",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421109.jpg",
	}
}

type CardDef2421110 struct{}

func (CardDef2421110) ID() string      { return "2421110" }
func (CardDef2421110) Name() string    { return "沙虫之饵" }
func (CardDef2421110) Kind() string    { return "道具" }
func (CardDef2421110) Element() string { return "地" }

func (CardDef2421110) Card() model.Card {
	return model.Card{
		Number:          "2421110",
		Type:            "道具",
		Name:            "沙虫之饵",
		Category:        "地",
		Tag:             "消耗品",
		Description:     "翻取1个入场花费大于5的地脉伙伴,如果是巨型沙虫则使其入场花费-2",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421110.jpg",
	}
}

type CardDef2421111 struct{}

func (CardDef2421111) ID() string      { return "2421111" }
func (CardDef2421111) Name() string    { return "沙漠护腿" }
func (CardDef2421111) Kind() string    { return "道具" }
func (CardDef2421111) Element() string { return "地" }

func (CardDef2421111) Card() model.Card {
	return model.Card{
		Number:          "2421111",
		Type:            "道具",
		Name:            "沙漠护腿",
		Category:        "地",
		Tag:             "装备-防具",
		Description:     "绝技:当1个友方单位受到2点及以上伤害时,使该伤害-2",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421111.jpg",
	}
}

type CardDef2421112 struct{}

func (CardDef2421112) ID() string      { return "2421112" }
func (CardDef2421112) Name() string    { return "秋枫宝钻" }
func (CardDef2421112) Kind() string    { return "道具" }
func (CardDef2421112) Element() string { return "地" }

func (CardDef2421112) Card() model.Card {
	return model.Card{
		Number:          "2421112",
		Type:            "道具",
		Name:            "秋枫宝钻",
		Category:        "地",
		Tag:             "装备-神器",
		Description:     "入场:放置2个标记物.回合技:移除1个标记物,重置你的1个地脉伙伴",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{"地": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\地\\2421112.jpg",
	}
}

type CardDef2501001 struct{}

func (CardDef2501001) ID() string      { return "2501001" }
func (CardDef2501001) Name() string    { return "桎梏" }
func (CardDef2501001) Kind() string    { return "道具" }
func (CardDef2501001) Element() string { return "光" }

func (CardDef2501001) Card() model.Card {
	return model.Card{
		Number:          "2501001",
		Type:            "道具",
		Name:            "桎梏",
		Category:        "光",
		Tag:             "衍生",
		Description:     "当你抽到这张牌时(起始手牌除外),必须将其展示并丢弃,之后你可以再抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2501001.jpg",
	}
}

type CardDef2511001 struct{}

func (CardDef2511001) ID() string      { return "2511001" }
func (CardDef2511001) Name() string    { return "万灵药" }
func (CardDef2511001) Kind() string    { return "道具" }
func (CardDef2511001) Element() string { return "光" }

func (CardDef2511001) Card() model.Card {
	return model.Card{
		Number:          "2511001",
		Type:            "道具",
		Name:            "万灵药",
		Category:        "光",
		Tag:             "传奇-消耗品-药剂",
		Description:     "回复1个友方单位所有生命,或抽4张牌,或获得5\\无,或重置你的1个技能",
		Quote:           "它曾历经数人之手,却无人舍得喝下一口",
		ElementsCost:    map[string]int{"光": 2, "无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2511001.jpg",
	}
}

type CardDef2511002 struct{}

func (CardDef2511002) ID() string      { return "2511002" }
func (CardDef2511002) Name() string    { return "辉之盾 闪耀" }
func (CardDef2511002) Kind() string    { return "道具" }
func (CardDef2511002) Element() string { return "光" }

func (CardDef2511002) Card() model.Card {
	return model.Card{
		Number:          "2511002",
		Type:            "道具",
		Name:            "辉之盾 闪耀",
		Category:        "光",
		Tag:             "传奇-装备-防具",
		Description:     "光环:你在防御时额外获得2\\威.诱发回合技:当你防御成功时,对法力范围内所有敌人造成晕眩1",
		Quote:           "真正的天才会在盾里嵌上一面镜子",
		ElementsCost:    map[string]int{"光": 5},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2511002.jpg",
	}
}

type CardDef2511101 struct{}

func (CardDef2511101) ID() string      { return "2511101" }
func (CardDef2511101) Name() string    { return "九霄辉迹" }
func (CardDef2511101) Kind() string    { return "道具" }
func (CardDef2511101) Element() string { return "光" }

func (CardDef2511101) Card() model.Card {
	return model.Card{
		Number:          "2511101",
		Type:            "道具",
		Name:            "九霄辉迹",
		Category:        "光",
		Tag:             "传奇-装备-武器",
		Description:     "绝技:双方将手牌全部丢弃,然后抽等量的牌",
		Quote:           "\"你看,我们也付出了巨大的牺牲\"",
		ElementsCost:    map[string]int{"光": 4, "气": 2},
		ElementsGain:    map[string]int{"光": 2, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2511101.jpg",
	}
}

type CardDef2511102 struct{}

func (CardDef2511102) ID() string      { return "2511102" }
func (CardDef2511102) Name() string    { return "五虹之环" }
func (CardDef2511102) Kind() string    { return "道具" }
func (CardDef2511102) Element() string { return "光" }

func (CardDef2511102) Card() model.Card {
	return model.Card{
		Number:          "2511102",
		Type:            "道具",
		Name:            "五虹之环",
		Category:        "光",
		Tag:             "传奇-装备-饰物",
		Description:     "回合技:花费1点\\火\\水\\地\\气\\光种类的元素,放置1个该种类的标记物.绑定技能:五虹之束",
		Quote:           "后来此环被遗落在某处遗迹的宝阁之内,又被某东方修士意外拾得,熠熠之辉不减当年",
		ElementsCost:    map[string]int{"光": 1, "地": 1, "气": 1, "水": 1, "火": 1},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"3501101"},
		OutputPath:      "output\\王权纷争\\道具\\光\\2511102.jpg",
	}
}

type CardDef2521001 struct{}

func (CardDef2521001) ID() string      { return "2521001" }
func (CardDef2521001) Name() string    { return "生命药剂" }
func (CardDef2521001) Kind() string    { return "道具" }
func (CardDef2521001) Element() string { return "光" }

func (CardDef2521001) Card() model.Card {
	return model.Card{
		Number:          "2521001",
		Type:            "道具",
		Name:            "生命药剂",
		Category:        "光",
		Tag:             "消耗品-药剂",
		Description:     "选择1个友方单位,使其回复2\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521001.jpg",
	}
}

type CardDef2521002 struct{}

func (CardDef2521002) ID() string      { return "2521002" }
func (CardDef2521002) Name() string    { return "庇护符文" }
func (CardDef2521002) Kind() string    { return "道具" }
func (CardDef2521002) Element() string { return "光" }

func (CardDef2521002) Card() model.Card {
	return model.Card{
		Number:          "2521002",
		Type:            "道具",
		Name:            "庇护符文",
		Category:        "光",
		Tag:             "消耗品-符文",
		Description:     "反制:当1个威力小于10的敌方法术命中时,将其无效",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521002.jpg",
	}
}

type CardDef2521003 struct{}

func (CardDef2521003) ID() string      { return "2521003" }
func (CardDef2521003) Name() string    { return "净化卷轴" }
func (CardDef2521003) Kind() string    { return "道具" }
func (CardDef2521003) Element() string { return "光" }

func (CardDef2521003) Card() model.Card {
	return model.Card{
		Number:          "2521003",
		Type:            "道具",
		Name:            "净化卷轴",
		Category:        "光",
		Tag:             "消耗品-卷轴",
		Description:     "移除1个友方卡牌所有负面状态或任意1个敌方卡牌所有标记物",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521003.jpg",
	}
}

type CardDef2521004 struct{}

func (CardDef2521004) ID() string      { return "2521004" }
func (CardDef2521004) Name() string    { return "神圣制裁卷轴" }
func (CardDef2521004) Kind() string    { return "道具" }
func (CardDef2521004) Element() string { return "光" }

func (CardDef2521004) Card() model.Card {
	return model.Card{
		Number:          "2521004",
		Type:            "道具",
		Name:            "神圣制裁卷轴",
		Category:        "光",
		Tag:             "消耗品-卷轴",
		Description:     "反制:敌方使用咒术时,无效敌人的那个技能",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521004.jpg",
	}
}

type CardDef2521005 struct{}

func (CardDef2521005) ID() string      { return "2521005" }
func (CardDef2521005) Name() string    { return "新生卷轴" }
func (CardDef2521005) Kind() string    { return "道具" }
func (CardDef2521005) Element() string { return "光" }

func (CardDef2521005) Card() model.Card {
	return model.Card{
		Number:          "2521005",
		Type:            "道具",
		Name:            "新生卷轴",
		Category:        "光",
		Tag:             "消耗品-卷轴",
		Description:     "选择你的一个弃牌堆中的光辉伙伴,支付其入场花费才能发动,将其复活",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521005.jpg",
	}
}

type CardDef2521006 struct{}

func (CardDef2521006) ID() string      { return "2521006" }
func (CardDef2521006) Name() string    { return "绿玉权杖" }
func (CardDef2521006) Kind() string    { return "道具" }
func (CardDef2521006) Element() string { return "光" }

func (CardDef2521006) Card() model.Card {
	return model.Card{
		Number:          "2521006",
		Type:            "道具",
		Name:            "绿玉权杖",
		Category:        "光",
		Tag:             "装备-武器",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521006.jpg",
	}
}

type CardDef2521007 struct{}

func (CardDef2521007) ID() string      { return "2521007" }
func (CardDef2521007) Name() string    { return "蓝晶灯盏" }
func (CardDef2521007) Kind() string    { return "道具" }
func (CardDef2521007) Element() string { return "光" }

func (CardDef2521007) Card() model.Card {
	return model.Card{
		Number:          "2521007",
		Type:            "道具",
		Name:            "蓝晶灯盏",
		Category:        "光",
		Tag:             "装备",
		Description:     "主动绝技:花费5\\光,负载+2\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521007.jpg",
	}
}

type CardDef2521008 struct{}

func (CardDef2521008) ID() string      { return "2521008" }
func (CardDef2521008) Name() string    { return "惩戒之箭卷轴" }
func (CardDef2521008) Kind() string    { return "道具" }
func (CardDef2521008) Element() string { return "光" }

func (CardDef2521008) Card() model.Card {
	return model.Card{
		Number:          "2521008",
		Type:            "道具",
		Name:            "惩戒之箭卷轴",
		Category:        "光",
		Tag:             "消耗品-法术卷轴-创造",
		Description:     "穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521008.jpg",
	}
}

type CardDef2521009 struct{}

func (CardDef2521009) ID() string      { return "2521009" }
func (CardDef2521009) Name() string    { return "光之刃卷轴" }
func (CardDef2521009) Kind() string    { return "道具" }
func (CardDef2521009) Element() string { return "光" }

func (CardDef2521009) Card() model.Card {
	return model.Card{
		Number:          "2521009",
		Type:            "道具",
		Name:            "光之刃卷轴",
		Category:        "光",
		Tag:             "消耗品-法术卷轴-创造",
		Description:     "范围:前排",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521009.jpg",
	}
}

type CardDef2521010 struct{}

func (CardDef2521010) ID() string      { return "2521010" }
func (CardDef2521010) Name() string    { return "闪耀水晶" }
func (CardDef2521010) Kind() string    { return "道具" }
func (CardDef2521010) Element() string { return "光" }

func (CardDef2521010) Card() model.Card {
	return model.Card{
		Number:          "2521010",
		Type:            "道具",
		Name:            "闪耀水晶",
		Category:        "光",
		Tag:             "装备-神器",
		Description:     "光环:你的光辉法术获得晕眩1",
		Quote:           "辉之盾终将破碎,而团结也难延续,所剩的闪耀水晶则变成了贵胄家族各自的珍藏",
		ElementsCost:    map[string]int{"光": 5},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521010.jpg",
	}
}

type CardDef2521011 struct{}

func (CardDef2521011) ID() string      { return "2521011" }
func (CardDef2521011) Name() string    { return "闪光符文" }
func (CardDef2521011) Kind() string    { return "道具" }
func (CardDef2521011) Element() string { return "光" }

func (CardDef2521011) Card() model.Card {
	return model.Card{
		Number:          "2521011",
		Type:            "道具",
		Name:            "闪光符文",
		Category:        "光",
		Tag:             "消耗品-符文",
		Description:     "反制:当敌方使用技能时,使所有前排敌人晕眩1",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521011.jpg",
	}
}

type CardDef2521012 struct{}

func (CardDef2521012) ID() string      { return "2521012" }
func (CardDef2521012) Name() string    { return "幻彩颜料" }
func (CardDef2521012) Kind() string    { return "道具" }
func (CardDef2521012) Element() string { return "光" }

func (CardDef2521012) Card() model.Card {
	return model.Card{
		Number:          "2521012",
		Type:            "道具",
		Name:            "幻彩颜料",
		Category:        "光",
		Tag:             "消耗品-药剂",
		Description:     "将你场上负载的最多4点\\光变为\\无",
		Quote:           "生命理应更加绚烂多彩!",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521012.jpg",
	}
}

type CardDef2521013 struct{}

func (CardDef2521013) ID() string      { return "2521013" }
func (CardDef2521013) Name() string    { return "防护结界卷轴" }
func (CardDef2521013) Kind() string    { return "道具" }
func (CardDef2521013) Element() string { return "光" }

func (CardDef2521013) Card() model.Card {
	return model.Card{
		Number:          "2521013",
		Type:            "道具",
		Name:            "防护结界卷轴",
		Category:        "光",
		Tag:             "消耗品-法术卷轴-聚能",
		Description:     "防御",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           8,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521013.jpg",
	}
}

type CardDef2521014 struct{}

func (CardDef2521014) ID() string      { return "2521014" }
func (CardDef2521014) Name() string    { return "祝福之杖" }
func (CardDef2521014) Kind() string    { return "道具" }
func (CardDef2521014) Element() string { return "光" }

func (CardDef2521014) Card() model.Card {
	return model.Card{
		Number:          "2521014",
		Type:            "道具",
		Name:            "祝福之杖",
		Category:        "光",
		Tag:             "装备-武器",
		Description:     "入场:放置3个标记物.主动回合技:消耗此卡并取除1个标记物才能发动,使1个友方单位+1\\血,然后你获得2点\\光元素",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\光\\2521014.jpg",
	}
}

type CardDef2521101 struct{}

func (CardDef2521101) ID() string      { return "2521101" }
func (CardDef2521101) Name() string    { return "赐福之孤星" }
func (CardDef2521101) Kind() string    { return "道具" }
func (CardDef2521101) Element() string { return "光" }

func (CardDef2521101) Card() model.Card {
	return model.Card{
		Number:          "2521101",
		Type:            "道具",
		Name:            "赐福之孤星",
		Category:        "光",
		Tag:             "消耗品",
		Description:     "使1个友方伙伴获得负载+1\\光和+1\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521101.jpg",
	}
}

type CardDef2521102 struct{}

func (CardDef2521102) ID() string      { return "2521102" }
func (CardDef2521102) Name() string    { return "月霞之尘" }
func (CardDef2521102) Kind() string    { return "道具" }
func (CardDef2521102) Element() string { return "光" }

func (CardDef2521102) Card() model.Card {
	return model.Card{
		Number:          "2521102",
		Type:            "道具",
		Name:            "月霞之尘",
		Category:        "光",
		Tag:             "消耗品-药剂",
		Description:     "摧毁敌方盖放的所有卡牌,或者使前排敌人失去隐蔽",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521102.jpg",
	}
}

type CardDef2521103 struct{}

func (CardDef2521103) ID() string      { return "2521103" }
func (CardDef2521103) Name() string    { return "红玛瑙圣杯" }
func (CardDef2521103) Kind() string    { return "道具" }
func (CardDef2521103) Element() string { return "光" }

func (CardDef2521103) Card() model.Card {
	return model.Card{
		Number:          "2521103",
		Type:            "道具",
		Name:            "红玛瑙圣杯",
		Category:        "光",
		Tag:             "装备-神器",
		Description:     "光环:如果你的场上集齐绿玉权杖,蓝晶灯盏和红玛瑙圣杯,这些卡牌均获得负载+1\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 4},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521103.jpg",
	}
}

type CardDef2521104 struct{}

func (CardDef2521104) ID() string      { return "2521104" }
func (CardDef2521104) Name() string    { return "黄金龙骨" }
func (CardDef2521104) Kind() string    { return "道具" }
func (CardDef2521104) Element() string { return "光" }

func (CardDef2521104) Card() model.Card {
	return model.Card{
		Number:          "2521104",
		Type:            "道具",
		Name:            "黄金龙骨",
		Category:        "光",
		Tag:             "装备",
		Description:     "主动:献祭此卡才能发动,抽2张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{"光": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521104.jpg",
	}
}

type CardDef2521105 struct{}

func (CardDef2521105) ID() string      { return "2521105" }
func (CardDef2521105) Name() string    { return "孤星守护者" }
func (CardDef2521105) Kind() string    { return "道具" }
func (CardDef2521105) Element() string { return "光" }

func (CardDef2521105) Card() model.Card {
	return model.Card{
		Number:          "2521105",
		Type:            "道具",
		Name:            "孤星守护者",
		Category:        "光",
		Tag:             "装备-防具",
		Description:     "入场:获得护盾2",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3, "无": 1},
		ElementsGain:    map[string]int{"光": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521105.jpg",
	}
}

type CardDef2521106 struct{}

func (CardDef2521106) ID() string      { return "2521106" }
func (CardDef2521106) Name() string    { return "沐光卷轴" }
func (CardDef2521106) Kind() string    { return "道具" }
func (CardDef2521106) Element() string { return "光" }

func (CardDef2521106) Card() model.Card {
	return model.Card{
		Number:          "2521106",
		Type:            "道具",
		Name:            "沐光卷轴",
		Category:        "光",
		Tag:             "消耗品-卷轴",
		Description:     "为所有友方单位回复2\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521106.jpg",
	}
}

type CardDef2521107 struct{}

func (CardDef2521107) ID() string      { return "2521107" }
func (CardDef2521107) Name() string    { return "百灵药P型" }
func (CardDef2521107) Kind() string    { return "道具" }
func (CardDef2521107) Element() string { return "光" }

func (CardDef2521107) Card() model.Card {
	return model.Card{
		Number:          "2521107",
		Type:            "道具",
		Name:            "百灵药P型",
		Category:        "光",
		Tag:             "消耗品-药剂",
		Description:     "造成1点伤害,回复1\\血,抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521107.jpg",
	}
}

type CardDef2521108 struct{}

func (CardDef2521108) ID() string      { return "2521108" }
func (CardDef2521108) Name() string    { return "议庭审判锤" }
func (CardDef2521108) Kind() string    { return "道具" }
func (CardDef2521108) Element() string { return "光" }

func (CardDef2521108) Card() model.Card {
	return model.Card{
		Number:          "2521108",
		Type:            "道具",
		Name:            "议庭审判锤",
		Category:        "光",
		Tag:             "装备-武器",
		Description:     "诱发回合技:当敌方使用1个法术攻击,将3张九霄印记洗入对方卡组",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2, "气": 1},
		ElementsGain:    map[string]int{"光": 1, "气": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2001102"},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521108.jpg",
	}
}

type CardDef2521109 struct{}

func (CardDef2521109) ID() string      { return "2521109" }
func (CardDef2521109) Name() string    { return "惩戒符文" }
func (CardDef2521109) Kind() string    { return "道具" }
func (CardDef2521109) Element() string { return "光" }

func (CardDef2521109) Card() model.Card {
	return model.Card{
		Number:          "2521109",
		Type:            "道具",
		Name:            "惩戒符文",
		Category:        "光",
		Tag:             "消耗品-符文",
		Description:     "反制:当对方在一回合内进行超过2次法术攻击后,对任意一个敌方伙伴造成2点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521109.jpg",
	}
}

type CardDef2521110 struct{}

func (CardDef2521110) ID() string      { return "2521110" }
func (CardDef2521110) Name() string    { return "天使之祈祷" }
func (CardDef2521110) Kind() string    { return "道具" }
func (CardDef2521110) Element() string { return "光" }

func (CardDef2521110) Card() model.Card {
	return model.Card{
		Number:          "2521110",
		Type:            "道具",
		Name:            "天使之祈祷",
		Category:        "光",
		Tag:             "消耗品-卷轴",
		Description:     "翻取1个光辉属性的精灵.反制:对方使用法术攻击后使用此卡则无需花费并使该精灵入场花费-1\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521110.jpg",
	}
}

type CardDef2521111 struct{}

func (CardDef2521111) ID() string      { return "2521111" }
func (CardDef2521111) Name() string    { return "神谕卷轴 荣耀" }
func (CardDef2521111) Kind() string    { return "道具" }
func (CardDef2521111) Element() string { return "光" }

func (CardDef2521111) Card() model.Card {
	return model.Card{
		Number:          "2521111",
		Type:            "道具",
		Name:            "神谕卷轴 荣耀",
		Category:        "光",
		Tag:             "消耗品-法术卷轴-神秘",
		Description:     "穿透.使用时必须选择你的1个生命和负载总和大于5的伙伴,此卡威力上升该数值",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          3,
		Life:            -1,
		Duration:        -1,
		Power:           0,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521111.jpg",
	}
}

type CardDef2521112 struct{}

func (CardDef2521112) ID() string      { return "2521112" }
func (CardDef2521112) Name() string    { return "神谕卷轴 团结" }
func (CardDef2521112) Kind() string    { return "道具" }
func (CardDef2521112) Element() string { return "光" }

func (CardDef2521112) Card() model.Card {
	return model.Card{
		Number:          "2521112",
		Type:            "道具",
		Name:            "神谕卷轴 团结",
		Category:        "光",
		Tag:             "消耗品-法术卷轴-神秘",
		Description:     "范围:方阵.你场上每有1个光辉单位此卡花费-1\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 9},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\光\\2521112.jpg",
	}
}

type CardDef2601001 struct{}

func (CardDef2601001) ID() string      { return "2601001" }
func (CardDef2601001) Name() string    { return "幻痛" }
func (CardDef2601001) Kind() string    { return "道具" }
func (CardDef2601001) Element() string { return "暗" }

func (CardDef2601001) Card() model.Card {
	return model.Card{
		Number:          "2601001",
		Type:            "道具",
		Name:            "幻痛",
		Category:        "暗",
		Tag:             "衍生-装备-神器",
		Description:     "诱发回合技:当敌方使用法术防御成功后,使用于防御和强化防御的法术虚弱2",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2601001.jpg",
	}
}

type CardDef2601002 struct{}

func (CardDef2601002) ID() string      { return "2601002" }
func (CardDef2601002) Name() string    { return "咒言书" }
func (CardDef2601002) Kind() string    { return "道具" }
func (CardDef2601002) Element() string { return "暗" }

func (CardDef2601002) Card() model.Card {
	return model.Card{
		Number:          "2601002",
		Type:            "道具",
		Name:            "咒言书",
		Category:        "暗",
		Tag:             "衍生-装备",
		Description:     "入场:使所有敌方法术虚弱1",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2601002.jpg",
	}
}

type CardDef2611001 struct{}

func (CardDef2611001) ID() string      { return "2611001" }
func (CardDef2611001) Name() string    { return "死灵魔石 虚无" }
func (CardDef2611001) Kind() string    { return "道具" }
func (CardDef2611001) Element() string { return "暗" }

func (CardDef2611001) Card() model.Card {
	return model.Card{
		Number:          "2611001",
		Type:            "道具",
		Name:            "死灵魔石 虚无",
		Category:        "暗",
		Tag:             "传奇-装备-神器",
		Description:     "诱发回合技:当1个友方伙伴死亡后,此卡获得负载+1\\暗",
		Quote:           "比起被永远囚禁在这,魂飞魄散显得那么安详",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2611001.jpg",
	}
}

type CardDef2611002 struct{}

func (CardDef2611002) ID() string      { return "2611002" }
func (CardDef2611002) Name() string    { return "与恶魔的契约书" }
func (CardDef2611002) Kind() string    { return "道具" }
func (CardDef2611002) Element() string { return "暗" }

func (CardDef2611002) Card() model.Card {
	return model.Card{
		Number:          "2611002",
		Type:            "道具",
		Name:            "与恶魔的契约书",
		Category:        "暗",
		Tag:             "传奇-消耗品-卷轴",
		Description:     "献祭1个友方单位然后消灭法力范围内1个敌方伙伴,二者每相差1点\\血必须额外支付2\\暗.此卡在打出后洗回卡组",
		Quote:           "放心用吧,毕竟现在被献祭的还不是你",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2611002.jpg",
	}
}

type CardDef2611101 struct{}

func (CardDef2611101) ID() string      { return "2611101" }
func (CardDef2611101) Name() string    { return "厄瑞波斯的魂链" }
func (CardDef2611101) Kind() string    { return "道具" }
func (CardDef2611101) Element() string { return "暗" }

func (CardDef2611101) Card() model.Card {
	return model.Card{
		Number:          "2611101",
		Type:            "道具",
		Name:            "厄瑞波斯的魂链",
		Category:        "暗",
		Tag:             "传奇-装备-神器",
		Description:     "诱发绝技:当敌方透支伙伴来使用法术时,标记那些伙伴和法术.诱发:每当那些标记的伙伴被消耗或透支,使那些标记的法术虚弱1",
		Quote:           "随着一次次挥舞,他的手臂愈发疲惫",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2611101.jpg",
	}
}

type CardDef2611102 struct{}

func (CardDef2611102) ID() string      { return "2611102" }
func (CardDef2611102) Name() string    { return "渡灵之烛" }
func (CardDef2611102) Kind() string    { return "道具" }
func (CardDef2611102) Element() string { return "暗" }

func (CardDef2611102) Card() model.Card {
	return model.Card{
		Number:          "2611102",
		Type:            "道具",
		Name:            "渡灵之烛",
		Category:        "暗",
		Tag:             "传奇-装备-神器",
		Description:     "光环:获得2个只能放置灵媒和神秘技能的槽位,你的其他种类的法术\\威减半(向上取整)",
		Quote:           "相信我,所能看见的少是一种幸福",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2611102.jpg",
	}
}

type CardDef2621001 struct{}

func (CardDef2621001) ID() string      { return "2621001" }
func (CardDef2621001) Name() string    { return "虚弱药剂" }
func (CardDef2621001) Kind() string    { return "道具" }
func (CardDef2621001) Element() string { return "暗" }

func (CardDef2621001) Card() model.Card {
	return model.Card{
		Number:          "2621001",
		Type:            "道具",
		Name:            "虚弱药剂",
		Category:        "暗",
		Tag:             "消耗品-药剂",
		Description:     "使敌方最多2个不同的法术虚弱2",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1, "暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621001.jpg",
	}
}

type CardDef2621002 struct{}

func (CardDef2621002) ID() string      { return "2621002" }
func (CardDef2621002) Name() string    { return "巫毒娃娃" }
func (CardDef2621002) Kind() string    { return "道具" }
func (CardDef2621002) Element() string { return "暗" }

func (CardDef2621002) Card() model.Card {
	return model.Card{
		Number:          "2621002",
		Type:            "道具",
		Name:            "巫毒娃娃",
		Category:        "暗",
		Tag:             "装备",
		Description:     "入场:在此卡上放置3个暗影标记物并选择法力范围内的2个伙伴,其一受到伤害时可以让另一者收到同等的伤害,并取除伤害数量的暗影标记物",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621002.jpg",
	}
}

type CardDef2621003 struct{}

func (CardDef2621003) ID() string      { return "2621003" }
func (CardDef2621003) Name() string    { return "杀戮本能" }
func (CardDef2621003) Kind() string    { return "道具" }
func (CardDef2621003) Element() string { return "暗" }

func (CardDef2621003) Card() model.Card {
	return model.Card{
		Number:          "2621003",
		Type:            "道具",
		Name:            "杀戮本能",
		Category:        "暗",
		Tag:             "消耗品-卷轴",
		Description:     "反制:当对手召唤1个伙伴时,对其造成2点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621003.jpg",
	}
}

type CardDef2621004 struct{}

func (CardDef2621004) ID() string      { return "2621004" }
func (CardDef2621004) Name() string    { return "暗影帷幕" }
func (CardDef2621004) Kind() string    { return "道具" }
func (CardDef2621004) Element() string { return "暗" }

func (CardDef2621004) Card() model.Card {
	return model.Card{
		Number:          "2621004",
		Type:            "道具",
		Name:            "暗影帷幕",
		Category:        "暗",
		Tag:             "装备",
		Description:     "诱发:敌方法术命中时,献祭此卡才能发动,这个回合你的暗影伙伴不会受到法术伤害,但你的人物会获得引魔",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621004.jpg",
	}
}

type CardDef2621005 struct{}

func (CardDef2621005) ID() string      { return "2621005" }
func (CardDef2621005) Name() string    { return "献祭符文" }
func (CardDef2621005) Kind() string    { return "道具" }
func (CardDef2621005) Element() string { return "暗" }

func (CardDef2621005) Card() model.Card {
	return model.Card{
		Number:          "2621005",
		Type:            "道具",
		Name:            "献祭符文",
		Category:        "暗",
		Tag:             "消耗品-符文",
		Description:     "反制:当1个伙伴死亡时,抽2张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621005.jpg",
	}
}

type CardDef2621006 struct{}

func (CardDef2621006) ID() string      { return "2621006" }
func (CardDef2621006) Name() string    { return "亡魂项链" }
func (CardDef2621006) Kind() string    { return "道具" }
func (CardDef2621006) Element() string { return "暗" }

func (CardDef2621006) Card() model.Card {
	return model.Card{
		Number:          "2621006",
		Type:            "道具",
		Name:            "亡魂项链",
		Category:        "暗",
		Tag:             "装备-饰物",
		Description:     "诱发回合技:当你的1个伙伴死亡时,获得1\\暗",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621006.jpg",
	}
}

type CardDef2621007 struct{}

func (CardDef2621007) ID() string      { return "2621007" }
func (CardDef2621007) Name() string    { return "安迪斯之镰" }
func (CardDef2621007) Kind() string    { return "道具" }
func (CardDef2621007) Element() string { return "暗" }

func (CardDef2621007) Card() model.Card {
	return model.Card{
		Number:          "2621007",
		Type:            "道具",
		Name:            "安迪斯之镰",
		Category:        "暗",
		Tag:             "装备-武器",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 5},
		ElementsGain:    map[string]int{"暗": 3},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621007.jpg",
	}
}

type CardDef2621008 struct{}

func (CardDef2621008) ID() string      { return "2621008" }
func (CardDef2621008) Name() string    { return "魂噬卷轴" }
func (CardDef2621008) Kind() string    { return "道具" }
func (CardDef2621008) Element() string { return "暗" }

func (CardDef2621008) Card() model.Card {
	return model.Card{
		Number:          "2621008",
		Type:            "道具",
		Name:            "魂噬卷轴",
		Category:        "暗",
		Tag:             "消耗品-法术卷轴-灵媒",
		Description:     "命中:将3层虚弱分配给敌方法术",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621008.jpg",
	}
}

type CardDef2621009 struct{}

func (CardDef2621009) ID() string      { return "2621009" }
func (CardDef2621009) Name() string    { return "暗冥弹卷轴" }
func (CardDef2621009) Kind() string    { return "道具" }
func (CardDef2621009) Element() string { return "暗" }

func (CardDef2621009) Card() model.Card {
	return model.Card{
		Number:          "2621009",
		Type:            "道具",
		Name:            "暗冥弹卷轴",
		Category:        "暗",
		Tag:             "消耗品-法术卷轴-聚能",
		Description:     "范围:溅射.",
		Quote:           "别人写好的就是比自己学的好用",
		ElementsCost:    map[string]int{"无": 1, "暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621009.jpg",
	}
}

type CardDef2621010 struct{}

func (CardDef2621010) ID() string      { return "2621010" }
func (CardDef2621010) Name() string    { return "拖入深渊" }
func (CardDef2621010) Kind() string    { return "道具" }
func (CardDef2621010) Element() string { return "暗" }

func (CardDef2621010) Card() model.Card {
	return model.Card{
		Number:          "2621010",
		Type:            "道具",
		Name:            "拖入深渊",
		Category:        "暗",
		Tag:             "消耗品-卷轴",
		Description:     "反制:当1个友方单位受到伤害且死亡后,对法力范围内的1个敌人造成等同于那个友方单位在本回合受到的全部伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621010.jpg",
	}
}

type CardDef2621011 struct{}

func (CardDef2621011) ID() string      { return "2621011" }
func (CardDef2621011) Name() string    { return "狂乱符文" }
func (CardDef2621011) Kind() string    { return "道具" }
func (CardDef2621011) Element() string { return "暗" }

func (CardDef2621011) Card() model.Card {
	return model.Card{
		Number:          "2621011",
		Type:            "道具",
		Name:            "狂乱符文",
		Category:        "暗",
		Tag:             "消耗品-符文",
		Description:     "反制:当1个具有攻击力的敌方伙伴消耗时,使那次消耗视为其对1个相邻单位的攻击",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621011.jpg",
	}
}

type CardDef2621012 struct{}

func (CardDef2621012) ID() string      { return "2621012" }
func (CardDef2621012) Name() string    { return "暗影披风" }
func (CardDef2621012) Kind() string    { return "道具" }
func (CardDef2621012) Element() string { return "暗" }

func (CardDef2621012) Card() model.Card {
	return model.Card{
		Number:          "2621012",
		Type:            "道具",
		Name:            "暗影披风",
		Category:        "暗",
		Tag:             "装备-防具",
		Description:     "入场:敌方下一次命中的敌方法术伤害为0",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621012.jpg",
	}
}

type CardDef2621013 struct{}

func (CardDef2621013) ID() string      { return "2621013" }
func (CardDef2621013) Name() string    { return "巫术指环" }
func (CardDef2621013) Kind() string    { return "道具" }
func (CardDef2621013) Element() string { return "暗" }

func (CardDef2621013) Card() model.Card {
	return model.Card{
		Number:          "2621013",
		Type:            "道具",
		Name:            "巫术指环",
		Category:        "暗",
		Tag:             "装备-饰物",
		Description:     "诱发回合技:当敌方1个法术受到虚弱时,使该虚弱层数+1",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621013.jpg",
	}
}

type CardDef2621014 struct{}

func (CardDef2621014) ID() string      { return "2621014" }
func (CardDef2621014) Name() string    { return "埋葬者" }
func (CardDef2621014) Kind() string    { return "道具" }
func (CardDef2621014) Element() string { return "暗" }

func (CardDef2621014) Card() model.Card {
	return model.Card{
		Number:          "2621014",
		Type:            "道具",
		Name:            "埋葬者",
		Category:        "暗",
		Tag:             "装备-武器",
		Description:     "入场:放置3个标记物.主动回合技:消耗此卡并取除1个标记物才能发动,从卡组上方将2张牌送去弃牌堆,然后获得2点\\暗",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\道具\\暗\\2621014.jpg",
	}
}

type CardDef2621101 struct{}

func (CardDef2621101) ID() string      { return "2621101" }
func (CardDef2621101) Name() string    { return "黑松木魔杖" }
func (CardDef2621101) Kind() string    { return "道具" }
func (CardDef2621101) Element() string { return "暗" }

func (CardDef2621101) Card() model.Card {
	return model.Card{
		Number:          "2621101",
		Type:            "道具",
		Name:            "黑松木魔杖",
		Category:        "暗",
		Tag:             "装备-武器",
		Description:     "光环:你以友方单位为目标释放法术的花费-1",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621101.jpg",
	}
}

type CardDef2621102 struct{}

func (CardDef2621102) ID() string      { return "2621102" }
func (CardDef2621102) Name() string    { return "血蔷薇诅咒" }
func (CardDef2621102) Kind() string    { return "道具" }
func (CardDef2621102) Element() string { return "暗" }

func (CardDef2621102) Card() model.Card {
	return model.Card{
		Number:          "2621102",
		Type:            "道具",
		Name:            "血蔷薇诅咒",
		Category:        "暗",
		Tag:             "消耗品-卷轴",
		Description:     "将敌方的一个法术变为敌方的一个伙伴(由敌方选择)的绑定技能,在其死亡后将该法术从游戏中移除",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621102.jpg",
	}
}

type CardDef2621103 struct{}

func (CardDef2621103) ID() string      { return "2621103" }
func (CardDef2621103) Name() string    { return "血蛊" }
func (CardDef2621103) Kind() string    { return "道具" }
func (CardDef2621103) Element() string { return "暗" }

func (CardDef2621103) Card() model.Card {
	return model.Card{
		Number:          "2621103",
		Type:            "道具",
		Name:            "血蛊",
		Category:        "暗",
		Tag:             "装备",
		Description:     "诱发:每当你的人物受到1点伤害,在此卡上放置1个标记物,最多6个.主动:献祭此卡才能发动,此卡每有2个标记物,本回合你的法术+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621103.jpg",
	}
}

type CardDef2621104 struct{}

func (CardDef2621104) ID() string      { return "2621104" }
func (CardDef2621104) Name() string    { return "献身契约" }
func (CardDef2621104) Kind() string    { return "道具" }
func (CardDef2621104) Element() string { return "暗" }

func (CardDef2621104) Card() model.Card {
	return model.Card{
		Number:          "2621104",
		Type:            "道具",
		Name:            "献身契约",
		Category:        "暗",
		Tag:             "装备",
		Description:     "诱发回合技:在你使用1个代赎法术后,必须对你的人物造成1点伤害并抽1张牌",
		Quote:           "很好,恶魔已经看上你了",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621104.jpg",
	}
}

type CardDef2621105 struct{}

func (CardDef2621105) ID() string      { return "2621105" }
func (CardDef2621105) Name() string    { return "红月吊坠" }
func (CardDef2621105) Kind() string    { return "道具" }
func (CardDef2621105) Element() string { return "暗" }

func (CardDef2621105) Card() model.Card {
	return model.Card{
		Number:          "2621105",
		Type:            "道具",
		Name:            "红月吊坠",
		Category:        "暗",
		Tag:             "装备-饰物",
		Description:     "主动:献祭此卡才能发动,使你的下一次红月持续时间+1",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621105.jpg",
	}
}

type CardDef2621106 struct{}

func (CardDef2621106) ID() string      { return "2621106" }
func (CardDef2621106) Name() string    { return "苦痛尖啸卷轴" }
func (CardDef2621106) Kind() string    { return "道具" }
func (CardDef2621106) Element() string { return "暗" }

func (CardDef2621106) Card() model.Card {
	return model.Card{
		Number:          "2621106",
		Type:            "道具",
		Name:            "苦痛尖啸卷轴",
		Category:        "暗",
		Tag:             "消耗品-卷轴",
		Description:     "本回合友方单位每受到1点伤害就选择1个没有虚弱的敌方法术,使那些法术虚弱2",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621106.jpg",
	}
}

type CardDef2621107 struct{}

func (CardDef2621107) ID() string      { return "2621107" }
func (CardDef2621107) Name() string    { return "诅咒魔盒" }
func (CardDef2621107) Kind() string    { return "道具" }
func (CardDef2621107) Element() string { return "暗" }

func (CardDef2621107) Card() model.Card {
	return model.Card{
		Number:          "2621107",
		Type:            "道具",
		Name:            "诅咒魔盒",
		Category:        "暗",
		Tag:             "装备-神器",
		Description:     "诱发:每当1个单位死亡,在此卡上放置1个标记物.主动回合技:移除最多3个标记物,使那个数量的敌方法术虚弱1",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621107.jpg",
	}
}

type CardDef2621108 struct{}

func (CardDef2621108) ID() string      { return "2621108" }
func (CardDef2621108) Name() string    { return "黑松棺木" }
func (CardDef2621108) Kind() string    { return "道具" }
func (CardDef2621108) Element() string { return "暗" }

func (CardDef2621108) Card() model.Card {
	return model.Card{
		Number:          "2621108",
		Type:            "道具",
		Name:            "黑松棺木",
		Category:        "暗",
		Tag:             "装备",
		Description:     "入场:从手牌丢弃最多2张入场花费小于5的暗影伙伴,立刻结算它们的遗言效果",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{"暗": 1},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621108.jpg",
	}
}

type CardDef2621109 struct{}

func (CardDef2621109) ID() string      { return "2621109" }
func (CardDef2621109) Name() string    { return "哀歌卷轴" }
func (CardDef2621109) Kind() string    { return "道具" }
func (CardDef2621109) Element() string { return "暗" }

func (CardDef2621109) Card() model.Card {
	return model.Card{
		Number:          "2621109",
		Type:            "道具",
		Name:            "哀歌卷轴",
		Category:        "暗",
		Tag:             "消耗品-卷轴",
		Description:     "翻取1个具有遗言的暗影伙伴,如果你的弃牌堆已有暗影伙伴则使翻取的卡牌入场花费-1\\暗",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621109.jpg",
	}
}

type CardDef2621110 struct{}

func (CardDef2621110) ID() string      { return "2621110" }
func (CardDef2621110) Name() string    { return "安迪斯的赠与" }
func (CardDef2621110) Kind() string    { return "道具" }
func (CardDef2621110) Element() string { return "暗" }

func (CardDef2621110) Card() model.Card {
	return model.Card{
		Number:          "2621110",
		Type:            "道具",
		Name:            "安迪斯的赠与",
		Category:        "暗",
		Tag:             "消耗品-药剂",
		Description:     "使1个友方单位获得负载+2\\暗,但会在回合结束时死亡",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621110.jpg",
	}
}

type CardDef2621111 struct{}

func (CardDef2621111) ID() string      { return "2621111" }
func (CardDef2621111) Name() string    { return "暗黑爆发卷轴" }
func (CardDef2621111) Kind() string    { return "道具" }
func (CardDef2621111) Element() string { return "暗" }

func (CardDef2621111) Card() model.Card {
	return model.Card{
		Number:          "2621111",
		Type:            "道具",
		Name:            "暗黑爆发卷轴",
		Category:        "暗",
		Tag:             "消耗品-卷轴",
		Description:     "你的弃牌堆有5个及以上暗影伙伴时才能使用,将那些暗影伙伴全部移出游戏,每1个获得2\\暗",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621111.jpg",
	}
}

type CardDef2621112 struct{}

func (CardDef2621112) ID() string      { return "2621112" }
func (CardDef2621112) Name() string    { return "灵魂法杖" }
func (CardDef2621112) Kind() string    { return "道具" }
func (CardDef2621112) Element() string { return "暗" }

func (CardDef2621112) Card() model.Card {
	return model.Card{
		Number:          "2621112",
		Type:            "道具",
		Name:            "灵魂法杖",
		Category:        "暗",
		Tag:             "装备-武器",
		Description:     "主动回合技:从你的弃牌堆将2张暗影伙伴移出游戏,给你的1个暗影法术放置1个灵魂标记物.每个灵魂标记物会使法术+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 5},
		ElementsGain:    map[string]int{"暗": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\道具\\暗\\2621112.jpg",
	}
}

type CardDef3001001 struct{}

func (CardDef3001001) ID() string      { return "3001001" }
func (CardDef3001001) Name() string    { return "破灭魔光" }
func (CardDef3001001) Kind() string    { return "技能" }
func (CardDef3001001) Element() string { return "无" }

func (CardDef3001001) Card() model.Card {
	return model.Card{
		Number:          "3001001",
		Type:            "技能",
		Name:            "破灭魔光",
		Category:        "无",
		Tag:             "衍生-法术-聚能",
		Description:     "范围:前排.无法强化或被强化",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1, "无": 1, "气": 1, "水": 1, "火": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           10,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3001001.jpg",
	}
}

type CardDef3001002 struct{}

func (CardDef3001002) ID() string      { return "3001002" }
func (CardDef3001002) Name() string    { return "纯净奥术" }
func (CardDef3001002) Kind() string    { return "技能" }
func (CardDef3001002) Element() string { return "无" }

func (CardDef3001002) Card() model.Card {
	return model.Card{
		Number:          "3001002",
		Type:            "技能",
		Name:            "纯净奥术",
		Category:        "无",
		Tag:             "衍生-咒术-聚能",
		Description:     "花费最多10点同种元素,使下一次该属性法术威力上升那个数值",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3001002.jpg",
	}
}

type CardDef3001101 struct{}

func (CardDef3001101) ID() string      { return "3001101" }
func (CardDef3001101) Name() string    { return "入局" }
func (CardDef3001101) Kind() string    { return "技能" }
func (CardDef3001101) Element() string { return "无" }

func (CardDef3001101) Card() model.Card {
	return model.Card{
		Number:          "3001101",
		Type:            "技能",
		Name:            "入局",
		Category:        "无",
		Tag:             "衍生-咒术-驱动",
		Description:     "速攻.为任意玩家召唤1个弃子(由你决定位置)",
		Quote:           "棋至困局,破者得道,乱者失机",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\无\\3001101.jpg",
	}
}

type CardDef3011101 struct{}

func (CardDef3011101) ID() string      { return "3011101" }
func (CardDef3011101) Name() string    { return "绝对纯净 奥能一心" }
func (CardDef3011101) Kind() string    { return "技能" }
func (CardDef3011101) Element() string { return "无" }

func (CardDef3011101) Card() model.Card {
	return model.Card{
		Number:          "3011101",
		Type:            "技能",
		Name:            "绝对纯净 奥能一心",
		Category:        "无",
		Tag:             "传奇-法术-聚能",
		Description:     "冷却2.范围:全场.此卡的学习和使用花费必须严格为奥术元素.使用时翻开你卡组最上方的卡牌直到出现非奥术卡牌,此卡\\威上升奥术卡牌的数量,然后洗混卡组",
		Quote:           "\"罗慕路斯,一个只会奥术的法师是走不远的\"",
		ElementsCost:    map[string]int{"无": 11},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 7},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           0,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\无\\3011101.jpg",
	}
}

type CardDef3021001 struct{}

func (CardDef3021001) ID() string      { return "3021001" }
func (CardDef3021001) Name() string    { return "移形换影" }
func (CardDef3021001) Kind() string    { return "技能" }
func (CardDef3021001) Element() string { return "无" }

func (CardDef3021001) Card() model.Card {
	return model.Card{
		Number:          "3021001",
		Type:            "技能",
		Name:            "移形换影",
		Category:        "无",
		Tag:             "咒术-幻变",
		Description:     "速攻.移动1个友方单位",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021001.jpg",
	}
}

type CardDef3021002 struct{}

func (CardDef3021002) ID() string      { return "3021002" }
func (CardDef3021002) Name() string    { return "预见" }
func (CardDef3021002) Kind() string    { return "技能" }
func (CardDef3021002) Element() string { return "无" }

func (CardDef3021002) Card() model.Card {
	return model.Card{
		Number:          "3021002",
		Type:            "技能",
		Name:            "预见",
		Category:        "无",
		Tag:             "咒术-灵媒",
		Description:     "查看牌堆顶3张牌,将其置于牌堆顶或牌堆底",
		Quote:           "这是否也是命运的一部分?",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021002.jpg",
	}
}

type CardDef3021003 struct{}

func (CardDef3021003) ID() string      { return "3021003" }
func (CardDef3021003) Name() string    { return "冥想" }
func (CardDef3021003) Kind() string    { return "技能" }
func (CardDef3021003) Element() string { return "无" }

func (CardDef3021003) Card() model.Card {
	return model.Card{
		Number:          "3021003",
		Type:            "技能",
		Name:            "冥想",
		Category:        "无",
		Tag:             "咒术-灵媒",
		Description:     "获得1\\无",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021003.jpg",
	}
}

type CardDef3021004 struct{}

func (CardDef3021004) ID() string      { return "3021004" }
func (CardDef3021004) Name() string    { return "刻印" }
func (CardDef3021004) Kind() string    { return "技能" }
func (CardDef3021004) Element() string { return "无" }

func (CardDef3021004) Card() model.Card {
	return model.Card{
		Number:          "3021004",
		Type:            "技能",
		Name:            "刻印",
		Category:        "无",
		Tag:             "咒术-代赎",
		Description:     "冷却2.丢弃1张手牌才能发动,从卡组检索1张卷轴或符文",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021004.jpg",
	}
}

type CardDef3021005 struct{}

func (CardDef3021005) ID() string      { return "3021005" }
func (CardDef3021005) Name() string    { return "奥术箭矢" }
func (CardDef3021005) Kind() string    { return "技能" }
func (CardDef3021005) Element() string { return "无" }

func (CardDef3021005) Card() model.Card {
	return model.Card{
		Number:          "3021005",
		Type:            "技能",
		Name:            "奥术箭矢",
		Category:        "无",
		Tag:             "咒术-创造",
		Description:     "对法力范围内1个单位造成1点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021005.jpg",
	}
}

type CardDef3021006 struct{}

func (CardDef3021006) ID() string      { return "3021006" }
func (CardDef3021006) Name() string    { return "洞察之眼" }
func (CardDef3021006) Kind() string    { return "技能" }
func (CardDef3021006) Element() string { return "无" }

func (CardDef3021006) Card() model.Card {
	return model.Card{
		Number:          "3021006",
		Type:            "技能",
		Name:            "洞察之眼",
		Category:        "无",
		Tag:             "咒术-灵媒",
		Description:     "速攻.冷却1.摧毁1张敌方盖放的卡牌",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021006.jpg",
	}
}

type CardDef3021007 struct{}

func (CardDef3021007) ID() string      { return "3021007" }
func (CardDef3021007) Name() string    { return "元素附魔" }
func (CardDef3021007) Kind() string    { return "技能" }
func (CardDef3021007) Element() string { return "无" }

func (CardDef3021007) Card() model.Card {
	return model.Card{
		Number:          "3021007",
		Type:            "技能",
		Name:            "元素附魔",
		Category:        "无",
		Tag:             "咒术-幻变",
		Description:     "使你的下一次法术获得1点任意负面效果(点燃、冻结、晕眩、石化、虚弱)",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021007.jpg",
	}
}

type CardDef3021008 struct{}

func (CardDef3021008) ID() string      { return "3021008" }
func (CardDef3021008) Name() string    { return "缴械" }
func (CardDef3021008) Kind() string    { return "技能" }
func (CardDef3021008) Element() string { return "无" }

func (CardDef3021008) Card() model.Card {
	return model.Card{
		Number:          "3021008",
		Type:            "技能",
		Name:            "缴械",
		Category:        "无",
		Tag:             "法术-驱动",
		Description:     "速攻.冷却1.命中:摧毁目标控制者的1个装备",
		Quote:           "除你武器!",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021008.jpg",
	}
}

type CardDef3021009 struct{}

func (CardDef3021009) ID() string      { return "3021009" }
func (CardDef3021009) Name() string    { return "昏睡" }
func (CardDef3021009) Kind() string    { return "技能" }
func (CardDef3021009) Element() string { return "无" }

func (CardDef3021009) Card() model.Card {
	return model.Card{
		Number:          "3021009",
		Type:            "技能",
		Name:            "昏睡",
		Category:        "无",
		Tag:             "法术-幻变",
		Description:     "速攻.命中:使目标伙伴晕眩1",
		Quote:           "如何卷赢室友?使用昏睡咒",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021009.jpg",
	}
}

type CardDef3021010 struct{}

func (CardDef3021010) ID() string      { return "3021010" }
func (CardDef3021010) Name() string    { return "解咒" }
func (CardDef3021010) Kind() string    { return "技能" }
func (CardDef3021010) Element() string { return "无" }

func (CardDef3021010) Card() model.Card {
	return model.Card{
		Number:          "3021010",
		Type:            "技能",
		Name:            "解咒",
		Category:        "无",
		Tag:             "咒术-幻变",
		Description:     "冷却1.诱发:当敌方使用防御型法术时才能使用此卡,将那个敌方法术无效",
		Quote:           "我破防了!",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021010.jpg",
	}
}

type CardDef3021011 struct{}

func (CardDef3021011) ID() string      { return "3021011" }
func (CardDef3021011) Name() string    { return "统御者的制裁" }
func (CardDef3021011) Kind() string    { return "技能" }
func (CardDef3021011) Element() string { return "无" }

func (CardDef3021011) Card() model.Card {
	return model.Card{
		Number:          "3021011",
		Type:            "技能",
		Name:            "统御者的制裁",
		Category:        "无",
		Tag:             "法术-神秘",
		Description:     "穿透.此卡的学习和使用花费必须为同种元素",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 9},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 4},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          3,
		Life:            -1,
		Duration:        -1,
		Power:           9,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021011.jpg",
	}
}

type CardDef3021012 struct{}

func (CardDef3021012) ID() string      { return "3021012" }
func (CardDef3021012) Name() string    { return "心炼" }
func (CardDef3021012) Kind() string    { return "技能" }
func (CardDef3021012) Element() string { return "无" }

func (CardDef3021012) Card() model.Card {
	return model.Card{
		Number:          "3021012",
		Type:            "技能",
		Name:            "心炼",
		Category:        "无",
		Tag:             "咒术-神秘",
		Description:     "冷却1.使你的1个法术永久获得+3\\威或者+1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\无\\3021012.jpg",
	}
}

type CardDef3021101 struct{}

func (CardDef3021101) ID() string      { return "3021101" }
func (CardDef3021101) Name() string    { return "奥术冲击" }
func (CardDef3021101) Kind() string    { return "技能" }
func (CardDef3021101) Element() string { return "无" }

func (CardDef3021101) Card() model.Card {
	return model.Card{
		Number:          "3021101",
		Type:            "技能",
		Name:            "奥术冲击",
		Category:        "无",
		Tag:             "法术-聚能",
		Description:     "如果此卡学习和使用花费的元素均为奥术元素,获得+1\\攻+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\无\\3021101.jpg",
	}
}

type CardDef3021102 struct{}

func (CardDef3021102) ID() string      { return "3021102" }
func (CardDef3021102) Name() string    { return "奥术屏障" }
func (CardDef3021102) Kind() string    { return "技能" }
func (CardDef3021102) Element() string { return "无" }

func (CardDef3021102) Card() model.Card {
	return model.Card{
		Number:          "3021102",
		Type:            "技能",
		Name:            "奥术屏障",
		Category:        "无",
		Tag:             "法术-创造",
		Description:     "防御",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\无\\3021102.jpg",
	}
}

type CardDef3021103 struct{}

func (CardDef3021103) ID() string      { return "3021103" }
func (CardDef3021103) Name() string    { return "奥能汲取" }
func (CardDef3021103) Kind() string    { return "技能" }
func (CardDef3021103) Element() string { return "无" }

func (CardDef3021103) Card() model.Card {
	return model.Card{
		Number:          "3021103",
		Type:            "技能",
		Name:            "奥能汲取",
		Category:        "无",
		Tag:             "咒术-神秘",
		Description:     "冷却1.抽2张牌.此卡使用花费的元素属性必须各不相同",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\无\\3021103.jpg",
	}
}

type CardDef3021104 struct{}

func (CardDef3021104) ID() string      { return "3021104" }
func (CardDef3021104) Name() string    { return "七神加护" }
func (CardDef3021104) Kind() string    { return "技能" }
func (CardDef3021104) Element() string { return "无" }

func (CardDef3021104) Card() model.Card {
	return model.Card{
		Number:          "3021104",
		Type:            "技能",
		Name:            "七神加护",
		Category:        "无",
		Tag:             "咒术-神秘",
		Description:     "异能:如果你的所有技能属性各不相同,使它们使用花费-1,若为法术再+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        2,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\无\\3021104.jpg",
	}
}

type CardDef3021105 struct{}

func (CardDef3021105) ID() string      { return "3021105" }
func (CardDef3021105) Name() string    { return "奥能净化" }
func (CardDef3021105) Kind() string    { return "技能" }
func (CardDef3021105) Element() string { return "无" }

func (CardDef3021105) Card() model.Card {
	return model.Card{
		Number:          "3021105",
		Type:            "技能",
		Name:            "奥能净化",
		Category:        "无",
		Tag:             "咒术-幻变",
		Description:     "速攻.冷却2.本回合友方卡牌不受负面状态影响(仍可处于负面状态)",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\无\\3021105.jpg",
	}
}

type CardDef3021106 struct{}

func (CardDef3021106) ID() string      { return "3021106" }
func (CardDef3021106) Name() string    { return "奥能流贯" }
func (CardDef3021106) Kind() string    { return "技能" }
func (CardDef3021106) Element() string { return "无" }

func (CardDef3021106) Card() model.Card {
	return model.Card{
		Number:          "3021106",
		Type:            "技能",
		Name:            "奥能流贯",
		Category:        "无",
		Tag:             "法术-聚能",
		Description:     "光环:如果你场上只有奥术卡牌,此卡\\威翻倍",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\无\\3021106.jpg",
	}
}

type CardDef3021107 struct{}

func (CardDef3021107) ID() string      { return "3021107" }
func (CardDef3021107) Name() string    { return "奥能护盾" }
func (CardDef3021107) Kind() string    { return "技能" }
func (CardDef3021107) Element() string { return "无" }

func (CardDef3021107) Card() model.Card {
	return model.Card{
		Number:          "3021107",
		Type:            "技能",
		Name:            "奥能护盾",
		Category:        "无",
		Tag:             "咒术-创造",
		Description:     "速攻.在下个回合开始时获得护盾1",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\无\\3021107.jpg",
	}
}

type CardDef3021108 struct{}

func (CardDef3021108) ID() string      { return "3021108" }
func (CardDef3021108) Name() string    { return "奥术封印" }
func (CardDef3021108) Kind() string    { return "技能" }
func (CardDef3021108) Element() string { return "无" }

func (CardDef3021108) Card() model.Card {
	return model.Card{
		Number:          "3021108",
		Type:            "技能",
		Name:            "奥术封印",
		Category:        "无",
		Tag:             "咒术-幻变",
		Description:     "速攻.冷却1.选择对方场上1个技能,对方在他的下个回合不能使用该技能,然后此卡的花费以后永久+2\\无",
		Quote:           "",
		ElementsCost:    map[string]int{"无": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"无": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\无\\3021108.jpg",
	}
}

type CardDef3101001 struct{}

func (CardDef3101001) ID() string      { return "3101001" }
func (CardDef3101001) Name() string    { return "火焰吐息" }
func (CardDef3101001) Kind() string    { return "技能" }
func (CardDef3101001) Element() string { return "火" }

func (CardDef3101001) Card() model.Card {
	return model.Card{
		Number:          "3101001",
		Type:            "技能",
		Name:            "火焰吐息",
		Category:        "火",
		Tag:             "衍生-法术-聚能",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3101001.jpg",
	}
}

type CardDef3101002 struct{}

func (CardDef3101002) ID() string      { return "3101002" }
func (CardDef3101002) Name() string    { return "万火合一术" }
func (CardDef3101002) Kind() string    { return "技能" }
func (CardDef3101002) Element() string { return "火" }

func (CardDef3101002) Card() model.Card {
	return model.Card{
		Number:          "3101002",
		Type:            "技能",
		Name:            "万火合一术",
		Category:        "火",
		Tag:             "衍生-法术-聚能",
		Description:     "光环:此卡每有5点\\威获得+1\\攻,包括强化获得的\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3101002.jpg",
	}
}

type CardDef3111101 struct{}

func (CardDef3111101) ID() string      { return "3111101" }
func (CardDef3111101) Name() string    { return "火炎炼狱" }
func (CardDef3111101) Kind() string    { return "技能" }
func (CardDef3111101) Element() string { return "火" }

func (CardDef3111101) Card() model.Card {
	return model.Card{
		Number:          "3111101",
		Type:            "技能",
		Name:            "火炎炼狱",
		Category:        "火",
		Tag:             "传奇-法术-幻变",
		Description:     "冷却1.范围:全场.每有8\\威获得点燃1.此卡被防御成功时,每1层点燃需额外4\\威进行防御,否则仍会生效",
		Quote:           "狄斯托德从地狱归来,手中是炽热的百炼战刃,胯下是嘶吼的千里炎驹,身后是无尽的万葬火海",
		ElementsCost:    map[string]int{"地": 2, "火": 7},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 2, "火": 4},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           8,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3111101.jpg",
	}
}

type CardDef3111102 struct{}

func (CardDef3111102) ID() string      { return "3111102" }
func (CardDef3111102) Name() string    { return "原初神炎 洛普修斯" }
func (CardDef3111102) Kind() string    { return "技能" }
func (CardDef3111102) Element() string { return "火" }

func (CardDef3111102) Card() model.Card {
	return model.Card{
		Number:          "3111102",
		Type:            "技能",
		Name:            "原初神炎 洛普修斯",
		Category:        "火",
		Tag:             "传奇-法术-神秘",
		Description:     "光环:无法被强化,无法用于强化,不受其他卡牌效果影响,不受负面效果影响.主动回合技:将你学习的一个火焰技能移出游戏,此卡获得+1\\攻+2\\威",
		Quote:           "\"努尔已经被战争和毁灭蒙蔽双眼,而洛普修斯将引领我们向往更辉煌的终点!\"",
		ElementsCost:    map[string]int{"火": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3111102.jpg",
	}
}

type CardDef3121001 struct{}

func (CardDef3121001) ID() string      { return "3121001" }
func (CardDef3121001) Name() string    { return "火球术" }
func (CardDef3121001) Kind() string    { return "技能" }
func (CardDef3121001) Element() string { return "火" }

func (CardDef3121001) Card() model.Card {
	return model.Card{
		Number:          "3121001",
		Type:            "技能",
		Name:            "火球术",
		Category:        "火",
		Tag:             "法术-创造",
		Description:     "",
		Quote:           "巫师学院里有一个传说,你可以看到五颜六色的火焰,蓝色的、绿色的、黑色的,唯独没有红色的",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121001.jpg",
	}
}

type CardDef3121002 struct{}

func (CardDef3121002) ID() string      { return "3121002" }
func (CardDef3121002) Name() string    { return "焚烧" }
func (CardDef3121002) Kind() string    { return "技能" }
func (CardDef3121002) Element() string { return "火" }

func (CardDef3121002) Card() model.Card {
	return model.Card{
		Number:          "3121002",
		Type:            "技能",
		Name:            "焚烧",
		Category:        "火",
		Tag:             "法术-幻变",
		Description:     "攻击时目标具有点燃则+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121002.jpg",
	}
}

type CardDef3121003 struct{}

func (CardDef3121003) ID() string      { return "3121003" }
func (CardDef3121003) Name() string    { return "炽热射线" }
func (CardDef3121003) Kind() string    { return "技能" }
func (CardDef3121003) Element() string { return "火" }

func (CardDef3121003) Card() model.Card {
	return model.Card{
		Number:          "3121003",
		Type:            "技能",
		Name:            "炽热射线",
		Category:        "火",
		Tag:             "法术-聚能",
		Description:     "点燃2",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121003.jpg",
	}
}

type CardDef3121004 struct{}

func (CardDef3121004) ID() string      { return "3121004" }
func (CardDef3121004) Name() string    { return "燃烧大地" }
func (CardDef3121004) Kind() string    { return "技能" }
func (CardDef3121004) Element() string { return "火" }

func (CardDef3121004) Card() model.Card {
	return model.Card{
		Number:          "3121004",
		Type:            "技能",
		Name:            "燃烧大地",
		Category:        "火",
		Tag:             "法术-驱动",
		Description:     "范围:前排",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121004.jpg",
	}
}

type CardDef3121005 struct{}

func (CardDef3121005) ID() string      { return "3121005" }
func (CardDef3121005) Name() string    { return "烈焰风暴" }
func (CardDef3121005) Kind() string    { return "技能" }
func (CardDef3121005) Element() string { return "火" }

func (CardDef3121005) Card() model.Card {
	return model.Card{
		Number:          "3121005",
		Type:            "技能",
		Name:            "烈焰风暴",
		Category:        "火",
		Tag:             "法术-创造",
		Description:     "范围:方阵.点燃1",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1, "火": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121005.jpg",
	}
}

type CardDef3121006 struct{}

func (CardDef3121006) ID() string      { return "3121006" }
func (CardDef3121006) Name() string    { return "陨石术" }
func (CardDef3121006) Kind() string    { return "技能" }
func (CardDef3121006) Element() string { return "火" }

func (CardDef3121006) Card() model.Card {
	return model.Card{
		Number:          "3121006",
		Type:            "技能",
		Name:            "陨石术",
		Category:        "火",
		Tag:             "法术-驱动",
		Description:     "穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "火": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 4},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          3,
		Life:            -1,
		Duration:        -1,
		Power:           8,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121006.jpg",
	}
}

type CardDef3121007 struct{}

func (CardDef3121007) ID() string      { return "3121007" }
func (CardDef3121007) Name() string    { return "激情之火" }
func (CardDef3121007) Kind() string    { return "技能" }
func (CardDef3121007) Element() string { return "火" }

func (CardDef3121007) Card() model.Card {
	return model.Card{
		Number:          "3121007",
		Type:            "技能",
		Name:            "激情之火",
		Category:        "火",
		Tag:             "咒术-聚能",
		Description:     "速攻.冷却1.异能:你的火焰法术命中时抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121007.jpg",
	}
}

type CardDef3121008 struct{}

func (CardDef3121008) ID() string      { return "3121008" }
func (CardDef3121008) Name() string    { return "火焰结界" }
func (CardDef3121008) Kind() string    { return "技能" }
func (CardDef3121008) Element() string { return "火" }

func (CardDef3121008) Card() model.Card {
	return model.Card{
		Number:          "3121008",
		Type:            "技能",
		Name:            "火焰结界",
		Category:        "火",
		Tag:             "咒术-创造",
		Description:     "冷却1.异能:你的火焰法术获得点燃1和+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121008.jpg",
	}
}

type CardDef3121009 struct{}

func (CardDef3121009) ID() string      { return "3121009" }
func (CardDef3121009) Name() string    { return "爆焰一击" }
func (CardDef3121009) Kind() string    { return "技能" }
func (CardDef3121009) Element() string { return "火" }

func (CardDef3121009) Card() model.Card {
	return model.Card{
		Number:          "3121009",
		Type:            "技能",
		Name:            "爆焰一击",
		Category:        "火",
		Tag:             "法术-聚能",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121009.jpg",
	}
}

type CardDef3121010 struct{}

func (CardDef3121010) ID() string      { return "3121010" }
func (CardDef3121010) Name() string    { return "岩浆爆发" }
func (CardDef3121010) Kind() string    { return "技能" }
func (CardDef3121010) Element() string { return "火" }

func (CardDef3121010) Card() model.Card {
	return model.Card{
		Number:          "3121010",
		Type:            "技能",
		Name:            "岩浆爆发",
		Category:        "火",
		Tag:             "法术-驱动",
		Description:     "范围:方阵.穿透.点燃1",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "火": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1, "火": 4},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121010.jpg",
	}
}

type CardDef3121011 struct{}

func (CardDef3121011) ID() string      { return "3121011" }
func (CardDef3121011) Name() string    { return "引燃" }
func (CardDef3121011) Kind() string    { return "技能" }
func (CardDef3121011) Element() string { return "火" }

func (CardDef3121011) Card() model.Card {
	return model.Card{
		Number:          "3121011",
		Type:            "技能",
		Name:            "引燃",
		Category:        "火",
		Tag:             "咒术-幻变",
		Description:     "速攻.使1个敌方单位点燃1",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121011.jpg",
	}
}

type CardDef3121012 struct{}

func (CardDef3121012) ID() string      { return "3121012" }
func (CardDef3121012) Name() string    { return "烈焰护体" }
func (CardDef3121012) Kind() string    { return "技能" }
func (CardDef3121012) Element() string { return "火" }

func (CardDef3121012) Card() model.Card {
	return model.Card{
		Number:          "3121012",
		Type:            "技能",
		Name:            "烈焰护体",
		Category:        "火",
		Tag:             "法术-创造",
		Description:     "防御",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121012.jpg",
	}
}

type CardDef3121013 struct{}

func (CardDef3121013) ID() string      { return "3121013" }
func (CardDef3121013) Name() string    { return "烈焰反噬" }
func (CardDef3121013) Kind() string    { return "技能" }
func (CardDef3121013) Element() string { return "火" }

func (CardDef3121013) Card() model.Card {
	return model.Card{
		Number:          "3121013",
		Type:            "技能",
		Name:            "烈焰反噬",
		Category:        "火",
		Tag:             "法术-驱动",
		Description:     "防御.若防御成功,对敌方人物造成点燃1",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121013.jpg",
	}
}

type CardDef3121014 struct{}

func (CardDef3121014) ID() string      { return "3121014" }
func (CardDef3121014) Name() string    { return "烈焰重燃" }
func (CardDef3121014) Kind() string    { return "技能" }
func (CardDef3121014) Element() string { return "火" }

func (CardDef3121014) Card() model.Card {
	return model.Card{
		Number:          "3121014",
		Type:            "技能",
		Name:            "烈焰重燃",
		Category:        "火",
		Tag:             "咒术-聚能",
		Description:     "本回合你每使用过1个火焰法术就获得1\\火",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121014.jpg",
	}
}

type CardDef3121015 struct{}

func (CardDef3121015) ID() string      { return "3121015" }
func (CardDef3121015) Name() string    { return "焚风" }
func (CardDef3121015) Kind() string    { return "技能" }
func (CardDef3121015) Element() string { return "火" }

func (CardDef3121015) Card() model.Card {
	return model.Card{
		Number:          "3121015",
		Type:            "技能",
		Name:            "焚风",
		Category:        "火",
		Tag:             "法术-驱动",
		Description:     "穿透.强化其他法术时使其获得穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1, "火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 1, "火": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           2,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\火\\3121015.jpg",
	}
}

type CardDef3121101 struct{}

func (CardDef3121101) ID() string      { return "3121101" }
func (CardDef3121101) Name() string    { return "唤灵术 火蛇" }
func (CardDef3121101) Kind() string    { return "技能" }
func (CardDef3121101) Element() string { return "火" }

func (CardDef3121101) Card() model.Card {
	return model.Card{
		Number:          "3121101",
		Type:            "技能",
		Name:            "唤灵术 火蛇",
		Category:        "火",
		Tag:             "法术-创造",
		Description:     "若防御成功对法力范围内1个敌人造成1点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3121101.jpg",
	}
}

type CardDef3121102 struct{}

func (CardDef3121102) ID() string      { return "3121102" }
func (CardDef3121102) Name() string    { return "雄狮之守护" }
func (CardDef3121102) Kind() string    { return "技能" }
func (CardDef3121102) Element() string { return "火" }

func (CardDef3121102) Card() model.Card {
	return model.Card{
		Number:          "3121102",
		Type:            "技能",
		Name:            "雄狮之守护",
		Category:        "火",
		Tag:             "法术-神秘",
		Description:     "防御.若防御成功使你的其他火焰法术永久+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3121102.jpg",
	}
}

type CardDef3121103 struct{}

func (CardDef3121103) ID() string      { return "3121103" }
func (CardDef3121103) Name() string    { return "祈祷之焰" }
func (CardDef3121103) Kind() string    { return "技能" }
func (CardDef3121103) Element() string { return "火" }

func (CardDef3121103) Card() model.Card {
	return model.Card{
		Number:          "3121103",
		Type:            "技能",
		Name:            "祈祷之焰",
		Category:        "火",
		Tag:             "咒术-神秘",
		Description:     "在此卡上放置3个标记物,或者取除所有标记物并召唤1个入场花费小于等于标记物数量的火焰伙伴,无需花费",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3121103.jpg",
	}
}

type CardDef3121104 struct{}

func (CardDef3121104) ID() string      { return "3121104" }
func (CardDef3121104) Name() string    { return "炎狱怒吼" }
func (CardDef3121104) Kind() string    { return "技能" }
func (CardDef3121104) Element() string { return "火" }

func (CardDef3121104) Card() model.Card {
	return model.Card{
		Number:          "3121104",
		Type:            "技能",
		Name:            "炎狱怒吼",
		Category:        "火",
		Tag:             "法术-聚能",
		Description:     "范围:纵列.使用时可以额外消耗1个友方火焰伙伴,此次\\威上升其入场花费的数值",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2, "火": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1, "火": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3121104.jpg",
	}
}

type CardDef3121105 struct{}

func (CardDef3121105) ID() string      { return "3121105" }
func (CardDef3121105) Name() string    { return "余火" }
func (CardDef3121105) Kind() string    { return "技能" }
func (CardDef3121105) Element() string { return "火" }

func (CardDef3121105) Card() model.Card {
	return model.Card{
		Number:          "3121105",
		Type:            "技能",
		Name:            "余火",
		Category:        "火",
		Tag:             "法术-创造",
		Description:     "不可用于防御.诱发:你的回合结束时,如果你没有剩余的\\火,可以立刻使用此卡且无需花费",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3121105.jpg",
	}
}

type CardDef3121106 struct{}

func (CardDef3121106) ID() string      { return "3121106" }
func (CardDef3121106) Name() string    { return "爆炎气焰" }
func (CardDef3121106) Kind() string    { return "技能" }
func (CardDef3121106) Element() string { return "火" }

func (CardDef3121106) Card() model.Card {
	return model.Card{
		Number:          "3121106",
		Type:            "技能",
		Name:            "爆炎气焰",
		Category:        "火",
		Tag:             "法术-创造",
		Description:     "穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1, "火": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 1, "火": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3121106.jpg",
	}
}

type CardDef3121107 struct{}

func (CardDef3121107) ID() string      { return "3121107" }
func (CardDef3121107) Name() string    { return "战争践踏" }
func (CardDef3121107) Kind() string    { return "技能" }
func (CardDef3121107) Element() string { return "火" }

func (CardDef3121107) Card() model.Card {
	return model.Card{
		Number:          "3121107",
		Type:            "技能",
		Name:            "战争践踏",
		Category:        "火",
		Tag:             "法术-驱动",
		Description:     "范围:前排.攻击时目标区域每有1个单位此法术-1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          4,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3121107.jpg",
	}
}

type CardDef3121108 struct{}

func (CardDef3121108) ID() string      { return "3121108" }
func (CardDef3121108) Name() string    { return "熔岩障壁" }
func (CardDef3121108) Kind() string    { return "技能" }
func (CardDef3121108) Element() string { return "火" }

func (CardDef3121108) Card() model.Card {
	return model.Card{
		Number:          "3121108",
		Type:            "技能",
		Name:            "熔岩障壁",
		Category:        "火",
		Tag:             "法术-创造",
		Description:     "防御",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1, "火": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3121108.jpg",
	}
}

type CardDef3121109 struct{}

func (CardDef3121109) ID() string      { return "3121109" }
func (CardDef3121109) Name() string    { return "烈焰闪" }
func (CardDef3121109) Kind() string    { return "技能" }
func (CardDef3121109) Element() string { return "火" }

func (CardDef3121109) Card() model.Card {
	return model.Card{
		Number:          "3121109",
		Type:            "技能",
		Name:            "烈焰闪",
		Category:        "火",
		Tag:             "法术-幻变",
		Description:     "命中:获得3\\火",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3121109.jpg",
	}
}

type CardDef3121110 struct{}

func (CardDef3121110) ID() string      { return "3121110" }
func (CardDef3121110) Name() string    { return "咒火" }
func (CardDef3121110) Kind() string    { return "技能" }
func (CardDef3121110) Element() string { return "火" }

func (CardDef3121110) Card() model.Card {
	return model.Card{
		Number:          "3121110",
		Type:            "技能",
		Name:            "咒火",
		Category:        "火",
		Tag:             "咒术-创造",
		Description:     "冷却1.翻取1个花费小于4的火焰法术卷轴,你可以立刻使用它,无需花费",
		Quote:           "",
		ElementsCost:    map[string]int{"火": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"火": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\火\\3121110.jpg",
	}
}

type CardDef3201001 struct{}

func (CardDef3201001) ID() string      { return "3201001" }
func (CardDef3201001) Name() string    { return "百川归海" }
func (CardDef3201001) Kind() string    { return "技能" }
func (CardDef3201001) Element() string { return "水" }

func (CardDef3201001) Card() model.Card {
	return model.Card{
		Number:          "3201001",
		Type:            "技能",
		Name:            "百川归海",
		Category:        "水",
		Tag:             "衍生-法术-聚能",
		Description:     "防御.若防御成功,获得等同于所有攻击法术的攻击力合计的\\水",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3201001.jpg",
	}
}

type CardDef3201002 struct{}

func (CardDef3201002) ID() string      { return "3201002" }
func (CardDef3201002) Name() string    { return "凛冬将至" }
func (CardDef3201002) Kind() string    { return "技能" }
func (CardDef3201002) Element() string { return "水" }

func (CardDef3201002) Card() model.Card {
	return model.Card{
		Number:          "3201002",
		Type:            "技能",
		Name:            "凛冬将至",
		Category:        "水",
		Tag:             "衍生-法术-幻变",
		Description:     "范围:溅射.穿透.冻结1.使用时必须移除嗜魔弓 凛冬上的5个水纹标记物",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 4},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           8,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3201002.jpg",
	}
}

type CardDef3211101 struct{}

func (CardDef3211101) ID() string      { return "3211101" }
func (CardDef3211101) Name() string    { return "心海迷离" }
func (CardDef3211101) Kind() string    { return "技能" }
func (CardDef3211101) Element() string { return "水" }

func (CardDef3211101) Card() model.Card {
	return model.Card{
		Number:          "3211101",
		Type:            "技能",
		Name:            "心海迷离",
		Category:        "水",
		Tag:             "传奇-法术-灵媒",
		Description:     "晕眩1.诱发:你每次使用此技能后直到回合结束获得+1\\威并改为任意范围(前排,纵列,方阵,溅射)",
		Quote:           "层层幻梦,层层迷离",
		ElementsCost:    map[string]int{"水": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\水\\3211101.jpg",
	}
}

type CardDef3211102 struct{}

func (CardDef3211102) ID() string      { return "3211102" }
func (CardDef3211102) Name() string    { return "龙吟雪域" }
func (CardDef3211102) Kind() string    { return "技能" }
func (CardDef3211102) Element() string { return "水" }

func (CardDef3211102) Card() model.Card {
	return model.Card{
		Number:          "3211102",
		Type:            "技能",
		Name:            "龙吟雪域",
		Category:        "水",
		Tag:             "传奇-咒术-幻变",
		Description:     "冷却2.异能:双方玩家的法术获得以下效果\"命中:选择法力范围内1个单位冻结1\".第5次触发这个效果后你可以召唤一个凛冰之龙",
		Quote:           "\"是的小姑娘,她知道你的一切,如果你真的能感动她,就尝试穿过这片无尽雪域吧\"",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        2,
		Power:           -1,
		Spawns:          []string{"1201101"},
		OutputPath:      "output\\王权纷争\\技能\\水\\3211102.jpg",
	}
}

type CardDef3221001 struct{}

func (CardDef3221001) ID() string      { return "3221001" }
func (CardDef3221001) Name() string    { return "冰雹术" }
func (CardDef3221001) Kind() string    { return "技能" }
func (CardDef3221001) Element() string { return "水" }

func (CardDef3221001) Card() model.Card {
	return model.Card{
		Number:          "3221001",
		Type:            "技能",
		Name:            "冰雹术",
		Category:        "水",
		Tag:             "法术-幻变",
		Description:     "范围:方阵",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221001.jpg",
	}
}

type CardDef3221002 struct{}

func (CardDef3221002) ID() string      { return "3221002" }
func (CardDef3221002) Name() string    { return "冰锥术" }
func (CardDef3221002) Kind() string    { return "技能" }
func (CardDef3221002) Element() string { return "水" }

func (CardDef3221002) Card() model.Card {
	return model.Card{
		Number:          "3221002",
		Type:            "技能",
		Name:            "冰锥术",
		Category:        "水",
		Tag:             "法术-创造",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221002.jpg",
	}
}

type CardDef3221003 struct{}

func (CardDef3221003) ID() string      { return "3221003" }
func (CardDef3221003) Name() string    { return "激冻寒流" }
func (CardDef3221003) Kind() string    { return "技能" }
func (CardDef3221003) Element() string { return "水" }

func (CardDef3221003) Card() model.Card {
	return model.Card{
		Number:          "3221003",
		Type:            "技能",
		Name:            "激冻寒流",
		Category:        "水",
		Tag:             "法术-聚能",
		Description:     "强化其他水纹法术时+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221003.jpg",
	}
}

type CardDef3221004 struct{}

func (CardDef3221004) ID() string      { return "3221004" }
func (CardDef3221004) Name() string    { return "寒冰屏障" }
func (CardDef3221004) Kind() string    { return "技能" }
func (CardDef3221004) Element() string { return "水" }

func (CardDef3221004) Card() model.Card {
	return model.Card{
		Number:          "3221004",
		Type:            "技能",
		Name:            "寒冰屏障",
		Category:        "水",
		Tag:             "法术-创造",
		Description:     "防御",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221004.jpg",
	}
}

type CardDef3221005 struct{}

func (CardDef3221005) ID() string      { return "3221005" }
func (CardDef3221005) Name() string    { return "玄冰阵" }
func (CardDef3221005) Kind() string    { return "技能" }
func (CardDef3221005) Element() string { return "水" }

func (CardDef3221005) Card() model.Card {
	return model.Card{
		Number:          "3221005",
		Type:            "技能",
		Name:            "玄冰阵",
		Category:        "水",
		Tag:             "法术-创造",
		Description:     "范围:溅射.冻结1",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221005.jpg",
	}
}

type CardDef3221006 struct{}

func (CardDef3221006) ID() string      { return "3221006" }
func (CardDef3221006) Name() string    { return "海啸" }
func (CardDef3221006) Kind() string    { return "技能" }
func (CardDef3221006) Element() string { return "水" }

func (CardDef3221006) Card() model.Card {
	return model.Card{
		Number:          "3221006",
		Type:            "技能",
		Name:            "海啸",
		Category:        "水",
		Tag:             "法术-驱动",
		Description:     "冷却1.范围:全场",
		Quote:           "参天的浪潮裹挟着天空,随后便是利齿,以及无尽的黑暗",
		ElementsCost:    map[string]int{"水": 8},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 5},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           9,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221006.jpg",
	}
}

type CardDef3221007 struct{}

func (CardDef3221007) ID() string      { return "3221007" }
func (CardDef3221007) Name() string    { return "水占术" }
func (CardDef3221007) Kind() string    { return "技能" }
func (CardDef3221007) Element() string { return "水" }

func (CardDef3221007) Card() model.Card {
	return model.Card{
		Number:          "3221007",
		Type:            "技能",
		Name:            "水占术",
		Category:        "水",
		Tag:             "咒术-灵媒",
		Description:     "冷却1.查看牌堆顶4张牌并检索其中1张水纹卡牌,其余按任意顺序置于牌堆顶或牌堆底",
		Quote:           "\"可能你不相信命运,但是命运似乎很相信你\"",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221007.jpg",
	}
}

type CardDef3221008 struct{}

func (CardDef3221008) ID() string      { return "3221008" }
func (CardDef3221008) Name() string    { return "冰封消解" }
func (CardDef3221008) Kind() string    { return "技能" }
func (CardDef3221008) Element() string { return "水" }

func (CardDef3221008) Card() model.Card {
	return model.Card{
		Number:          "3221008",
		Type:            "技能",
		Name:            "冰封消解",
		Category:        "水",
		Tag:             "咒术-幻变",
		Description:     "冷却1.诱发:当对方使用法术时可以使用此卡,使其中1个\\威变为0",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221008.jpg",
	}
}

type CardDef3221009 struct{}

func (CardDef3221009) ID() string      { return "3221009" }
func (CardDef3221009) Name() string    { return "冰霜利刃" }
func (CardDef3221009) Kind() string    { return "技能" }
func (CardDef3221009) Element() string { return "水" }

func (CardDef3221009) Card() model.Card {
	return model.Card{
		Number:          "3221009",
		Type:            "技能",
		Name:            "冰霜利刃",
		Category:        "水",
		Tag:             "法术-创造",
		Description:     "此卡攻击或强化攻击时+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           2,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221009.jpg",
	}
}

type CardDef3221010 struct{}

func (CardDef3221010) ID() string      { return "3221010" }
func (CardDef3221010) Name() string    { return "水幻影" }
func (CardDef3221010) Kind() string    { return "技能" }
func (CardDef3221010) Element() string { return "水" }

func (CardDef3221010) Card() model.Card {
	return model.Card{
		Number:          "3221010",
		Type:            "技能",
		Name:            "水幻影",
		Category:        "水",
		Tag:             "咒术-创造",
		Description:     "冷却1.选择1个本回合你召唤的水纹伙伴,召唤1个只有1\\血的复制.",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1, "水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 1, "水": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221010.jpg",
	}
}

type CardDef3221011 struct{}

func (CardDef3221011) ID() string      { return "3221011" }
func (CardDef3221011) Name() string    { return "幽影寒锋" }
func (CardDef3221011) Kind() string    { return "技能" }
func (CardDef3221011) Element() string { return "水" }

func (CardDef3221011) Card() model.Card {
	return model.Card{
		Number:          "3221011",
		Type:            "技能",
		Name:            "幽影寒锋",
		Category:        "水",
		Tag:             "法术-聚能",
		Description:     "穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1, "水": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221011.jpg",
	}
}

type CardDef3221012 struct{}

func (CardDef3221012) ID() string      { return "3221012" }
func (CardDef3221012) Name() string    { return "霜冻射线" }
func (CardDef3221012) Kind() string    { return "技能" }
func (CardDef3221012) Element() string { return "水" }

func (CardDef3221012) Card() model.Card {
	return model.Card{
		Number:          "3221012",
		Type:            "技能",
		Name:            "霜冻射线",
		Category:        "水",
		Tag:             "法术-聚能",
		Description:     "冻结2",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221012.jpg",
	}
}

type CardDef3221013 struct{}

func (CardDef3221013) ID() string      { return "3221013" }
func (CardDef3221013) Name() string    { return "猎潮" }
func (CardDef3221013) Kind() string    { return "技能" }
func (CardDef3221013) Element() string { return "水" }

func (CardDef3221013) Card() model.Card {
	return model.Card{
		Number:          "3221013",
		Type:            "技能",
		Name:            "猎潮",
		Category:        "水",
		Tag:             "法术-驱动",
		Description:     "",
		Quote:           "归来的勇士跃向无底的巨口,鳞甲深处是它跳动的心脏",
		ElementsCost:    map[string]int{"气": 1, "水": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          3,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221013.jpg",
	}
}

type CardDef3221014 struct{}

func (CardDef3221014) ID() string      { return "3221014" }
func (CardDef3221014) Name() string    { return "坚冰领域" }
func (CardDef3221014) Kind() string    { return "技能" }
func (CardDef3221014) Element() string { return "水" }

func (CardDef3221014) Card() model.Card {
	return model.Card{
		Number:          "3221014",
		Type:            "技能",
		Name:            "坚冰领域",
		Category:        "水",
		Tag:             "法术-创造",
		Description:     "防御.若防御成功,使所有前排敌人冻结1",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221014.jpg",
	}
}

type CardDef3221015 struct{}

func (CardDef3221015) ID() string      { return "3221015" }
func (CardDef3221015) Name() string    { return "暴风雪" }
func (CardDef3221015) Kind() string    { return "技能" }
func (CardDef3221015) Element() string { return "水" }

func (CardDef3221015) Card() model.Card {
	return model.Card{
		Number:          "3221015",
		Type:            "技能",
		Name:            "暴风雪",
		Category:        "水",
		Tag:             "咒术-幻变",
		Description:     "冷却1.异能:你的水纹和大气法术获得冻结1和+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1, "水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 1, "水": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\水\\3221015.jpg",
	}
}

type CardDef3221101 struct{}

func (CardDef3221101) ID() string      { return "3221101" }
func (CardDef3221101) Name() string    { return "踏浪术" }
func (CardDef3221101) Kind() string    { return "技能" }
func (CardDef3221101) Element() string { return "水" }

func (CardDef3221101) Card() model.Card {
	return model.Card{
		Number:          "3221101",
		Type:            "技能",
		Name:            "踏浪术",
		Category:        "水",
		Tag:             "咒术-驱动",
		Description:     "冷却2.异能:在你使用1个水纹法术后,重置你的另一个水纹法术,但使其下一次花费多X+1\\水,X为本回合触发此效果的次数",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\水\\3221101.jpg",
	}
}

type CardDef3221102 struct{}

func (CardDef3221102) ID() string      { return "3221102" }
func (CardDef3221102) Name() string    { return "唤灵术 蛟龙" }
func (CardDef3221102) Kind() string    { return "技能" }
func (CardDef3221102) Element() string { return "水" }

func (CardDef3221102) Card() model.Card {
	return model.Card{
		Number:          "3221102",
		Type:            "技能",
		Name:            "唤灵术 蛟龙",
		Category:        "水",
		Tag:             "法术-创造",
		Description:     "若防御成功对法力范围内所有敌人造成1点伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\水\\3221102.jpg",
	}
}

type CardDef3221103 struct{}

func (CardDef3221103) ID() string      { return "3221103" }
func (CardDef3221103) Name() string    { return "水镜壁" }
func (CardDef3221103) Kind() string    { return "技能" }
func (CardDef3221103) Element() string { return "水" }

func (CardDef3221103) Card() model.Card {
	return model.Card{
		Number:          "3221103",
		Type:            "技能",
		Name:            "水镜壁",
		Category:        "水",
		Tag:             "法术-创造",
		Description:     "防御.若防御成功获得护盾1",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\水\\3221103.jpg",
	}
}

type CardDef3221104 struct{}

func (CardDef3221104) ID() string      { return "3221104" }
func (CardDef3221104) Name() string    { return "水遁术" }
func (CardDef3221104) Kind() string    { return "技能" }
func (CardDef3221104) Element() string { return "水" }

func (CardDef3221104) Card() model.Card {
	return model.Card{
		Number:          "3221104",
		Type:            "技能",
		Name:            "水遁术",
		Category:        "水",
		Tag:             "咒术-幻变",
		Description:     "冷却1.速攻.使1个没有隐蔽的单位获得隐蔽2",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\水\\3221104.jpg",
	}
}

type CardDef3221105 struct{}

func (CardDef3221105) ID() string      { return "3221105" }
func (CardDef3221105) Name() string    { return "腐蚀之流" }
func (CardDef3221105) Kind() string    { return "技能" }
func (CardDef3221105) Element() string { return "水" }

func (CardDef3221105) Card() model.Card {
	return model.Card{
		Number:          "3221105",
		Type:            "技能",
		Name:            "腐蚀之流",
		Category:        "水",
		Tag:             "法术-幻变",
		Description:     "命中:随机弃置敌方1张手牌",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1, "水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 1, "水": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\水\\3221105.jpg",
	}
}

type CardDef3221106 struct{}

func (CardDef3221106) ID() string      { return "3221106" }
func (CardDef3221106) Name() string    { return "暗流涌动" }
func (CardDef3221106) Kind() string    { return "技能" }
func (CardDef3221106) Element() string { return "水" }

func (CardDef3221106) Card() model.Card {
	return model.Card{
		Number:          "3221106",
		Type:            "技能",
		Name:            "暗流涌动",
		Category:        "水",
		Tag:             "法术-驱动",
		Description:     "可以攻击隐蔽的单位(无论在哪),攻击隐蔽单位时+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\水\\3221106.jpg",
	}
}

type CardDef3221107 struct{}

func (CardDef3221107) ID() string      { return "3221107" }
func (CardDef3221107) Name() string    { return "海龙卷" }
func (CardDef3221107) Kind() string    { return "技能" }
func (CardDef3221107) Element() string { return "水" }

func (CardDef3221107) Card() model.Card {
	return model.Card{
		Number:          "3221107",
		Type:            "技能",
		Name:            "海龙卷",
		Category:        "水",
		Tag:             "法术-驱动",
		Description:     "范围:方阵.",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1, "水": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 1, "水": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\水\\3221107.jpg",
	}
}

type CardDef3221108 struct{}

func (CardDef3221108) ID() string      { return "3221108" }
func (CardDef3221108) Name() string    { return "六瓣雪花" }
func (CardDef3221108) Kind() string    { return "技能" }
func (CardDef3221108) Element() string { return "水" }

func (CardDef3221108) Card() model.Card {
	return model.Card{
		Number:          "3221108",
		Type:            "技能",
		Name:            "六瓣雪花",
		Category:        "水",
		Tag:             "法术-创造",
		Description:     "冻结1.速攻.此卡的冻结对人物无效",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\水\\3221108.jpg",
	}
}

type CardDef3221109 struct{}

func (CardDef3221109) ID() string      { return "3221109" }
func (CardDef3221109) Name() string    { return "波纹斩" }
func (CardDef3221109) Kind() string    { return "技能" }
func (CardDef3221109) Element() string { return "水" }

func (CardDef3221109) Card() model.Card {
	return model.Card{
		Number:          "3221109",
		Type:            "技能",
		Name:            "波纹斩",
		Category:        "水",
		Tag:             "法术-创造",
		Description:     "光环:如果本回合你已经使用过波纹斩,此卡获得+2\\威和范围:前排",
		Quote:           "",
		ElementsCost:    map[string]int{"水": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"水": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\水\\3221109.jpg",
	}
}

type CardDef3221110 struct{}

func (CardDef3221110) ID() string      { return "3221110" }
func (CardDef3221110) Name() string    { return "劫掠之潮" }
func (CardDef3221110) Kind() string    { return "技能" }
func (CardDef3221110) Element() string { return "水" }

func (CardDef3221110) Card() model.Card {
	return model.Card{
		Number:          "3221110",
		Type:            "技能",
		Name:            "劫掠之潮",
		Category:        "水",
		Tag:             "法术-驱动",
		Description:     "范围:前排.命中:每命中1个单位随机弃置敌方1张手牌并且你抽1张牌",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1, "水": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 1, "水": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\水\\3221110.jpg",
	}
}

type CardDef3301001 struct{}

func (CardDef3301001) ID() string      { return "3301001" }
func (CardDef3301001) Name() string    { return "风暴之怒" }
func (CardDef3301001) Kind() string    { return "技能" }
func (CardDef3301001) Element() string { return "气" }

func (CardDef3301001) Card() model.Card {
	return model.Card{
		Number:          "3301001",
		Type:            "技能",
		Name:            "风暴之怒",
		Category:        "气",
		Tag:             "衍生-咒术-驱动",
		Description:     "异能:展示你的所有手牌,每张使你的大气法术+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3301001.jpg",
	}
}

type CardDef3311101 struct{}

func (CardDef3311101) ID() string      { return "3311101" }
func (CardDef3311101) Name() string    { return "苍穹幻韵" }
func (CardDef3311101) Kind() string    { return "技能" }
func (CardDef3311101) Element() string { return "气" }

func (CardDef3311101) Card() model.Card {
	return model.Card{
		Number:          "3311101",
		Type:            "技能",
		Name:            "苍穹幻韵",
		Category:        "气",
		Tag:             "传奇-法术-幻变",
		Description:     "使用时等同于释放你学习的另一个驱动或聚能法术",
		Quote:           "无形之风为意,浮云之影成形",
		ElementsCost:    map[string]int{"光": 2, "气": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 1, "气": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           0,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3311101.jpg",
	}
}

type CardDef3311102 struct{}

func (CardDef3311102) ID() string      { return "3311102" }
func (CardDef3311102) Name() string    { return "星落之银叶" }
func (CardDef3311102) Kind() string    { return "技能" }
func (CardDef3311102) Element() string { return "气" }

func (CardDef3311102) Card() model.Card {
	return model.Card{
		Number:          "3311102",
		Type:            "技能",
		Name:            "星落之银叶",
		Category:        "气",
		Tag:             "传奇-咒术-创造",
		Description:     "诱发:每当你弃牌时,从中选择1张放在此卡下方.使用时将此卡下方1张牌洗回卡组,然后抽1张牌",
		Quote:           "\"有人说,我们是坠落的神灵弃子,但我想,我们是受眷顾的苍茫众生\"",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3311102.jpg",
	}
}

type CardDef3321001 struct{}

func (CardDef3321001) ID() string      { return "3321001" }
func (CardDef3321001) Name() string    { return "闪电链" }
func (CardDef3321001) Kind() string    { return "技能" }
func (CardDef3321001) Element() string { return "气" }

func (CardDef3321001) Card() model.Card {
	return model.Card{
		Number:          "3321001",
		Type:            "技能",
		Name:            "闪电链",
		Category:        "气",
		Tag:             "法术-创造",
		Description:     "可以额外选择1个无视范围的目标",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321001.jpg",
	}
}

type CardDef3321002 struct{}

func (CardDef3321002) ID() string      { return "3321002" }
func (CardDef3321002) Name() string    { return "雷击" }
func (CardDef3321002) Kind() string    { return "技能" }
func (CardDef3321002) Element() string { return "气" }

func (CardDef3321002) Card() model.Card {
	return model.Card{
		Number:          "3321002",
		Type:            "技能",
		Name:            "雷击",
		Category:        "气",
		Tag:             "法术-驱动",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321002.jpg",
	}
}

type CardDef3321003 struct{}

func (CardDef3321003) ID() string      { return "3321003" }
func (CardDef3321003) Name() string    { return "静电脉冲" }
func (CardDef3321003) Kind() string    { return "技能" }
func (CardDef3321003) Element() string { return "气" }

func (CardDef3321003) Card() model.Card {
	return model.Card{
		Number:          "3321003",
		Type:            "技能",
		Name:            "静电脉冲",
		Category:        "气",
		Tag:             "法术-聚能",
		Description:     "晕眩1",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321003.jpg",
	}
}

type CardDef3321004 struct{}

func (CardDef3321004) ID() string      { return "3321004" }
func (CardDef3321004) Name() string    { return "雷闪" }
func (CardDef3321004) Kind() string    { return "技能" }
func (CardDef3321004) Element() string { return "气" }

func (CardDef3321004) Card() model.Card {
	return model.Card{
		Number:          "3321004",
		Type:            "技能",
		Name:            "雷闪",
		Category:        "气",
		Tag:             "法术-驱动",
		Description:     "穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "气": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321004.jpg",
	}
}

type CardDef3321005 struct{}

func (CardDef3321005) ID() string      { return "3321005" }
func (CardDef3321005) Name() string    { return "气旋波" }
func (CardDef3321005) Kind() string    { return "技能" }
func (CardDef3321005) Element() string { return "气" }

func (CardDef3321005) Card() model.Card {
	return model.Card{
		Number:          "3321005",
		Type:            "技能",
		Name:            "气旋波",
		Category:        "气",
		Tag:             "法术-创造",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321005.jpg",
	}
}

type CardDef3321006 struct{}

func (CardDef3321006) ID() string      { return "3321006" }
func (CardDef3321006) Name() string    { return "雷暴术" }
func (CardDef3321006) Kind() string    { return "技能" }
func (CardDef3321006) Element() string { return "气" }

func (CardDef3321006) Card() model.Card {
	return model.Card{
		Number:          "3321006",
		Type:            "技能",
		Name:            "雷暴术",
		Category:        "气",
		Tag:             "法术-聚能",
		Description:     "范围:方阵,晕眩1",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "气": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321006.jpg",
	}
}

type CardDef3321007 struct{}

func (CardDef3321007) ID() string      { return "3321007" }
func (CardDef3321007) Name() string    { return "源力之风" }
func (CardDef3321007) Kind() string    { return "技能" }
func (CardDef3321007) Element() string { return "气" }

func (CardDef3321007) Card() model.Card {
	return model.Card{
		Number:          "3321007",
		Type:            "技能",
		Name:            "源力之风",
		Category:        "气",
		Tag:             "咒术-灵媒",
		Description:     "冷却1.补充手牌至手牌上限,每补1张花费1\\气",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321007.jpg",
	}
}

type CardDef3321008 struct{}

func (CardDef3321008) ID() string      { return "3321008" }
func (CardDef3321008) Name() string    { return "风洞" }
func (CardDef3321008) Kind() string    { return "技能" }
func (CardDef3321008) Element() string { return "气" }

func (CardDef3321008) Card() model.Card {
	return model.Card{
		Number:          "3321008",
		Type:            "技能",
		Name:            "风洞",
		Category:        "气",
		Tag:             "咒术-创造",
		Description:     "冷却1.诱发:当敌方的非范围法术命中时可以使用此卡,将其无效",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321008.jpg",
	}
}

type CardDef3321009 struct{}

func (CardDef3321009) ID() string      { return "3321009" }
func (CardDef3321009) Name() string    { return "宇宙飓风" }
func (CardDef3321009) Kind() string    { return "技能" }
func (CardDef3321009) Element() string { return "气" }

func (CardDef3321009) Card() model.Card {
	return model.Card{
		Number:          "3321009",
		Type:            "技能",
		Name:            "宇宙飓风",
		Category:        "气",
		Tag:             "法术-驱动",
		Description:     "穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 8},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 5},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          4,
		Life:            -1,
		Duration:        -1,
		Power:           9,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321009.jpg",
	}
}

type CardDef3321010 struct{}

func (CardDef3321010) ID() string      { return "3321010" }
func (CardDef3321010) Name() string    { return "涡旋屏障" }
func (CardDef3321010) Kind() string    { return "技能" }
func (CardDef3321010) Element() string { return "气" }

func (CardDef3321010) Card() model.Card {
	return model.Card{
		Number:          "3321010",
		Type:            "技能",
		Name:            "涡旋屏障",
		Category:        "气",
		Tag:             "法术-创造",
		Description:     "防御",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321010.jpg",
	}
}

type CardDef3321011 struct{}

func (CardDef3321011) ID() string      { return "3321011" }
func (CardDef3321011) Name() string    { return "撕裂长空" }
func (CardDef3321011) Kind() string    { return "技能" }
func (CardDef3321011) Element() string { return "气" }

func (CardDef3321011) Card() model.Card {
	return model.Card{
		Number:          "3321011",
		Type:            "技能",
		Name:            "撕裂长空",
		Category:        "气",
		Tag:             "法术-聚能",
		Description:     "范围:纵列",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 7},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 4},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321011.jpg",
	}
}

type CardDef3321012 struct{}

func (CardDef3321012) ID() string      { return "3321012" }
func (CardDef3321012) Name() string    { return "空天感应" }
func (CardDef3321012) Kind() string    { return "技能" }
func (CardDef3321012) Element() string { return "气" }

func (CardDef3321012) Card() model.Card {
	return model.Card{
		Number:          "3321012",
		Type:            "技能",
		Name:            "空天感应",
		Category:        "气",
		Tag:             "咒术-灵媒",
		Description:     "冷却1.速攻.异能:如果你的法术目标或区域中包含了非前排单位,使其获得+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321012.jpg",
	}
}

type CardDef3321013 struct{}

func (CardDef3321013) ID() string      { return "3321013" }
func (CardDef3321013) Name() string    { return "霹雳惊雷" }
func (CardDef3321013) Kind() string    { return "技能" }
func (CardDef3321013) Element() string { return "气" }

func (CardDef3321013) Card() model.Card {
	return model.Card{
		Number:          "3321013",
		Type:            "技能",
		Name:            "霹雳惊雷",
		Category:        "气",
		Tag:             "法术-驱动",
		Description:     "速攻.穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           2,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321013.jpg",
	}
}

type CardDef3321014 struct{}

func (CardDef3321014) ID() string      { return "3321014" }
func (CardDef3321014) Name() string    { return "引雷" }
func (CardDef3321014) Kind() string    { return "技能" }
func (CardDef3321014) Element() string { return "气" }

func (CardDef3321014) Card() model.Card {
	return model.Card{
		Number:          "3321014",
		Type:            "技能",
		Name:            "引雷",
		Category:        "气",
		Tag:             "咒术-驱动",
		Description:     "冷却1.丢弃1张手牌,使1个敌方伙伴晕眩1",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321014.jpg",
	}
}

type CardDef3321015 struct{}

func (CardDef3321015) ID() string      { return "3321015" }
func (CardDef3321015) Name() string    { return "静电屏障" }
func (CardDef3321015) Kind() string    { return "技能" }
func (CardDef3321015) Element() string { return "气" }

func (CardDef3321015) Card() model.Card {
	return model.Card{
		Number:          "3321015",
		Type:            "技能",
		Name:            "静电屏障",
		Category:        "气",
		Tag:             "法术-创造",
		Description:     "防御.若防御失败,使1个前排敌人晕眩1",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\气\\3321015.jpg",
	}
}

type CardDef3321101 struct{}

func (CardDef3321101) ID() string      { return "3321101" }
func (CardDef3321101) Name() string    { return "急速涡旋" }
func (CardDef3321101) Kind() string    { return "技能" }
func (CardDef3321101) Element() string { return "气" }

func (CardDef3321101) Card() model.Card {
	return model.Card{
		Number:          "3321101",
		Type:            "技能",
		Name:            "急速涡旋",
		Category:        "气",
		Tag:             "法术-驱动",
		Description:     "速攻",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3321101.jpg",
	}
}

type CardDef3321102 struct{}

func (CardDef3321102) ID() string      { return "3321102" }
func (CardDef3321102) Name() string    { return "肃杀之风" }
func (CardDef3321102) Kind() string    { return "技能" }
func (CardDef3321102) Element() string { return "气" }

func (CardDef3321102) Card() model.Card {
	return model.Card{
		Number:          "3321102",
		Type:            "技能",
		Name:            "肃杀之风",
		Category:        "气",
		Tag:             "法术-驱动",
		Description:     "范围:方阵.光环:此卡威力上升双方玩家手牌数量的差值",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3321102.jpg",
	}
}

type CardDef3321103 struct{}

func (CardDef3321103) ID() string      { return "3321103" }
func (CardDef3321103) Name() string    { return "雷霆万钧" }
func (CardDef3321103) Kind() string    { return "技能" }
func (CardDef3321103) Element() string { return "气" }

func (CardDef3321103) Card() model.Card {
	return model.Card{
		Number:          "3321103",
		Type:            "技能",
		Name:            "雷霆万钧",
		Category:        "气",
		Tag:             "法术-聚能",
		Description:     "范围:前排",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2, "气": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 1, "气": 4},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          3,
		Life:            -1,
		Duration:        -1,
		Power:           9,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3321103.jpg",
	}
}

type CardDef3321104 struct{}

func (CardDef3321104) ID() string      { return "3321104" }
func (CardDef3321104) Name() string    { return "收势" }
func (CardDef3321104) Kind() string    { return "技能" }
func (CardDef3321104) Element() string { return "气" }

func (CardDef3321104) Card() model.Card {
	return model.Card{
		Number:          "3321104",
		Type:            "技能",
		Name:            "收势",
		Category:        "气",
		Tag:             "法术-聚能",
		Description:     "防御.若防御成功,你下一个用于攻击的法术+3\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3321104.jpg",
	}
}

type CardDef3321105 struct{}

func (CardDef3321105) ID() string      { return "3321105" }
func (CardDef3321105) Name() string    { return "风卷残云" }
func (CardDef3321105) Kind() string    { return "技能" }
func (CardDef3321105) Element() string { return "气" }

func (CardDef3321105) Card() model.Card {
	return model.Card{
		Number:          "3321105",
		Type:            "技能",
		Name:            "风卷残云",
		Category:        "气",
		Tag:             "咒术-幻变",
		Description:     "速攻.冷却1.异能:每当有单位受到伤害,如果其生命为1,将其消灭",
		Quote:           "一不做,二不休",
		ElementsCost:    map[string]int{"气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3321105.jpg",
	}
}

type CardDef3321106 struct{}

func (CardDef3321106) ID() string      { return "3321106" }
func (CardDef3321106) Name() string    { return "紫电穿空" }
func (CardDef3321106) Kind() string    { return "技能" }
func (CardDef3321106) Element() string { return "气" }

func (CardDef3321106) Card() model.Card {
	return model.Card{
		Number:          "3321106",
		Type:            "技能",
		Name:            "紫电穿空",
		Category:        "气",
		Tag:             "法术-聚能",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3321106.jpg",
	}
}

type CardDef3321107 struct{}

func (CardDef3321107) ID() string      { return "3321107" }
func (CardDef3321107) Name() string    { return "屏息凝神" }
func (CardDef3321107) Kind() string    { return "技能" }
func (CardDef3321107) Element() string { return "气" }

func (CardDef3321107) Card() model.Card {
	return model.Card{
		Number:          "3321107",
		Type:            "技能",
		Name:            "屏息凝神",
		Category:        "气",
		Tag:             "咒术-灵媒",
		Description:     "速攻.异能:如果本回合你除回合开始外没有抽牌,你的大气法术+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3321107.jpg",
	}
}

type CardDef3321108 struct{}

func (CardDef3321108) ID() string      { return "3321108" }
func (CardDef3321108) Name() string    { return "唤灵术 苍鹰" }
func (CardDef3321108) Kind() string    { return "技能" }
func (CardDef3321108) Element() string { return "气" }

func (CardDef3321108) Card() model.Card {
	return model.Card{
		Number:          "3321108",
		Type:            "技能",
		Name:            "唤灵术 苍鹰",
		Category:        "气",
		Tag:             "法术-创造",
		Description:     "入场:使1个友方大气法术下一次使用时+1\\攻+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3321108.jpg",
	}
}

type CardDef3321109 struct{}

func (CardDef3321109) ID() string      { return "3321109" }
func (CardDef3321109) Name() string    { return "银叶旋风" }
func (CardDef3321109) Kind() string    { return "技能" }
func (CardDef3321109) Element() string { return "气" }

func (CardDef3321109) Card() model.Card {
	return model.Card{
		Number:          "3321109",
		Type:            "技能",
		Name:            "银叶旋风",
		Category:        "气",
		Tag:             "法术-驱动",
		Description:     "光环:如果当回合有卡被送去弃牌堆,此卡\\威变为6",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3321109.jpg",
	}
}

type CardDef3321110 struct{}

func (CardDef3321110) ID() string      { return "3321110" }
func (CardDef3321110) Name() string    { return "气蕴成流" }
func (CardDef3321110) Kind() string    { return "技能" }
func (CardDef3321110) Element() string { return "气" }

func (CardDef3321110) Card() model.Card {
	return model.Card{
		Number:          "3321110",
		Type:            "技能",
		Name:            "气蕴成流",
		Category:        "气",
		Tag:             "法术-幻变",
		Description:     "入场:你学习的下一个大气法术获得速攻",
		Quote:           "",
		ElementsCost:    map[string]int{"气": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"气": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\气\\3321110.jpg",
	}
}

type CardDef3411101 struct{}

func (CardDef3411101) ID() string      { return "3411101" }
func (CardDef3411101) Name() string    { return "时岁轮转" }
func (CardDef3411101) Kind() string    { return "技能" }
func (CardDef3411101) Element() string { return "地" }

func (CardDef3411101) Card() model.Card {
	return model.Card{
		Number:          "3411101",
		Type:            "技能",
		Name:            "时岁轮转",
		Category:        "地",
		Tag:             "传奇-咒术-幻变",
		Description:     "冷却2.异能:双方玩家不能打出、召唤任何卡牌,不能学习、使用任何法术,不能使用卡牌攻击.此卡的学习和使用花费必须严格为地脉和奥术元素",
		Quote:           "\"你看起来,还和我记忆中一样\"",
		ElementsCost:    map[string]int{"地": 1, "无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1, "无": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        2,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3411101.jpg",
	}
}

type CardDef3411102 struct{}

func (CardDef3411102) ID() string      { return "3411102" }
func (CardDef3411102) Name() string    { return "蔽天阵 血沙" }
func (CardDef3411102) Kind() string    { return "技能" }
func (CardDef3411102) Element() string { return "地" }

func (CardDef3411102) Card() model.Card {
	return model.Card{
		Number:          "3411102",
		Type:            "技能",
		Name:            "蔽天阵 血沙",
		Category:        "地",
		Tag:             "传奇-法术-创造",
		Description:     "冷却1.范围:全场.主动回合技:双方各自选择支付自己场上单位的最多3点负载或生命值(由你先进行),在此卡上放置差值数量的标记物.每一点标记物此卡获得+3\\威和+1\\攻",
		Quote:           "遮沙蔽风了",
		ElementsCost:    map[string]int{"地": 5, "暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 4, "暗": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3411102.jpg",
	}
}

type CardDef3421001 struct{}

func (CardDef3421001) ID() string      { return "3421001" }
func (CardDef3421001) Name() string    { return "森林的庇护" }
func (CardDef3421001) Kind() string    { return "技能" }
func (CardDef3421001) Element() string { return "地" }

func (CardDef3421001) Card() model.Card {
	return model.Card{
		Number:          "3421001",
		Type:            "技能",
		Name:            "森林的庇护",
		Category:        "地",
		Tag:             "法术-神秘",
		Description:     "防御.精通1:改为4\\威;精通3:改为6\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           2,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421001.jpg",
	}
}

type CardDef3421002 struct{}

func (CardDef3421002) ID() string      { return "3421002" }
func (CardDef3421002) Name() string    { return "石化缠绕" }
func (CardDef3421002) Kind() string    { return "技能" }
func (CardDef3421002) Element() string { return "地" }

func (CardDef3421002) Card() model.Card {
	return model.Card{
		Number:          "3421002",
		Type:            "技能",
		Name:            "石化缠绕",
		Category:        "地",
		Tag:             "法术-幻变",
		Description:     "石化1.",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421002.jpg",
	}
}

type CardDef3421003 struct{}

func (CardDef3421003) ID() string      { return "3421003" }
func (CardDef3421003) Name() string    { return "裂地重击" }
func (CardDef3421003) Kind() string    { return "技能" }
func (CardDef3421003) Element() string { return "地" }

func (CardDef3421003) Card() model.Card {
	return model.Card{
		Number:          "3421003",
		Type:            "技能",
		Name:            "裂地重击",
		Category:        "地",
		Tag:             "法术-聚能",
		Description:     "精通1,3:获得+1\\威和+1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421003.jpg",
	}
}

type CardDef3421004 struct{}

func (CardDef3421004) ID() string      { return "3421004" }
func (CardDef3421004) Name() string    { return "再生之力" }
func (CardDef3421004) Kind() string    { return "技能" }
func (CardDef3421004) Element() string { return "地" }

func (CardDef3421004) Card() model.Card {
	return model.Card{
		Number:          "3421004",
		Type:            "技能",
		Name:            "再生之力",
		Category:        "地",
		Tag:             "咒术-幻变",
		Description:     "重置你的1张地脉伙伴",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421004.jpg",
	}
}

type CardDef3421005 struct{}

func (CardDef3421005) ID() string      { return "3421005" }
func (CardDef3421005) Name() string    { return "岩石壁垒" }
func (CardDef3421005) Kind() string    { return "技能" }
func (CardDef3421005) Element() string { return "地" }

func (CardDef3421005) Card() model.Card {
	return model.Card{
		Number:          "3421005",
		Type:            "技能",
		Name:            "岩石壁垒",
		Category:        "地",
		Tag:             "法术-创造",
		Description:     "防御",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421005.jpg",
	}
}

type CardDef3421006 struct{}

func (CardDef3421006) ID() string      { return "3421006" }
func (CardDef3421006) Name() string    { return "天崩地裂" }
func (CardDef3421006) Kind() string    { return "技能" }
func (CardDef3421006) Element() string { return "地" }

func (CardDef3421006) Card() model.Card {
	return model.Card{
		Number:          "3421006",
		Type:            "技能",
		Name:            "天崩地裂",
		Category:        "地",
		Tag:             "法术-驱动",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 5, "气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          3,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421006.jpg",
	}
}

type CardDef3421007 struct{}

func (CardDef3421007) ID() string      { return "3421007" }
func (CardDef3421007) Name() string    { return "大地震" }
func (CardDef3421007) Kind() string    { return "技能" }
func (CardDef3421007) Element() string { return "地" }

func (CardDef3421007) Card() model.Card {
	return model.Card{
		Number:          "3421007",
		Type:            "技能",
		Name:            "大地震",
		Category:        "地",
		Tag:             "法术-驱动",
		Description:     "范围:方阵.晕眩1",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421007.jpg",
	}
}

type CardDef3421008 struct{}

func (CardDef3421008) ID() string      { return "3421008" }
func (CardDef3421008) Name() string    { return "联合施法" }
func (CardDef3421008) Kind() string    { return "技能" }
func (CardDef3421008) Element() string { return "地" }

func (CardDef3421008) Card() model.Card {
	return model.Card{
		Number:          "3421008",
		Type:            "技能",
		Name:            "联合施法",
		Category:        "地",
		Tag:             "法术-灵媒",
		Description:     "强化其他法术时使其+1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 2, "无": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1, "无": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421008.jpg",
	}
}

type CardDef3421009 struct{}

func (CardDef3421009) ID() string      { return "3421009" }
func (CardDef3421009) Name() string    { return "惧怖之颜" }
func (CardDef3421009) Kind() string    { return "技能" }
func (CardDef3421009) Element() string { return "地" }

func (CardDef3421009) Card() model.Card {
	return model.Card{
		Number:          "3421009",
		Type:            "技能",
		Name:            "惧怖之颜",
		Category:        "地",
		Tag:             "咒术-幻变",
		Description:     "穿透.冷却1.使1个敌人石化2",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421009.jpg",
	}
}

type CardDef3421010 struct{}

func (CardDef3421010) ID() string      { return "3421010" }
func (CardDef3421010) Name() string    { return "大地穿刺" }
func (CardDef3421010) Kind() string    { return "技能" }
func (CardDef3421010) Element() string { return "地" }

func (CardDef3421010) Card() model.Card {
	return model.Card{
		Number:          "3421010",
		Type:            "技能",
		Name:            "大地穿刺",
		Category:        "地",
		Tag:             "法术-驱动",
		Description:     "穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421010.jpg",
	}
}

type CardDef3421011 struct{}

func (CardDef3421011) ID() string      { return "3421011" }
func (CardDef3421011) Name() string    { return "自然生长" }
func (CardDef3421011) Kind() string    { return "技能" }
func (CardDef3421011) Element() string { return "地" }

func (CardDef3421011) Card() model.Card {
	return model.Card{
		Number:          "3421011",
		Type:            "技能",
		Name:            "自然生长",
		Category:        "地",
		Tag:             "咒术-幻变",
		Description:     "选择你的1个横置状态且负载小于4的地脉伙伴,使其获得负载+1\\地",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421011.jpg",
	}
}

type CardDef3421012 struct{}

func (CardDef3421012) ID() string      { return "3421012" }
func (CardDef3421012) Name() string    { return "石破天惊" }
func (CardDef3421012) Kind() string    { return "技能" }
func (CardDef3421012) Element() string { return "地" }

func (CardDef3421012) Card() model.Card {
	return model.Card{
		Number:          "3421012",
		Type:            "技能",
		Name:            "石破天惊",
		Category:        "地",
		Tag:             "法术-驱动",
		Description:     "穿透.光环:你的伙伴每负载1点\\地获得+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 7},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421012.jpg",
	}
}

type CardDef3421013 struct{}

func (CardDef3421013) ID() string      { return "3421013" }
func (CardDef3421013) Name() string    { return "大地共鸣" }
func (CardDef3421013) Kind() string    { return "技能" }
func (CardDef3421013) Element() string { return "地" }

func (CardDef3421013) Card() model.Card {
	return model.Card{
		Number:          "3421013",
		Type:            "技能",
		Name:            "大地共鸣",
		Category:        "地",
		Tag:             "法术-灵媒",
		Description:     "冷却1.范围:全场.光环:你的场上每有1个负载或生命大于3的伙伴,获得+1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 8},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 5},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           9,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421013.jpg",
	}
}

type CardDef3421014 struct{}

func (CardDef3421014) ID() string      { return "3421014" }
func (CardDef3421014) Name() string    { return "千里流沙" }
func (CardDef3421014) Kind() string    { return "技能" }
func (CardDef3421014) Element() string { return "地" }

func (CardDef3421014) Card() model.Card {
	return model.Card{
		Number:          "3421014",
		Type:            "技能",
		Name:            "千里流沙",
		Category:        "地",
		Tag:             "法术-驱动",
		Description:     "范围:方阵.冷却1.若本卡攻击未命中,无需冷却",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421014.jpg",
	}
}

type CardDef3421015 struct{}

func (CardDef3421015) ID() string      { return "3421015" }
func (CardDef3421015) Name() string    { return "急袭沙暴" }
func (CardDef3421015) Kind() string    { return "技能" }
func (CardDef3421015) Element() string { return "地" }

func (CardDef3421015) Card() model.Card {
	return model.Card{
		Number:          "3421015",
		Type:            "技能",
		Name:            "急袭沙暴",
		Category:        "地",
		Tag:             "咒术-驱动",
		Description:     "冷却2.速攻.异能:双方所有原始威力小于5的法术-2\\攻-2\\威(最低为0)",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1, "气": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        2,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\地\\3421015.jpg",
	}
}

type CardDef3421101 struct{}

func (CardDef3421101) ID() string      { return "3421101" }
func (CardDef3421101) Name() string    { return "森之洞察" }
func (CardDef3421101) Kind() string    { return "技能" }
func (CardDef3421101) Element() string { return "地" }

func (CardDef3421101) Card() model.Card {
	return model.Card{
		Number:          "3421101",
		Type:            "技能",
		Name:            "森之洞察",
		Category:        "地",
		Tag:             "咒术-灵媒",
		Description:     "冷却1.你场上每有1个地脉伙伴就抽1张牌(最多5张)然后将抽牌数量的手牌洗回卡组",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3421101.jpg",
	}
}

type CardDef3421102 struct{}

func (CardDef3421102) ID() string      { return "3421102" }
func (CardDef3421102) Name() string    { return "苍岚之刃" }
func (CardDef3421102) Kind() string    { return "技能" }
func (CardDef3421102) Element() string { return "地" }

func (CardDef3421102) Card() model.Card {
	return model.Card{
		Number:          "3421102",
		Type:            "技能",
		Name:            "苍岚之刃",
		Category:        "地",
		Tag:             "法术-创造",
		Description:     "范围:前排",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4, "气": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 2, "气": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3421102.jpg",
	}
}

type CardDef3421103 struct{}

func (CardDef3421103) ID() string      { return "3421103" }
func (CardDef3421103) Name() string    { return "巨岩崩落" }
func (CardDef3421103) Kind() string    { return "技能" }
func (CardDef3421103) Element() string { return "地" }

func (CardDef3421103) Card() model.Card {
	return model.Card{
		Number:          "3421103",
		Type:            "技能",
		Name:            "巨岩崩落",
		Category:        "地",
		Tag:             "法术-驱动",
		Description:     "无法用于强化",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           8,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3421103.jpg",
	}
}

type CardDef3421104 struct{}

func (CardDef3421104) ID() string      { return "3421104" }
func (CardDef3421104) Name() string    { return "自然回响" }
func (CardDef3421104) Kind() string    { return "技能" }
func (CardDef3421104) Element() string { return "地" }

func (CardDef3421104) Card() model.Card {
	return model.Card{
		Number:          "3421104",
		Type:            "技能",
		Name:            "自然回响",
		Category:        "地",
		Tag:             "法术-灵媒",
		Description:     "主动回合技:移除友方卡牌负载的1\\地,重置此卡且下次施放+2\\威,可额外选择一个目标(可以相同)",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3421104.jpg",
	}
}

type CardDef3421105 struct{}

func (CardDef3421105) ID() string      { return "3421105" }
func (CardDef3421105) Name() string    { return "苍老之触" }
func (CardDef3421105) Kind() string    { return "技能" }
func (CardDef3421105) Element() string { return "地" }

func (CardDef3421105) Card() model.Card {
	return model.Card{
		Number:          "3421105",
		Type:            "技能",
		Name:            "苍老之触",
		Category:        "地",
		Tag:             "法术-幻变",
		Description:     "命中:如果目标为伙伴,使其失去所有负载",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3, "暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 1, "暗": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3421105.jpg",
	}
}

type CardDef3421106 struct{}

func (CardDef3421106) ID() string      { return "3421106" }
func (CardDef3421106) Name() string    { return "腐朽侵蚀" }
func (CardDef3421106) Kind() string    { return "技能" }
func (CardDef3421106) Element() string { return "地" }

func (CardDef3421106) Card() model.Card {
	return model.Card{
		Number:          "3421106",
		Type:            "技能",
		Name:            "腐朽侵蚀",
		Category:        "地",
		Tag:             "法术-幻变",
		Description:     "范围:前排.精通3,6:获得+1\\攻和+1\\威.诱发:此法术命中敌人并结算后,对所有敌方法术造成此卡\\攻数量的虚弱然后立刻触发下一次精通.",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4, "暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 3, "暗": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3421106.jpg",
	}
}

type CardDef3421107 struct{}

func (CardDef3421107) ID() string      { return "3421107" }
func (CardDef3421107) Name() string    { return "破土而出" }
func (CardDef3421107) Kind() string    { return "技能" }
func (CardDef3421107) Element() string { return "地" }

func (CardDef3421107) Card() model.Card {
	return model.Card{
		Number:          "3421107",
		Type:            "技能",
		Name:            "破土而出",
		Category:        "地",
		Tag:             "法术-驱动",
		Description:     "冷却1.范围:溅射.精通1,2:额外无视范围选择一个溅射目标.",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 4},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3421107.jpg",
	}
}

type CardDef3421108 struct{}

func (CardDef3421108) ID() string      { return "3421108" }
func (CardDef3421108) Name() string    { return "御守石阵" }
func (CardDef3421108) Kind() string    { return "技能" }
func (CardDef3421108) Element() string { return "地" }

func (CardDef3421108) Card() model.Card {
	return model.Card{
		Number:          "3421108",
		Type:            "技能",
		Name:            "御守石阵",
		Category:        "地",
		Tag:             "法术-创造",
		Description:     "范围:方阵.用于防御时此卡使用花费-1",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3421108.jpg",
	}
}

type CardDef3421109 struct{}

func (CardDef3421109) ID() string      { return "3421109" }
func (CardDef3421109) Name() string    { return "石化死光" }
func (CardDef3421109) Kind() string    { return "技能" }
func (CardDef3421109) Element() string { return "地" }

func (CardDef3421109) Card() model.Card {
	return model.Card{
		Number:          "3421109",
		Type:            "技能",
		Name:            "石化死光",
		Category:        "地",
		Tag:             "法术-聚能",
		Description:     "石化3",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3421109.jpg",
	}
}

type CardDef3421110 struct{}

func (CardDef3421110) ID() string      { return "3421110" }
func (CardDef3421110) Name() string    { return "粉碎石破" }
func (CardDef3421110) Kind() string    { return "技能" }
func (CardDef3421110) Element() string { return "地" }

func (CardDef3421110) Card() model.Card {
	return model.Card{
		Number:          "3421110",
		Type:            "技能",
		Name:            "粉碎石破",
		Category:        "地",
		Tag:             "法术-驱动",
		Description:     "如果目标\\血大于2,此卡本次+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"地": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\地\\3421110.jpg",
	}
}

type CardDef3501001 struct{}

func (CardDef3501001) ID() string      { return "3501001" }
func (CardDef3501001) Name() string    { return "团结的希望" }
func (CardDef3501001) Kind() string    { return "技能" }
func (CardDef3501001) Element() string { return "光" }

func (CardDef3501001) Card() model.Card {
	return model.Card{
		Number:          "3501001",
		Type:            "技能",
		Name:            "团结的希望",
		Category:        "光",
		Tag:             "衍生-咒术-神秘",
		Description:     "从卡组上方开始将翻开5张牌,检索其中1张光辉伙伴,之后重洗卡组",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3501001.jpg",
	}
}

type CardDef3501101 struct{}

func (CardDef3501101) ID() string      { return "3501101" }
func (CardDef3501101) Name() string    { return "五虹之束" }
func (CardDef3501101) Kind() string    { return "技能" }
func (CardDef3501101) Element() string { return "光" }

func (CardDef3501101) Card() model.Card {
	return model.Card{
		Number:          "3501101",
		Type:            "技能",
		Name:            "五虹之束",
		Category:        "光",
		Tag:             "衍生-法术-聚能",
		Description:     "使用时移除五虹之环上任意数量标记物,根据种类获得以下效果:\\火:+2\\攻;\\水:额外选择1个目标;\\地:+3\\威;\\气:获得穿透;\\光:使用花费-2;全部:额外使\\威翻倍",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3501101.jpg",
	}
}

type CardDef3511010 struct{}

func (CardDef3511010) ID() string      { return "3511010" }
func (CardDef3511010) Name() string    { return "破晓" }
func (CardDef3511010) Kind() string    { return "技能" }
func (CardDef3511010) Element() string { return "光" }

func (CardDef3511010) Card() model.Card {
	return model.Card{
		Number:          "3511010",
		Type:            "技能",
		Name:            "破晓",
		Category:        "光",
		Tag:             "传奇-法术-驱动",
		Description:     "如果攻击目标为敌方伙伴,将同时命中所有与之属性相同的敌人.诱发:你每召唤1个负载有光的伙伴此卡获得永久+1\\威.此卡仅当\\威大于8时才能用于攻击",
		Quote:           "\"结束了?我们...胜利了吗?\"",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3511010.jpg",
	}
}

type CardDef3511101 struct{}

func (CardDef3511101) ID() string      { return "3511101" }
func (CardDef3511101) Name() string    { return "神辉驭空" }
func (CardDef3511101) Kind() string    { return "技能" }
func (CardDef3511101) Element() string { return "光" }

func (CardDef3511101) Card() model.Card {
	return model.Card{
		Number:          "3511101",
		Type:            "技能",
		Name:            "神辉驭空",
		Category:        "光",
		Tag:             "传奇-法术-神秘",
		Description:     "范围:前排或纵列.光环:对方每有1张手牌,此卡+1\\威.命中:可以选择1个玩家将全部手牌丢弃然后抽牌至手牌上限.",
		Quote:           "金色的光芒散落在九霄议庭,仿佛他们真的是被选中的",
		ElementsCost:    map[string]int{"光": 4, "气": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 3, "气": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3511101.jpg",
	}
}

type CardDef3511102 struct{}

func (CardDef3511102) ID() string      { return "3511102" }
func (CardDef3511102) Name() string    { return "绝境之光 孤星闪耀" }
func (CardDef3511102) Kind() string    { return "技能" }
func (CardDef3511102) Element() string { return "光" }

func (CardDef3511102) Card() model.Card {
	return model.Card{
		Number:          "3511102",
		Type:            "技能",
		Name:            "绝境之光 孤星闪耀",
		Category:        "光",
		Tag:             "传奇-法术-聚能",
		Description:     "穿透.范围:溅射.仅在你场上单位比对方少时才能使用.光环:此卡\\威上升你场上光辉伙伴生命值和负载最高者的合计数值",
		Quote:           "这一次,再也没有值得信任的伙伴了",
		ElementsCost:    map[string]int{"光": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 4},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3511102.jpg",
	}
}

type CardDef3521001 struct{}

func (CardDef3521001) ID() string      { return "3521001" }
func (CardDef3521001) Name() string    { return "治疗术" }
func (CardDef3521001) Kind() string    { return "技能" }
func (CardDef3521001) Element() string { return "光" }

func (CardDef3521001) Card() model.Card {
	return model.Card{
		Number:          "3521001",
		Type:            "技能",
		Name:            "治疗术",
		Category:        "光",
		Tag:             "咒术-聚能",
		Description:     "回复2\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521001.jpg",
	}
}

type CardDef3521002 struct{}

func (CardDef3521002) ID() string      { return "3521002" }
func (CardDef3521002) Name() string    { return "神圣之火" }
func (CardDef3521002) Kind() string    { return "技能" }
func (CardDef3521002) Element() string { return "光" }

func (CardDef3521002) Card() model.Card {
	return model.Card{
		Number:          "3521002",
		Type:            "技能",
		Name:            "神圣之火",
		Category:        "光",
		Tag:             "法术-创造",
		Description:     "此卡对友方单位不造成伤害,改为移除所有负面状态",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521002.jpg",
	}
}

type CardDef3521003 struct{}

func (CardDef3521003) ID() string      { return "3521003" }
func (CardDef3521003) Name() string    { return "神圣防护罩" }
func (CardDef3521003) Kind() string    { return "技能" }
func (CardDef3521003) Element() string { return "光" }

func (CardDef3521003) Card() model.Card {
	return model.Card{
		Number:          "3521003",
		Type:            "技能",
		Name:            "神圣防护罩",
		Category:        "光",
		Tag:             "法术-创造",
		Description:     "防御",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521003.jpg",
	}
}

type CardDef3521004 struct{}

func (CardDef3521004) ID() string      { return "3521004" }
func (CardDef3521004) Name() string    { return "闪光魔术" }
func (CardDef3521004) Kind() string    { return "技能" }
func (CardDef3521004) Element() string { return "光" }

func (CardDef3521004) Card() model.Card {
	return model.Card{
		Number:          "3521004",
		Type:            "技能",
		Name:            "闪光魔术",
		Category:        "光",
		Tag:             "法术-幻变",
		Description:     "晕眩1",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521004.jpg",
	}
}

type CardDef3521005 struct{}

func (CardDef3521005) ID() string      { return "3521005" }
func (CardDef3521005) Name() string    { return "星陨" }
func (CardDef3521005) Kind() string    { return "技能" }
func (CardDef3521005) Element() string { return "光" }

func (CardDef3521005) Card() model.Card {
	return model.Card{
		Number:          "3521005",
		Type:            "技能",
		Name:            "星陨",
		Category:        "光",
		Tag:             "法术-驱动",
		Description:     "穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3, "地": 1, "气": 1, "水": 1, "火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2, "无": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          3,
		Life:            -1,
		Duration:        -1,
		Power:           8,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521005.jpg",
	}
}

type CardDef3521006 struct{}

func (CardDef3521006) ID() string      { return "3521006" }
func (CardDef3521006) Name() string    { return "光辉斩裂" }
func (CardDef3521006) Kind() string    { return "技能" }
func (CardDef3521006) Element() string { return "光" }

func (CardDef3521006) Card() model.Card {
	return model.Card{
		Number:          "3521006",
		Type:            "技能",
		Name:            "光辉斩裂",
		Category:        "光",
		Tag:             "法术-聚能",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521006.jpg",
	}
}

type CardDef3521007 struct{}

func (CardDef3521007) ID() string      { return "3521007" }
func (CardDef3521007) Name() string    { return "希望呼唤" }
func (CardDef3521007) Kind() string    { return "技能" }
func (CardDef3521007) Element() string { return "光" }

func (CardDef3521007) Card() model.Card {
	return model.Card{
		Number:          "3521007",
		Type:            "技能",
		Name:            "希望呼唤",
		Category:        "光",
		Tag:             "咒术-神秘",
		Description:     "从卡组上方开始将翻到的第1张光辉伙伴抽取,之后重洗卡组",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521007.jpg",
	}
}

type CardDef3521008 struct{}

func (CardDef3521008) ID() string      { return "3521008" }
func (CardDef3521008) Name() string    { return "光辉波动" }
func (CardDef3521008) Kind() string    { return "技能" }
func (CardDef3521008) Element() string { return "光" }

func (CardDef3521008) Card() model.Card {
	return model.Card{
		Number:          "3521008",
		Type:            "技能",
		Name:            "光辉波动",
		Category:        "光",
		Tag:             "法术-聚能",
		Description:     "范围:前排.晕眩1",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521008.jpg",
	}
}

type CardDef3521009 struct{}

func (CardDef3521009) ID() string      { return "3521009" }
func (CardDef3521009) Name() string    { return "幻彩流光" }
func (CardDef3521009) Kind() string    { return "技能" }
func (CardDef3521009) Element() string { return "光" }

func (CardDef3521009) Card() model.Card {
	return model.Card{
		Number:          "3521009",
		Type:            "技能",
		Name:            "幻彩流光",
		Category:        "光",
		Tag:             "法术-幻变",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "地": 1, "气": 1, "水": 1, "火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521009.jpg",
	}
}

type CardDef3521011 struct{}

func (CardDef3521011) ID() string      { return "3521011" }
func (CardDef3521011) Name() string    { return "光之庇护" }
func (CardDef3521011) Kind() string    { return "技能" }
func (CardDef3521011) Element() string { return "光" }

func (CardDef3521011) Card() model.Card {
	return model.Card{
		Number:          "3521011",
		Type:            "技能",
		Name:            "光之庇护",
		Category:        "光",
		Tag:             "咒术-神秘",
		Description:     "速攻.冷却2.选择1个伙伴,直到下个回合结束防止所有致命伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521011.jpg",
	}
}

type CardDef3521012 struct{}

func (CardDef3521012) ID() string      { return "3521012" }
func (CardDef3521012) Name() string    { return "长虹贯日" }
func (CardDef3521012) Kind() string    { return "技能" }
func (CardDef3521012) Element() string { return "光" }

func (CardDef3521012) Card() model.Card {
	return model.Card{
		Number:          "3521012",
		Type:            "技能",
		Name:            "长虹贯日",
		Category:        "光",
		Tag:             "法术-驱动",
		Description:     "范围:纵列",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2, "地": 1, "气": 1, "水": 1, "火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 1, "地": 1, "气": 1, "水": 1, "火": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521012.jpg",
	}
}

type CardDef3521013 struct{}

func (CardDef3521013) ID() string      { return "3521013" }
func (CardDef3521013) Name() string    { return "月之辉" }
func (CardDef3521013) Kind() string    { return "技能" }
func (CardDef3521013) Element() string { return "光" }

func (CardDef3521013) Card() model.Card {
	return model.Card{
		Number:          "3521013",
		Type:            "技能",
		Name:            "月之辉",
		Category:        "光",
		Tag:             "法术-神秘",
		Description:     "用于防御或强化防御时+2威力",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521013.jpg",
	}
}

type CardDef3521014 struct{}

func (CardDef3521014) ID() string      { return "3521014" }
func (CardDef3521014) Name() string    { return "光之祝福" }
func (CardDef3521014) Kind() string    { return "技能" }
func (CardDef3521014) Element() string { return "光" }

func (CardDef3521014) Card() model.Card {
	return model.Card{
		Number:          "3521014",
		Type:            "技能",
		Name:            "光之祝福",
		Category:        "光",
		Tag:             "咒术-神秘",
		Description:     "冷却1.使1个友方伙伴获得+1\\血和负载+1\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521014.jpg",
	}
}

type CardDef3521015 struct{}

func (CardDef3521015) ID() string      { return "3521015" }
func (CardDef3521015) Name() string    { return "寂灭之光" }
func (CardDef3521015) Kind() string    { return "技能" }
func (CardDef3521015) Element() string { return "光" }

func (CardDef3521015) Card() model.Card {
	return model.Card{
		Number:          "3521015",
		Type:            "技能",
		Name:            "寂灭之光",
		Category:        "光",
		Tag:             "法术-聚能",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 7},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 4},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          3,
		Life:            -1,
		Duration:        -1,
		Power:           9,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\光\\3521015.jpg",
	}
}

type CardDef3521101 struct{}

func (CardDef3521101) ID() string      { return "3521101" }
func (CardDef3521101) Name() string    { return "福音" }
func (CardDef3521101) Kind() string    { return "技能" }
func (CardDef3521101) Element() string { return "光" }

func (CardDef3521101) Card() model.Card {
	return model.Card{
		Number:          "3521101",
		Type:            "技能",
		Name:            "福音",
		Category:        "光",
		Tag:             "法术-神秘",
		Description:     "诱发:每当你的一个光辉伙伴被消耗,本卡使用费用-1\\光,直到你的回合结束",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 9},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           8,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3521101.jpg",
	}
}

type CardDef3521102 struct{}

func (CardDef3521102) ID() string      { return "3521102" }
func (CardDef3521102) Name() string    { return "神助" }
func (CardDef3521102) Kind() string    { return "技能" }
func (CardDef3521102) Element() string { return "光" }

func (CardDef3521102) Card() model.Card {
	return model.Card{
		Number:          "3521102",
		Type:            "技能",
		Name:            "神助",
		Category:        "光",
		Tag:             "法术-神秘",
		Description:     "光环:此卡用于强化神秘法术时+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3521102.jpg",
	}
}

type CardDef3521103 struct{}

func (CardDef3521103) ID() string      { return "3521103" }
func (CardDef3521103) Name() string    { return "光铸飞弹" }
func (CardDef3521103) Kind() string    { return "技能" }
func (CardDef3521103) Element() string { return "光" }

func (CardDef3521103) Card() model.Card {
	return model.Card{
		Number:          "3521103",
		Type:            "技能",
		Name:            "光铸飞弹",
		Category:        "光",
		Tag:             "法术-创造",
		Description:     "穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3521103.jpg",
	}
}

type CardDef3521104 struct{}

func (CardDef3521104) ID() string      { return "3521104" }
func (CardDef3521104) Name() string    { return "罪责" }
func (CardDef3521104) Kind() string    { return "技能" }
func (CardDef3521104) Element() string { return "光" }

func (CardDef3521104) Card() model.Card {
	return model.Card{
		Number:          "3521104",
		Type:            "技能",
		Name:            "罪责",
		Category:        "光",
		Tag:             "法术-神秘",
		Description:     "入场:选择1个伙伴种类,此卡仅在攻击该种类伙伴时获得穿透和+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3521104.jpg",
	}
}

type CardDef3521105 struct{}

func (CardDef3521105) ID() string      { return "3521105" }
func (CardDef3521105) Name() string    { return "流光之束" }
func (CardDef3521105) Kind() string    { return "技能" }
func (CardDef3521105) Element() string { return "光" }

func (CardDef3521105) Card() model.Card {
	return model.Card{
		Number:          "3521105",
		Type:            "技能",
		Name:            "流光之束",
		Category:        "光",
		Tag:             "法术-聚能",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 1, "气": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3521105.jpg",
	}
}

type CardDef3521106 struct{}

func (CardDef3521106) ID() string      { return "3521106" }
func (CardDef3521106) Name() string    { return "归心" }
func (CardDef3521106) Kind() string    { return "技能" }
func (CardDef3521106) Element() string { return "光" }

func (CardDef3521106) Card() model.Card {
	return model.Card{
		Number:          "3521106",
		Type:            "技能",
		Name:            "归心",
		Category:        "光",
		Tag:             "法术-神秘",
		Description:     "范围:方阵.光环:你的场上每有一个光辉伙伴获得+1\\威,每有一个其他属性伙伴-1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 4},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           2,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3521106.jpg",
	}
}

type CardDef3521107 struct{}

func (CardDef3521107) ID() string      { return "3521107" }
func (CardDef3521107) Name() string    { return "虹彩之壁" }
func (CardDef3521107) Kind() string    { return "技能" }
func (CardDef3521107) Element() string { return "光" }

func (CardDef3521107) Card() model.Card {
	return model.Card{
		Number:          "3521107",
		Type:            "技能",
		Name:            "虹彩之壁",
		Category:        "光",
		Tag:             "法术-创造",
		Description:     "防御",
		Quote:           "",
		ElementsCost:    map[string]int{"地": 1, "气": 1, "水": 1, "火": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3521107.jpg",
	}
}

type CardDef3521108 struct{}

func (CardDef3521108) ID() string      { return "3521108" }
func (CardDef3521108) Name() string    { return "恩典" }
func (CardDef3521108) Kind() string    { return "技能" }
func (CardDef3521108) Element() string { return "光" }

func (CardDef3521108) Card() model.Card {
	return model.Card{
		Number:          "3521108",
		Type:            "技能",
		Name:            "恩典",
		Category:        "光",
		Tag:             "咒术-神秘",
		Description:     "冷却1.使1个受伤的友方伙伴回复2\\血,如果使生命值回满则使其+1\\血,负载+1\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3521108.jpg",
	}
}

type CardDef3521109 struct{}

func (CardDef3521109) ID() string      { return "3521109" }
func (CardDef3521109) Name() string    { return "重整旗鼓" }
func (CardDef3521109) Kind() string    { return "技能" }
func (CardDef3521109) Element() string { return "光" }

func (CardDef3521109) Card() model.Card {
	return model.Card{
		Number:          "3521109",
		Type:            "技能",
		Name:            "重整旗鼓",
		Category:        "光",
		Tag:             "咒术-幻变",
		Description:     "诱发:当敌方法术命中后才能使用此卡,使1个友方伙伴负载+1\\光+1\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3521109.jpg",
	}
}

type CardDef3521110 struct{}

func (CardDef3521110) ID() string      { return "3521110" }
func (CardDef3521110) Name() string    { return "光灵汲取" }
func (CardDef3521110) Kind() string    { return "技能" }
func (CardDef3521110) Element() string { return "光" }

func (CardDef3521110) Card() model.Card {
	return model.Card{
		Number:          "3521110",
		Type:            "技能",
		Name:            "光灵汲取",
		Category:        "光",
		Tag:             "法术-灵媒",
		Description:     "命中:使1个友方光辉伙伴获得负载+1\\光",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\光\\3521110.jpg",
	}
}

type CardDef3601101 struct{}

func (CardDef3601101) ID() string      { return "3601101" }
func (CardDef3601101) Name() string    { return "鲜血盛宴" }
func (CardDef3601101) Kind() string    { return "技能" }
func (CardDef3601101) Element() string { return "暗" }

func (CardDef3601101) Card() model.Card {
	return model.Card{
		Number:          "3601101",
		Type:            "技能",
		Name:            "鲜血盛宴",
		Category:        "暗",
		Tag:             "衍生-法术-代赎",
		Description:     "只能用于攻击友方单位.命中:获得2\\暗,或使你的人物回复1\\血.主动:你可以花费1\\暗将此卡变为人物的绑定技能",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3601101.jpg",
	}
}

type CardDef3611101 struct{}

func (CardDef3611101) ID() string      { return "3611101" }
func (CardDef3611101) Name() string    { return "红月" }
func (CardDef3611101) Kind() string    { return "技能" }
func (CardDef3611101) Element() string { return "暗" }

func (CardDef3611101) Card() model.Card {
	return model.Card{
		Number:          "3611101",
		Type:            "技能",
		Name:            "红月",
		Category:        "暗",
		Tag:             "传奇-咒术-神秘",
		Description:     "冷却2.异能:使你的暗影法术+2\\威",
		Quote:           "酒杯中的黑血平静如镜,倒映出这个诡谲的世界:低头的安迪斯与瑟薇安娜,花园的荆棘和松树,深邃的夜空,以及夜空中倒悬的红月",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3611101.jpg",
	}
}

type CardDef3611102 struct{}

func (CardDef3611102) ID() string      { return "3611102" }
func (CardDef3611102) Name() string    { return "厄瑞波斯之爪" }
func (CardDef3611102) Kind() string    { return "技能" }
func (CardDef3611102) Element() string { return "暗" }

func (CardDef3611102) Card() model.Card {
	return model.Card{
		Number:          "3611102",
		Type:            "技能",
		Name:            "厄瑞波斯之爪",
		Category:        "暗",
		Tag:             "传奇-法术-神秘",
		Description:     "敌方有三个及以上虚弱法术才能学习.光环:敌方法术每有1层虚弱此卡+1\\威.诱发:每次使用此卡后使最多3个不同的敌方法术虚弱1",
		Quote:           "嘘...它已经来了",
		ElementsCost:    map[string]int{"暗": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3611102.jpg",
	}
}

type CardDef3621001 struct{}

func (CardDef3621001) ID() string      { return "3621001" }
func (CardDef3621001) Name() string    { return "暗影冲击" }
func (CardDef3621001) Kind() string    { return "技能" }
func (CardDef3621001) Element() string { return "暗" }

func (CardDef3621001) Card() model.Card {
	return model.Card{
		Number:          "3621001",
		Type:            "技能",
		Name:            "暗影冲击",
		Category:        "暗",
		Tag:             "法术-聚能",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 1},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621001.jpg",
	}
}

type CardDef3621002 struct{}

func (CardDef3621002) ID() string      { return "3621002" }
func (CardDef3621002) Name() string    { return "噬血" }
func (CardDef3621002) Kind() string    { return "技能" }
func (CardDef3621002) Element() string { return "暗" }

func (CardDef3621002) Card() model.Card {
	return model.Card{
		Number:          "3621002",
		Type:            "技能",
		Name:            "噬血",
		Category:        "暗",
		Tag:             "法术-代赎",
		Description:     "命中:使1个友方单位+2\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621002.jpg",
	}
}

type CardDef3621003 struct{}

func (CardDef3621003) ID() string      { return "3621003" }
func (CardDef3621003) Name() string    { return "死亡收割" }
func (CardDef3621003) Kind() string    { return "技能" }
func (CardDef3621003) Element() string { return "暗" }

func (CardDef3621003) Card() model.Card {
	return model.Card{
		Number:          "3621003",
		Type:            "技能",
		Name:            "死亡收割",
		Category:        "暗",
		Tag:             "法术-驱动",
		Description:     "范围:前排",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621003.jpg",
	}
}

type CardDef3621004 struct{}

func (CardDef3621004) ID() string      { return "3621004" }
func (CardDef3621004) Name() string    { return "暗影箭" }
func (CardDef3621004) Kind() string    { return "技能" }
func (CardDef3621004) Element() string { return "暗" }

func (CardDef3621004) Card() model.Card {
	return model.Card{
		Number:          "3621004",
		Type:            "技能",
		Name:            "暗影箭",
		Category:        "暗",
		Tag:             "法术-创造",
		Description:     "穿透",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621004.jpg",
	}
}

type CardDef3621005 struct{}

func (CardDef3621005) ID() string      { return "3621005" }
func (CardDef3621005) Name() string    { return "暗冥弹" }
func (CardDef3621005) Kind() string    { return "技能" }
func (CardDef3621005) Element() string { return "暗" }

func (CardDef3621005) Card() model.Card {
	return model.Card{
		Number:          "3621005",
		Type:            "技能",
		Name:            "暗冥弹",
		Category:        "暗",
		Tag:             "法术-创造",
		Description:     "",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 3},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          3,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621005.jpg",
	}
}

type CardDef3621006 struct{}

func (CardDef3621006) ID() string      { return "3621006" }
func (CardDef3621006) Name() string    { return "死魂之噬" }
func (CardDef3621006) Kind() string    { return "技能" }
func (CardDef3621006) Element() string { return "暗" }

func (CardDef3621006) Card() model.Card {
	return model.Card{
		Number:          "3621006",
		Type:            "技能",
		Name:            "死魂之噬",
		Category:        "暗",
		Tag:             "法术-灵媒",
		Description:     "命中:将3层虚弱分配给敌方法术",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621006.jpg",
	}
}

type CardDef3621007 struct{}

func (CardDef3621007) ID() string      { return "3621007" }
func (CardDef3621007) Name() string    { return "安迪斯的惩罚" }
func (CardDef3621007) Kind() string    { return "技能" }
func (CardDef3621007) Element() string { return "暗" }

func (CardDef3621007) Card() model.Card {
	return model.Card{
		Number:          "3621007",
		Type:            "技能",
		Name:            "安迪斯的惩罚",
		Category:        "暗",
		Tag:             "法术-神秘",
		Description:     "诱发:每当友方单位受到1点伤害时,下一次此技能获得+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621007.jpg",
	}
}

type CardDef3621008 struct{}

func (CardDef3621008) ID() string      { return "3621008" }
func (CardDef3621008) Name() string    { return "亡者之怒" }
func (CardDef3621008) Kind() string    { return "技能" }
func (CardDef3621008) Element() string { return "暗" }

func (CardDef3621008) Card() model.Card {
	return model.Card{
		Number:          "3621008",
		Type:            "技能",
		Name:            "亡者之怒",
		Category:        "暗",
		Tag:             "法术-聚能",
		Description:     "诱发:每当1个伙伴死亡后,此法术获得永久+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621008.jpg",
	}
}

type CardDef3621009 struct{}

func (CardDef3621009) ID() string      { return "3621009" }
func (CardDef3621009) Name() string    { return "虚弱诅咒" }
func (CardDef3621009) Kind() string    { return "技能" }
func (CardDef3621009) Element() string { return "暗" }

func (CardDef3621009) Card() model.Card {
	return model.Card{
		Number:          "3621009",
		Type:            "技能",
		Name:            "虚弱诅咒",
		Category:        "暗",
		Tag:             "咒术-灵媒",
		Description:     "速攻.使1个敌方法术虚弱2",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621009.jpg",
	}
}

type CardDef3621010 struct{}

func (CardDef3621010) ID() string      { return "3621010" }
func (CardDef3621010) Name() string    { return "血魔爆" }
func (CardDef3621010) Kind() string    { return "技能" }
func (CardDef3621010) Element() string { return "暗" }

func (CardDef3621010) Card() model.Card {
	return model.Card{
		Number:          "3621010",
		Type:            "技能",
		Name:            "血魔爆",
		Category:        "暗",
		Tag:             "咒术-代赎",
		Description:     "冷却1.献祭你的1个前排暗影伙伴才能发动此卡,对法力范围内1个敌人造成该伙伴生命值的伤害",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621010.jpg",
	}
}

type CardDef3621011 struct{}

func (CardDef3621011) ID() string      { return "3621011" }
func (CardDef3621011) Name() string    { return "次元爆诞" }
func (CardDef3621011) Kind() string    { return "技能" }
func (CardDef3621011) Element() string { return "暗" }

func (CardDef3621011) Card() model.Card {
	return model.Card{
		Number:          "3621011",
		Type:            "技能",
		Name:            "次元爆诞",
		Category:        "暗",
		Tag:             "法术-聚能",
		Description:     "穿透.范围:方阵",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "暗": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 1, "暗": 4},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621011.jpg",
	}
}

type CardDef3621012 struct{}

func (CardDef3621012) ID() string      { return "3621012" }
func (CardDef3621012) Name() string    { return "回魂术" }
func (CardDef3621012) Kind() string    { return "技能" }
func (CardDef3621012) Element() string { return "暗" }

func (CardDef3621012) Card() model.Card {
	return model.Card{
		Number:          "3621012",
		Type:            "技能",
		Name:            "回魂术",
		Category:        "暗",
		Tag:             "咒术-幻变",
		Description:     "冷却1.从你的弃牌堆将最多2个伙伴移回手牌",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 1},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621012.jpg",
	}
}

type CardDef3621013 struct{}

func (CardDef3621013) ID() string      { return "3621013" }
func (CardDef3621013) Name() string    { return "亡灵护壁" }
func (CardDef3621013) Kind() string    { return "技能" }
func (CardDef3621013) Element() string { return "暗" }

func (CardDef3621013) Card() model.Card {
	return model.Card{
		Number:          "3621013",
		Type:            "技能",
		Name:            "亡灵护壁",
		Category:        "暗",
		Tag:             "法术-驱动",
		Description:     "防御.光环:如果当回合或上个回合有友方单位死亡,此法术+2\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621013.jpg",
	}
}

type CardDef3621014 struct{}

func (CardDef3621014) ID() string      { return "3621014" }
func (CardDef3621014) Name() string    { return "业障" }
func (CardDef3621014) Kind() string    { return "技能" }
func (CardDef3621014) Element() string { return "暗" }

func (CardDef3621014) Card() model.Card {
	return model.Card{
		Number:          "3621014",
		Type:            "技能",
		Name:            "业障",
		Category:        "暗",
		Tag:             "法术-幻变",
		Description:     "防御.若防御成功,使敌方攻击法术虚弱2",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          0,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621014.jpg",
	}
}

type CardDef3621015 struct{}

func (CardDef3621015) ID() string      { return "3621015" }
func (CardDef3621015) Name() string    { return "虹吸" }
func (CardDef3621015) Kind() string    { return "技能" }
func (CardDef3621015) Element() string { return "暗" }

func (CardDef3621015) Card() model.Card {
	return model.Card{
		Number:          "3621015",
		Type:            "技能",
		Name:            "虹吸",
		Category:        "暗",
		Tag:             "咒术-聚能",
		Description:     "冷却2.诱发:当敌方法术命中时可以使用此卡,将即将造成的伤害改为对目标回复生命值",
		Quote:           "",
		ElementsCost:    map[string]int{"光": 1, "暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"光": 1, "暗": 2},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\技能\\暗\\3621015.jpg",
	}
}

type CardDef3621101 struct{}

func (CardDef3621101) ID() string      { return "3621101" }
func (CardDef3621101) Name() string    { return "歃血" }
func (CardDef3621101) Kind() string    { return "技能" }
func (CardDef3621101) Element() string { return "暗" }

func (CardDef3621101) Card() model.Card {
	return model.Card{
		Number:          "3621101",
		Type:            "技能",
		Name:            "歃血",
		Category:        "暗",
		Tag:             "法术-代赎",
		Description:     "诱发:如果此法术对友方单位造成伤害,获得2\\暗并使下一次获得+2\\威和+1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 1},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           2,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3621101.jpg",
	}
}

type CardDef3621102 struct{}

func (CardDef3621102) ID() string      { return "3621102" }
func (CardDef3621102) Name() string    { return "报应" }
func (CardDef3621102) Kind() string    { return "技能" }
func (CardDef3621102) Element() string { return "暗" }

func (CardDef3621102) Card() model.Card {
	return model.Card{
		Number:          "3621102",
		Type:            "技能",
		Name:            "报应",
		Category:        "暗",
		Tag:             "法术-代赎",
		Description:     "光环:本回合以及上个回合你的人物每受到1点伤害,此技能本次+1\\攻",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3621102.jpg",
	}
}

type CardDef3621103 struct{}

func (CardDef3621103) ID() string      { return "3621103" }
func (CardDef3621103) Name() string    { return "血魂斩" }
func (CardDef3621103) Kind() string    { return "技能" }
func (CardDef3621103) Element() string { return "暗" }

func (CardDef3621103) Card() model.Card {
	return model.Card{
		Number:          "3621103",
		Type:            "技能",
		Name:            "血魂斩",
		Category:        "暗",
		Tag:             "法术-代赎",
		Description:     "诱发:此卡用于攻击时,对你的人物造成1点伤害.命中:为你的人物回复2\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           6,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3621103.jpg",
	}
}

type CardDef3621104 struct{}

func (CardDef3621104) ID() string      { return "3621104" }
func (CardDef3621104) Name() string    { return "血蔷薇咒印" }
func (CardDef3621104) Kind() string    { return "技能" }
func (CardDef3621104) Element() string { return "暗" }

func (CardDef3621104) Card() model.Card {
	return model.Card{
		Number:          "3621104",
		Type:            "技能",
		Name:            "血蔷薇咒印",
		Category:        "暗",
		Tag:             "法术-灵媒",
		Description:     "入场:标记场上任意1个敌方单位.诱发:若该单位在你的下个回合结束前死亡,此卡变为你的人物绑定技能且使用花费-1",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 3},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3621104.jpg",
	}
}

type CardDef3621105 struct{}

func (CardDef3621105) ID() string      { return "3621105" }
func (CardDef3621105) Name() string    { return "恐吓" }
func (CardDef3621105) Kind() string    { return "技能" }
func (CardDef3621105) Element() string { return "暗" }

func (CardDef3621105) Card() model.Card {
	return model.Card{
		Number:          "3621105",
		Type:            "技能",
		Name:            "恐吓",
		Category:        "暗",
		Tag:             "法术-灵媒",
		Description:     "光环:此卡\\威和\\攻+X,X为敌方具有虚弱的法术数量且最多为2",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           3,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3621105.jpg",
	}
}

type CardDef3621106 struct{}

func (CardDef3621106) ID() string      { return "3621106" }
func (CardDef3621106) Name() string    { return "红月吞噬" }
func (CardDef3621106) Kind() string    { return "技能" }
func (CardDef3621106) Element() string { return "暗" }

func (CardDef3621106) Card() model.Card {
	return model.Card{
		Number:          "3621106",
		Type:            "技能",
		Name:            "红月吞噬",
		Category:        "暗",
		Tag:             "法术-神秘",
		Description:     "命中:如果目标为伙伴则将其消灭,如果红月生效可以使1个友方暗影单位获得被消灭伙伴剩余生命值的\\血",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 5},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           7,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3621106.jpg",
	}
}

type CardDef3621107 struct{}

func (CardDef3621107) ID() string      { return "3621107" }
func (CardDef3621107) Name() string    { return "意志侵蚀" }
func (CardDef3621107) Kind() string    { return "技能" }
func (CardDef3621107) Element() string { return "暗" }

func (CardDef3621107) Card() model.Card {
	return model.Card{
		Number:          "3621107",
		Type:            "技能",
		Name:            "意志侵蚀",
		Category:        "暗",
		Tag:             "法术-灵媒",
		Description:     "可额外选择1个目标(可以相同).光环:红月生效期间此卡获得穿透和+1\\威",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3621107.jpg",
	}
}

type CardDef3621108 struct{}

func (CardDef3621108) ID() string      { return "3621108" }
func (CardDef3621108) Name() string    { return "月影" }
func (CardDef3621108) Kind() string    { return "技能" }
func (CardDef3621108) Element() string { return "暗" }

func (CardDef3621108) Card() model.Card {
	return model.Card{
		Number:          "3621108",
		Type:            "技能",
		Name:            "月影",
		Category:        "暗",
		Tag:             "法术-神秘",
		Description:     "诱发回合技:当你的法术与虚弱法术或被虚弱法术强化的法术战斗后,可以重置此卡",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 4},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 2},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          2,
		Life:            -1,
		Duration:        -1,
		Power:           4,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3621108.jpg",
	}
}

type CardDef3621109 struct{}

func (CardDef3621109) ID() string      { return "3621109" }
func (CardDef3621109) Name() string    { return "裂魂尖啸" }
func (CardDef3621109) Kind() string    { return "技能" }
func (CardDef3621109) Element() string { return "暗" }

func (CardDef3621109) Card() model.Card {
	return model.Card{
		Number:          "3621109",
		Type:            "技能",
		Name:            "裂魂尖啸",
		Category:        "暗",
		Tag:             "法术-灵媒",
		Description:     "晕眩1.范围:溅射.诱发:此卡作为主法术战斗后,使战斗的敌方法术(包括强化)获得虚弱1",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 6},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{"暗": 3},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          1,
		Life:            -1,
		Duration:        -1,
		Power:           5,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3621109.jpg",
	}
}

type CardDef3621110 struct{}

func (CardDef3621110) ID() string      { return "3621110" }
func (CardDef3621110) Name() string    { return "鲜血滋养" }
func (CardDef3621110) Kind() string    { return "技能" }
func (CardDef3621110) Element() string { return "暗" }

func (CardDef3621110) Card() model.Card {
	return model.Card{
		Number:          "3621110",
		Type:            "技能",
		Name:            "鲜血滋养",
		Category:        "暗",
		Tag:             "咒术-代赎",
		Description:     "将弃牌堆1张暗影卡牌移出游戏,获得2\\暗",
		Quote:           "",
		ElementsCost:    map[string]int{"暗": 2},
		ElementsGain:    map[string]int{},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            -1,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\技能\\暗\\3621110.jpg",
	}
}

type CardDef4011001 struct{}

func (CardDef4011001) ID() string      { return "4011001" }
func (CardDef4011001) Name() string    { return "\"南境百灵\" 斯卡尔蒂 罗佳" }
func (CardDef4011001) Kind() string    { return "人物" }
func (CardDef4011001) Element() string { return "无" }

func (CardDef4011001) Card() model.Card {
	return model.Card{
		Number:          "4011001",
		Type:            "人物",
		Name:            "\"南境百灵\" 斯卡尔蒂 罗佳",
		Category:        "无",
		Tag:             "",
		Description:     "主动回合技:丢弃1张手牌才能发动,获得2点所丢弃卡牌属性种类的元素,这个效果对于每个奥术以外的属性一局只能使用1次",
		Quote:           "发表获奖感言时,这位闪耀的歌星热泪盈眶:\"尽管我未能成功毕业,但我要特别感谢我曾经的地脉法术教师耶伦尔,在所有人不看好我的情况下,他却说我是一个发展均衡,充满潜力的孩子\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"无": 2},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\无\\4011001.jpg",
	}
}

type CardDef4011002 struct{}

func (CardDef4011002) ID() string      { return "4011002" }
func (CardDef4011002) Name() string    { return "\"无面\"" }
func (CardDef4011002) Kind() string    { return "人物" }
func (CardDef4011002) Element() string { return "无" }

func (CardDef4011002) Card() model.Card {
	return model.Card{
		Number:          "4011002",
		Type:            "人物",
		Name:            "\"无面\"",
		Category:        "无",
		Tag:             "",
		Description:     "诱发:当你打出或学习1张与你场上原有卡牌属性相同的卡牌后,你受到1点伤害",
		Quote:           "凡有所相,皆是虚妄",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"无": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\无\\4011002.jpg",
	}
}

type CardDef4011101 struct{}

func (CardDef4011101) ID() string      { return "4011101" }
func (CardDef4011101) Name() string    { return "纯净灵体 奥希斯" }
func (CardDef4011101) Kind() string    { return "人物" }
func (CardDef4011101) Element() string { return "无" }

func (CardDef4011101) Card() model.Card {
	return model.Card{
		Number:          "4011101",
		Type:            "人物",
		Name:            "纯净灵体 奥希斯",
		Category:        "无",
		Tag:             "",
		Description:     "诱发:当你使1张奥术以外的卡牌入场时,使你的所有法术获得虚弱2",
		Quote:           "\"你觉得你创造了我?我倒是觉得恰恰相反\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"无": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\无\\4011101.jpg",
	}
}

type CardDef4011102 struct{}

func (CardDef4011102) ID() string      { return "4011102" }
func (CardDef4011102) Name() string    { return "大法师 罗慕路斯" }
func (CardDef4011102) Kind() string    { return "人物" }
func (CardDef4011102) Element() string { return "无" }

func (CardDef4011102) Card() model.Card {
	return model.Card{
		Number:          "4011102",
		Type:            "人物",
		Name:            "大法师 罗慕路斯",
		Category:        "无",
		Tag:             "",
		Description:     "此卡提供的元素只能严格当做\\无使用",
		Quote:           "\"看啊,多么完美而纯粹的造物!\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"无": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\无\\4011102.jpg",
	}
}

type CardDef4111001 struct{}

func (CardDef4111001) ID() string      { return "4111001" }
func (CardDef4111001) Name() string    { return "掌门 龙卷火" }
func (CardDef4111001) Kind() string    { return "人物" }
func (CardDef4111001) Element() string { return "火" }

func (CardDef4111001) Card() model.Card {
	return model.Card{
		Number:          "4111001",
		Type:            "人物",
		Name:            "掌门 龙卷火",
		Category:        "火",
		Tag:             "",
		Description:     "入场:将一张衍生卡牌万火合一术置于你的技能池",
		Quote:           "\"我最强大的弟子,我们本可一起称霸...很遗憾你选择与我为敌\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"火": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"3101002"},
		OutputPath:      "output\\基础包\\人物\\火\\4111001.jpg",
	}
}

type CardDef4111002 struct{}

func (CardDef4111002) ID() string      { return "4111002" }
func (CardDef4111002) Name() string    { return "女巫 维兰德" }
func (CardDef4111002) Kind() string    { return "人物" }
func (CardDef4111002) Element() string { return "火" }

func (CardDef4111002) Card() model.Card {
	return model.Card{
		Number:          "4111002",
		Type:            "人物",
		Name:            "女巫 维兰德",
		Category:        "火",
		Tag:             "",
		Description:     "主动回合技:你的人物获得点燃1,然后直到本回合结束,将此卡负载的1\\火变为1\\无",
		Quote:           "\"别瞎说,我们可不敢烧女巫,是她自己烧的\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"火": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\火\\4111002.jpg",
	}
}

type CardDef4111003 struct{}

func (CardDef4111003) ID() string      { return "4111003" }
func (CardDef4111003) Name() string    { return "大祭司 梵天" }
func (CardDef4111003) Kind() string    { return "人物" }
func (CardDef4111003) Element() string { return "火" }

func (CardDef4111003) Card() model.Card {
	return model.Card{
		Number:          "4111003",
		Type:            "人物",
		Name:            "大祭司 梵天",
		Category:        "火",
		Tag:             "",
		Description:     "主动绝技:本回合内每当你的火焰法术命中,此卡永久获得负载+1\\火",
		Quote:           "\"我,梵天,带领全火焰门派,将追随我们唯一的王.愿祭坛之火庇佑她的王国\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"火": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\火\\4111003.jpg",
	}
}

type CardDef4111101 struct{}

func (CardDef4111101) ID() string      { return "4111101" }
func (CardDef4111101) Name() string    { return "首席顾问 费林" }
func (CardDef4111101) Kind() string    { return "人物" }
func (CardDef4111101) Element() string { return "火" }

func (CardDef4111101) Card() model.Card {
	return model.Card{
		Number:          "4111101",
		Type:            "人物",
		Name:            "首席顾问 费林",
		Category:        "火",
		Tag:             "",
		Description:     "主动绝技:献祭1个友方火焰伙伴才能发动,你的下1次入场的火焰卡牌入场花费减少献祭卡牌入场花费的元素",
		Quote:           "\"不用担心,我想师兄会认同我们的,毕竟,我们必须以王国的利益为重\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"火": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\火\\4111101.jpg",
	}
}

type CardDef4111102 struct{}

func (CardDef4111102) ID() string      { return "4111102" }
func (CardDef4111102) Name() string    { return "大将军 克兰" }
func (CardDef4111102) Kind() string    { return "人物" }
func (CardDef4111102) Element() string { return "火" }

func (CardDef4111102) Card() model.Card {
	return model.Card{
		Number:          "4111102",
		Type:            "人物",
		Name:            "大将军 克兰",
		Category:        "火",
		Tag:             "",
		Description:     "诱发回合技:当双方法术成功防御时才能发动,翻取1张火焰卡牌,然后你丢弃1张手牌",
		Quote:           "\"我给你了一次又一次的机会,可惜你依然选择违抗我\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"火": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\火\\4111102.jpg",
	}
}

type CardDef4211001 struct{}

func (CardDef4211001) ID() string      { return "4211001" }
func (CardDef4211001) Name() string    { return "\"浪之人\" 巴特尔" }
func (CardDef4211001) Kind() string    { return "人物" }
func (CardDef4211001) Element() string { return "水" }

func (CardDef4211001) Card() model.Card {
	return model.Card{
		Number:          "4211001",
		Type:            "人物",
		Name:            "\"浪之人\" 巴特尔",
		Category:        "水",
		Tag:             "",
		Description:     "主动绝技:展示你的1张手牌,其属性永久变为水,入场花费和负载的元素全部变为等量的\\水",
		Quote:           "\"把这个带上吧\"人鱼手捧一小块蓝色的晶石,\"它会带给你...好运\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"水": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\水\\4211001.jpg",
	}
}

type CardDef4211002 struct{}

func (CardDef4211002) ID() string      { return "4211002" }
func (CardDef4211002) Name() string    { return "大贤者 沃尔波特" }
func (CardDef4211002) Kind() string    { return "人物" }
func (CardDef4211002) Element() string { return "水" }

func (CardDef4211002) Card() model.Card {
	return model.Card{
		Number:          "4211002",
		Type:            "人物",
		Name:            "大贤者 沃尔波特",
		Category:        "水",
		Tag:             "",
		Description:     "入场:将1张衍生卡牌百川归海置于你的技能池",
		Quote:           "\"反对我?你不是第一个,但确实是现在唯一的一个\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"水": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"3201001"},
		OutputPath:      "output\\基础包\\人物\\水\\4211002.jpg",
	}
}

type CardDef4211003 struct{}

func (CardDef4211003) ID() string      { return "4211003" }
func (CardDef4211003) Name() string    { return "凛冬城主 水晶心" }
func (CardDef4211003) Kind() string    { return "人物" }
func (CardDef4211003) Element() string { return "水" }

func (CardDef4211003) Card() model.Card {
	return model.Card{
		Number:          "4211003",
		Type:            "人物",
		Name:            "凛冬城主 水晶心",
		Category:        "水",
		Tag:             "",
		Description:     "主动绝技:在本回合剩余时间内,你技能区内的法术获得\"冻结1\"",
		Quote:           "\"卡姆陛下,我想凛冬城下的无数冰雕已经足够说明了,这里不需要一个国王\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"水": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\水\\4211003.jpg",
	}
}

type CardDef4211101 struct{}

func (CardDef4211101) ID() string      { return "4211101" }
func (CardDef4211101) Name() string    { return "海神之使 珊瑚 贝莉 " }
func (CardDef4211101) Kind() string    { return "人物" }
func (CardDef4211101) Element() string { return "水" }

func (CardDef4211101) Card() model.Card {
	return model.Card{
		Number:          "4211101",
		Type:            "人物",
		Name:            "海神之使 珊瑚 贝莉 ",
		Category:        "水",
		Tag:             "",
		Description:     "诱发:本局游戏你第一次使用法术攻击时,使该法术永久+3\\威",
		Quote:           "\"别忘了父亲的话,每个人出生便被赋予了属于自己和家族的责任\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"水": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\水\\4211101.jpg",
	}
}

type CardDef4211102 struct{}

func (CardDef4211102) ID() string      { return "4211102" }
func (CardDef4211102) Name() string    { return "凛冰魔巫 索菲娅" }
func (CardDef4211102) Kind() string    { return "人物" }
func (CardDef4211102) Element() string { return "水" }

func (CardDef4211102) Card() model.Card {
	return model.Card{
		Number:          "4211102",
		Type:            "人物",
		Name:            "凛冰魔巫 索菲娅",
		Category:        "水",
		Tag:             "",
		Description:     "光环:此卡不受冻结效果影响.绝技:移除场上任意1个单位的1层冻结,对其造成2点伤害",
		Quote:           "\"终于,你找到了我\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"水": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\水\\4211102.jpg",
	}
}

type CardDef4311001 struct{}

func (CardDef4311001) ID() string      { return "4311001" }
func (CardDef4311001) Name() string    { return "雷术士 肃" }
func (CardDef4311001) Kind() string    { return "人物" }
func (CardDef4311001) Element() string { return "气" }

func (CardDef4311001) Card() model.Card {
	return model.Card{
		Number:          "4311001",
		Type:            "人物",
		Name:            "雷术士 肃",
		Category:        "气",
		Tag:             "",
		Description:     "主动绝技:丢弃2张大气手牌才能发动,对任意1名敌人造成1点伤害",
		Quote:           "真正的法术是与生俱来的,可惜这点天分与神的力量相比还是太渺小了",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"气": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\气\\4311001.jpg",
	}
}

type CardDef4311002 struct{}

func (CardDef4311002) ID() string      { return "4311002" }
func (CardDef4311002) Name() string    { return "\"渡鸦\" 睿文" }
func (CardDef4311002) Kind() string    { return "人物" }
func (CardDef4311002) Element() string { return "气" }

func (CardDef4311002) Card() model.Card {
	return model.Card{
		Number:          "4311002",
		Type:            "人物",
		Name:            "\"渡鸦\" 睿文",
		Category:        "气",
		Tag:             "",
		Description:     "你的起始手牌数与换牌机会+1",
		Quote:           "凌晨,他们冲入囚室,但是房间内空无一人,地上只有一根黑色的羽毛",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"气": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\气\\4311002.jpg",
	}
}

type CardDef4311003 struct{}

func (CardDef4311003) ID() string      { return "4311003" }
func (CardDef4311003) Name() string    { return "掌门 穆伶" }
func (CardDef4311003) Kind() string    { return "人物" }
func (CardDef4311003) Element() string { return "气" }

func (CardDef4311003) Card() model.Card {
	return model.Card{
		Number:          "4311003",
		Type:            "人物",
		Name:            "掌门 穆伶",
		Category:        "气",
		Tag:             "",
		Description:     "主动绝技:选择你的法力范围内双方各1个伙伴,花费它们入场费用差值的\\气,将它们移回各自手牌",
		Quote:           "\"时过境迁,世上不再会有我们的立足之地\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"气": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\气\\4311003.jpg",
	}
}

type CardDef4311101 struct{}

func (CardDef4311101) ID() string      { return "4311101" }
func (CardDef4311101) Name() string    { return "司天魔巫 索兰德" }
func (CardDef4311101) Kind() string    { return "人物" }
func (CardDef4311101) Element() string { return "气" }

func (CardDef4311101) Card() model.Card {
	return model.Card{
		Number:          "4311101",
		Type:            "人物",
		Name:            "司天魔巫 索兰德",
		Category:        "气",
		Tag:             "",
		Description:     "光环:你的驱动和聚能法术永久+1\\威,你不能从技能池学习其他种类的法术",
		Quote:           "",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"气": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\气\\4311101.jpg",
	}
}

type CardDef4311102 struct{}

func (CardDef4311102) ID() string      { return "4311102" }
func (CardDef4311102) Name() string    { return "布雾者 弗格" }
func (CardDef4311102) Kind() string    { return "人物" }
func (CardDef4311102) Element() string { return "气" }

func (CardDef4311102) Card() model.Card {
	return model.Card{
		Number:          "4311102",
		Type:            "人物",
		Name:            "布雾者 弗格",
		Category:        "气",
		Tag:             "",
		Description:     "绝技:双方各自召唤的下1个伙伴入场时获得隐蔽2",
		Quote:           "你确定你战胜的是索兰德?还是一片很像他的云?",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"气": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\气\\4311102.jpg",
	}
}

type CardDef4411001 struct{}

func (CardDef4411001) ID() string      { return "4411001" }
func (CardDef4411001) Name() string    { return "森林隐士 白须" }
func (CardDef4411001) Kind() string    { return "人物" }
func (CardDef4411001) Element() string { return "地" }

func (CardDef4411001) Card() model.Card {
	return model.Card{
		Number:          "4411001",
		Type:            "人物",
		Name:            "森林隐士 白须",
		Category:        "地",
		Tag:             "",
		Description:     "诱发:在你的首个回合的抽牌阶段,你可以用检索1张地属性野兽、植物或精灵来代替抽牌",
		Quote:           "人们相信,白须能听懂动物说的话.这可能也就是为什么他在海边洗脚时会戴上耳塞",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"地": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\地\\4411001.jpg",
	}
}

type CardDef4411002 struct{}

func (CardDef4411002) ID() string      { return "4411002" }
func (CardDef4411002) Name() string    { return "大法师 安德鲁" }
func (CardDef4411002) Kind() string    { return "人物" }
func (CardDef4411002) Element() string { return "地" }

func (CardDef4411002) Card() model.Card {
	return model.Card{
		Number:          "4411002",
		Type:            "人物",
		Name:            "大法师 安德鲁",
		Category:        "地",
		Tag:             "",
		Description:     "入场:将1张衍生卡牌灵兽 辛柯置于你的卡组",
		Quote:           "走过千山万水,遍历浮生世事,友谊、爱情,越是珍视越是脆弱,倒不如重返自然落得个风轻云淡",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"地": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"1401002"},
		OutputPath:      "output\\基础包\\人物\\地\\4411002.jpg",
	}
}

type CardDef4411003 struct{}

func (CardDef4411003) ID() string      { return "4411003" }
func (CardDef4411003) Name() string    { return "麦吉教授" }
func (CardDef4411003) Kind() string    { return "人物" }
func (CardDef4411003) Element() string { return "地" }

func (CardDef4411003) Card() model.Card {
	return model.Card{
		Number:          "4411003",
		Type:            "人物",
		Name:            "麦吉教授",
		Category:        "地",
		Tag:             "",
		Description:     "光环:你的打出或学习的第一张原始费用大于5的卡牌费用减少2\\地",
		Quote:           "《地理学入门》成功将没人听得懂的东西,弄成了没人愿意听的东西",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"地": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\地\\4411003.jpg",
	}
}

type CardDef4411101 struct{}

func (CardDef4411101) ID() string      { return "4411101" }
func (CardDef4411101) Name() string    { return "翡翠男爵 杰德 拜利兰" }
func (CardDef4411101) Kind() string    { return "人物" }
func (CardDef4411101) Element() string { return "地" }

func (CardDef4411101) Card() model.Card {
	return model.Card{
		Number:          "4411101",
		Type:            "人物",
		Name:            "翡翠男爵 杰德 拜利兰",
		Category:        "地",
		Tag:             "",
		Description:     "光环:当你的护盾小于3时不会在回合结束时减少",
		Quote:           "拜利兰家族可能多少对永生有着独特的渴望,杰德也一样,不过他选择了大家都能接受的方式",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"地": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\地\\4411101.jpg",
	}
}

type CardDef4411102 struct{}

func (CardDef4411102) ID() string      { return "4411102" }
func (CardDef4411102) Name() string    { return "秋枫领主 狄利克雷" }
func (CardDef4411102) Kind() string    { return "人物" }
func (CardDef4411102) Element() string { return "地" }

func (CardDef4411102) Card() model.Card {
	return model.Card{
		Number:          "4411102",
		Type:            "人物",
		Name:            "秋枫领主 狄利克雷",
		Category:        "地",
		Tag:             "",
		Description:     "光环:当你的地脉卡牌透支时,你可以获得剩余的元素的2倍(当你透支超过所需费用时,你不能继续透支)",
		Quote:           "我与休伯特的意见并不相同,我认为人活着比死了更有价值",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"地": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\地\\4411102.jpg",
	}
}

type CardDef4511001 struct{}

func (CardDef4511001) ID() string      { return "4511001" }
func (CardDef4511001) Name() string    { return "圣使 玛丽斯 南森埃尔" }
func (CardDef4511001) Kind() string    { return "人物" }
func (CardDef4511001) Element() string { return "光" }

func (CardDef4511001) Card() model.Card {
	return model.Card{
		Number:          "4511001",
		Type:            "人物",
		Name:            "圣使 玛丽斯 南森埃尔",
		Category:        "光",
		Tag:             "",
		Description:     "诱发绝技:当敌方将要造成伤害时可以发动,直到你的下个回合结束,你的每个单位每次受到对方伤害,获得2\\光",
		Quote:           "\"我从未迷茫,因为唯一的道路早已被照亮\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"光": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\光\\4511001.jpg",
	}
}

type CardDef4511002 struct{}

func (CardDef4511002) ID() string      { return "4511002" }
func (CardDef4511002) Name() string    { return "神之眷子 爱里默" }
func (CardDef4511002) Kind() string    { return "人物" }
func (CardDef4511002) Element() string { return "光" }

func (CardDef4511002) Card() model.Card {
	return model.Card{
		Number:          "4511002",
		Type:            "人物",
		Name:            "神之眷子 爱里默",
		Category:        "光",
		Tag:             "",
		Description:     "入场:将5张衍生卡牌桎梏置于对手的牌组,当全部被解除(进入弃牌堆)时你的人物获得主动绝技:移除场上最多3张人物牌以外的任意卡牌",
		Quote:           "桎梏下的鸟儿,是否终能有飞翔的一天?",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"光": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2501001"},
		OutputPath:      "output\\基础包\\人物\\光\\4511002.jpg",
	}
}

type CardDef4511003 struct{}

func (CardDef4511003) ID() string      { return "4511003" }
func (CardDef4511003) Name() string    { return "骑士团长 蕾曦娅" }
func (CardDef4511003) Kind() string    { return "人物" }
func (CardDef4511003) Element() string { return "光" }

func (CardDef4511003) Card() model.Card {
	return model.Card{
		Number:          "4511003",
		Type:            "人物",
		Name:            "骑士团长 蕾曦娅",
		Category:        "光",
		Tag:             "",
		Description:     "入场:如果你的技能池中有\"希望呼唤\",用1张衍生卡牌\"团结的希望\"将其替换",
		Quote:           "\"选择与被选择,都是一种奢侈\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"光": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"3501001"},
		OutputPath:      "output\\基础包\\人物\\光\\4511003.jpg",
	}
}

type CardDef4511101 struct{}

func (CardDef4511101) ID() string      { return "4511101" }
func (CardDef4511101) Name() string    { return "庇护者 西瓦尔" }
func (CardDef4511101) Kind() string    { return "人物" }
func (CardDef4511101) Element() string { return "光" }

func (CardDef4511101) Card() model.Card {
	return model.Card{
		Number:          "4511101",
		Type:            "人物",
		Name:            "庇护者 西瓦尔",
		Category:        "光",
		Tag:             "",
		Description:     "诱发绝技:当友方单位在一个回合内受到3点及以上总伤害时才能发动,直到下个回合结束所有友方单位不会受到任何伤害",
		Quote:           "\"我不关心你们间的恩怨,但是在这里,没有人可以开战\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"光": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\光\\4511101.jpg",
	}
}

type CardDef4511102 struct{}

func (CardDef4511102) ID() string      { return "4511102" }
func (CardDef4511102) Name() string    { return "救赎者 伊芙 秋枫" }
func (CardDef4511102) Kind() string    { return "人物" }
func (CardDef4511102) Element() string { return "光" }

func (CardDef4511102) Card() model.Card {
	return model.Card{
		Number:          "4511102",
		Type:            "人物",
		Name:            "救赎者 伊芙 秋枫",
		Category:        "光",
		Tag:             "",
		Description:     "主动绝技:敌方场上单位数量比我方多时才能发动,选择1个受伤的友方伙伴,花费X\\光使该伙伴获得+X\\血和负载+X\\光,X为你场上受伤单位数量",
		Quote:           "只要你探寻,必得指引;只要你祈求,必得救赎",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"光": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\光\\4511102.jpg",
	}
}

type CardDef4611001 struct{}

func (CardDef4611001) ID() string      { return "4611001" }
func (CardDef4611001) Name() string    { return "暗影学者 爱莉斯" }
func (CardDef4611001) Kind() string    { return "人物" }
func (CardDef4611001) Element() string { return "暗" }

func (CardDef4611001) Card() model.Card {
	return model.Card{
		Number:          "4611001",
		Type:            "人物",
		Name:            "暗影学者 爱莉斯",
		Category:        "暗",
		Tag:             "",
		Description:     "诱发回合技:当1个你的伙伴死亡后,使你的1个法术+1\\威",
		Quote:           "\"这本叫'群屿编年史'的,真是我见过最富想象力的东西\"——\"观察者\" 欧柯茹",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"暗": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\暗\\4611001.jpg",
	}
}

type CardDef4611002 struct{}

func (CardDef4611002) ID() string      { return "4611002" }
func (CardDef4611002) Name() string    { return "芙雅夫人" }
func (CardDef4611002) Kind() string    { return "人物" }
func (CardDef4611002) Element() string { return "暗" }

func (CardDef4611002) Card() model.Card {
	return model.Card{
		Number:          "4611002",
		Type:            "人物",
		Name:            "芙雅夫人",
		Category:        "暗",
		Tag:             "",
		Description:     "主动绝技:选择你的1个竖置状态的伙伴,使其攻击和负载翻倍,但会在消耗或透支后死亡",
		Quote:           "\"你是想现在就为我服务,还是死后再为我服务?\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"暗": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\基础包\\人物\\暗\\4611002.jpg",
	}
}

type CardDef4611003 struct{}

func (CardDef4611003) ID() string      { return "4611003" }
func (CardDef4611003) Name() string    { return "咒言师 结影" }
func (CardDef4611003) Kind() string    { return "人物" }
func (CardDef4611003) Element() string { return "暗" }

func (CardDef4611003) Card() model.Card {
	return model.Card{
		Number:          "4611003",
		Type:            "人物",
		Name:            "咒言师 结影",
		Category:        "暗",
		Tag:             "",
		Description:     "入场:将三张衍生卡牌咒言书洗入你的卡组",
		Quote:           "\"书是人的养料,人也可以是书的养料\"",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"暗": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "基础包",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"2601002"},
		OutputPath:      "output\\基础包\\人物\\暗\\4611003.jpg",
	}
}

type CardDef4611101 struct{}

func (CardDef4611101) ID() string      { return "4611101" }
func (CardDef4611101) Name() string    { return "鲜血伯爵 休伯特 黑松" }
func (CardDef4611101) Kind() string    { return "人物" }
func (CardDef4611101) Element() string { return "暗" }

func (CardDef4611101) Card() model.Card {
	return model.Card{
		Number:          "4611101",
		Type:            "人物",
		Name:            "鲜血伯爵 休伯特 黑松",
		Category:        "暗",
		Tag:             "",
		Description:     "入场:将1张衍生卡牌鲜血盛宴加入你的技能池",
		Quote:           "只要鲜血和黑魔法还在滋养我的家族,黑松就会在这片土地上屹立不倒",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"暗": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{"3601101"},
		OutputPath:      "output\\王权纷争\\人物\\暗\\4611101.jpg",
	}
}

type CardDef4611102 struct{}

func (CardDef4611102) ID() string      { return "4611102" }
func (CardDef4611102) Name() string    { return "灾厄玫瑰 多姆" }
func (CardDef4611102) Kind() string    { return "人物" }
func (CardDef4611102) Element() string { return "暗" }

func (CardDef4611102) Card() model.Card {
	return model.Card{
		Number:          "4611102",
		Type:            "人物",
		Name:            "灾厄玫瑰 多姆",
		Category:        "暗",
		Tag:             "",
		Description:     "入场:从双方卡组上方将4张牌送去弃牌堆",
		Quote:           "\"我就是家里最年长的\",少年看着手中的花苞说道",
		ElementsCost:    map[string]int{},
		ElementsGain:    map[string]int{"暗": 4},
		ElementsExpense: map[string]int{},
		VersionNum:      "1",
		VersionName:     "王权纷争",
		Attack:          -1,
		Life:            6,
		Duration:        -1,
		Power:           -1,
		Spawns:          []string{},
		OutputPath:      "output\\王权纷争\\人物\\暗\\4611102.jpg",
	}
}

var compiledCardDefinitions = []CardDefinition{
	CardDef1001101{},
	CardDef1011001{},
	CardDef1011002{},
	CardDef1011003{},
	CardDef1011101{},
	CardDef1011102{},
	CardDef1011103{},
	CardDef1021001{},
	CardDef1021002{},
	CardDef1021003{},
	CardDef1021004{},
	CardDef1021005{},
	CardDef1021006{},
	CardDef1021007{},
	CardDef1021008{},
	CardDef1021009{},
	CardDef1021010{},
	CardDef1021011{},
	CardDef1021012{},
	CardDef1021013{},
	CardDef1021014{},
	CardDef1021015{},
	CardDef1021016{},
	CardDef1021017{},
	CardDef1021018{},
	CardDef1021101{},
	CardDef1021102{},
	CardDef1021103{},
	CardDef1021104{},
	CardDef1021105{},
	CardDef1021106{},
	CardDef1021107{},
	CardDef1021108{},
	CardDef1021109{},
	CardDef1021110{},
	CardDef1021111{},
	CardDef1021112{},
	CardDef1021113{},
	CardDef1021114{},
	CardDef1021115{},
	CardDef1111001{},
	CardDef1111002{},
	CardDef1111003{},
	CardDef1111101{},
	CardDef1111102{},
	CardDef1111103{},
	CardDef1121001{},
	CardDef1121002{},
	CardDef1121003{},
	CardDef1121004{},
	CardDef1121005{},
	CardDef1121006{},
	CardDef1121007{},
	CardDef1121008{},
	CardDef1121009{},
	CardDef1121010{},
	CardDef1121011{},
	CardDef1121012{},
	CardDef1121013{},
	CardDef1121014{},
	CardDef1121015{},
	CardDef1121016{},
	CardDef1121101{},
	CardDef1121102{},
	CardDef1121103{},
	CardDef1121104{},
	CardDef1121105{},
	CardDef1121106{},
	CardDef1121107{},
	CardDef1121108{},
	CardDef1121109{},
	CardDef1121110{},
	CardDef1121111{},
	CardDef1121112{},
	CardDef1121113{},
	CardDef1121114{},
	CardDef1121115{},
	CardDef1201101{},
	CardDef1211001{},
	CardDef1211002{},
	CardDef1211003{},
	CardDef1211101{},
	CardDef1211102{},
	CardDef1211103{},
	CardDef1221001{},
	CardDef1221002{},
	CardDef1221003{},
	CardDef1221004{},
	CardDef1221005{},
	CardDef1221006{},
	CardDef1221007{},
	CardDef1221008{},
	CardDef1221009{},
	CardDef1221010{},
	CardDef1221011{},
	CardDef1221012{},
	CardDef1221013{},
	CardDef1221014{},
	CardDef1221015{},
	CardDef1221016{},
	CardDef1221101{},
	CardDef1221102{},
	CardDef1221103{},
	CardDef1221104{},
	CardDef1221105{},
	CardDef1221106{},
	CardDef1221107{},
	CardDef1221108{},
	CardDef1221109{},
	CardDef1221110{},
	CardDef1221111{},
	CardDef1221112{},
	CardDef1221113{},
	CardDef1221114{},
	CardDef1221115{},
	CardDef1311001{},
	CardDef1311002{},
	CardDef1311003{},
	CardDef1311101{},
	CardDef1311102{},
	CardDef1311103{},
	CardDef1321001{},
	CardDef1321002{},
	CardDef1321003{},
	CardDef1321004{},
	CardDef1321005{},
	CardDef1321006{},
	CardDef1321007{},
	CardDef1321008{},
	CardDef1321009{},
	CardDef1321010{},
	CardDef1321011{},
	CardDef1321012{},
	CardDef1321013{},
	CardDef1321014{},
	CardDef1321015{},
	CardDef1321016{},
	CardDef1321101{},
	CardDef1321102{},
	CardDef1321103{},
	CardDef1321104{},
	CardDef1321105{},
	CardDef1321106{},
	CardDef1321107{},
	CardDef1321108{},
	CardDef1321109{},
	CardDef1321110{},
	CardDef1321111{},
	CardDef1321112{},
	CardDef1321113{},
	CardDef1321114{},
	CardDef1321115{},
	CardDef1401001{},
	CardDef1401002{},
	CardDef1401101{},
	CardDef1411001{},
	CardDef1411002{},
	CardDef1411003{},
	CardDef1411101{},
	CardDef1411102{},
	CardDef1411103{},
	CardDef1421001{},
	CardDef1421002{},
	CardDef1421003{},
	CardDef1421004{},
	CardDef1421005{},
	CardDef1421006{},
	CardDef1421007{},
	CardDef1421008{},
	CardDef1421009{},
	CardDef1421010{},
	CardDef1421011{},
	CardDef1421012{},
	CardDef1421013{},
	CardDef1421014{},
	CardDef1421015{},
	CardDef1421016{},
	CardDef1421101{},
	CardDef1421102{},
	CardDef1421103{},
	CardDef1421104{},
	CardDef1421105{},
	CardDef1421106{},
	CardDef1421107{},
	CardDef1421108{},
	CardDef1421109{},
	CardDef1421110{},
	CardDef1421111{},
	CardDef1421112{},
	CardDef1421113{},
	CardDef1421114{},
	CardDef1421115{},
	CardDef1501001{},
	CardDef1511001{},
	CardDef1511002{},
	CardDef1511003{},
	CardDef1511101{},
	CardDef1511102{},
	CardDef1511103{},
	CardDef1521001{},
	CardDef1521002{},
	CardDef1521003{},
	CardDef1521004{},
	CardDef1521005{},
	CardDef1521006{},
	CardDef1521007{},
	CardDef1521008{},
	CardDef1521009{},
	CardDef1521010{},
	CardDef1521011{},
	CardDef1521012{},
	CardDef1521013{},
	CardDef1521014{},
	CardDef1521015{},
	CardDef1521016{},
	CardDef1521101{},
	CardDef1521102{},
	CardDef1521103{},
	CardDef1521104{},
	CardDef1521105{},
	CardDef1521106{},
	CardDef1521107{},
	CardDef1521108{},
	CardDef1521109{},
	CardDef1521110{},
	CardDef1521111{},
	CardDef1521112{},
	CardDef1521113{},
	CardDef1521114{},
	CardDef1521115{},
	CardDef1601101{},
	CardDef1611001{},
	CardDef1611002{},
	CardDef1611003{},
	CardDef1611101{},
	CardDef1611102{},
	CardDef1611103{},
	CardDef1621001{},
	CardDef1621002{},
	CardDef1621003{},
	CardDef1621004{},
	CardDef1621005{},
	CardDef1621006{},
	CardDef1621007{},
	CardDef1621008{},
	CardDef1621009{},
	CardDef1621010{},
	CardDef1621011{},
	CardDef1621012{},
	CardDef1621013{},
	CardDef1621014{},
	CardDef1621015{},
	CardDef1621016{},
	CardDef1621101{},
	CardDef1621102{},
	CardDef1621103{},
	CardDef1621104{},
	CardDef1621105{},
	CardDef1621106{},
	CardDef1621107{},
	CardDef1621108{},
	CardDef1621109{},
	CardDef1621110{},
	CardDef1621111{},
	CardDef1621112{},
	CardDef1621113{},
	CardDef1621114{},
	CardDef1621115{},
	CardDef2001101{},
	CardDef2001102{},
	CardDef2011001{},
	CardDef2011002{},
	CardDef2011003{},
	CardDef2011101{},
	CardDef2011102{},
	CardDef2021001{},
	CardDef2021002{},
	CardDef2021003{},
	CardDef2021004{},
	CardDef2021005{},
	CardDef2021006{},
	CardDef2021007{},
	CardDef2021008{},
	CardDef2021009{},
	CardDef2021010{},
	CardDef2021011{},
	CardDef2021012{},
	CardDef2021013{},
	CardDef2021014{},
	CardDef2021015{},
	CardDef2021016{},
	CardDef2021017{},
	CardDef2021018{},
	CardDef2021019{},
	CardDef2021020{},
	CardDef2021021{},
	CardDef2021022{},
	CardDef2021023{},
	CardDef2021101{},
	CardDef2021102{},
	CardDef2021103{},
	CardDef2021104{},
	CardDef2021105{},
	CardDef2021106{},
	CardDef2021107{},
	CardDef2021108{},
	CardDef2021109{},
	CardDef2021110{},
	CardDef2021111{},
	CardDef2021112{},
	CardDef2021113{},
	CardDef2021114{},
	CardDef2021115{},
	CardDef2021116{},
	CardDef2111001{},
	CardDef2111002{},
	CardDef2111101{},
	CardDef2111102{},
	CardDef2121001{},
	CardDef2121002{},
	CardDef2121003{},
	CardDef2121004{},
	CardDef2121005{},
	CardDef2121006{},
	CardDef2121007{},
	CardDef2121008{},
	CardDef2121009{},
	CardDef2121010{},
	CardDef2121011{},
	CardDef2121012{},
	CardDef2121013{},
	CardDef2121014{},
	CardDef2121101{},
	CardDef2121102{},
	CardDef2121103{},
	CardDef2121104{},
	CardDef2121105{},
	CardDef2121106{},
	CardDef2121107{},
	CardDef2121108{},
	CardDef2121109{},
	CardDef2121110{},
	CardDef2121111{},
	CardDef2121112{},
	CardDef2201101{},
	CardDef2201102{},
	CardDef2201103{},
	CardDef2211001{},
	CardDef2211002{},
	CardDef2211101{},
	CardDef2211102{},
	CardDef2221001{},
	CardDef2221002{},
	CardDef2221003{},
	CardDef2221004{},
	CardDef2221005{},
	CardDef2221006{},
	CardDef2221007{},
	CardDef2221008{},
	CardDef2221009{},
	CardDef2221010{},
	CardDef2221011{},
	CardDef2221012{},
	CardDef2221013{},
	CardDef2221014{},
	CardDef2221101{},
	CardDef2221102{},
	CardDef2221103{},
	CardDef2221104{},
	CardDef2221105{},
	CardDef2221106{},
	CardDef2221107{},
	CardDef2221108{},
	CardDef2221109{},
	CardDef2221110{},
	CardDef2221111{},
	CardDef2221112{},
	CardDef2311001{},
	CardDef2311002{},
	CardDef2311101{},
	CardDef2311102{},
	CardDef2321001{},
	CardDef2321002{},
	CardDef2321003{},
	CardDef2321004{},
	CardDef2321005{},
	CardDef2321006{},
	CardDef2321007{},
	CardDef2321008{},
	CardDef2321009{},
	CardDef2321010{},
	CardDef2321011{},
	CardDef2321012{},
	CardDef2321013{},
	CardDef2321014{},
	CardDef2321101{},
	CardDef2321102{},
	CardDef2321103{},
	CardDef2321104{},
	CardDef2321105{},
	CardDef2321106{},
	CardDef2321107{},
	CardDef2321108{},
	CardDef2321109{},
	CardDef2321110{},
	CardDef2321111{},
	CardDef2321112{},
	CardDef2411001{},
	CardDef2411002{},
	CardDef2411101{},
	CardDef2411102{},
	CardDef2421001{},
	CardDef2421002{},
	CardDef2421003{},
	CardDef2421004{},
	CardDef2421005{},
	CardDef2421006{},
	CardDef2421007{},
	CardDef2421008{},
	CardDef2421009{},
	CardDef2421010{},
	CardDef2421011{},
	CardDef2421012{},
	CardDef2421013{},
	CardDef2421014{},
	CardDef2421101{},
	CardDef2421102{},
	CardDef2421103{},
	CardDef2421104{},
	CardDef2421105{},
	CardDef2421106{},
	CardDef2421107{},
	CardDef2421108{},
	CardDef2421109{},
	CardDef2421110{},
	CardDef2421111{},
	CardDef2421112{},
	CardDef2501001{},
	CardDef2511001{},
	CardDef2511002{},
	CardDef2511101{},
	CardDef2511102{},
	CardDef2521001{},
	CardDef2521002{},
	CardDef2521003{},
	CardDef2521004{},
	CardDef2521005{},
	CardDef2521006{},
	CardDef2521007{},
	CardDef2521008{},
	CardDef2521009{},
	CardDef2521010{},
	CardDef2521011{},
	CardDef2521012{},
	CardDef2521013{},
	CardDef2521014{},
	CardDef2521101{},
	CardDef2521102{},
	CardDef2521103{},
	CardDef2521104{},
	CardDef2521105{},
	CardDef2521106{},
	CardDef2521107{},
	CardDef2521108{},
	CardDef2521109{},
	CardDef2521110{},
	CardDef2521111{},
	CardDef2521112{},
	CardDef2601001{},
	CardDef2601002{},
	CardDef2611001{},
	CardDef2611002{},
	CardDef2611101{},
	CardDef2611102{},
	CardDef2621001{},
	CardDef2621002{},
	CardDef2621003{},
	CardDef2621004{},
	CardDef2621005{},
	CardDef2621006{},
	CardDef2621007{},
	CardDef2621008{},
	CardDef2621009{},
	CardDef2621010{},
	CardDef2621011{},
	CardDef2621012{},
	CardDef2621013{},
	CardDef2621014{},
	CardDef2621101{},
	CardDef2621102{},
	CardDef2621103{},
	CardDef2621104{},
	CardDef2621105{},
	CardDef2621106{},
	CardDef2621107{},
	CardDef2621108{},
	CardDef2621109{},
	CardDef2621110{},
	CardDef2621111{},
	CardDef2621112{},
	CardDef3001001{},
	CardDef3001002{},
	CardDef3001101{},
	CardDef3011101{},
	CardDef3021001{},
	CardDef3021002{},
	CardDef3021003{},
	CardDef3021004{},
	CardDef3021005{},
	CardDef3021006{},
	CardDef3021007{},
	CardDef3021008{},
	CardDef3021009{},
	CardDef3021010{},
	CardDef3021011{},
	CardDef3021012{},
	CardDef3021101{},
	CardDef3021102{},
	CardDef3021103{},
	CardDef3021104{},
	CardDef3021105{},
	CardDef3021106{},
	CardDef3021107{},
	CardDef3021108{},
	CardDef3101001{},
	CardDef3101002{},
	CardDef3111101{},
	CardDef3111102{},
	CardDef3121001{},
	CardDef3121002{},
	CardDef3121003{},
	CardDef3121004{},
	CardDef3121005{},
	CardDef3121006{},
	CardDef3121007{},
	CardDef3121008{},
	CardDef3121009{},
	CardDef3121010{},
	CardDef3121011{},
	CardDef3121012{},
	CardDef3121013{},
	CardDef3121014{},
	CardDef3121015{},
	CardDef3121101{},
	CardDef3121102{},
	CardDef3121103{},
	CardDef3121104{},
	CardDef3121105{},
	CardDef3121106{},
	CardDef3121107{},
	CardDef3121108{},
	CardDef3121109{},
	CardDef3121110{},
	CardDef3201001{},
	CardDef3201002{},
	CardDef3211101{},
	CardDef3211102{},
	CardDef3221001{},
	CardDef3221002{},
	CardDef3221003{},
	CardDef3221004{},
	CardDef3221005{},
	CardDef3221006{},
	CardDef3221007{},
	CardDef3221008{},
	CardDef3221009{},
	CardDef3221010{},
	CardDef3221011{},
	CardDef3221012{},
	CardDef3221013{},
	CardDef3221014{},
	CardDef3221015{},
	CardDef3221101{},
	CardDef3221102{},
	CardDef3221103{},
	CardDef3221104{},
	CardDef3221105{},
	CardDef3221106{},
	CardDef3221107{},
	CardDef3221108{},
	CardDef3221109{},
	CardDef3221110{},
	CardDef3301001{},
	CardDef3311101{},
	CardDef3311102{},
	CardDef3321001{},
	CardDef3321002{},
	CardDef3321003{},
	CardDef3321004{},
	CardDef3321005{},
	CardDef3321006{},
	CardDef3321007{},
	CardDef3321008{},
	CardDef3321009{},
	CardDef3321010{},
	CardDef3321011{},
	CardDef3321012{},
	CardDef3321013{},
	CardDef3321014{},
	CardDef3321015{},
	CardDef3321101{},
	CardDef3321102{},
	CardDef3321103{},
	CardDef3321104{},
	CardDef3321105{},
	CardDef3321106{},
	CardDef3321107{},
	CardDef3321108{},
	CardDef3321109{},
	CardDef3321110{},
	CardDef3411101{},
	CardDef3411102{},
	CardDef3421001{},
	CardDef3421002{},
	CardDef3421003{},
	CardDef3421004{},
	CardDef3421005{},
	CardDef3421006{},
	CardDef3421007{},
	CardDef3421008{},
	CardDef3421009{},
	CardDef3421010{},
	CardDef3421011{},
	CardDef3421012{},
	CardDef3421013{},
	CardDef3421014{},
	CardDef3421015{},
	CardDef3421101{},
	CardDef3421102{},
	CardDef3421103{},
	CardDef3421104{},
	CardDef3421105{},
	CardDef3421106{},
	CardDef3421107{},
	CardDef3421108{},
	CardDef3421109{},
	CardDef3421110{},
	CardDef3501001{},
	CardDef3501101{},
	CardDef3511010{},
	CardDef3511101{},
	CardDef3511102{},
	CardDef3521001{},
	CardDef3521002{},
	CardDef3521003{},
	CardDef3521004{},
	CardDef3521005{},
	CardDef3521006{},
	CardDef3521007{},
	CardDef3521008{},
	CardDef3521009{},
	CardDef3521011{},
	CardDef3521012{},
	CardDef3521013{},
	CardDef3521014{},
	CardDef3521015{},
	CardDef3521101{},
	CardDef3521102{},
	CardDef3521103{},
	CardDef3521104{},
	CardDef3521105{},
	CardDef3521106{},
	CardDef3521107{},
	CardDef3521108{},
	CardDef3521109{},
	CardDef3521110{},
	CardDef3601101{},
	CardDef3611101{},
	CardDef3611102{},
	CardDef3621001{},
	CardDef3621002{},
	CardDef3621003{},
	CardDef3621004{},
	CardDef3621005{},
	CardDef3621006{},
	CardDef3621007{},
	CardDef3621008{},
	CardDef3621009{},
	CardDef3621010{},
	CardDef3621011{},
	CardDef3621012{},
	CardDef3621013{},
	CardDef3621014{},
	CardDef3621015{},
	CardDef3621101{},
	CardDef3621102{},
	CardDef3621103{},
	CardDef3621104{},
	CardDef3621105{},
	CardDef3621106{},
	CardDef3621107{},
	CardDef3621108{},
	CardDef3621109{},
	CardDef3621110{},
	CardDef4011001{},
	CardDef4011002{},
	CardDef4011101{},
	CardDef4011102{},
	CardDef4111001{},
	CardDef4111002{},
	CardDef4111003{},
	CardDef4111101{},
	CardDef4111102{},
	CardDef4211001{},
	CardDef4211002{},
	CardDef4211003{},
	CardDef4211101{},
	CardDef4211102{},
	CardDef4311001{},
	CardDef4311002{},
	CardDef4311003{},
	CardDef4311101{},
	CardDef4311102{},
	CardDef4411001{},
	CardDef4411002{},
	CardDef4411003{},
	CardDef4411101{},
	CardDef4411102{},
	CardDef4511001{},
	CardDef4511002{},
	CardDef4511003{},
	CardDef4511101{},
	CardDef4511102{},
	CardDef4611001{},
	CardDef4611002{},
	CardDef4611003{},
	CardDef4611101{},
	CardDef4611102{},
}
