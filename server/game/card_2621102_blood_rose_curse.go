package game

type Card2621102BloodRoseCurse struct{ AlwaysActive }

func (Card2621102BloodRoseCurse) ID() string { return "2621102" }

func (Card2621102BloodRoseCurse) Name() string { return "血蔷薇诅咒" }

func (Card2621102BloodRoseCurse) OnUseItem(ctx *EffectContext) error {
	if ctx == nil || ctx.Engine == nil {
		return nil
	}
	opponentID := ctx.OpponentID
	spellCandidates := ctx.Engine.enemySkills(ctx.PlayerID, func(skill *CardInstance) bool {
		return skill != nil && skill.Card != nil && isSpellLikeCard(skill.Card)
	})
	if len(spellCandidates) == 0 {
		return nil
	}
	hostCandidates := ctx.Engine.enemyUnits(ctx.PlayerID, false, func(card *CardInstance) bool {
		return card != nil && card.Card != nil && card.Card.IsCompanion()
	})
	if len(hostCandidates) == 0 {
		return nil
	}
	ctx.Engine.SetPendingAction(ctx.PlayerID, "blood_rose_curse_spell",
		"血蔷薇诅咒:选择敌方1个法术", spellCandidates, 1, 1,
		func(selected []string) {
			skillID := firstSelected(selected)
			skill, skillIndex := findSkillSlotByInstance(ctx.Engine.State.Players[opponentID], skillID)
			if skill == nil || skill.Card == nil || !isSpellLikeCard(skill.Card) {
				return
			}
			ctx.Engine.SetPendingAction(opponentID, "blood_rose_curse_host",
				"血蔷薇诅咒:选择你的1个伙伴作为绑定宿主", hostCandidates, 1, 1,
				func(hostSelected []string) {
					host, zone := ctx.Engine.findFriendlyCandidate(opponentID, firstSelected(hostSelected))
					if host == nil || zone != "unit" || host.Card == nil || !host.Card.IsCompanion() {
						return
					}
					currentSkill, currentIndex := findSkillSlotByInstance(ctx.Engine.State.Players[opponentID], skillID)
					if currentSkill != skill || currentIndex != skillIndex {
						return
					}
					ctx.Engine.State.Players[opponentID].Skills[skillIndex] = nil
					skill.SlotIndex = -1
					markTransferredBoundSkill(skill)
					host.BoundSkills = append(host.BoundSkills, skill)
					ctx.Engine.emit(GameEvent{
						Type:   "blood_rose_curse_bind",
						Player: -1,
						Data: map[string]any{
							"player": ctx.PlayerID,
							"skill":  cardToInfo(skill),
							"host":   cardToInfo(host),
						},
					})
				})
		})
	return nil
}

var _ OnUseItemBehavior = Card2621102BloodRoseCurse{}
