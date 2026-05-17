package game

import "fmt"

type Card1521005TwinAngel struct{ AlwaysActive }

func (Card1521005TwinAngel) ID() string   { return "1521005" }
func (Card1521005TwinAngel) Name() string { return "双生天使" }

func (Card1521005TwinAngel) OnEnter(ctx *EffectContext) error {
	card := getCardDB()["1501001"]
	if card == nil {
		return fmt.Errorf("missing twin angel token card 1501001")
	}
	instance := NewCardInstance(card, ctx.PlayerID, ctx.Engine.State.TurnNumber)
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	ps.Hand = append(ps.Hand, instance)
	ctx.Engine.emit(GameEvent{
		Type:   "effect_trigger",
		Player: ctx.PlayerID,
		Data: map[string]any{
			"source": cardToInfo(ctx.Source),
			"effect": "create_card_in_hand",
			"card":   cardToInfo(instance),
		},
	})
	return nil
}
