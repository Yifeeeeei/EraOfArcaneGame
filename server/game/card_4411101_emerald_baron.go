package game

type Card4411101EmeraldBaron struct{ AlwaysActive }

func (Card4411101EmeraldBaron) ID() string   { return "4411101" }
func (Card4411101EmeraldBaron) Name() string { return "翡翠男爵 杰德 拜利兰" }
func (Card4411101EmeraldBaron) PreventsShieldDecay(_ *CardInstance, shield int) bool {
	return shield < 3
}
