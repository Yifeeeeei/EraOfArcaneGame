package game

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRegistryConcurrentLazyPublication(t *testing.T) {
	r := NewEffectRegistry()
	var constructed atomic.Int32
	r.RegisterBehaviorFactory("1021006", func() CardBehavior {
		constructed.Add(1)
		return Card1021006Grocer{}
	})
	var workers sync.WaitGroup
	for i := 0; i < 64; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for j := 0; j < 10; j++ {
				if !r.HasEffect("1021006", TriggerOnEnter) || len(r.GetAllEffects("1021006")) != 1 {
					t.Error("observed an incomplete card adapter")
				}
				if r.GetBehavior("1021006") == nil {
					t.Error("missing behavior")
				}
			}
		}()
	}
	workers.Wait()
	if got := constructed.Load(); got != 1 {
		t.Fatalf("factory constructed %d times", got)
	}
}

func TestCounterExtensionNeedsNoEngineCardList(t *testing.T) {
	previous := globalRegistry
	t.Cleanup(func() { globalRegistry = previous })
	globalRegistry = NewEffectRegistry()
	globalRegistry.RegisterBehaviorFactory("extension_counter", func() CardBehavior { return extensionCounter{} })
	if !isCounterTrapCard("extension_counter") || !counterTrapHasTrigger("extension_counter", TriggerOnDraw) {
		t.Fatal("counter discovery still depends on a central card-number list")
	}
	if counterTrapHasTrigger("extension_counter", TriggerOnConsume) {
		t.Fatal("counter subscribed to an undeclared event")
	}
}

type extensionCounter struct{ AlwaysActive }

func (extensionCounter) ID() string                                 { return "extension_counter" }
func (extensionCounter) Name() string                               { return "extension counter" }
func (extensionCounter) CounterTriggers() []EffectTrigger           { return []EffectTrigger{TriggerOnDraw} }
func (extensionCounter) CanTriggerCounter(ctx *CounterContext) bool { return ctx.Event.DrawCount >= 2 }
