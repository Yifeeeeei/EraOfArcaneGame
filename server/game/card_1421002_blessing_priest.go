package game

type Card1421002BlessingPriest struct{}

func (Card1421002BlessingPriest) ID() string   { return "1421002" }
func (Card1421002BlessingPriest) Name() string { return "祝祷祭师" }
func (Card1421002BlessingPriest) ProtectsAdjacentFromNegativeStatus() bool {
	return true
}
