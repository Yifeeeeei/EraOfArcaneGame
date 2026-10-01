package game

type Card2321104ThunderlightCrown struct{ AlwaysActive }

func (Card2321104ThunderlightCrown) ID() string { return "2321104" }

func (Card2321104ThunderlightCrown) Name() string { return "雷光头冠" }

func (Card2321104ThunderlightCrown) IsPrayerAbility() bool { return true }

func (Card2321104ThunderlightCrown) OnPerTurn(ctx *EffectContext) error {
	ctx.Engine.addNextTaggedSpellPowerBonus(ctx.PlayerID, "聚能", 1)
	return nil
}
