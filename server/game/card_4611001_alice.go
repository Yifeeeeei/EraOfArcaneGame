package game

type Card4611001Alice struct{ AlwaysActive }

func (Card4611001Alice) ID() string   { return "4611001" }
func (Card4611001Alice) Name() string { return "暗影学者 爱莉斯" }
func (Card4611001Alice) OnFriendlyDeath(*EffectContext) error {
	return nil
}
