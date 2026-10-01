package game

type Card3621007AndisPunishment struct{ AlwaysActive }

func (Card3621007AndisPunishment) ID() string { return "3621007" }

func (Card3621007AndisPunishment) Name() string { return "安迪斯的惩罚" }

func (Card3621007AndisPunishment) DamageScope() DamageScope { return DamageFriendly }

func (Card3621007AndisPunishment) OnDamaged(ctx *EffectContext, event DamageEvent) error {
	damagedPlayer := event.Target.OwnerID
	if damagedPlayer == ctx.PlayerID {
		amount := event.Amount
		ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
			Type:             TempModSkillPowerBonus,
			SourceCardNumber: ctx.Source.Card.Number,
			SourceName:       ctx.Source.Card.Name,
			TargetInstanceID: ctx.Source.InstanceID,
			Amount:           max(amount, 0),
			RemainingUses:    1,
			ExpiresTurn:      ctx.Engine.State.TurnNumber + 2,
		})
	}
	return nil
}
