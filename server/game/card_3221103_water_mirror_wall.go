package game

type Card3221103WaterMirrorWall struct{ AlwaysActive }

func (Card3221103WaterMirrorWall) ID() string { return "3221103" }

func (Card3221103WaterMirrorWall) Name() string { return "水镜壁" }

func (Card3221103WaterMirrorWall) OnDefend(ctx *EffectContext) error {
	if ctx.ExtraData == nil {
		return nil
	}
	success, _ := ctx.ExtraData["defense_success"].(bool)
	if !success {
		return nil
	}
	ctx.Engine.gainPlayerShield(ctx.PlayerID, 1)
	return nil
}
