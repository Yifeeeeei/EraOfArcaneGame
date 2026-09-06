package game

type Card1421111RockWallMonster struct{ AlwaysActive }

func (Card1421111RockWallMonster) ID() string { return "1421111" }

func (Card1421111RockWallMonster) Name() string { return "岩壁魔怪" }

func (Card1421111RockWallMonster) ModifyDamageAmount(ctx *EffectContext, amount int) int {
	if ctx == nil || ctx.Engine == nil || amount <= 1 || playerHasLearnedSpell(ctx.Engine.State.Players[ctx.PlayerID]) {
		return amount
	}
	return 1
}

func playerHasLearnedSpell(ps *PlayerState) bool {
	if ps == nil {
		return false
	}
	for _, skill := range ps.Skills {
		if skill != nil && skill.Card != nil && hasCardTag(skill.Card, "法术") {
			return true
		}
	}
	return false
}

var _ DamageAmountModifier = Card1421111RockWallMonster{}
