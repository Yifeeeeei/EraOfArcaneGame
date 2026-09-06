package game

type Card1221103WinterfellArcher struct{ AlwaysActive }

func (Card1221103WinterfellArcher) ID() string { return "1221103" }

func (Card1221103WinterfellArcher) Name() string { return "凛冬城射手" }

func (Card1221103WinterfellArcher) CanAttackFromNonFront() bool { return true }
