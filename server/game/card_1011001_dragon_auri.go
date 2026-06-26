package game

type Card1011001DragonAuri struct{ AlwaysActive }

func (Card1011001DragonAuri) ID() string   { return "1011001" }
func (Card1011001DragonAuri) Name() string { return "魔龙 奥瑞" }

func (Card1011001DragonAuri) OnEnter(ctx *EffectContext) error {
	bindSkillToHost(ctx, "3001001")
	return nil
}
