package game

import "fmt"

type Card1211101MistKingNalanti struct{ AlwaysActive }

func (Card1211101MistKingNalanti) ID() string   { return "1211101" }
func (Card1211101MistKingNalanti) Name() string { return "雾之国主 那兰提" }

func (Card1211101MistKingNalanti) OnEnter(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			card := ps.Units[col][row]
			if card != nil && card.Position != nil && !ctx.Engine.hasActiveStealth(card) {
				ctx.Engine.grantStealth(card, 2)
			}
		}
	}
	return nil
}

type Card1221102MistMage struct{ AlwaysActive }

func (Card1221102MistMage) ID() string   { return "1221102" }
func (Card1221102MistMage) Name() string { return "雾之国法师" }

func (Card1221102MistMage) OnPerTurn(ctx *EffectContext) error {
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card != ctx.Source && card.Position != nil && !ctx.Engine.hasActiveStealth(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "mist_mage_stealth", "雾之国法师:选择1个没有隐蔽的其他友方单位隐蔽2", candidates, 1, 1, func(selected []string) {
		target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
		if target != nil {
			ctx.Engine.grantStealth(target, 2)
		}
	})
	return nil
}

type Card1221105MistDancer struct{ AlwaysActive }

func (Card1221105MistDancer) ID() string   { return "1221105" }
func (Card1221105MistDancer) Name() string { return "雾之国舞女" }

func (Card1221105MistDancer) OnEnter(ctx *EffectContext) error {
	candidates := companionSpellRangeCandidates(ctx, false)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "mist_dancer_stealth", "雾之国舞女:选择法力范围内1个伙伴隐蔽2", candidates, 1, 1, func(selected []string) {
		target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
		if target != nil {
			ctx.Engine.grantStealth(target, 2)
		}
	})
	return nil
}

type Card1221109MistPhantom struct{ AlwaysActive }

func (Card1221109MistPhantom) ID() string   { return "1221109" }
func (Card1221109MistPhantom) Name() string { return "雾霭幽魂" }

func (Card1221109MistPhantom) OnEnter(ctx *EffectContext) error {
	ctx.Engine.grantStealth(ctx.Source, 3)
	return nil
}

type Card1321104MistWeaver struct{ AlwaysActive }

func (Card1321104MistWeaver) ID() string   { return "1321104" }
func (Card1321104MistWeaver) Name() string { return "织雾者" }

func (Card1321104MistWeaver) OnEnter(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return card.Position != nil && !ctx.Engine.hasActiveStealth(card)
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "mist_weaver_stealth", "织雾者:选择任意1个敌方单位隐蔽2", candidates, 1, 1, func(selected []string) {
		target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
		if target != nil {
			ctx.Engine.grantStealth(target, 2)
		}
	})
	return nil
}

type Card1421114GiantSandworm struct{ AlwaysActive }

func (Card1421114GiantSandworm) ID() string   { return "1421114" }
func (Card1421114GiantSandworm) Name() string { return "巨型沙虫" }

func (Card1421114GiantSandworm) OnDamaged(ctx *EffectContext) error {
	ctx.Engine.grantStealth(ctx.Source, 1)
	return nil
}

type Card2021103MistPotion struct{ AlwaysActive }

func (Card2021103MistPotion) ID() string   { return "2021103" }
func (Card2021103MistPotion) Name() string { return "幻雾药剂" }

func (Card2021103MistPotion) OnUseItem(ctx *EffectContext) error {
	candidates := companionSpellRangeCandidates(ctx, false)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "mist_potion_stealth", "幻雾药剂:选择法力范围内1个伙伴隐蔽2", candidates, 1, 1, func(selected []string) {
		target := selectedUnitFromCandidates(ctx.Engine, selected, candidates)
		if target != nil {
			ctx.Engine.grantStealth(target, 2)
		}
	})
	return nil
}

type Card3221104WaterEscape struct{ AlwaysActive }

func (Card3221104WaterEscape) ID() string   { return "3221104" }
func (Card3221104WaterEscape) Name() string { return "水遁术" }

func (Card3221104WaterEscape) AllowsFriendlySpellTarget() bool { return true }

func (Card3221104WaterEscape) ValidateSpellTarget(ctx *EffectContext, _ SpellTarget, target *CardInstance) error {
	if target == nil {
		return nil
	}
	if ctx != nil && ctx.Engine != nil && ctx.Engine.hasActiveStealth(target) {
		return fmt.Errorf("target already has stealth")
	}
	return nil
}

func (Card3221104WaterEscape) SpellHitStatuses(*EffectContext) map[string]int {
	return map[string]int{StatusStealth: 2}
}

type Card3221106Undercurrent struct{ AlwaysActive }

func (Card3221106Undercurrent) ID() string   { return "3221106" }
func (Card3221106Undercurrent) Name() string { return "暗流涌动" }

func (Card3221106Undercurrent) AllowsStealthSpellTarget() bool { return true }

func (Card3221106Undercurrent) ModifySkillContribution(ctx *EffectContext, stats *SpellStats) {
	target, _ := ctx.ExtraData["spell_target_unit"].(*CardInstance)
	if target != nil && ctx.Engine.hasActiveStealth(target) {
		stats.PowerBonus += 2
	}
}

type Card4311102MistmakerFug struct{ AlwaysActive }

func (Card4311102MistmakerFug) ID() string   { return "4311102" }
func (Card4311102MistmakerFug) Name() string { return "布雾者 弗格" }

func (Card4311102MistmakerFug) OnUltimate(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].NextCompanionStealth = 2
	ctx.Engine.State.Players[ctx.OpponentID].NextCompanionStealth = 2
	return nil
}

func companionSpellRangeCandidates(ctx *EffectContext, requireNoStealth bool) []map[string]any {
	matches := func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() && (!requireNoStealth || !ctx.Engine.hasActiveStealth(card))
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, matches)
	candidates = append(candidates, ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return matches(card) && card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
	})...)
	return candidates
}

var _ FriendlySpellTargetBehavior = Card3221104WaterEscape{}
var _ SpellTargetValidationBehavior = Card3221104WaterEscape{}
var _ SpellHitStatusBehavior = Card3221104WaterEscape{}
var _ StealthSpellTargetBehavior = Card3221106Undercurrent{}
var _ SkillContributionModifier = Card3221106Undercurrent{}
