package game

type Card2601002Spellbook struct{}

func (Card2601002Spellbook) ID() string   { return "2601002" }
func (Card2601002Spellbook) Name() string { return "咒言书" }
func (Card2601002Spellbook) OnEnter(ctx *EffectContext) error {
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	for _, skill := range opponent.Skills {
		if skill != nil {
			skill.Statuses[StatusWeaken]++
		}
	}
	return nil
}
