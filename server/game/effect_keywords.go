package game

import (
	"strings"
)

// Keyword flags - detected from card description text
const (
	KW_Rush       = "速攻"   // enters vertical (can act immediately)
	KW_Pierce     = "穿透"   // ignores spell range, can target any unit
	KW_Temporary  = "临时"   // removed at end of turn
	KW_Stealth    = "隐蔽"   // can't be targeted by enemy, doesn't block spell range
	KW_Taunt      = "引魔"   // always in enemy spell range
	KW_Ritual     = "仪式"   // special ritual mechanic
	KW_Terrain    = "地形"   // terrain item mechanic
	KW_Artifact   = "法宝"   // artifact equipment mechanic
	KW_Sacrifice  = "献祭"   // sacrifice mechanic
	KW_Defense    = "防御"   // defense-only spell
	KW_Counter    = "反制"   // counter item
	KW_Aura       = "异能"   // aura/continuous effect
	KW_Devour     = "吞噬"   // must sacrifice a friendly unit on enter
	KW_Shield     = "护盾"   // shield mechanic
	KW_Charge     = "充能"   // charge counter mechanic
	KW_Overload   = "过载"   // overload trigger
	KW_Bind       = "绑定"   // bind/attach mechanic
	KW_Shielding  = "屏蔽"   // blocks enemy pierce
	KW_Cannon     = "火炮"   // cannon mechanic
	KW_ElemLock   = "元素限制" // element restriction
)

// HasKeyword checks if a card's description contains a keyword
func HasKeyword(description string, keyword string) bool {
	return strings.Contains(description, keyword)
}

// HasAnyKeyword checks if a card has any of the given keywords
func HasAnyKeyword(description string, keywords ...string) bool {
	for _, kw := range keywords {
		if strings.Contains(description, kw) {
			return true
		}
	}
	return false
}

// ParseCooldown extracts cooldown value from description (e.g., "冷却2" → 2)
// Only matches "冷却" at the start of the description or after a period/comma
func ParseCooldown(description string) int {
	// Look for "冷却" that is NOT preceded by characters that negate it
	idx := -1
	searchFrom := 0
	for {
		i := strings.Index(description[searchFrom:], "冷却")
		if i == -1 {
			break
		}
		pos := searchFrom + i
		// Check if it's a standalone keyword (at start, or after '.', ',', '。', ' ')
		if pos == 0 {
			idx = pos
			break
		}
		prevRune := []rune(description[:pos])
		if len(prevRune) > 0 {
			last := prevRune[len(prevRune)-1]
			if last == '.' || last == '。' || last == ',' || last == '，' || last == ' ' || last == '\n' {
				idx = pos
				break
			}
		}
		searchFrom = pos + len("冷却")
	}

	if idx == -1 {
		return 0
	}
	// Parse the number after "冷却"
	rest := description[idx+len("冷却"):]
	if len(rest) == 0 {
		return 1
	}
	n := 0
	for _, r := range rest {
		if r >= '0' && r <= '9' {
			n = n*10 + int(r-'0')
		} else {
			break
		}
	}
	if n == 0 {
		return 1
	}
	return n
}

// ApplyKeywordOnEnter applies keyword effects when a card enters the field
func (e *Engine) ApplyKeywordOnEnter(card *CardInstance) {
	desc := card.Card.Description

	// 速攻 (Rush): enters vertical instead of horizontal
	if HasKeyword(desc, KW_Rush) {
		card.IsHorizontal = false
	}

	// 隐蔽 (Stealth): add stealth status
	// Check for 隐蔽X pattern
	if HasKeyword(desc, KW_Stealth) {
		stealthLayers := parseNumberAfter(desc, KW_Stealth)
		if stealthLayers == 0 {
			stealthLayers = 1
		}
		card.Statuses["隐蔽"] = stealthLayers
	}

	// 引魔 (Taunt): mark as always targetable
	if HasKeyword(desc, KW_Taunt) {
		card.Statuses["引魔"] = 1
	}

	// 护盾 (Shield): parse shield value
	if HasKeyword(desc, KW_Shield) {
		shieldVal := parseNumberAfter(desc, KW_Shield)
		if shieldVal > 0 {
			card.Statuses["护盾"] = shieldVal
		}
	}
}

// ApplyKeywordOnSkillUse applies keyword effects when a skill is used
func (e *Engine) ApplyKeywordOnSkillUse(skill *CardInstance) {
	desc := skill.Card.Description

	// 冷却 (Cooldown): add cooldown stacks after use
	cd := ParseCooldown(desc)
	if cd > 0 {
		skill.Statuses[StatusCooldown] = cd
	}
}

// IsInSpellRange checks if a target position is in the caster's spell range
// Default spell range: all friendly units + enemy front row
// 穿透 (Pierce): can target any unit
// 隐蔽 (Stealth): can't be targeted unless caster has 穿透
// 引魔 (Taunt): always in range
// 屏蔽 (Shielding): blocks enemy pierce
func (e *Engine) IsInSpellRange(casterID int, targetCol, targetRow int, hasPierce bool) bool {
	opponent := e.State.Players[1-casterID]

	// Check if opponent has shielding active (屏蔽)
	if hasPierce {
		opponentCards := e.getAllFieldCards(opponent)
		for _, c := range opponentCards {
			if HasKeyword(c.Card.Description, KW_Shielding) && c.Statuses[StatusPetrify] == 0 {
				hasPierce = false // shielding blocks pierce
				break
			}
		}
	}

	if hasPierce {
		return true // pierce can target anything
	}

	// Default: enemy front row only
	frontRow := opponent.GetFrontRow()
	if frontRow == -1 {
		return true // no enemy units, any position is valid
	}

	target := opponent.Units[targetCol][targetRow]
	if target == nil {
		return false
	}

	// 引魔 (Taunt) units are always in range
	if target.Statuses["引魔"] > 0 {
		return true
	}

	// 隐蔽 (Stealth) units can't be targeted
	if target.Statuses["隐蔽"] > 0 {
		return false
	}

	return targetRow == frontRow
}

// IsInAttackRange checks if a unit can attack a target
// Only front row units can attack, and only enemy front row
func (e *Engine) IsInAttackRange(attackerID int, attacker *CardInstance, targetCol, targetRow int) bool {
	ps := e.State.Players[attackerID]
	opponent := e.State.Players[1-attackerID]

	// Attacker must be in front row
	frontRow := ps.GetFrontRow()
	if attacker.Position == nil || attacker.Position.Row != frontRow {
		return false
	}

	// Target must be in opponent's front row
	enemyFront := opponent.GetFrontRow()
	if enemyFront == -1 {
		return false
	}

	target := opponent.Units[targetCol][targetRow]
	if target == nil {
		return false
	}

	// 隐蔽 units can't be attacked unless attacker has 穿透
	if target.Statuses["隐蔽"] > 0 {
		if !HasKeyword(attacker.Card.Description, KW_Pierce) {
			return false
		}
	}

	return targetRow == enemyFront
}

// HandleTemporaryUnits removes 临时 (temporary) units at end of turn
func (e *Engine) HandleTemporaryUnits(ps *PlayerState) {
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			unit := ps.Units[col][row]
			if unit != nil && HasKeyword(unit.Card.Description, KW_Temporary) {
				e.destroyUnit(unit, ps.PlayerID)
			}
		}
	}
}

// HandleShieldDecay reduces shield by 1 at end of turn
func (e *Engine) HandleShieldDecay(ps *PlayerState) {
	allCards := e.getAllFieldCards(ps)
	for _, card := range allCards {
		if card.Statuses["护盾"] > 0 {
			card.Statuses["护盾"]--
		}
		// Stealth decays too
		if card.Statuses["隐蔽"] > 0 {
			card.Statuses["隐蔽"]--
		}
	}
}

// ApplyShieldDamage reduces damage by shield amount
// Returns remaining damage after shield
func ApplyShieldDamage(target *CardInstance, damage int) int {
	shield := target.Statuses["护盾"]
	if shield <= 0 {
		return damage
	}
	if damage <= shield {
		target.Statuses["护盾"] -= 1 // shield absorbs hit, lose 1 layer
		return 0
	}
	target.Statuses["护盾"] = 0
	return damage - shield
}

// parseNumberAfter extracts a number immediately following a keyword in text
func parseNumberAfter(text, keyword string) int {
	idx := strings.Index(text, keyword)
	if idx == -1 {
		return 0
	}
	rest := text[idx+len(keyword):]
	n := 0
	for _, r := range rest {
		if r >= '0' && r <= '9' {
			n = n*10 + int(r-'0')
		} else {
			break
		}
	}
	return n
}
