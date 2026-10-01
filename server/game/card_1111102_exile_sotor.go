package game

type Card1111102ExileSotor struct{ AlwaysActive }

func (Card1111102ExileSotor) ID() string   { return "1111102" }
func (Card1111102ExileSotor) Name() string { return "\"流放者\" 索拓尔" }
func (Card1111102ExileSotor) ExpandSpellTargets(ctx *EffectContext, target SpellTarget, extraTargets []SpellTarget) []SpellTarget {
	e, playerID := ctx.Engine, ctx.PlayerID
	if e == nil || target.Type != "unit" || !target.Position.Valid() || playerID < 0 || playerID >= len(e.State.Players) {
		return extraTargets
	}
	targetOwnerID := 1 - playerID
	if target.OwnerID != nil {
		targetOwnerID = *target.OwnerID
	}
	if targetOwnerID < 0 || targetOwnerID >= len(e.State.Players) {
		return extraTargets
	}
	for _, delta := range []struct{ col, row int }{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		extra := SpellTarget{Type: "unit", Position: Position{Col: target.Position.Col + delta.col, Row: target.Position.Row + delta.row}}
		if !extra.Position.Valid() || e.State.Players[targetOwnerID].Units[extra.Position.Col][extra.Position.Row] == nil {
			continue
		}
		ownerID := targetOwnerID
		extra.OwnerID = &ownerID
		if spellTargetsContain(extraTargets, extra) {
			continue
		}
		extraTargets = append(extraTargets, extra)
	}
	return extraTargets
}
