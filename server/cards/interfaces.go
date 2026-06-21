package cards

// HeroCard marks a hero/person card definition.
type HeroCard interface {
	CardDefinition
	isHeroCard()
}

// CompanionCard marks a companion card definition.
type CompanionCard interface {
	CardDefinition
	isCompanionCard()
}

// SkillCard marks a skill card definition.
type SkillCard interface {
	CardDefinition
	isSkillCard()
}

// ItemCard marks any item card definition.
type ItemCard interface {
	CardDefinition
	isItemCard()
}

// EquipmentCard marks an item that is equipped into an equipment slot.
type EquipmentCard interface {
	ItemCard
	isEquipmentCard()
}

// WeaponCard marks an equipment item whose card tag includes weapon behavior.
type WeaponCard interface {
	EquipmentCard
	isWeaponCard()
}

// ArmorCard marks an equipment item whose card tag includes armor behavior.
type ArmorCard interface {
	EquipmentCard
	isArmorCard()
}

// AccessoryCard marks an equipment item whose card tag includes accessory behavior.
type AccessoryCard interface {
	EquipmentCard
	isAccessoryCard()
}

// ArtifactCard marks an equipment item whose card tag includes artifact behavior.
type ArtifactCard interface {
	EquipmentCard
	isArtifactCard()
}

// ConsumableCard marks an item that is used from hand and then discarded.
type ConsumableCard interface {
	ItemCard
	isConsumableCard()
}

// TerrainCard marks an item placed onto the battlefield grid.
type TerrainCard interface {
	ItemCard
	isTerrainCard()
}
