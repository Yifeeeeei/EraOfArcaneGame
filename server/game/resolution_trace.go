package game

const resolutionTraceCapacity = 512

// ResolutionTraceEntry is a private debugging record, never a player event.
// It describes execution order and wait relationships without copying state.
type ResolutionTraceEntry struct {
	Sequence uint64 `json:"sequence"`
	Action   uint64 `json:"action"`
	Frame    uint64 `json:"frame,omitempty"`
	Parent   uint64 `json:"parent,omitempty"`
	Kind     string `json:"kind"`
	Name     string `json:"name,omitempty"`
	Step     int    `json:"step,omitempty"`
	Turn     int    `json:"turn"`
}

func (e *Engine) traceResolution(kind, name string, frame *resolutionFrame) {
	e.traceSequence++
	entry := ResolutionTraceEntry{Sequence: e.traceSequence, Action: e.actionSequence, Kind: kind, Name: name, Turn: e.State.TurnNumber}
	if frame != nil {
		entry.Frame = frame.id
		entry.Parent = frame.parent
		entry.Step = frame.next
	}
	e.resolutionTrace[(e.traceSequence-1)%resolutionTraceCapacity] = entry
}

func (e *Engine) newResolutionFrame(name string, steps []func()) *resolutionFrame {
	e.frameSequence++
	frame := &resolutionFrame{id: e.frameSequence, name: name, steps: steps, baselineSpell: e.State.PendingSpell}
	if e.activeFrame != nil {
		frame.parent = e.activeFrame.id
	}
	e.traceResolution("frame_start", name, frame)
	return frame
}

// DebugResolutionTrace returns a bounded copy in execution order. Do not add
// this private diagnostic stream to public or player state serialization.
func (e *Engine) DebugResolutionTrace() []ResolutionTraceEntry {
	e.mu.Lock()
	defer e.mu.Unlock()
	count := min(e.traceSequence, uint64(resolutionTraceCapacity))
	result := make([]ResolutionTraceEntry, 0, count)
	for seq := e.traceSequence - count; seq < e.traceSequence; seq++ {
		result = append(result, e.resolutionTrace[seq%resolutionTraceCapacity])
	}
	return result
}
