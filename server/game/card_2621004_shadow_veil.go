package game

type Card2621004ShadowVeil struct{ AlwaysActive }

func (Card2621004ShadowVeil) ID() string { return "2621004" }

func (Card2621004ShadowVeil) Name() string { return "暗影帷幕" }

func (Card2621004ShadowVeil) OnSpellHit(ctx *EffectContext) error {
	if isEnemySpellCast(ctx) {
		ctx.Engine.State.Players[ctx.PlayerID].Hero.Statuses["引魔"] = 1
	}
	return nil
}
