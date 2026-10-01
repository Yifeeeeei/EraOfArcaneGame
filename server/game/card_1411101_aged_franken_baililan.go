package game

import (
	"eraofarcane/model"
)

type Card1411101AgedFrankenBaililan struct{ AlwaysActive }

func (Card1411101AgedFrankenBaililan) ID() string { return "1411101" }

func (Card1411101AgedFrankenBaililan) Name() string { return "苍老者 弗兰肯 拜利兰" }

func (Card1411101AgedFrankenBaililan) IsPrayerAbility() bool { return true }

func (Card1411101AgedFrankenBaililan) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	load := dragonBloodTreantReducibleLoad(ctx.Source)
	candidates := make([]map[string]any, 0)
	for _, elem := range model.AllElements {
		if load[elem] <= 0 {
			continue
		}
		candidates = append(candidates, map[string]any{
			"instance_id": elem,
			"name":        elem,
			"zone":        "choice",
			"side":        "own",
		})
	}
	removeLoad := func(elem string) {
		if elem == "" || reducibleElementLoad(ctx.Source, elem) <= 0 {
			return
		}
		ctx.Engine.reduceCardElementLoadWithTriggers(ctx.PlayerID, ctx.Source, elem, 1, ctx.Source)
	}
	if len(candidates) == 0 {
		return nil
	}
	if len(candidates) == 1 {
		elem, _ := candidates[0]["instance_id"].(string)
		removeLoad(elem)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "aged_franken_prayer_load",
		"苍老者 弗兰肯 拜利兰:选择失去1点负载", candidates, 1, 1,
		func(selected []string) {
			removeLoad(firstSelected(selected))
		})
	return nil
}

func (Card1411101AgedFrankenBaililan) OnLoadLoss(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	lossTarget, _ := ctx.ExtraData["load_loss_target"].(*CardInstance)
	if lossTarget != ctx.Source || !ctx.Engine.cardStillOnField(ctx.Source) {
		return nil
	}
	candidates := ctx.Engine.enemySkills(ctx.PlayerID, nil)
	if len(candidates) == 0 {
		return nil
	}
	weakenSkill := func(selection string) {
		target, _ := ctx.Engine.findFriendlyCandidate(ctx.OpponentID, selection)
		if target == nil || target.Card == nil || !target.Card.IsSkill() {
			return
		}
		target.PowerBonus -= 2
		ctx.Engine.emit(GameEvent{
			Type:   "aged_franken_weaken_spell",
			Player: -1,
			Data: map[string]any{
				"player": ctx.PlayerID,
				"source": cardToInfo(ctx.Source),
				"target": cardToInfo(target),
				"power":  -2,
			},
		})
	}
	if len(candidates) == 1 {
		id, _ := candidates[0]["instance_id"].(string)
		weakenSkill(id)
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "aged_franken_weaken_spell",
		"苍老者 弗兰肯 拜利兰:选择1个敌方法术永久-2威", candidates, 1, 1,
		func(selected []string) {
			weakenSkill(firstSelected(selected))
		})
	return nil
}

var _ PrayerAbility = Card1411101AgedFrankenBaililan{}

var _ PerTurnAbility = Card1411101AgedFrankenBaililan{}

var _ OnLoadLossBehavior = Card1411101AgedFrankenBaililan{}
