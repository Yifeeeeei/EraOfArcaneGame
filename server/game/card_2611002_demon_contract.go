package game

import (
	"fmt"

	"eraofarcane/model"
)

type Card2611002DemonContract struct{ AlwaysActive }

func (Card2611002DemonContract) ID() string   { return "2611002" }
func (Card2611002DemonContract) Name() string { return "与恶魔的契约书" }

func (Card2611002DemonContract) OnUseItem(ctx *EffectContext) error {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	sacrifices := ctx.Engine.friendlyUnits(ctx.PlayerID, true, func(card *CardInstance) bool {
		return len(demonContractPayableTargets(ctx, ps, card)) > 0
	})
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
			targets := demonContractPayableTargets(ctx, ps, sacrifice)
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
					cost := demonContractExtraCost(sacrifice, target)
					if len(cost) == 0 {
						resolveDemonContract(ctx, sacrifice, target)
						return
					}
					candidate := candidateInfo(target, "unit", "enemy")
					ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "demon_contract_payment",
						"与恶魔的契约书:支付生命差额费用", []map[string]any{candidate}, 1, 1, cost, false,
						func(selected []string, data map[string]any) error {
							if len(selected) == 0 || selected[0] != target.InstanceID {
								return fmt.Errorf("invalid demon contract target")
							}
							if sacrifice == nil || target == nil || sacrifice.CurrentLife <= 0 || target.CurrentLife <= 0 {
								return fmt.Errorf("demon contract target no longer valid")
							}
							if !ctx.Engine.payCostForAction(ps, cost, ActionMessage{Data: data}) {
								return fmt.Errorf("invalid demon contract payment")
							}
							resolveDemonContract(ctx, sacrifice, target)
							return nil
						})
				})
		})
	return nil
}

func demonContractPayableTargets(ctx *EffectContext, ps *PlayerState, sacrifice *CardInstance) []map[string]any {
	if ctx == nil || ctx.Engine == nil || ps == nil || sacrifice == nil {
		return nil
	}
	return ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card.Card.IsCompanion() &&
			card.Position != nil &&
			ctx.Engine.IsInSpellRange(ctx.PlayerID, card.Position.Col, card.Position.Row, false) &&
			ctx.Engine.canPayCost(ps, demonContractExtraCost(sacrifice, card))
	})
}

func (e *Engine) demonContractHasPayablePathAfterEntryCost(playerID int, contract *CardInstance) bool {
	if contract == nil || contract.Card == nil {
		return false
	}
	ps := e.State.Players[playerID]
	entryCost := e.effectiveCardPlayCost(ps, contract)
	opponent := e.State.Players[1-playerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			sacrifice := ps.Units[col][row]
			if sacrifice == nil {
				continue
			}
			for targetCol := 0; targetCol < 3; targetCol++ {
				for targetRow := 0; targetRow < 3; targetRow++ {
					target := opponent.Units[targetCol][targetRow]
					if target == nil || target.Card == nil || !target.Card.IsCompanion() || target.Position == nil {
						continue
					}
					if !e.IsInSpellRange(playerID, target.Position.Col, target.Position.Row, false) {
						continue
					}
					totalCost := mergeElementCosts(entryCost, demonContractExtraCost(sacrifice, target))
					if _, ok := calculateElementPaymentWithOptions(ps.Elements, totalCost, e.playerHasLightWildcard(ps)); ok {
						return true
					}
				}
			}
		}
	}
	return false
}

func resolveDemonContract(ctx *EffectContext, sacrifice *CardInstance, target *CardInstance) {
	if ctx == nil || ctx.Engine == nil || sacrifice == nil || target == nil {
		return
	}
	ctx.Engine.destroyUnitWithCause(sacrifice, ctx.PlayerID, DeathCauseSacrifice)
	ctx.Engine.destroyUnit(target, ctx.OpponentID)
	shuffleDemonContractIntoDeck(ctx)
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
