package game

type Card1311103SkyCouncilSpeaker struct{ AlwaysActive }

func (Card1311103SkyCouncilSpeaker) ID() string   { return "1311103" }
func (Card1311103SkyCouncilSpeaker) Name() string { return "九霄议庭言主 麦阿提" }
func (Card1311103SkyCouncilSpeaker) EnforcesHandLimit(source *CardInstance, affected int) bool {
	return source.OwnerID != affected
}

func (Card1311103SkyCouncilSpeaker) HandLimitGrant(source *CardInstance, playerID int) HandLimitGrant {
	if source.OwnerID == playerID {
		return HandLimitGrant{}
	}
	return HandLimitGrant{Group: "sky_council_speaker", Delta: -1}
}
