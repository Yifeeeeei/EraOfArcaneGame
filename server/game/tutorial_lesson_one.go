package game

import "fmt"

// This is an authored lesson, not an AI policy or a second rules engine. Every
// write below is an ordinary validated game action, including all payments.
func (e *Engine) firstLessonSteps() []tutorialStep {
	steps := []tutorialStep{}
	add := func(step tutorialStep) {
		step.id = fmt.Sprintf("first-%03d", len(steps)+1)
		steps = append(steps, step)
	}
	read := func(chapter, title, text, focus, card string) {
		add(tutorialStep{chapter: chapter, title: title, text: text, instruction: "看清这一步说明后，点击“继续”。", focus: focus, card: card})
	}
	move := func(player int, chapter, title, text, instruction, focus, card, action string, data func(*Engine) map[string]any) {
		add(tutorialStep{player: player, chapter: chapter, title: title, text: text, instruction: instruction, focus: focus, card: card, action: action, data: data})
	}
	fixed := func(data map[string]any) func(*Engine) map[string]any {
		return func(*Engine) map[string]any { return data }
	}
	cardData := func(role string, payment map[string]int) func(*Engine) map[string]any {
		return func(e *Engine) map[string]any {
			d := map[string]any{"instance_id": e.tutorial.roles[role]}
			if len(payment) > 0 {
				d["payment"] = payment
			}
			return d
		}
	}
	consume := func(role, name, text string) {
		move(0, "取得元素", "消耗"+name, text, "点击高亮卡牌，选择“消耗”。", "card", role, "consume", cardData(role, nil))
	}
	learn := func(role, title, text string, payment map[string]int) {
		move(0, "学习技能", title, text, "打开技能池，选择高亮技能并点击“学习”，按提示支付。", "skills", role, "learn_skill", cardData(role, payment))
	}
	cast := func(role, title, text string, cost int) {
		move(0, "法术进攻", title, text, "选择高亮法术施放，点击敌方中央人物为目标；不添加强化，按提示支付并确认。", "cast", role, "cast_spell", func(e *Engine) map[string]any {
			return map[string]any{"instance_id": e.tutorial.roles[role], "target_type": "unit", "target_col": 1, "target_row": 1, "target_owner": 1, "payment": map[string]int{"火": cost}}
		})
	}
	end := func(text string) {
		move(0, "回合流程", "结束自己的回合", text, "点击“结束回合”。", "end_turn", "", "end_turn", fixed(map[string]any{}))
	}
	op := func(title, role, action string, data func(*Engine) map[string]any) {
		move(1, "对手行动", title, "观察对手按同样的规则行动；轮到需要你决定时，教学会暂停。", "等待对手完成这一步。", "", role, action, data)
	}
	opConsume := func() { op("对手取得元素", "enemy", "consume", cardData("enemy", nil)) }
	opEnd := func(discard bool) {
		op("对手结束回合", "", "end_turn", fixed(map[string]any{}))
		if discard {
			op("对手按上限弃牌", "", "resolve_action", func(e *Engine) map[string]any {
				pa := e.State.PendingAction
				if pa == nil || pa.PlayerID != 1 || len(pa.Candidates) == 0 {
					return map[string]any{}
				}
				return map[string]any{"selected": []string{pa.Candidates[0]["instance_id"].(string)}}
			})
		}
	}
	opCast := func() {
		op("对手使用冰锥术", "ice", "cast_spell", func(e *Engine) map[string]any {
			return map[string]any{"instance_id": e.tutorial.roles["ice"], "target_type": "unit", "target_col": 1, "target_row": 1, "target_owner": 0, "payment": map[string]int{"水": 1}}
		})
	}
	opNoDefend := func() { op("对手选择不防御", "", "no_defend", fixed(map[string]any{})) }
	defend := func() {
		move(0, "法术防御", "按四个步骤完成防御", "右侧防御面板分为①防御、②强化、③透支、④付款。烈焰护体威力3，足以抵挡冰锥术；本次跳过强化，透支火焰面甲支付1火。选择期间不会立即扣费，到最后确认才一起结算。", "按右侧当前操作提示依次选择；卡牌和下一步按钮会用金框标出。", "defense", "shield", "defend", func(e *Engine) map[string]any {
			return map[string]any{"skill_ids": []string{e.tutorial.roles["shield"]}, "overexert_ids": []string{e.tutorial.roles["mask"]}, "payment": map[string]int{"火": 1}}
		})
	}

	read("认识对局", "你的目标：让敌方人物生命归零", "这是一场完整的教学对局。双方人物从6生命开始。将敌方人物的生命降到0就获胜；本教程会一直带你操作到胜利。人物特殊能力留到后续课程。", "hero", "hero")
	read("认识对局", "先认识双方战场", "双方各有3×3单位区，人物从中央开始。伙伴会放在己方空格；你的伙伴本教程放在人物后方。人物与伙伴都叫单位。", "board", "")
	read("认识对局", "手牌来自主卡组", "正常构筑有1个人物、30张主卡和10张技能。主卡组中的伙伴与道具会抽到手里。起手4张，以后每个自己的回合通常抽1张。暂不用的牌在本教程会收起。", "hand", "apprentice")
	read("认识对局", "技能不是从手牌抽来的", "技能池是可学习的技能，技能区是已学习的技能。先付学习费用放入技能区，以后使用再付使用费用。技能区初始有5个位置，本教程只用其中3个。", "skills", "")
	read("认识牌面", "读懂一张伙伴牌", "巫师的学徒：名称帮助识别，类型是伙伴，属性是奥术，标签是巫师。入场费用1奥术，负载1奥术，生命1，没有额外效果。牌号、版本和引言不必记忆。", "card", "apprentice")
	read("认识牌面", "费用与负载是不同的数字", "费用是打出这张牌要花什么；负载是它在场上被消耗时能产生什么。元素池才是你现在实际可以花的元素，负载不会自动加入池中。", "elements", "apprentice")
	read("认识元素", "先认识火焰与奥术", "七种元素是火焰、水纹、大气、地脉、光辉、暗影、奥术；图标旁数字表示数量。本教程主要用火与奥术。火不能直接替代水；奥术费用可用任意元素支付，池中的奥术也可支付普通的火等需求。", "elements", "")
	move(0, "准备开局", "保留这次起手", "正常对局可以整手重抽一次。本教程为了固定流程保留起手。双方保留后进入第一回合，先手仍然会抽1张。", "点击“保留手牌”。", "mulligan", "", "mulligan", fixed(map[string]any{"keep": true}))
	op("对手保留起手", "", "mulligan", fixed(map[string]any{"keep": true}))
	read("第一次回合", "先手首回合人物负载折半", "你已抽到第5张手牌。梵天印刷负载是4火，但先手首回合消耗人物只取得一半，即2火；其他回合会恢复完整负载。", "hero", "hero")
	consume("hero", "人物", "消耗把竖置卡变为横置，并取得负载。卡牌仍在场上，不是弃置；同一张横置牌不能再消耗。本次元素池会从0变为2火。")
	read("第一次回合", "用元素支付伙伴入场费用", "现在池中有2火。学徒需要1奥术，可以用1火支付。伙伴默认横置入场，所以刚进场还不能再拿它取得负载。", "hand", "apprentice")
	move(0, "第一次回合", "召唤巫师的学徒", "从手牌选择学徒，花1火支付其1奥术入场费用。把它放在人物后方的左侧空格。", "选择高亮手牌“召唤”，再点高亮的己方后排左格并确认。", "summon", "apprentice", "summon", func(e *Engine) map[string]any {
		return map[string]any{"instance_id": e.tutorial.roles["apprentice"], "col": 0, "row": 2, "payment": map[string]int{"火": 1}}
	})
	read("第一次回合", "道具中的装备进入装备区", "火焰面甲是道具，子类是装备／饰物。入场花费1火，负载1火；它没有生命，也没有额外效果。装备放在装备区，默认同样横置入场。", "hand", "mask")
	move(0, "第一次回合", "装备火焰面甲", "用池中最后1火装备面甲。它以后也能消耗取得负载；本回合刚进场时先横置。", "选择高亮面甲并点击“装备”。", "equip", "mask", "equip", cardData("mask", map[string]int{"火": 1}))
	read("回合流程", "自己的回合结束才重置", "本回合的元素已花完。结束时会先检查手牌上限，再尝试重置自己的卡牌，处理状态与结束效果，清空元素后换人。本局还没有状态效果。", "end_turn", "")
	end("点击结束后，看人物、学徒和面甲从横置恢复为竖置，再进入对手回合。")
	read("观察对手", "轮到对手主动行动", "对手已经抽1张。一般只有当前玩家能召唤、学习或主动施法；你可以随时看公开卡牌。对方将消耗人物取得4水，学习冰锥术，再结束回合。", "opponent", "")
	opConsume()
	op("对手学习冰锥术", "ice", "learn_skill", cardData("ice", map[string]int{"水": 2}))
	opEnd(false)
	read("取得元素", "第二回合：场上卡牌提供负载", "人物本次可以取得完整4火，面甲能取得1火，学徒能取得1奥术。依次消耗，元素会在池中累积到5火和1奥术；它们不会自己自动产费。", "elements", "")
	consume("hero", "人物", "本次不再折半，取得4火。")
	consume("mask", "面甲", "装备也能提供负载。取得1火，与人物产生的元素累积。")
	consume("apprentice", "学徒", "学徒取得1奥术。本教程每次花费前都会明确指出如何支付。")
	read("学习技能", "学习费、使用费、威力和攻击", "爆焰一击学习费5火，使用费3火。威力7用于比较能否突破防御；攻击2用于命中后的伤害，通常不会因为威力是7就造成7伤害。缺少数值也不等于印刷为0。", "skills", "blast")
	learn("blast", "学习爆焰一击", "本步支付5火，保留池中1奥术。技能从技能池进入技能区，默认横置，本回合暂时不能使用。", map[string]int{"火": 5})
	read("学习技能", "新技能先横置", "已经学会不等于现在能用。技能要等自己的回合结束重置后，才能在合法时机支付使用费用施放。使用后仍留在技能区，不会像一次性道具那样用掉。", "card", "blast")
	learn("shield", "用奥术学习防御法术", "烈焰护体学习1火、使用1火、威力3、攻击0。印刷的独立关键词“防御”表示它只能用于防御或强化防御。本步用最后1奥术支付1火学习费。", map[string]int{"无": 1})
	end("两张新学技能会在自己的结束流程重置。接下来将在对手的法术窗口中使用烈焰护体。")
	opConsume()
	opCast()
	read("法术防御", "先看清来袭法术", "冰锥术以你的人物为目标：威力3，攻击1。若不防御会造成1伤害。现在用威力也是3的烈焰护体，达到对方威力即可成功防御。", "incoming", "ice")
	defend()
	read("法术防御", "防御成功，人物没有失去生命", "双方威力相等，已成功挡住。烈焰护体和透支的面甲都横置。透支只为这一笔费用提供元素，不触发普通消耗，也不会把元素留在池里。", "card", "mask")
	opEnd(true)
	read("主动进攻", "透支的卡还没有重置", "现在是你的第三回合。面甲在对手回合被透支，必须等你自己的回合结束才重置，所以本回合不能再用它产费。改用人物4火与学徒1奥术。", "card", "mask")
	consume("hero", "人物", "取得4火，用于本回合的施法与学习。")
	consume("apprentice", "学徒", "取得1奥术，稍后可以补足火球术的学习费用。")
	cast("blast", "第一次主动施法", "已学爆焰一击现在竖置，可以花3火使用，不需要再付5火学费。对手只有中央人物，它就在前排，是普通法术的合法目标。", 3)
	opNoDefend()
	read("主动进攻", "威力7，实际伤害是2", "没有被成功防御就叫法术命中。对手本次不防御，爆焰一击以攻击2造成伤害，敌方生命从6降到4。技能使用后横置，下次可重置后还能使用。", "card", "blast")
	learn("fireball", "学习下一张普通攻击法术", "火球术学习2火，使用1火，威力3、攻击1。本步用剩余1火和1奥术支付2火学习费；它横置进入技能区。", map[string]int{"火": 1, "无": 1})
	end("本回合新学的火球术会重置，之后可以反复使用，无需重复付学习费用。")
	opConsume()
	opCast()
	read("法术防御", "看看不防御会发生什么", "本次仍是威力3、攻击1的冰锥术。现在为了认识这个选项，我们选择不防御；本次按攻击1扣减人物生命。", "incoming", "ice")
	move(0, "法术防御", "选择不防御", "不使用防御主法术，本次攻击会命中。这里只演示选项的结果，不要求判断战术。", "点击“不防御”。", "no_defend", "", "no_defend", fixed(map[string]any{}))
	read("法术防御", "人物生命从6变为5", "这是冰锥术的1点攻击伤害。没有造成其他状态，后续你仍可以按流程行动。", "hero", "hero")
	opEnd(true)
	read("连续行动", "手牌6张，可以先正常出牌", "本回合抽牌后手牌有6张。默认上限5张是在结束回合的弃牌阶段检查，并非刚抽到第6张就立即弃牌。本回合先完成两次普通施法。", "hand", "")
	consume("hero", "人物", "取得4火。")
	consume("mask", "面甲", "面甲已在你的上一次回合结束重置，现在可以取得1火。")
	consume("apprentice", "学徒", "取得1奥术，本回合合计5火和1奥术。")
	cast("fireball", "使用已经学会的火球术", "只需1火使用费；选择同一个敌方人物，不添加强化。对手将选择不防御。", 1)
	opNoDefend()
	read("连续行动", "出牌阶段不只允许做一件事", "火球术造成1伤害，敌方生命4变3。你还有元素，也有另一张竖置技能，可以继续执行一次独立的施法。强化等其他用法留到下一课。", "skills", "blast")
	cast("blast", "再使用一次爆焰一击", "使用费3火，命中造成2伤害；这是另一笔独立支付、另一次独立法术交互。", 3)
	opNoDefend()
	read("回合结束", "这次完整处理弃牌和剩余元素", "敌方人物剩1生命。你现在仍有6张手牌，池中剩1火和1奥术。结束时先弃1张到上限5张，随后重置并清空剩余元素。", "end_turn", "")
	end("点击结束回合后，先完成右侧出现的“弃牌至手牌上限”选择。")
	move(0, "回合结束", "弃掉手里的第二张面甲", "这是弃至上限的正式阶段。指定弃掉的是手中的第二张面甲，不是场上已经装备的那张。", "在右侧只选择高亮的面甲，然后确认弃牌。", "pending", "spare_mask", "resolve_action", func(e *Engine) map[string]any {
		return map[string]any{"selected": []string{e.tutorial.roles["spare_mask"]}}
	})
	read("回合结束", "弃牌、重置、清空，再换人", "手牌从6变5，自己的卡牌已重置，剩余1火和1奥术已清空。状态与结束效果本教程没有额外内容，真实流程仍已经过这些阶段。现在轮到对手。", "elements", "")
	opConsume()
	opCast()
	defend()
	opEnd(true)
	read("最后一击", "完成真正的胜利条件", "敌方人物剩1生命。你抽牌后暂时6张手牌，仍可先出牌。使用火球术造成最后1伤害，即刻获胜；不用先结束回合或先弃牌。", "hero", "enemy")
	consume("hero", "人物", "取得4火，本次火球术只需要其中1火。")
	cast("fireball", "让敌方人物生命归零", "支付1火，选择敌方人物。对手不防御后，火球术的攻击1会完成击杀。", 1)
	opNoDefend()
	return steps
}
