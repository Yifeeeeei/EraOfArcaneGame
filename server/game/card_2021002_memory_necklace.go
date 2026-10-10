package game

type Card2021002MemoryNecklace struct{ AlwaysActive }

func (Card2021002MemoryNecklace) ID() string { return "2021002" }

func (Card2021002MemoryNecklace) Name() string { return "记忆项链" }

func (Card2021002MemoryNecklace) OnEquip(ctx *EffectContext) error {
	ctx.Source.Statuses["技能槽位+1"] = 1
	return nil
}

func (Card2021002MemoryNecklace) SlotGrant(*CardInstance) SlotGrant {
	return SlotGrant{Group: "memory_necklace", SkillSlots: 1}
}
