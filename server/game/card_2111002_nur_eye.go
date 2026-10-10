package game

import (
	"eraofarcane/model"
)

type Card2111002NurEye struct{ AlwaysActive }

func (Card2111002NurEye) ID() string { return "2111002" }

func (Card2111002NurEye) Name() string { return "努尔之眼" }

func (Card2111002NurEye) IsPrayerAbility() bool { return true }

func (Card2111002NurEye) DamageScope() DamageScope { return DamageAny }

func (Card2111002NurEye) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	if event.IsFire() {
		ctx.Source.Statuses[nurEyeFireMark]++
		ctx.Source.Statuses["火焰标记"]++
	}
	return nil
}

func (Card2111002NurEye) OnPerTurn(ctx *EffectContext) error {
	newMarkers := ctx.Source.Statuses[nurEyeFireMark]
	ctx.Source.Statuses[nurEyeFireMark] = 0
	ctx.Source.Statuses["火焰标记"] = 0
	switch newMarkers {
	case 0:
		ctx.Engine.discardFriendlyCandidate(ctx.PlayerID, ctx.Source.InstanceID)
		return nil
	case 1:
		ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementFire: 2})
	case 2:
		ctx.Engine.addNextElementSpellPowerBonus(ctx.PlayerID, model.ElementFire, 2)
	case 3:
		ctx.Engine.addNextElementSpellDamageBonus(ctx.PlayerID, model.ElementFire, 1)
	default:
		candidates := append(ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil), ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)...)
		ctx.Engine.SetPendingAction(ctx.PlayerID, "nur_eye_fire_damage", "努尔之眼:选择1个单位造成2点火焰伤害", candidates, 1, 1, func(selected []string) {
			for _, ps := range ctx.Engine.State.Players {
				target := ctx.Engine.findUnitOnGrid(ps, firstSelected(selected))
				if target != nil {
					ctx.Engine.ApplyDamage(DamageRequest{Target: target, Amount: 2, Kind: "effect", Element: model.ElementFire, Source: ctx.Source})
					return
				}
			}
		})
	}
	return nil
}
