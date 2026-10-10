package game

// resolutionFrame owns the remaining steps of a sequential effect. Choices and
// spells hold waiters, not wrappers around card callbacks. The baseline spell
// belongs to the caller: a before-hit choice must resume that spell rather than
// wait for the very spell its next step is supposed to finish.
//
// This is a continuation protocol, not a replacement for simultaneous damage
// and death batches. Steps still execute within the existing resolution scope.
type resolutionFrame struct {
	id            uint64
	parent        uint64
	name          string
	steps         []func()
	next          int
	baselineSpell *SpellCast
	waitingAction *PendingAction
	waitingSpell  *SpellCast
	done          bool
}

// runResolution executes sequential steps until one opens a choice or a new
// spell. Unlike Combine, later steps wait for the complete interaction.
func (e *Engine) runResolution(name string, steps ...func()) *resolutionFrame {
	frame := e.newResolutionFrame(name, steps)
	e.resumeResolution(frame)
	return frame
}

func (e *Engine) continueAfterPendingAction(continueFn func()) {
	if e.State.PendingAction == nil || continueFn == nil {
		return
	}
	e.addActionContinuation(e.State.PendingAction, "after choice", continueFn)
}

func (e *Engine) addActionContinuation(action *PendingAction, name string, continueFn func()) {
	if action == nil || continueFn == nil {
		return
	}
	frame := e.newResolutionFrame(name, []func(){continueFn})
	frame.waitingAction = action
	e.traceResolution("wait_choice", action.Type, frame)
	action.resolutions = append(action.resolutions, frame)
}

func (e *Engine) resumeResolution(frame *resolutionFrame) {
	if frame == nil || frame.done || frame.waitingAction != nil || frame.waitingSpell != nil {
		return
	}
	for {
		if e.State.Phase == PhaseGameOver {
			frame.done = true
			frame.steps = nil
			e.traceResolution("frame_stopped", frame.name, frame)
			return
		}
		if action := e.State.PendingAction; action != nil {
			frame.waitingAction = action
			e.traceResolution("wait_choice", action.Type, frame)
			action.resolutions = append(action.resolutions, frame)
			return
		}
		if spell := e.State.PendingSpell; spell != nil && spell != frame.baselineSpell {
			frame.waitingSpell = spell
			e.traceResolution("wait_spell", frame.name, frame)
			spell.resolutions = append(spell.resolutions, frame)
			return
		}
		if frame.next == len(frame.steps) {
			e.traceResolution("frame_complete", frame.name, frame)
			frame.done = true
			frame.steps = nil
			return
		}
		step := frame.steps[frame.next]
		frame.next++ // commit the step before it can create a nested interaction
		if step != nil {
			e.traceResolution("step", frame.name, frame)
			func() {
				previous := e.activeFrame
				e.activeFrame = frame
				defer func() { e.activeFrame = previous }()
				step()
			}()
		}
	}
}

// completeActionResolutions is called only after a successful callback or a
// queued action becoming unavailable. A failed choice retains its waiters for
// retry. Existing queued choices are activated before parents are resumed.
func (e *Engine) completeActionResolutions(action *PendingAction) {
	waiters := action.resolutions
	action.resolutions = nil
	for _, frame := range waiters {
		frame.waitingAction = nil
		e.resumeResolution(frame)
	}
}

func (e *Engine) completePendingSpell(spell *SpellCast) bool {
	if spell == nil || e.State.PendingSpell != spell {
		return false
	}
	e.State.PendingSpell = nil
	waiters := spell.resolutions
	spell.resolutions = nil
	continued := false
	for _, frame := range waiters {
		continued = continued || (!frame.done && frame.next < len(frame.steps))
		frame.waitingSpell = nil
		e.resumeResolution(frame)
	}
	// A frame whose final step started this spell has nothing left to run.
	// Let spell callers restore main phase in that case, but never overwrite
	// a replacement spell or game-over established during completion.
	return continued || e.State.PendingSpell != nil || e.State.Phase == PhaseGameOver
}

// replacePendingSpell preserves the parent interaction when a counter reflects
// a spell. start must validate before mutation on error, as other action
// preparation helpers do. A failed replacement leaves the original waiters
// intact so the player can retry the choice.
func (e *Engine) replacePendingSpell(start func() error) error {
	original := e.State.PendingSpell
	e.State.PendingSpell = nil
	if err := start(); err != nil {
		e.State.PendingSpell = original
		return err
	}
	if original == nil {
		return nil
	}
	clearFiveRainbowBeamSelection(original.Skill)
	waiters := original.resolutions
	original.resolutions = nil
	for _, frame := range waiters {
		frame.waitingSpell = nil
		e.resumeResolution(frame)
	}
	return nil
}

func (e *Engine) beginResolution() {
	e.resolutionDepth++
}

func (e *Engine) endResolution() {
	if e.resolutionDepth > 0 {
		e.resolutionDepth--
	}
	if e.resolutionDepth == 0 {
		e.resolvePendingDeaths()
	}
}
