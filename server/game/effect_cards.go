package game

import (
	"log"
	"math/rand"
	"strings"
)

// RegisterAllCardEffects registers manually-coded card effects.
// Complex cards that can't be auto-parsed need explicit handlers here.
func RegisterAllCardEffects() {
	r := GetEffectRegistry()

	registerHeroEffects(r)
	registerCompanionEffects(r)
	registerSkillEffects(r)
	registerItemEffects(r)

	log.Printf("[Effects] Registered manual card effects")
}

// =============================================
// HERO EFFECTS (人物)
// =============================================
func registerHeroEffects(r *EffectRegistry) {

	// 4011001 "南境百灵" 斯卡尔蒂 罗佳
	// 回合技:丢弃1张手牌,获得2点该卡牌属性种类的元素
	r.RegisterActive("4011001", TriggerPerTurn, func(ctx *EffectContext) error {
		// This is an active ability, handled by "use_ability" action
		// Implementation: discard a hand card, gain 2 of its element
		return nil
	})

	// 4011002 "无面"
	// 每当你让1张与你场上原有卡牌属性相同的卡牌入场,你受到1点伤害
	r.Register("4011002", TriggerOnUnitEnter, func(ctx *EffectContext) error {
		if ctx.Target == nil || ctx.Source == nil {
			return nil
		}
		newCard := ctx.Target
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		// Check if any other card on field shares the same category
		allCards := ctx.Engine.getAllFieldCards(ps)
		for _, c := range allCards {
			if c == newCard {
				continue
			}
			if c.Card.Category == newCard.Card.Category {
				// Deal 1 damage to hero
				if ps.Hero != nil {
					ctx.Engine.dealDamage(ps.Hero, 1, ctx.PlayerID)
					ctx.Engine.emit(GameEvent{
						Type:   "effect_trigger",
						Player: -1,
						Data: map[string]any{
							"source": cardToInfo(ctx.Source),
							"effect": "same_element_penalty",
							"damage": 1,
						},
					})
				}
				return nil
			}
		}
		return nil
	})

	// 4011101 纯净灵体 奥希斯
	// 你每使1张奥术以外的卡牌入场,使你的所有法术获得虚弱2
	r.Register("4011101", TriggerOnUnitEnter, func(ctx *EffectContext) error {
		if ctx.Target == nil {
			return nil
		}
		if ctx.Target.Card.Category != "无" {
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			for i := 0; i < 5; i++ {
				if ps.Skills[i] != nil {
					ps.Skills[i].Statuses[StatusWeaken] += 2
				}
			}
		}
		return nil
	})

	// 4011102 大法师 罗慕路斯
	// 此卡提供的元素只能严格当做\无使用
	// (This is a passive constraint, handled during element payment - needs engine support)

	// 4311001 雷术士 肃
	// 绝技:丢弃2张大气手牌,对任意1名敌人造成1点伤害
	r.RegisterActive("4311001", TriggerUltimate, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		// Find 2 air hand cards to discard
		discarded := 0
		for i := len(ps.Hand) - 1; i >= 0 && discarded < 2; i-- {
			if ps.Hand[i].Card.Category == "气" {
				ps.Graveyard = append(ps.Graveyard, ps.Hand[i])
				ps.Hand = append(ps.Hand[:i], ps.Hand[i+1:]...)
				discarded++
			}
		}
		if discarded < 2 {
			return nil // not enough air cards
		}
		// Deal 1 damage to random enemy
		opponent := ctx.Engine.State.Players[ctx.OpponentID]
		target := findAnyUnit(opponent)
		if target != nil {
			ctx.Engine.dealDamage(target, 1, ctx.OpponentID)
		}
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4311002 "渡鸦" 睿文
	// 你的起始手牌数与换牌机会+1 (handled during setup/mulligan)

	// 4311003 掌门 穆伶
	// 绝技:选择法力范围内双方各1个伙伴,花费差值气,移回手牌
	r.RegisterActive("4311003", TriggerUltimate, func(ctx *EffectContext) error {
		// Complex - needs target selection from UI
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4111001 掌门 龙卷火
	// 游戏开始前,将万火合一术置于你的技能池 (handled during setup)

	// 4111002 女巫 维兰德
	// 回合技:你获得点燃1,然后本回合将此卡负载的1\火变为1\无
	r.RegisterActive("4111002", TriggerPerTurn, func(ctx *EffectContext) error {
		ctx.Source.Statuses[StatusBurn] += 1
		// Temporarily change 1 fire to 1 arcane in load (simplified)
		return nil
	})

	// 4111003 大祭司 梵天
	// 绝技:本回合每当你的火焰法术命中,此卡获得负载+1\火
	r.RegisterActive("4111003", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		// Set a flag that fire spells hitting grant +1 fire load
		return nil
	})

	// 4211101 海神之使 珊瑚 贝莉
	// 本局游戏你第一次使用法术攻击时,使该法术永久+3\威
	r.Register("4211101", TriggerOnSpellCast, func(ctx *EffectContext) error {
		if ctx.Target != nil && ctx.Target.Card.IsSkill() {
			// First spell only - check if already triggered
			if ctx.Source.Statuses["first_spell_used"] == 0 {
				ctx.Source.Statuses["first_spell_used"] = 1
				// Boost the spell's power (simplified)
			}
		}
		return nil
	})

	// 4211102 凛冰魔巫 索菲娅
	// 此卡不受冻结效果影响
	// 绝技:移除场上任意1个单位的1层冻结,对其造成2点伤害
	r.RegisterActive("4211102", TriggerUltimate, func(ctx *EffectContext) error {
		if ctx.Target != nil {
			if ctx.Target.Statuses[StatusFreeze] > 0 {
				ctx.Target.Statuses[StatusFreeze]--
			}
			ctx.Engine.dealDamage(ctx.Target, 2, ctx.Target.OwnerID)
		}
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4411001 森林隐士 白须
	// 检索1张地属性野兽/植物/精灵代替首回合抽牌
	r.Register("4411001", TriggerOnTurnStart, func(ctx *EffectContext) error {
		if ctx.Engine.State.TurnNumber == 1 {
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			// Search deck for earth beast/plant/elf
			for i, c := range ps.Deck {
				if c.Card.Category == "地" && hasTag(c.Card.Tag, "野兽", "植物", "精灵") {
					ps.Hand = append(ps.Hand, c)
					ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
					shuffleDeck(ps.Deck)
					ctx.Engine.emit(GameEvent{
						Type:   "effect_trigger",
						Player: ctx.PlayerID,
						Data: map[string]any{
							"source": cardToInfo(ctx.Source),
							"effect": "search",
							"card":   cardToInfo(c),
						},
					})
					break
				}
			}
		}
		return nil
	})

	// 4411003 麦吉教授
	// 你的打出或学习的第一张原始费用大于5的卡牌费用减少2\地
	// (Passive - needs engine cost modification support)

	// 4611001 暗影学者 爱莉斯
	// 回合技:当1个你的伙伴死亡,使你的1个法术+1\威
	r.Register("4611001", TriggerOnFriendlyDeath, func(ctx *EffectContext) error {
		// Boost a random friendly spell's power
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		for i := 0; i < 5; i++ {
			if ps.Skills[i] != nil && ps.Skills[i].Card.Power > 0 {
				// Simplified: can't permanently boost power in current model
				// We'd need a PowerBoost field
				break
			}
		}
		return nil
	})

	// 4611002 芙雅夫人
	// 绝技:使你的1个伙伴攻击和负载翻倍,但会在回合结束时死亡
	r.RegisterActive("4611002", TriggerUltimate, func(ctx *EffectContext) error {
		if ctx.Target != nil {
			ctx.Target.CurrentAttack *= 2
			// Mark for death at end of turn
			ctx.Target.Statuses["临时"] = 1
		}
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4511001 圣使 玛丽斯 南森埃尔
	// 绝技:当敌方造成伤害时,直到下回合结束,每当友方单位受伤获得2光
	r.RegisterActive("4511001", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4311101 司天魔巫 索兰德
	// 绝技:使驱动和聚能法术永久+1\威,其他法术永久-2\威
	r.RegisterActive("4311101", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4311102 布雾者 弗格
	// 绝技:双方各自召唤的下1个伙伴入场时获得隐蔽2
	r.RegisterActive("4311102", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4311201 星术士 莉莉安
	// 回合技:将此卡移动至另一位置
	r.RegisterActive("4311201", TriggerPerTurn, func(ctx *EffectContext) error {
		// Needs target position from UI
		return nil
	})

	// 4311202 舞女 特拉奇
	// 绝技:花费2\气,双方各抽2张牌,或双方随机弃2张牌
	r.RegisterActive("4311202", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4411101 翡翠男爵 杰德 拜利兰
	// 当你的护盾小于3时不会在回合结束时减少 (passive - shield decay check)

	// 4411102 秋枫领主 狄利克雷
	// 当你的地脉卡牌透支时,你可以获得剩余的元素 (passive - overdraft modifier)

	// 4211003 凛冬城主 水晶心
	// 绝技:本回合剩余时间内,技能区内法术获得冻结1
	r.RegisterActive("4211003", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4611003 咒言师 结影
	// 游戏开始前,将三张咒言书洗入卡组 (handled during setup)

	// 4611101 鲜血伯爵 休伯特 黑松
	// 游戏开始前,将鲜血盛宴置于技能池 (handled during setup)

	// 4611102 灾厄玫瑰 多姆
	// 游戏开始前,从双方卡组上方4张牌送去弃牌堆
	// (handled during setup)

	// 4511301 洛迪总督
	// 单位上限+2,额外伙伴可共用位置 (passive - needs engine support)

	// 4411201 林地之王 希松
	// 对方一回合内召唤超2个时,每超1个-1\地 (passive listener)
	r.Register("4411201", TriggerOnUnitEnter, func(ctx *EffectContext) error {
		// Count opponent summons this turn (simplified)
		return nil
	})

	// 4411202 自然学家 多萝西
	// 召唤第一个野兽/植物/精灵后抽1张牌
	r.Register("4411202", TriggerOnUnitEnter, func(ctx *EffectContext) error {
		if ctx.Target == nil {
			return nil
		}
		tag := ctx.Target.Card.Tag
		if hasTag(tag, "野兽", "植物", "精灵") {
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			// Check if this is the first of its type this game
			key := "first_" + getFirstMatchingTag(tag, "野兽", "植物", "精灵")
			if ctx.Source.Statuses[key] == 0 {
				ctx.Source.Statuses[key] = 1
				drawn := ps.DrawCards(1)
				if len(drawn) > 0 {
					ctx.Engine.emit(GameEvent{
						Type:   "draw_card",
						Player: ctx.PlayerID,
						Data:   map[string]any{"card": cardToInfo(drawn[0])},
					})
				}
			}
		}
		return nil
	})

	// 4211201 深渊归来者 多里安
	// 打出的第1张伙伴牌入场花费减少1\水或1\暗
	// (passive cost modifier - needs engine support)

	// 4511201 "光影交错" 亚基尔
	// 游戏开始后,将"相生境"置于手牌 (handled during setup)

	// 4611201 灵媒师 梅
	// 前两个学习花费>3的灵媒/神秘技能花费减1 (passive cost modifier)

	// 4011201 驭龙者 阿卡托斯
	// 游戏开始后,将驭龙呼唤加入手牌.每控制1条龙,负载+1\无(最多+2)
	// (handled during setup + passive load modifier)

	// 4011202 卡珊教授
	// 游戏开始前选择属性,负载只能当该属性,额外技能携带量
	// (handled during setup)

	// 4011301 诺叩少爷
	// 来自原本技能池的技能无法造成伤害 (passive damage restriction)

	// 4111201 炎狱灵法师 鲁多克拉克
	// 法术战斗威力超3点时获得护盾1和1\火
	r.Register("4111201", TriggerOnDefend, func(ctx *EffectContext) error {
		// Check power difference during spell combat
		return nil
	})

	// 4111101 首席顾问 费林
	// 绝技:献祭1个火焰伙伴,下1个火焰卡牌减费
	r.RegisterActive("4111101", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4111102 大将军 克兰
	// 绝技:选择火焰法术,未命中时获得其攻的火
	r.RegisterActive("4111102", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4211001 "浪之人" 巴特尔
	// 绝技:下一张手牌属性视为水
	r.RegisterActive("4211001", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4211002 大贤者 沃尔波特
	// 游戏开始前,将百川归海置于技能 (handled during setup)

	// 4511002 神之眷子 爱里默
	// 游戏开始前将桎梏置于对手牌组 (handled during setup)

	// 4511003 骑士团长 蕾曦娅
	// 游戏开始前,替换技能池中的"希望呼唤" (handled during setup)

	// 4511101 庇护者 西瓦尔
	// 绝技:当受到3+伤害,下回合结束前所有友方不受伤
	r.RegisterActive("4511101", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4511102 救赎者 伊芙 秋枫
	// 游戏开始前从卡组移出1张光辉牌 (handled during setup)

	// 4311301 东海炼气士 元鸿
	// 游戏开始前将3张符纸置于卡组 (handled during setup)

	// 4311302 云梦仙子 玉玲
	// 绝技:将手牌/场上2张牌洗回卡组,抽2张牌
	r.RegisterActive("4311302", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 4611202 观者 月曦
	// 每当使用暗影技能,查看卡组顶1张.回合技:可以重洗
	r.Register("4611202", TriggerOnSpellCast, func(ctx *EffectContext) error {
		// Check if spell is shadow type
		return nil
	})
}

// =============================================
// COMPANION EFFECTS (伙伴)
// =============================================
func registerCompanionEffects(r *EffectRegistry) {
	// 1621112 谧语精灵猎手 - 遗言:对任意1个敌人造成1点伤害
	r.Register("1621112", TriggerOnDeath, func(ctx *EffectContext) error {
		opponent := ctx.Engine.State.Players[ctx.OpponentID]
		target := findAnyUnit(opponent)
		if target != nil {
			ctx.Engine.dealDamage(target, 1, ctx.OpponentID)
		}
		return nil
	})

	// 1621113 谧语精灵祭司 - 遗言:使1个友方伙伴负载+1\暗
	r.Register("1621113", TriggerOnDeath, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		for col := 0; col < 3; col++ {
			for row := 0; row < 3; row++ {
				unit := ps.Units[col][row]
				if unit != nil && unit != ctx.Source && unit.Card.IsCompanion() {
					if unit.Card.ElementsGain == nil {
						unit.Card.ElementsGain = make(map[string]int)
					}
					unit.Card.ElementsGain["暗"]++
					return nil
				}
			}
		}
		return nil
	})

	// 1621114 灵魂共生体 - 遗言:给最多2个法术放灵魂标记物,每个+2威
	r.Register("1621114", TriggerOnDeath, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		count := 0
		for i := 0; i < 5 && count < 2; i++ {
			if ps.Skills[i] != nil {
				ps.Skills[i].Statuses["灵魂标记"] += 1
				count++
			}
		}
		return nil
	})

	// 1621103 鲜血傀儡 - 入场:对你的人物造成2点伤害
	r.Register("1621103", TriggerOnEnter, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		if ps.Hero != nil {
			ctx.Engine.dealDamage(ps.Hero, 2, ctx.PlayerID)
		}
		return nil
	})

	// 1221004 寒霜傀儡 - 入场:对法力范围内1个敌方伙伴冻结1
	r.Register("1221004", TriggerOnEnter, func(ctx *EffectContext) error {
		if ctx.Target != nil {
			ctx.Target.Statuses[StatusFreeze] += 1
		} else {
			// Auto-target front row enemy
			opponent := ctx.Engine.State.Players[ctx.OpponentID]
			frontRow := opponent.GetFrontRow()
			if frontRow >= 0 {
				for col := 0; col < 3; col++ {
					if opponent.Units[col][frontRow] != nil {
						opponent.Units[col][frontRow].Statuses[StatusFreeze] += 1
						break
					}
				}
			}
		}
		return nil
	})

	// 1321308 浮空飞行员 - 入场:充能1
	r.Register("1321308", TriggerOnEnter, func(ctx *EffectContext) error {
		ctx.Engine.addCharge(ctx.PlayerID, 1)
		return nil
	})

	// 1321306 太极宗师 - 回合技:将场上或手牌1张卡洗回卡组,抽1张牌
	r.RegisterActive("1321306", TriggerPerTurn, func(ctx *EffectContext) error {
		// Needs target selection from UI
		return nil
	})

	// 1221207 重构体 - 回合技:在你检索1张卡牌后,获得负载+1水
	r.RegisterActive("1221207", TriggerPerTurn, func(ctx *EffectContext) error {
		return nil
	})

	// 1221208 狩猎者C3型巡洋舰 - 回合技:移动至其他位置
	r.RegisterActive("1221208", TriggerPerTurn, func(ctx *EffectContext) error {
		// Needs target position from UI
		return nil
	})

	// 1321301 门派外聘技师 - 绝技:本回合每当洗回卡组1张牌,充能1
	r.RegisterActive("1321301", TriggerUltimate, func(ctx *EffectContext) error {
		ctx.Source.UltimateUsed = true
		return nil
	})

	// 1321309 太极传人 - 入场:从卡组检索1张调和、化劲、收势
	r.Register("1321309", TriggerOnEnter, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		targets := []string{"调和", "化劲", "收势"}
		for i, c := range ps.Deck {
			for _, name := range targets {
				if strings.Contains(c.Card.Name, name) {
					ps.Hand = append(ps.Hand, c)
					ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
					shuffleDeck(ps.Deck)
					ctx.Engine.emit(GameEvent{
						Type:   "effect_trigger",
						Player: ctx.PlayerID,
						Data: map[string]any{
							"source": cardToInfo(ctx.Source),
							"effect": "search",
							"card":   cardToInfo(c),
						},
					})
					return nil
				}
			}
		}
		return nil
	})

	// 1221203 深渊布雾魔 - 入场:检索1张水属性或深渊地形牌
	r.Register("1221203", TriggerOnEnter, func(ctx *EffectContext) error {
		ps := ctx.Engine.State.Players[ctx.PlayerID]
		for i, c := range ps.Deck {
			if c.Card.Category == "水" || (HasKeyword(c.Card.Description, "深渊") && HasKeyword(c.Card.Description, "地形")) {
				ps.Hand = append(ps.Hand, c)
				ps.Deck = append(ps.Deck[:i], ps.Deck[i+1:]...)
				shuffleDeck(ps.Deck)
				ctx.Engine.emit(GameEvent{
					Type:   "effect_trigger",
					Player: ctx.PlayerID,
					Data: map[string]any{
						"source": cardToInfo(ctx.Source),
						"effect": "search",
						"card":   cardToInfo(c),
					},
				})
				break
			}
		}
		return nil
	})

	// 1321304 屠魔兵器 云霄井阑 - 引魔.你的具有威力的卡牌获得穿透
	r.Register("1321304", TriggerPassive, func(ctx *EffectContext) error {
		// Passive: all power-bearing cards get pierce
		// This is checked during spell targeting
		return nil
	})

	// 1611103 鲜血贵公子 罗伯特 黑松
	// 友方单位受友方伤害放标记,友方因友方死亡放2标记
	// 回合技:移除3标记,+1血或负载+1暗或+1攻
	r.Register("1611103", TriggerOnFriendlyDeath, func(ctx *EffectContext) error {
		ctx.Source.Statuses["鲜血标记"] += 2
		return nil
	})
	r.RegisterActive("1611103", TriggerPerTurn, func(ctx *EffectContext) error {
		if ctx.Source.Statuses["鲜血标记"] >= 3 {
			ctx.Source.Statuses["鲜血标记"] -= 3
			// Choose randomly: +1 life, +1 shadow load, or +1 attack
			choices := []int{0, 1, 2}
			choice := choices[rand.Intn(3)]
			switch choice {
			case 0:
				ctx.Source.CurrentLife++
			case 1:
				if ctx.Source.Card.ElementsGain == nil {
					ctx.Source.Card.ElementsGain = make(map[string]int)
				}
				ctx.Source.Card.ElementsGain["暗"]++
			case 2:
				ctx.Source.CurrentAttack++
			}
		}
		return nil
	})
}

// =============================================
// SKILL EFFECTS (技能)
// =============================================
func registerSkillEffects(r *EffectRegistry) {
	// 3521005 星陨 - 穿透 (simple keyword, handled by keyword system)

	// 3421102 苍岚之刃 - 范围:前排 (targets all enemy front row)
	// This is a targeting modifier, checked during spell resolution

	// 3421103 巨岩崩落 - 无法用于强化
	// This restricts the spell from being used as boost

	// 3621206 混沌吞噬 - 献祭你的1个生物:回复其生命值的生命,或获得光和暗
	r.Register("3621206", TriggerOnSpellHit, func(ctx *EffectContext) error {
		// Sacrifice mechanic needs UI target selection
		return nil
	})

	// 3621301 充能回收 - 冷却1.异能:每当1个单位死亡,充能1
	r.Register("3621301", TriggerOnFriendlyDeath, func(ctx *EffectContext) error {
		ctx.Engine.addCharge(ctx.PlayerID, 1)
		return nil
	})
	r.Register("3621301", TriggerOnEnemyDeath, func(ctx *EffectContext) error {
		ctx.Engine.addCharge(ctx.PlayerID, 1)
		return nil
	})
}

// =============================================
// ITEM EFFECTS (道具)
// =============================================
func registerItemEffects(r *EffectRegistry) {
	// 2521101 赐福之孤星 - 使1个友方伙伴获得负载+1\光和+1\血
	r.Register("2521101", TriggerOnEquip, func(ctx *EffectContext) error {
		// Needs target selection
		if ctx.Target != nil {
			if ctx.Target.Card.ElementsGain == nil {
				ctx.Target.Card.ElementsGain = make(map[string]int)
			}
			ctx.Target.Card.ElementsGain["光"]++
			ctx.Target.CurrentLife++
		}
		return nil
	})

	// 2521104 黄金龙骨 - 献祭:抽2张牌
	r.Register("2521104", TriggerOnEquip, func(ctx *EffectContext) error {
		// Sacrifice trigger - when used as sacrifice, draw 2
		return nil
	})

	// 2511101 九霄辉迹 - 绝技:双方将手牌全部丢弃,然后抽等量的牌
	r.RegisterActive("2511101", TriggerUltimate, func(ctx *EffectContext) error {
		for i := 0; i < 2; i++ {
			ps := ctx.Engine.State.Players[i]
			handSize := len(ps.Hand)
			// Discard all
			ps.Graveyard = append(ps.Graveyard, ps.Hand...)
			ps.Hand = make([]*CardInstance, 0)
			// Draw same amount
			drawn := ps.DrawCards(handSize)
			for _, c := range drawn {
				ctx.Engine.emit(GameEvent{
					Type:   "draw_card",
					Player: i,
					Data:   map[string]any{"card": cardToInfo(c)},
				})
			}
		}
		ctx.Source.UltimateUsed = true
		return nil
	})
}

// =============================================
// HELPER FUNCTIONS
// =============================================

func hasTag(tag string, tags ...string) bool {
	for _, t := range tags {
		if strings.Contains(tag, t) {
			return true
		}
	}
	return false
}

func getFirstMatchingTag(tag string, tags ...string) string {
	for _, t := range tags {
		if strings.Contains(tag, t) {
			return t
		}
	}
	return ""
}
