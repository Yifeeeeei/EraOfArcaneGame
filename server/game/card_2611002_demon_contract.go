package game

import "eraofarcane/model"

type Card2611002DemonContract struct{ AlwaysActive }

func (Card2611002DemonContract) ID() string   { return "2611002" }
func (Card2611002DemonContract) Name() string { return "与恶魔的契约书" }

func (Card2611002DemonContract) OnUseItem(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	sacrifices := ctx.Engine.friendlyUnits(ctx.PlayerID, true, nil)
	if len(sacrifices) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "demon_contract_sacrifice", "选择1个友方单位作为献祭对象", sacrifices, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			sacrifice := ctx.Engine.findFieldCardByInstance(ps, selected[0])
			if sacrifice == nil {
				return
			}
			targets := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
				return card.Card.IsCompanion() &&
					card.Position != nil &&
					ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false) &&
					ps.CanPayCost(demonContractExtraCost(sacrifice, card))
			})
			if len(targets) == 0 {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "demon_contract_destroy", "选择1个敌方伙伴并支付生命差额费用", targets, 1, 1,
				func(selected []string) {
					if len(selected) == 0 {
						return
					}
					target := ctx.Engine.findFieldCardByInstance(ctx.Engine.State.Players[ctx.OpponentID], selected[0])
					if target == nil {
						return
					}
					if target.Position == nil || !ctx.Engine.IsInSpellRange(ctx.PlayerID, target.Position.Col, target.Position.Row, false) {
						return
					}
					if !ps.PayCost(demonContractExtraCost(sacrifice, target)) {
						return
					}
					ctx.Engine.destroyUnitWithCause(sacrifice, ctx.PlayerID, DeathCauseSacrifice)
					ctx.Engine.destroyUnit(target, ctx.OpponentID)
					shuffleDemonContractIntoDeck(ctx)
				})
		})
	return nil
}

func demonContractExtraCost(sacrifice *CardInstance, target *CardInstance) map[string]int {
	if sacrifice == nil || target == nil {
		return map[string]int{}
	}
	diff := abs(sacrifice.CurrentLife - target.CurrentLife)
	if diff <= 0 {
		return map[string]int{}
	}
	return map[string]int{model.ElementShadow: diff * 2}
}

func shuffleDemonContractIntoDeck(ctx *EffectContext) {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	for i, card := range ps.Graveyard {
		if card != nil && card.InstanceID == ctx.Source.InstanceID {
			ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
			ps.Deck = append(ps.Deck, card)
			ctx.Engine.shuffleDeck(ctx.PlayerID)
			return
		}
	}
}
