package game

type Card3211101MindSeaMaze struct{ AlwaysActive }

func (Card3211101MindSeaMaze) ID() string { return "3211101" }

func (Card3211101MindSeaMaze) Name() string { return "心海迷离" }

func (Card3211101MindSeaMaze) OnSkillUseFinished(ctx *EffectContext) {
	if ctx == nil || ctx.Engine == nil || ctx.Source == nil || !ctx.Engine.cardStillOnField(ctx.Source) {
		return
	}
	ctx.Engine.addTemporaryModifier(ctx.PlayerID, TemporaryModifier{
		Type:              TempModSkillPowerBonus,
		SourceCardNumber:  ctx.Source.Card.Number,
		SourceName:        ctx.Source.Card.Name,
		TargetInstanceID:  ctx.Source.InstanceID,
		Amount:            1,
		ExpiresAtTurnEnd:  true,
		ExpiresOnPlayerID: ctx.Engine.State.CurrentTurn,
	})
	areas := []SpellArea{SpellAreaFrontRow, SpellAreaColumn, SpellAreaSquare, SpellAreaSplashCross}
	labels := []string{"前排", "纵列", "方阵", "溅射"}
	candidates := []map[string]any{}
	for i, area := range areas {
		candidates = append(candidates, map[string]any{"instance_id": string(area), "name": labels[i], "zone": "choice"})
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "mind_sea_maze_next_area", "心海迷离:选择直到回合结束的新范围", candidates, 1, 1, func(selected []string) {
		if !ctx.Engine.cardStillOnField(ctx.Source) {
			return
		}
		for _, area := range areas {
			if string(area) == firstSelected(selected) {
				ctx.Source.mindSeaRange = area
				ctx.Source.mindSeaRangeTurn = ctx.Engine.turnOrdinal()
				return
			}
		}
	})
}

func (Card3211101MindSeaMaze) ModifySpellArea(ctx *EffectContext, area *SpellArea) {
	if ctx.Target == ctx.Source && ctx.Source.mindSeaRange != "" && ctx.Source.mindSeaRangeTurn == ctx.Engine.turnOrdinal() {
		*area = ctx.Source.mindSeaRange
	}
}

var _ SkillUseFinishedBehavior = Card3211101MindSeaMaze{}
