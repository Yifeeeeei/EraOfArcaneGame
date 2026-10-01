package game

type Card1411002KnowledgeTreeDeepRoot struct{ AlwaysActive }

func (Card1411002KnowledgeTreeDeepRoot) ID() string { return "1411002" }

func (Card1411002KnowledgeTreeDeepRoot) Name() string { return "\"知识古树\" 深耕" }

func (Card1411002KnowledgeTreeDeepRoot) OnEnter(ctx *EffectContext) error {
	ctx.Engine.advanceAllMasteryToMax(ctx.Engine.State.Players[ctx.PlayerID])
	return nil
}
