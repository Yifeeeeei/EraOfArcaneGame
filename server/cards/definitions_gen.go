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

var compiledCardDefinitions = []CardDefinition{
	CardDef1011001{},
	CardDef1011002{},
	CardDef1011003{},
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
	CardDef1111001{},
	CardDef1111002{},
	CardDef1111003{},
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
	CardDef1211001{},
	CardDef1211002{},
	CardDef1211003{},
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
	CardDef1311001{},
	CardDef1311002{},
	CardDef1311003{},
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
	CardDef1401001{},
	CardDef1401002{},
	CardDef1411001{},
	CardDef1411002{},
	CardDef1411003{},
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
	CardDef1501001{},
	CardDef1511001{},
	CardDef1511002{},
	CardDef1511003{},
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
	CardDef1611001{},
	CardDef1611002{},
	CardDef1611003{},
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
	CardDef2011001{},
	CardDef2011002{},
	CardDef2011003{},
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
	CardDef2111001{},
	CardDef2111002{},
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
	CardDef2211001{},
	CardDef2211002{},
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
	CardDef2311001{},
	CardDef2311002{},
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
	CardDef2411001{},
	CardDef2411002{},
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
	CardDef2501001{},
	CardDef2511001{},
	CardDef2511002{},
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
	CardDef2601001{},
	CardDef2601002{},
	CardDef2611001{},
	CardDef2611002{},
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
	CardDef3001001{},
	CardDef3001002{},
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
	CardDef3101001{},
	CardDef3101002{},
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
	CardDef3201001{},
	CardDef3201002{},
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
	CardDef3301001{},
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
	CardDef3501001{},
	CardDef3511010{},
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
	CardDef4011001{},
	CardDef4011002{},
	CardDef4111001{},
	CardDef4111002{},
	CardDef4111003{},
	CardDef4211001{},
	CardDef4211002{},
	CardDef4211003{},
	CardDef4311001{},
	CardDef4311002{},
	CardDef4311003{},
	CardDef4411001{},
	CardDef4411002{},
	CardDef4411003{},
	CardDef4511001{},
	CardDef4511002{},
	CardDef4511003{},
	CardDef4611001{},
	CardDef4611002{},
	CardDef4611003{},
}
