package game

import "fmt"

type Card1321001RavenMessenger struct{}

func (Card1321001RavenMessenger) ID() string   { return "1321001" }
func (Card1321001RavenMessenger) Name() string { return "渡鸦信使" }
func (Card1321001RavenMessenger) OnPerTurn(ctx *EffectContext) error {
	if ctx.Source.IsHorizontal {
		return fmt.Errorf("渡鸦信使已经横置")
	}
	ctx.Source.IsHorizontal = true
	return DrawCards(1)(ctx)
}
