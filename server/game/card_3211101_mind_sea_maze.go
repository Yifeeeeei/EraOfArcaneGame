package game

type Card3211101MindSeaMaze struct{ AlwaysActive }

func (Card3211101MindSeaMaze) ID() string { return "3211101" }

func (Card3211101MindSeaMaze) Name() string { return "心海迷离" }

func (Card3211101MindSeaMaze) OnSpellCast(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || !isFriendlySpellCast(ctx) {
		return nil
	}
	ctx.Source.Statuses[mindSeaMazeAnyRangeUntilStatus] = ctx.Engine.State.TurnNumber
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:             TempModSkillPowerBonus,
		SourceCardNumber: ctx.Source.Card.Number,
		SourceName:       ctx.Source.Card.Name,
		TargetInstanceID: ctx.Source.InstanceID,
		Amount:           1,
		ExpiresTurn:      ctx.Engine.State.TurnNumber + 1,
	})
	return nil
}

func (Card3211101MindSeaMaze) ModifySpellArea(ctx *EffectContext, area *SpellArea) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || ctx.Target != ctx.Source || area == nil {
		return
	}
	if ctx.Source.Statuses[mindSeaMazeAnyRangeUntilStatus] >= ctx.Engine.State.TurnNumber {
		*area = SpellAreaAll
	}
}

var _ OnSpellCastBehavior = Card3211101MindSeaMaze{}

var _ SpellAreaModifier = Card3211101MindSeaMaze{}

const mindSeaMazeAnyRangeUntilStatus = "mind_sea_maze_any_range_until"
