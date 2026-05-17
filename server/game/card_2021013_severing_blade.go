package game

import "fmt"

type Card2021013SeveringBlade struct{ AlwaysActive }

func (Card2021013SeveringBlade) ID() string   { return "2021013" }
func (Card2021013SeveringBlade) Name() string { return "断绝之刃" }

func (Card2021013SeveringBlade) ModifySpellStats(ctx *EffectContext, stats *SpellStats) {
	if ctx.ExtraData["purpose"] != string(skillPurposeAttack) {
		return
	}
	stats.PowerBonus += 2
}

func (Card2021013SeveringBlade) ValidateSkillUse(ctx *EffectContext, skill *CardInstance, purpose skillPurpose) error {
	if purpose == skillPurposeDefend {
		return fmt.Errorf("断绝之刃使你的法术不能用于防御")
	}
	return nil
}
