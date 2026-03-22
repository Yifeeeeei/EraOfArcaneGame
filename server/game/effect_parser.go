package game

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

// AutoParseAndRegister scans all cards in the database and registers
// auto-parseable effects based on description patterns.
// This handles common simple patterns; complex cards need manual registration.
func AutoParseAndRegister() {
	db := getCardDB()
	if db == nil {
		return
	}

	registry := GetEffectRegistry()

	for number, card := range db {
		desc := card.Description
		if desc == "" {
			continue
		}

		// Skip cards that already have manually registered effects
		if len(registry.GetAllEffects(number)) > 0 {
			continue
		}

		parseAndRegisterCard(registry, number, desc)
	}

	log.Printf("[EffectParser] Auto-parsed card effects from database")
}

func parseAndRegisterCard(registry *EffectRegistry, number string, desc string) {
	// --- Entry effects (入场:) ---
	entryMatch := extractAfterKeyword(desc, "入场:")
	if entryMatch != "" {
		parseEntryEffect(registry, number, entryMatch)
	}

	// --- Deathrattle effects (遗言:) ---
	deathMatch := extractAfterKeyword(desc, "遗言:")
	if deathMatch != "" {
		parseDeathEffect(registry, number, deathMatch)
	}

	// --- Simple standalone keywords (no colon needed) ---
	// These are handled by ApplyKeywordOnEnter, no registry needed
}

// parseEntryEffect parses "入场:XXX" patterns
func parseEntryEffect(registry *EffectRegistry, number string, effectText string) {
	// Pattern: "充能X" - gain X charge
	if match := matchPattern(effectText, `充能(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnEnter, func(ctx *EffectContext) error {
			ctx.Engine.addCharge(ctx.PlayerID, amount)
			ctx.Engine.emit(GameEvent{
				Type:   "effect_trigger",
				Player: -1,
				Data: map[string]any{
					"source":  cardToInfo(ctx.Source),
					"effect":  "charge",
					"amount":  amount,
				},
			})
			return nil
		})
		return
	}

	// Pattern: "抽X张牌" or "抽1张牌" - draw cards
	if match := matchPattern(effectText, `抽(\d+)张牌`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnEnter, func(ctx *EffectContext) error {
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			drawn := ps.DrawCards(amount)
			for _, c := range drawn {
				ctx.Engine.emit(GameEvent{
					Type:   "draw_card",
					Player: ctx.PlayerID,
					Data:   map[string]any{"card": cardToInfo(c)},
				})
			}
			return nil
		})
		return
	}

	// Pattern: "对法力范围内1个敌方伙伴造成X点伤害" - deal damage on enter
	if match := matchPattern(effectText, `造成(\d+)点伤害`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnEnter, func(ctx *EffectContext) error {
			// Auto-target: find a target in spell range
			if ctx.Target != nil {
				ctx.Engine.dealDamage(ctx.Target, amount, ctx.OpponentID)
				ctx.Engine.emit(GameEvent{
					Type:   "effect_trigger",
					Player: -1,
					Data: map[string]any{
						"source":  cardToInfo(ctx.Source),
						"effect":  "damage",
						"amount":  amount,
						"target":  cardToInfo(ctx.Target),
					},
				})
			}
			return nil
		})
		return
	}

	// Pattern: "对X冻结Y" or "冻结X" - apply freeze on enter
	if match := matchPattern(effectText, `冻结(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registerStatusOnEnter(registry, number, StatusFreeze, amount)
		return
	}

	// Pattern: "点燃X" - apply burn on enter
	if match := matchPattern(effectText, `点燃(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registerStatusOnEnter(registry, number, StatusBurn, amount)
		return
	}

	// Pattern: "眩晕X" or "晕眩X" - apply stun on enter
	if match := matchPattern(effectText, `[眩晕|晕眩](\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registerStatusOnEnter(registry, number, StatusStun, amount)
		return
	}

	// Pattern: "石化X" - apply petrify on enter
	if match := matchPattern(effectText, `石化(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registerStatusOnEnter(registry, number, StatusPetrify, amount)
		return
	}

	// Pattern: "对你的人物造成X点伤害" - self-hero damage on enter
	if match := matchPattern(effectText, `对你的人物造成(\d+)点伤害`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnEnter, func(ctx *EffectContext) error {
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			if ps.Hero != nil {
				ctx.Engine.dealDamage(ps.Hero, amount, ctx.PlayerID)
			}
			return nil
		})
		return
	}
}

// parseDeathEffect parses "遗言:XXX" patterns
func parseDeathEffect(registry *EffectRegistry, number string, effectText string) {
	// Pattern: "对任意1个敌人造成X点伤害" - deal damage on death
	if match := matchPattern(effectText, `造成(\d+)点伤害`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnDeath, func(ctx *EffectContext) error {
			// Find a random enemy unit to hit
			opponent := ctx.Engine.State.Players[ctx.OpponentID]
			target := findAnyUnit(opponent)
			if target != nil {
				ctx.Engine.dealDamage(target, amount, ctx.OpponentID)
				ctx.Engine.emit(GameEvent{
					Type:   "effect_trigger",
					Player: -1,
					Data: map[string]any{
						"source":  cardToInfo(ctx.Source),
						"effect":  "deathrattle_damage",
						"amount":  amount,
						"target":  cardToInfo(target),
					},
				})
			}
			return nil
		})
		return
	}

	// Pattern: "抽X张牌" - draw on death
	if match := matchPattern(effectText, `抽(\d+)张牌`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnDeath, func(ctx *EffectContext) error {
			ps := ctx.Engine.State.Players[ctx.PlayerID]
			drawn := ps.DrawCards(amount)
			for _, c := range drawn {
				ctx.Engine.emit(GameEvent{
					Type:   "draw_card",
					Player: ctx.PlayerID,
					Data:   map[string]any{"card": cardToInfo(c)},
				})
			}
			return nil
		})
		return
	}

	// Pattern: "充能X" - gain charge on death
	if match := matchPattern(effectText, `充能(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnDeath, func(ctx *EffectContext) error {
			ctx.Engine.addCharge(ctx.PlayerID, amount)
			return nil
		})
		return
	}
}

// registerStatusOnEnter registers an effect that applies a status to a target on entry
func registerStatusOnEnter(registry *EffectRegistry, number string, status string, amount int) {
	registry.Register(number, TriggerOnEnter, func(ctx *EffectContext) error {
		if ctx.Target != nil {
			ctx.Target.Statuses[status] += amount
			ctx.Engine.emit(GameEvent{
				Type:   "effect_trigger",
				Player: -1,
				Data: map[string]any{
					"source": cardToInfo(ctx.Source),
					"effect": "apply_status",
					"status": status,
					"amount": amount,
					"target": cardToInfo(ctx.Target),
				},
			})
		}
		return nil
	})
}

// Helper: extract text after a keyword until period/comma or end
func extractAfterKeyword(desc, keyword string) string {
	idx := strings.Index(desc, keyword)
	if idx == -1 {
		return ""
	}
	rest := desc[idx+len(keyword):]
	// Find end: period, comma followed by non-digit, or end of string
	end := len(rest)
	for i, r := range rest {
		if r == '.' || r == '。' {
			end = i
			break
		}
	}
	return strings.TrimSpace(rest[:end])
}

// Helper: match regex pattern and return groups
func matchPattern(text, pattern string) []string {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(text)
	if len(m) == 0 {
		return nil
	}
	return m
}

// Helper: find any alive unit on a player's field
func findAnyUnit(ps *PlayerState) *CardInstance {
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil {
				return ps.Units[col][row]
			}
		}
	}
	return nil
}
