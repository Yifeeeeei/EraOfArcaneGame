package game

// A learned skill normally runs defense effects only when used as a main
// spell. An observer explicitly opts into battles involving other spells.
type SpellBattleObserverBehavior interface {
	ObservesSpellBattles() bool
}

func opposingBattleSpells(ctx *EffectContext) []*CardInstance {
	attack, _ := ctx.ExtraData["attack_skill"].(*CardInstance)
	if attack == nil {
		return nil
	}
	if attack.OwnerID == ctx.PlayerID {
		return append(append([]*CardInstance{}, spellInstancesFromData(ctx.ExtraData, "defense_skills")...), spellInstancesFromData(ctx.ExtraData, "defense_boosts")...)
	}
	return append([]*CardInstance{attack}, spellInstancesFromData(ctx.ExtraData, "boost_skills")...)
}
