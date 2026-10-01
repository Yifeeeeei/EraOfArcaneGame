package game

import (
	"sort"
)

// pendingDeath records a unit that should die after the current resolution scope.
type pendingDeath struct {
	unit      *CardInstance
	ownerID   int
	deathData map[string]any
}

type pendingDeathTrigger struct {
	unit      *CardInstance
	ownerID   int
	deathData map[string]any
}

const (
	DeathCauseSacrifice = "sacrifice"
	DeathCauseDevour    = "devour"
)

func (e *Engine) queueDeath(unit *CardInstance, ownerID int) {
	e.queueDeathWithData(unit, ownerID, nil)
}

func (e *Engine) queueDeathWithData(unit *CardInstance, ownerID int, deathData map[string]any) {
	if unit == nil {
		return
	}
	for i := range e.deathQueue {
		if e.deathQueue[i].unit == unit {
			if len(deathData) > 0 && len(e.deathQueue[i].deathData) == 0 {
				e.deathQueue[i].deathData = cloneExtraData(deathData)
			}
			return
		}
	}
	e.deathQueue = append(e.deathQueue, pendingDeath{unit: unit, ownerID: ownerID, deathData: cloneExtraData(deathData)})
}

func (e *Engine) removeQueuedDeath(unit *CardInstance) {
	if unit == nil || len(e.deathQueue) == 0 {
		return
	}
	kept := e.deathQueue[:0]
	for _, pending := range e.deathQueue {
		if pending.unit != unit {
			kept = append(kept, pending)
		}
	}
	e.deathQueue = kept
}

func (e *Engine) resolvePendingDeaths() {
	if e.resolvingDeaths {
		return
	}
	e.resolvingDeaths = true

	for len(e.deathQueue) > 0 {
		pending := e.deathQueue[0]
		e.deathQueue = e.deathQueue[1:]
		if pending.unit == nil || pending.unit.CurrentLife > 0 || !e.unitInOwnerGrid(pending.unit, pending.ownerID) {
			continue
		}
		e.destroyUnitWithData(pending.unit, pending.ownerID, pending.deathData)
	}

	e.resolvingDeaths = false
	e.resolveQueuedDeathTriggers()
	e.checkWinCondition()
}

func (e *Engine) queueDeathTriggers(unit *CardInstance, ownerID int) {
	e.queueDeathTriggersWithData(unit, ownerID, nil)
}

func (e *Engine) queueDeathTriggersWithData(unit *CardInstance, ownerID int, deathData map[string]any) {
	if unit == nil {
		return
	}
	e.deathTriggers = append(e.deathTriggers, pendingDeathTrigger{unit: unit, ownerID: ownerID, deathData: cloneExtraData(deathData)})
}

func (e *Engine) resolveQueuedDeathTriggers() {
	if len(e.deathTriggers) == 0 {
		return
	}
	e.beginResolution()
	defer e.endResolution()

	sort.SliceStable(e.deathTriggers, func(i, j int) bool {
		leftCurrent := e.deathTriggers[i].ownerID == e.State.CurrentTurn
		rightCurrent := e.deathTriggers[j].ownerID == e.State.CurrentTurn
		return leftCurrent && !rightCurrent
	})
	for len(e.deathTriggers) > 0 {
		pending := e.deathTriggers[0]
		e.deathTriggers = e.deathTriggers[1:]
		e.resolveDeathTriggersWithData(pending.unit, pending.ownerID, pending.deathData)
	}
}

func (e *Engine) resolveDeathTriggers(unit *CardInstance, ownerID int) {
	e.resolveDeathTriggersWithData(unit, ownerID, nil)
}

func (e *Engine) resolveDeathTriggersWithData(unit *CardInstance, ownerID int, deathData map[string]any) {
	if unit == nil || unit.Card == nil {
		return
	}

	// Trigger 遗言 (on death) effects
	e.triggerEffects(TriggerOnDeath, unit, nil, deathData)

	// Notify friendly cards about the death
	e.triggerFieldEffectsWithData(TriggerOnFriendlyDeath, ownerID, unit, deathData)

	// Notify enemy cards about the death
	e.triggerFieldEffectsWithData(TriggerOnEnemyDeath, 1-ownerID, unit, deathData)
}

func (e *Engine) unitInOwnerGrid(unit *CardInstance, ownerID int) bool {
	if unit == nil || ownerID < 0 || ownerID >= len(e.State.Players) {
		return false
	}
	ps := e.State.Players[ownerID]
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] == unit {
				return true
			}
		}
	}
	return false
}

// destroyUnit removes a unit from the field and sends it to graveyard
func (e *Engine) destroyUnit(unit *CardInstance, ownerID int) {
	e.destroyUnitWithData(unit, ownerID, nil)
}

func (e *Engine) destroyUnitWithCause(unit *CardInstance, ownerID int, cause string) {
	data := map[string]any{}
	if cause != "" {
		data["death_cause"] = cause
	}
	e.destroyUnitWithData(unit, ownerID, data)
}

func (e *Engine) destroyUnitWithData(unit *CardInstance, ownerID int, deathData map[string]any) {
	e.removeQueuedDeath(unit)
	ps := e.State.Players[ownerID]

	// Remove from grid
	if unit.Position != nil {
		ps.Units[unit.Position.Col][unit.Position.Row] = nil
	}

	// Printed/generated bound skills live only while their host is on the
	// battlefield. Learned skills turned into bound skills have their own exile
	// rule and are handled before clearing the host.
	e.releaseUnderCardsToGraveyard(ownerID, unit)
	e.exileTransferredBoundSkills(ownerID, unit)
	unit.BoundSkills = nil

	// Add to graveyard
	e.addToGraveyard(ownerID, unit)

	e.emit(GameEvent{
		Type:   "unit_destroyed",
		Player: -1,
		Data: map[string]any{
			"player": ownerID,
			"card":   cardToInfo(unit),
		},
	})

	if e.resolvingDeaths || e.resolutionDepth > 0 {
		e.queueDeathTriggersWithData(unit, ownerID, deathData)
	} else {
		e.resolveDeathTriggersWithData(unit, ownerID, deathData)
	}

	// Check if hero died. During queued death resolution, wait until all pending
	// deaths are processed so simultaneous hero deaths can become a draw.
	if unit.Card.IsHero() && !e.resolvingDeaths && e.resolutionDepth == 0 {
		e.checkWinCondition()
	}
}
