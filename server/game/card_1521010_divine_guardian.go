package game

type Card1521010DivineGuardian struct{ AlwaysActive }

func (Card1521010DivineGuardian) ID() string { return "1521010" }

func (Card1521010DivineGuardian) Name() string { return "神护者" }

func (Card1521010DivineGuardian) HasNegativeStatusImmunity() bool {
	return true
}
