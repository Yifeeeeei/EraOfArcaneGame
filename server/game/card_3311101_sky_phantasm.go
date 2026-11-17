package game

type Card3311101SkyPhantasm struct{ AlwaysActive }

func (Card3311101SkyPhantasm) ID() string               { return "3311101" }
func (Card3311101SkyPhantasm) Name() string             { return "苍穹幻韵" }
func (Card3311101SkyPhantasm) AutomaticSpellCopy() bool { return false }
func (Card3311101SkyPhantasm) SpellCopies(ctx *EffectContext) []*CardInstance {
	var result []*CardInstance
	for _, skill := range ctx.Engine.State.Players[ctx.PlayerID].Skills {
		if skill == nil || skill == ctx.Source || skill.Card == nil {
			continue
		}
		if hasCardTag(skill.Card, "驱动") || hasCardTag(skill.Card, "聚能") {
			result = append(result, skill)
		}
	}
	return result
}
