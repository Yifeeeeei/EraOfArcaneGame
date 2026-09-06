package game

type Card3021010Dispel struct{ AlwaysActive }

func (Card3021010Dispel) ID() string { return "3021010" }

func (Card3021010Dispel) Name() string { return "解咒" }

func (Card3021010Dispel) CanCancelDefenseSpell(*CardInstance) bool { return true }
