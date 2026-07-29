package game

const boneKnightRebornStatus = "失去遗言"

type Card1621011BoneKnight struct{ AlwaysActive }

func (Card1621011BoneKnight) ID() string   { return "1621011" }
func (Card1621011BoneKnight) Name() string { return "白骨骑士" }

func (Card1621011BoneKnight) HasActiveDeathrattle(card *CardInstance) bool {
	return card != nil && card.Statuses[boneKnightRebornStatus] <= 0
}

func (Card1621011BoneKnight) OnDeath(ctx *EffectContext) error {
	if ctx.Source == nil || ctx.Source.Position == nil {
		return nil
	}
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	pos := *ctx.Source.Position
	if !pos.Valid() || ps.Units[pos.Col][pos.Row] != nil {
		return nil
	}
	for i, card := range ps.Graveyard {
		if card == ctx.Source {
			ps.Graveyard = append(ps.Graveyard[:i], ps.Graveyard[i+1:]...)
			break
		}
	}
	ctx.Source.CurrentLife = ctx.Source.Card.Life
	ctx.Source.Statuses[boneKnightRebornStatus] = 1
	ctx.Source.IsHorizontal = true
	ctx.Source.EnterTurn = ctx.Engine.State.TurnNumber
	ps.Units[pos.Col][pos.Row] = ctx.Source
	ctx.Engine.ApplySummonModifiersOnEnter(ctx.Source)
	ctx.Engine.emit(GameEvent{Type: "summon", Player: ctx.PlayerID, Data: map[string]any{
		"player":   ctx.PlayerID,
		"card":     cardToInfo(ctx.Source),
		"position": pos,
	}})
	return nil
}
