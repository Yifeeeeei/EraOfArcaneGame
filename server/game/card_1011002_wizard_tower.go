package game

type Card1011002WizardTower struct{ AlwaysActive }

func (Card1011002WizardTower) ID() string { return "1011002" }

func (Card1011002WizardTower) Name() string { return "巫师之塔 通天阁" }

func (Card1011002WizardTower) HasGlobalSpellRange() bool {
	return true
}

func (Card1011002WizardTower) OnEnter(ctx *EffectContext) error {
	ctx.Source.Statuses["全场法力范围"] = 1
	return nil
}
