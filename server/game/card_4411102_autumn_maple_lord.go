package game

import "eraofarcane/model"

type Card4411102AutumnMapleLord struct{ AlwaysActive }

func (Card4411102AutumnMapleLord) ID() string   { return "4411102" }
func (Card4411102AutumnMapleLord) Name() string { return "秋枫领主 狄利克雷" }
func (Card4411102AutumnMapleLord) OnOverexertPayment(ctx *EffectContext, event OverexertPaymentEvent) {
	e, ps := ctx.Engine, ctx.Engine.State.Players[ctx.PlayerID]
	units, payment, poolSpent := event.Units, event.Payment, event.PoolSpent
	overexertSpent := make(map[string]int)
	for elem, amount := range payment {
		if amount <= 0 {
			continue
		}
		overexertSpent[elem] = amount - poolSpent[elem]
	}
	reward := make(map[string]int)
	for _, unit := range units {
		if unit == nil || unit.Card == nil || unit.Card.Category != model.ElementEarth {
			continue
		}
		for elem, amount := range e.effectiveElementsGain(unit) {
			remaining := amount
			if remaining <= 0 {
				continue
			}
			used := min(overexertSpent[elem], remaining)
			overexertSpent[elem] -= used
			remaining -= used
			if remaining > 0 {
				reward[elem] += remaining * 2
			}
		}
	}
	for elem, amount := range reward {
		if amount > 0 {
			ps.Elements[elem] += amount
		}
	}
	if len(reward) > 0 {
		e.emit(GameEvent{
			Type:   "autumn_maple_lord_overexert_reward",
			Player: ps.PlayerID,
			Data: map[string]any{
				"player": ps.PlayerID,
				"source": cardToInfo(ctx.Source),
				"reward": reward,
			},
		})
	}
}
