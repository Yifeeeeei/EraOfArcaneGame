package game

type Card2021017TravelPack struct{ AlwaysActive }

func (Card2021017TravelPack) ID() string { return "2021017" }

func (Card2021017TravelPack) Name() string { return "旅行行囊" }

func (Card2021017TravelPack) OnEquip(ctx *EffectContext) error {
	ctx.Source.Statuses["道具槽位+3"] = 1
	return nil
}

func (Card2021017TravelPack) SlotGrant(*CardInstance) SlotGrant {
	return SlotGrant{Group: "travel_pack", EquipmentSlots: 3}
}
