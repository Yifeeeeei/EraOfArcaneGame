package game

import "eraofarcane/model"

type Card4311001Su struct{ AlwaysActive }

func (Card4311001Su) ID() string   { return "4311001" }
func (Card4311001Su) Name() string { return "雷术士 肃" }
func (Card4311001Su) OnUltimate(ctx *EffectContext) error {
	airCards := ctx.Engine.friendlyHandCards(ctx.PlayerID, func(card *CardInstance) bool {
		return card.Card.Category == model.ElementAir
	})
	if len(airCards) < 2 {
		return nil
	}
	targets := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(targets) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "su_discard_air",
		"雷术士 肃:选择2张大气手牌作为费用", airCards, 2, 2,
		func(selected []string) {
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			for _, id := range selected {
				for i, card := range ps.Hand {
					if card == nil || card.InstanceID != id || card.Card.Category != model.ElementAir {
						continue
					}
					ctx.Engine.discardHandCardAt(ctx.PlayerID, i)
					break
				}
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "su_damage_enemy",
				"雷术士 肃:选择任意1名敌人造成1点伤害", targets, 1, 1,
				func(selected []string) {
					target := selectedUnitFromCandidates(ctx.Engine, selected, targets)
					if target != nil {
						ctx.Engine.dealDamage(target, 1, ctx.OpponentID)
					}
				})
		})
	return nil
}
