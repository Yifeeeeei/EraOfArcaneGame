package game

import "eraofarcane/model"

func (Card2121102FireCloudFan) SpellTargetGrant(ctx *EffectContext, skill *CardInstance, target SpellTarget) SpellTargetGrant {
	e, playerID := ctx.Engine, ctx.PlayerID
	if e == nil || skill == nil || skill.Card == nil || target.Type != "unit" || !target.Position.Valid() {
		return SpellTargetGrant{}
	}
	if skill.Card.Category != model.ElementFire && skill.Card.Category != model.ElementAir {
		return SpellTargetGrant{}
	}
	opponent := e.State.Players[1-playerID]
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 || target.Position.Row <= frontRow {
		return SpellTargetGrant{}
	}
	frontOfTarget := target.Position.Row - 1
	return SpellTargetGrant{IgnoreRange: frontOfTarget >= 0 && opponent.Units[target.Position.Col][frontOfTarget] == nil}
}

type Card2121102FireCloudFan struct{ AlwaysActive }

func (Card2121102FireCloudFan) ID() string   { return "2121102" }
func (Card2121102FireCloudFan) Name() string { return "火云扇" }
