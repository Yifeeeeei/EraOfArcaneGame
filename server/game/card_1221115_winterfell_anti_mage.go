package game

import (
	"eraofarcane/model"
)

type Card1221115WinterfellAntiMage struct{ AlwaysActive }

func (Card1221115WinterfellAntiMage) ID() string { return "1221115" }

func (Card1221115WinterfellAntiMage) Name() string { return "凛冬城御魔师" }

func (Card1221115WinterfellAntiMage) OnPrayer(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for _, skill := range ps.Skills {
		if skill == nil || skill.Card == nil || !skill.Card.IsSkill() {
			continue
		}
		ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
			Type:             TempModNextSkillUseCostMinus,
			SourceCardNumber: ctx.Source.Card.Number,
			SourceName:       ctx.Source.Card.Name,
			TargetInstanceID: skill.InstanceID,
			Element:          model.ElementWater,
			Amount:           1,
			RemainingUses:    1,
		})
	}
	return nil
}
