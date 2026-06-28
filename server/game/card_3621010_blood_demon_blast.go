package game

import "eraofarcane/model"

type Card3621010BloodDemonBlast struct{ AlwaysActive }

func (Card3621010BloodDemonBlast) ID() string   { return "3621010" }
func (Card3621010BloodDemonBlast) Name() string { return "血魔爆" }

func (Card3621010BloodDemonBlast) OnSpellCast(ctx *EffectContext) error {
	if !isSpellBeingCast(ctx) {
		return nil
	}
	frontRow := ctx.Engine.State.Players[ctx.PlayerID].GetFrontRow()
	if frontRow < 0 {
		return nil
	}
	candidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() && card.Card.Category == model.ElementShadow && card.Position != nil && card.Position.Row == frontRow
	})
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_demon_blast",
		"血魔爆:献祭1个前排暗影伙伴", candidates, 1, 1,
		func(selected []string) {
			player := ctx.Engine.State.Players[ctx.PlayerID]
			unit := ctx.Engine.findFieldCardByInstance(player, firstSelected(selected))
			if unit == nil || unit.Position == nil || unit.Position.Row != frontRow {
				return
			}
			damage := max(unit.CurrentLife, 0)
			ctx.Engine.destroyUnitWithCause(unit, ctx.PlayerID, DeathCauseSacrifice)
			if damage <= 0 {
				return
			}
			targets := ctx.Engine.enemyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
				return card.Position != nil && ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false)
			})
			if len(targets) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_demon_blast_target", "血魔爆:选择法力范围内1个敌人造成伤害", targets, 1, 1, func(targetSelected []string) {
				target := selectedUnitFromCandidates(ctx.Engine, targetSelected, targets)
				if target != nil {
					ctx.Engine.dealDamageWithExtra(target, damage, ctx.OpponentID, map[string]any{"attacker": ctx.PlayerID})
				}
			})
		})
	return nil
}
