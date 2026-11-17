package game

import (
	"fmt"
)

type Card1321001RavenMessenger struct{ AlwaysActive }

func (Card1321001RavenMessenger) AbilityRequiresTap(trigger EffectTrigger) bool {
	return trigger == TriggerPerTurn
}

func (Card1321001RavenMessenger) ID() string { return "1321001" }

func (Card1321001RavenMessenger) Name() string { return "渡鸦信使" }

func (Card1321001RavenMessenger) PerTurnLabel(*CardInstance) string {
	return "消耗抽牌"
}

func (Card1321001RavenMessenger) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source.IsHorizontal {
		return fmt.Errorf("渡鸦信使已经横置")
	}
	if !ctx.Engine.canConsumeCard(ctx.Source) {
		return fmt.Errorf("渡鸦信使不能支付消耗费用")
	}
	ctx.Engine.runResolution("raven messenger consume and draw",
		func() {
			ctx.Engine.consumeCardForEffectWithTriggers(ctx.PlayerID, ctx.Source, nil, ctx.Source.Card.Number)
		},
		func() { _ = DrawCards(1)(ctx) },
	)
	return nil
}

func (Card1321001RavenMessenger) RepeatableActiveAbility() bool { return true }
