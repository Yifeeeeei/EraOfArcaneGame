package game

import (
	"fmt"
)

type Card2321101ThunderChain struct{ AlwaysActive }

func (Card2321101ThunderChain) ID() string { return "2321101" }

func (Card2321101ThunderChain) Name() string { return "雷之链" }

func (Card2321101ThunderChain) PerTurnLabel(*CardInstance) string {
	return "充能"
}

func (Card2321101ThunderChain) OnPerTurn(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil {
		return nil
	}
	if ctx.Source.IsHorizontal {
		return fmt.Errorf("雷之链需要竖置才能消耗")
	}
	ctx.Source.IsHorizontal = true
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModNextDriveSpellExtraTarget,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		RemainingUses:    1,
	})
	return nil
}

var _ PerTurnAbility = Card2321101ThunderChain{}
