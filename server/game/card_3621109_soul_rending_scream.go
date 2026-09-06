package game

type Card3621109SoulRendingScream struct{ AlwaysActive }

func (Card3621109SoulRendingScream) ID() string { return "3621109" }

func (Card3621109SoulRendingScream) Name() string { return "裂魂尖啸" }

func (Card3621109SoulRendingScream) OnDefend(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.ExtraData == nil {
		return nil
	}
	attackSkill, _ := ctx.ExtraData["attack_skill"].(*CardInstance)
	if attackSkill != ctx.Source {
		return nil
	}
	defenseSkills, _ := ctx.ExtraData["defense_skills"].([]*CardInstance)
	defenseBoosts, _ := ctx.ExtraData["defense_boosts"].([]*CardInstance)
	weakened := 0
	for _, skill := range append(defenseSkills, defenseBoosts...) {
		if skill == nil || skill.OwnerID == ctx.PlayerID {
			continue
		}
		if ctx.Engine.addStatus(skill, StatusWeaken, 1) {
			weakened++
		}
	}
	if weakened > 0 {
		ctx.Engine.emit(GameEvent{
			Type:   "soul_rending_scream_weaken",
			Player: -1,
			Data: map[string]any{
				"player": ctx.PlayerID,
				"source": cardToInfo(ctx.Source),
				"count":  weakened,
			},
		})
	}
	return nil
}

var _ OnDefendBehavior = Card3621109SoulRendingScream{}
