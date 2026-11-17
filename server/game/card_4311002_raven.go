package game

type Card4311002Raven struct{ AlwaysActive }

func (Card4311002Raven) ID() string { return "4311002" }

func (Card4311002Raven) Name() string { return "\"渡鸦\" 睿文" }

func (Card4311002Raven) OpeningHandBonus(*CardInstance) int { return 1 }

func (Card4311002Raven) HandLimitGrant(source *CardInstance, playerID int) HandLimitGrant {
	if source.OwnerID != playerID {
		return HandLimitGrant{}
	}
	return HandLimitGrant{Group: "raven", Delta: 1}
}
