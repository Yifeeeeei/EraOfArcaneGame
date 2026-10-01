package game

import (
	"eraofarcane/model"
)

type Card2111101DivineFireStaffCrimsonSky struct{ AlwaysActive }

func (Card2111101DivineFireStaffCrimsonSky) ID() string { return "2111101" }

func (Card2111101DivineFireStaffCrimsonSky) Name() string { return "神火杖 赤空" }

func (Card2111101DivineFireStaffCrimsonSky) OnSpellHit(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	attacker, _ := ctx.ExtraData["attacker"].(int)
	if attacker != ctx.PlayerID || ctx.Target.Card.Category != model.ElementFire || !isSpellLikeCard(ctx.Target.Card) || !ctx.Engine.canConsumeCard(ctx.Source) {
		return nil
	}
	staff := ctx.Source
	skill := ctx.Target
	ctx.Engine.SetPendingAction(ctx.PlayerID, "divine_fire_staff_empower_spell",
		"神火杖 赤空:是否消耗此卡永久强化命中的火焰法术", []map[string]any{candidateInfo(staff, "equipment", "own")}, 0, 1,
		func(selected []string) {
			if len(selected) == 0 || !ctx.Engine.canConsumeCard(staff) || !ctx.Engine.cardStillOnField(staff) {
				return
			}
			if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) || skill.Card.Category != model.ElementFire {
				return
			}
			ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, staff, ctx.Engine.effectiveElementsGain(staff), "2111101")
			skill.PowerBonus++
			skill.Statuses[permanentPierceStatus] = 1
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source":  cardToInfo(staff),
				"target":  cardToInfo(skill),
				"effect":  "permanent_fire_spell_empower",
				"power":   1,
				"pierce":  true,
				"consume": true,
			}})
		})
	return nil
}

const permanentPierceStatus = "永久穿透"
