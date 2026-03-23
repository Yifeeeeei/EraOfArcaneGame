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
	parsed := 0

	for number, card := range db {
		desc := card.Description
		if desc == "" {
			continue
		}

		// Skip cards that already have manually registered effects
		if len(registry.GetAllEffects(number)) > 0 {
			continue
		}

		if parseAndRegisterCard(registry, number, desc) {
			parsed++
		}
	}

	log.Printf("[EffectParser] Auto-parsed %d card effects from database", parsed)
}

// parseAndRegisterCard parses description and registers all matching effects.
// Returns true if any effects were registered.
func parseAndRegisterCard(registry *EffectRegistry, number string, desc string) bool {
	registered := false

	// --- Entry effects (入场:) ---
	entryMatch := extractAfterKeyword(desc, "入场:")
	if entryMatch != "" {
		if parseEntryEffects(registry, number, entryMatch) {
			registered = true
		}
	}

	// --- Deathrattle effects (遗言:) ---
	deathMatch := extractAfterKeyword(desc, "遗言:")
	if deathMatch != "" {
		if parseDeathEffects(registry, number, deathMatch) {
			registered = true
		}
	}

	// --- 祈咒 effects (回合开始触发) ---
	prayerMatch := extractAfterKeyword(desc, "祈咒:")
	if prayerMatch != "" {
		if parsePrayerEffects(registry, number, prayerMatch) {
			registered = true
		}
	}

	// --- Simple standalone keywords (no colon needed) ---
	// These are handled by ApplyKeywordOnEnter, no registry needed

	return registered
}

// parseEntryEffects parses "入场:XXX" patterns. Multiple effects can be registered.
func parseEntryEffects(registry *EffectRegistry, number string, effectText string) bool {
	registered := false

	// Pattern: "充能X" - gain X charge
	if match := matchPattern(effectText, `充能(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnEnter, func(ctx *EffectContext) error {
			ctx.Engine.addCharge(ctx.PlayerID, amount)
			ctx.Engine.emit(GameEvent{
				Type:   "effect_trigger",
				Player: -1,
				Data: map[string]any{
					"source": cardToInfo(ctx.Source),
					"effect": "charge",
					"amount": amount,
				},
			})
			return nil
		})
		registered = true
	}

	// Pattern: "抽X张牌" - draw cards
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
		registered = true
	}

	// Pattern: "造成X点伤害" - deal damage on enter (to a target or auto-target)
	if !registered { // Only if no other entry effect matched — damage is often part of complex effects
		if match := matchPattern(effectText, `造成(\d+)点伤害`); match != nil {
			amount, _ := strconv.Atoi(match[1])
			registry.Register(number, TriggerOnEnter, func(ctx *EffectContext) error {
				if ctx.Target != nil {
					ctx.Engine.dealDamage(ctx.Target, amount, ctx.OpponentID)
					ctx.Engine.emit(GameEvent{
						Type:   "effect_trigger",
						Player: -1,
						Data: map[string]any{
							"source": cardToInfo(ctx.Source),
							"effect": "damage",
							"amount": amount,
							"target": cardToInfo(ctx.Target),
						},
					})
				} else {
					// Auto-target: find front row enemy
					opponent := ctx.Engine.State.Players[ctx.OpponentID]
					target := findFrontRowUnit(opponent)
					if target != nil {
						ctx.Engine.dealDamage(target, amount, ctx.OpponentID)
						ctx.Engine.emit(GameEvent{
							Type:   "effect_trigger",
							Player: -1,
							Data: map[string]any{
								"source": cardToInfo(ctx.Source),
								"effect": "damage",
								"amount": amount,
								"target": cardToInfo(target),
							},
						})
					}
				}
				return nil
			})
			registered = true
		}
	}

	// Pattern: "冻结X" - apply freeze on enter
	if match := matchPattern(effectText, `冻结(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registerStatusOnEnter(registry, number, StatusFreeze, amount)
		registered = true
	}

	// Pattern: "点燃X" - apply burn on enter
	if match := matchPattern(effectText, `点燃(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registerStatusOnEnter(registry, number, StatusBurn, amount)
		registered = true
	}

	// Pattern: "眩晕X" - apply stun on enter
	if match := matchPattern(effectText, `(?:眩晕|晕眩)(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registerStatusOnEnter(registry, number, StatusStun, amount)
		registered = true
	}

	// Pattern: "石化X" - apply petrify on enter
	if match := matchPattern(effectText, `石化(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registerStatusOnEnter(registry, number, StatusPetrify, amount)
		registered = true
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
		registered = true
	}

	// Pattern: "护盾X" - gain shield on enter (for self)
	if match := matchPattern(effectText, `护盾(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnEnter, func(ctx *EffectContext) error {
			ctx.Source.Statuses["护盾"] += amount
			ctx.Engine.emit(GameEvent{
				Type:   "effect_trigger",
				Player: -1,
				Data: map[string]any{
					"source": cardToInfo(ctx.Source),
					"effect": "gain_shield",
					"amount": amount,
				},
			})
			return nil
		})
		registered = true
	}

	// Pattern: "隐蔽" (without number, treated as 1 layer) or "隐蔽X"
	if strings.Contains(effectText, "隐蔽") {
		amount := 1
		if match := matchPattern(effectText, `隐蔽(\d+)`); match != nil {
			amount, _ = strconv.Atoi(match[1])
		}
		registry.Register(number, TriggerOnEnter, func(ctx *EffectContext) error {
			ctx.Source.Statuses["隐蔽"] += amount
			ctx.Engine.emit(GameEvent{
				Type:   "effect_trigger",
				Player: -1,
				Data: map[string]any{
					"source": cardToInfo(ctx.Source),
					"effect": "gain_stealth",
					"amount": amount,
				},
			})
			return nil
		})
		registered = true
	}

	// Pattern: "虚弱X" - apply weaken on enter
	if match := matchPattern(effectText, `虚弱(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registerStatusOnEnter(registry, number, StatusWeaken, amount)
		registered = true
	}

	return registered
}

// parseDeathEffects parses "遗言:XXX" patterns. Multiple effects allowed.
func parseDeathEffects(registry *EffectRegistry, number string, effectText string) bool {
	registered := false

	// Pattern: "造成X点伤害" - deal damage on death
	if match := matchPattern(effectText, `造成(\d+)点伤害`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnDeath, func(ctx *EffectContext) error {
			opponent := ctx.Engine.State.Players[ctx.OpponentID]
			target := findAnyUnit(opponent)
			if target != nil {
				ctx.Engine.dealDamage(target, amount, ctx.OpponentID)
				ctx.Engine.emit(GameEvent{
					Type:   "effect_trigger",
					Player: -1,
					Data: map[string]any{
						"source": cardToInfo(ctx.Source),
						"effect": "deathrattle_damage",
						"amount": amount,
						"target": cardToInfo(target),
					},
				})
			}
			return nil
		})
		registered = true
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
		registered = true
	}

	// Pattern: "充能X" - gain charge on death
	if match := matchPattern(effectText, `充能(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnDeath, func(ctx *EffectContext) error {
			ctx.Engine.addCharge(ctx.PlayerID, amount)
			return nil
		})
		registered = true
	}

	// Pattern: "冻结X" - apply freeze on death to enemy
	if match := matchPattern(effectText, `冻结(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnDeath, func(ctx *EffectContext) error {
			opponent := ctx.Engine.State.Players[ctx.OpponentID]
			target := findFrontRowUnit(opponent)
			if target != nil {
				target.Statuses[StatusFreeze] += amount
				ctx.Engine.emit(GameEvent{
					Type:   "effect_trigger",
					Player: -1,
					Data: map[string]any{
						"source": cardToInfo(ctx.Source),
						"effect": "deathrattle_status",
						"status": StatusFreeze,
						"amount": amount,
						"target": cardToInfo(target),
					},
				})
			}
			return nil
		})
		registered = true
	}

	// Pattern: "点燃X" - apply burn on death
	if match := matchPattern(effectText, `点燃(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnDeath, func(ctx *EffectContext) error {
			opponent := ctx.Engine.State.Players[ctx.OpponentID]
			target := findFrontRowUnit(opponent)
			if target != nil {
				target.Statuses[StatusBurn] += amount
				ctx.Engine.emit(GameEvent{
					Type:   "effect_trigger",
					Player: -1,
					Data: map[string]any{
						"source": cardToInfo(ctx.Source),
						"effect": "deathrattle_status",
						"status": StatusBurn,
						"amount": amount,
						"target": cardToInfo(target),
					},
				})
			}
			return nil
		})
		registered = true
	}

	return registered
}

// parsePrayerEffects parses "祈咒:XXX" patterns (triggers at turn start)
func parsePrayerEffects(registry *EffectRegistry, number string, effectText string) bool {
	registered := false

	// Pattern: "抽X张牌" - draw on each turn start
	if match := matchPattern(effectText, `抽(\d+)张牌`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnTurnStart, func(ctx *EffectContext) error {
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
		registered = true
	}

	// Pattern: "充能X" - gain charge on turn start
	if match := matchPattern(effectText, `充能(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnTurnStart, func(ctx *EffectContext) error {
			ctx.Engine.addCharge(ctx.PlayerID, amount)
			return nil
		})
		registered = true
	}

	// Pattern: "造成X点伤害" - deal damage on turn start
	if match := matchPattern(effectText, `造成(\d+)点伤害`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		isSelf := strings.Contains(effectText, "对你") || strings.Contains(effectText, "自身")
		registry.Register(number, TriggerOnTurnStart, func(ctx *EffectContext) error {
			if isSelf {
				ctx.Engine.dealDamage(ctx.Source, amount, ctx.PlayerID)
			} else {
				opponent := ctx.Engine.State.Players[ctx.OpponentID]
				target := findFrontRowUnit(opponent)
				if target != nil {
					ctx.Engine.dealDamage(target, amount, ctx.OpponentID)
				}
			}
			return nil
		})
		registered = true
	}

	// Pattern: "点燃X" - apply burn on turn start
	if match := matchPattern(effectText, `点燃(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnTurnStart, func(ctx *EffectContext) error {
			opponent := ctx.Engine.State.Players[ctx.OpponentID]
			target := findFrontRowUnit(opponent)
			if target != nil {
				target.Statuses[StatusBurn] += amount
			}
			return nil
		})
		registered = true
	}

	// Pattern: "冻结X" - apply freeze on turn start
	if match := matchPattern(effectText, `冻结(\d+)`); match != nil {
		amount, _ := strconv.Atoi(match[1])
		registry.Register(number, TriggerOnTurnStart, func(ctx *EffectContext) error {
			opponent := ctx.Engine.State.Players[ctx.OpponentID]
			target := findFrontRowUnit(opponent)
			if target != nil {
				target.Statuses[StatusFreeze] += amount
			}
			return nil
		})
		registered = true
	}

	return registered
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
		} else {
			// Auto-target: front row enemy
			opponent := ctx.Engine.State.Players[ctx.OpponentID]
			target := findFrontRowUnit(opponent)
			if target != nil {
				target.Statuses[status] += amount
				ctx.Engine.emit(GameEvent{
					Type:   "effect_trigger",
					Player: -1,
					Data: map[string]any{
						"source": cardToInfo(ctx.Source),
						"effect": "apply_status",
						"status": status,
						"amount": amount,
						"target": cardToInfo(target),
					},
				})
			}
		}
		return nil
	})
}

// Helper: extract text after a keyword until period/newline/semicolon or end
func extractAfterKeyword(desc, keyword string) string {
	idx := strings.Index(desc, keyword)
	if idx == -1 {
		return ""
	}
	rest := desc[idx+len(keyword):]
	// Find end: period, fullstop, semicolon, or next keyword section
	end := len(rest)
	for i, r := range rest {
		if r == '.' || r == '。' || r == ';' || r == '；' {
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

// Helper: find a unit in the front row of a player's field
func findFrontRowUnit(ps *PlayerState) *CardInstance {
	frontRow := ps.GetFrontRow()
	if frontRow < 0 {
		return nil
	}
	for col := 0; col < 3; col++ {
		if ps.Units[col][frontRow] != nil {
			return ps.Units[col][frontRow]
		}
	}
	return nil
}
