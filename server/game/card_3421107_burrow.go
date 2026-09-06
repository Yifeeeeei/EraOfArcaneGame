package game

import "fmt"

type Card3421107Burrow struct{ AlwaysActive }

func (Card3421107Burrow) ID() string { return "3421107" }

func (Card3421107Burrow) Name() string { return "破土而出" }

func (Card3421107Burrow) MasteryMax() int { return 2 }

func (Card3421107Burrow) OnMastery(*EffectContext, int) error { return nil }

var _ MasteryBehavior = Card3421107Burrow{}

func (Card3421107Burrow) PrepareSpellCast(ctx *EffectContext, target SpellTarget, action ActionMessage) (SpellCastOptions, error) {
	return SpellCastOptions{ExtraTargets: func(bool) ([]SpellTarget, error) {
		return ctx.Engine.burrowExtraTargetsFromAction(ctx.PlayerID, ctx.Source, target, action)
	}}, nil
}

func (e *Engine) burrowExtraTargetsFromAction(playerID int, skill *CardInstance, mainTarget SpellTarget, action ActionMessage) ([]SpellTarget, error) {
	if skill == nil || skill.Card == nil || skill.Card.Number != "3421107" {
		return nil, nil
	}
	raw, ok := action.Data["extra_targets"].([]any)
	if !ok || len(raw) == 0 {
		return nil, nil
	}
	maxTargets := skill.Statuses[StatusMastery]
	if maxTargets > 2 {
		maxTargets = 2
	}
	if maxTargets <= 0 {
		return nil, fmt.Errorf("burrow extra targets require mastery")
	}
	if len(raw) > maxTargets {
		return nil, fmt.Errorf("too many extra targets for burrow")
	}
	result := make([]SpellTarget, 0, len(raw))
	seen := map[Position]bool{mainTarget.Position: true}
	for _, entry := range raw {
		data, ok := entry.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid burrow extra target")
		}
		pos, err := requiredBoardPosition(data, "col", "row")
		if err != nil {
			return nil, fmt.Errorf("invalid burrow extra target")
		}
		extra := SpellTarget{Type: "unit", Position: pos}
		if err := e.validateSpellExtraTarget(playerID, extra); err != nil {
			return nil, err
		}
		if seen[extra.Position] {
			return nil, fmt.Errorf("duplicate burrow extra target")
		}
		seen[extra.Position] = true
		result = append(result, extra)
	}
	return result, nil
}
