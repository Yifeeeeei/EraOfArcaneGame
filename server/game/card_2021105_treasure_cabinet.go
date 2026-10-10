package game

type Card2021105TreasureCabinet struct{ AlwaysActive }

func (Card2021105TreasureCabinet) ID() string   { return "2021105" }
func (Card2021105TreasureCabinet) Name() string { return "珍宝柜" }
func (Card2021105TreasureCabinet) SlotGrant(*CardInstance) SlotGrant {
	return SlotGrant{Group: "treasure_cabinet", EquipmentSlots: 1, DuplicateEquipmentSubtypes: true}
}
