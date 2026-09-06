package game

import ()

type Card1621112WhisperElfHunter struct{ AlwaysActive }

func (Card1621112WhisperElfHunter) ID() string { return "1621112" }

func (Card1621112WhisperElfHunter) Name() string { return "谧语精灵猎手" }

func (Card1621112WhisperElfHunter) OnDeath(ctx *EffectContext) error {
	candidates := ctx.Engine.enemyUnits(ctx.PlayerID, true, nil)
	if len(candidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "whisper_elf_hunter_damage",
		"谧语精灵猎手:选择1个敌人造成1点伤害", candidates, 1, 1,
		func(selected []string) {
			target := ctx.Engine.findUnitByInstanceID(firstSelected(selected))
			if target != nil {
				ctx.DealDamage(target, 1)
			}
		})
	return nil
}
