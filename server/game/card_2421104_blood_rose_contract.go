package game

import (
	"eraofarcane/model"
)

type Card2421104BloodRoseContract struct{ AlwaysActive }

func (Card2421104BloodRoseContract) ID() string { return "2421104" }

func (Card2421104BloodRoseContract) Name() string { return "血蔷薇契约" }

func (Card2421104BloodRoseContract) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	spellCandidates := ctx.Engine.friendlySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card)
	})
	if len(spellCandidates) == 0 {
		return nil
	}
	hostCandidates := ctx.Engine.friendlyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion() &&
			(card.Card.Category == model.ElementEarth || card.Card.Category == model.ElementShadow)
	})
	if len(hostCandidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_rose_contract_spell",
		"血蔷薇契约:选择要绑定的己方法术", spellCandidates, 1, 1,
		func(selected []string) {
			skillID := firstSelected(selected)
			skill, skillIndex := findSkillSlotByInstance(ctx.Engine.State.Players[ctx.PlayerID], skillID)
			if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
				return
			}
			ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_rose_contract_host",
				"血蔷薇契约:选择地脉或暗影伙伴作为绑定宿主", hostCandidates, 1, 1,
				func(hostSelected []string) {
					host, zone := ctx.Engine.findFriendlyCandidate(ctx.PlayerID, firstSelected(hostSelected))
					if host == nil || zone != "unit" || host.Card == nil || !host.Card.IsCompanion() ||
						(host.Card.Category != model.ElementEarth && host.Card.Category != model.ElementShadow) {
						return
					}
					currentSkill, currentIndex := findSkillSlotByInstance(ctx.Engine.State.Players[ctx.PlayerID], skillID)
					if currentSkill != skill || currentIndex != skillIndex {
						return
					}
					bonus := ctx.Engine.totalLoad(host)
					ctx.Engine.State.Players[ctx.PlayerID].Skills[skillIndex] = nil
					skill.SlotIndex = -1
					markTransferredBoundSkill(skill)
					skill.PowerBonus += bonus
					host.BoundSkills = append(host.BoundSkills, skill)
					ctx.Engine.emit(GameEvent{
						Type:   "blood_rose_contract_bind",
						Player: -1,
						Data: map[string]any{
							"player":      ctx.PlayerID,
							"skill":       cardToInfo(skill),
							"host":        cardToInfo(host),
							"power_bonus": bonus,
						},
					})
				})
		})
	return nil
}

var _ OnUseItemBehavior = Card2421104BloodRoseContract{}
