package game

import "fmt"

type Card4611001Alice struct{ AlwaysActive }

func (Card4611001Alice) ID() string   { return "4611001" }
func (Card4611001Alice) Name() string { return "暗影学者 爱莉斯" }

func (Card4611001Alice) OnFriendlyDeath(ctx *EffectContext) error {
	if ctx == nil || ctx.Source == nil || ctx.Target == nil || ctx.Target.Card == nil {
		return nil
	}
	if ctx.Source.UsedThisTurn > 0 || !ctx.Target.Card.IsCompanion() || ctx.Target.Card.IsHero() {
		return nil
	}
	candidates := friendlySkillCandidates(ctx.Engine.State.Players[ctx.PlayerID])
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "alice_boost_spell",
		fmt.Sprintf("%s: 选择1个你的法术+1威", ctx.Source.Card.Name), candidates, 1, 1,
		func(selected []string) {
			if len(selected) == 0 {
				return
			}
			skill := ctx.Engine.findSkill(ctx.Engine.State.Players[ctx.PlayerID], selected[0])
			if skill == nil {
				return
			}
			skill.PowerBonus++
			ctx.Source.UsedThisTurn++
			ctx.Engine.emit(GameEvent{Type: "effect_trigger", Player: ctx.PlayerID, Data: map[string]any{
				"source": cardToInfo(ctx.Source),
				"target": cardToInfo(skill),
				"effect": "modify_spell_power",
				"amount": 1,
			}})
		})
	return nil
}

func friendlySkillCandidates(ps *PlayerState) []map[string]any {
	candidates := make([]map[string]any, 0, len(ps.Skills))
	for _, skill := range ps.Skills {
		if skill == nil || skill.Card == nil {
			continue
		}
		candidates = append(candidates, cardToInfo(skill))
	}
	return candidates
}
