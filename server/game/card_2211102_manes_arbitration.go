package game

import (
	"eraofarcane/model"
	"fmt"
)

type Card2211102ManesArbitration struct{ AlwaysActive }

func (Card2211102ManesArbitration) ID() string { return "2211102" }

func (Card2211102ManesArbitration) Name() string { return "玛涅斯之予夺" }

func (Card2211102ManesArbitration) OnEquip(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	choices := []map[string]any{
		{"instance_id": "all_water_power", "name": "此后学习的所有水纹法术+2威", "zone": "choice", "side": "own"},
	}
	waterSkills := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && skill.Card.IsSkill() && skill.Card.Category == model.ElementWater
	})
	if len(waterSkills) > 0 {
		choices = append(choices, map[string]any{
			"instance_id": "one_water_power_attack",
			"name":        "选择1个水纹法术+3威+1攻并禁止再学习水纹",
			"zone":        "choice",
			"side":        "own",
		})
	}
	ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "manes_arbitration_choice",
		"玛涅斯之予夺:选择本局游戏持续的效果", choices, 1, 1,
		nil, false, func(selected []string, _ map[string]any) error {
			switch firstSelected(selected) {
			case "all_water_power":
				ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
					Type:             TempModSkillPowerBonus,
					SourceCardNumber: "2211102",
					SourceName:       "玛涅斯之予夺",
					Element:          model.ElementWater,
					Amount:           2,
				})
			case "one_water_power_attack":
				currentWaterSkills := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
					return skill != nil && skill.Card != nil && skill.Card.IsSkill() && skill.Card.Category == model.ElementWater
				})
				if len(currentWaterSkills) == 0 {
					return fmt.Errorf("no water skill to empower")
				}
				ctx.Engine.SetPendingActionWithError(ctx.PlayerID, "manes_arbitration_skill",
					"玛涅斯之予夺:选择1个水纹法术", currentWaterSkills, 1, 1,
					nil, false, func(skillSelected []string, _ map[string]any) error {
						skill := findFriendlySkillIncludingBound(ctx.Engine, ctx.PlayerID, firstSelected(skillSelected))
						if skill == nil || skill.Card == nil || skill.Card.Category != model.ElementWater {
							return fmt.Errorf("invalid Manes target")
						}
						skill.PowerBonus += 3
						skill.AttackBonus++
						ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
							Type:             TempModCannotLearnElementSkill,
							SourceCardNumber: "2211102",
							SourceName:       "玛涅斯之予夺",
							Element:          model.ElementWater,
						})
						return nil
					})
			default:
				return fmt.Errorf("invalid Manes choice")
			}
			return nil
		})
	return nil
}
