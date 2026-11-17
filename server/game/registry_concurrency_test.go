package game

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Use an isolated registry so this exercises first-time initialization even when
// earlier tests have already initialized the process-wide production catalog.
// Global-registry fixtures are sequential; all workers finish before restoration.
func TestRegisterAllCardEffectsConcurrentInitializationAndReads(t *testing.T) {
	previous := globalRegistry
	t.Cleanup(func() { globalRegistry = previous })
	r := NewEffectRegistry()
	globalRegistry = r

	start := make(chan struct{})
	var workers sync.WaitGroup
	for i := 0; i < 64; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			RegisterAllCardEffects()
			if GetEffectRegistry() != r {
				t.Error("registration replaced the shared registry")
				return
			}
			r.mu.RLock()
			count := len(r.behaviorFactories)
			entry := r.behaviorFactories["1021006"]
			r.mu.RUnlock()
			if count != len(baseSetBehaviorFactories) || entry == nil {
				t.Error("registration returned before the complete catalog was available")
				return
			}
			behavior := behaviorForNumber("1021006")
			for j := 0; j < 10; j++ {
				RegisterAllCardEffects()
				if behaviorForNumber("1021006") != behavior || !r.HasEffect("1021006", TriggerOnEnter) || len(r.GetAllEffects("1021006")) != 1 {
					t.Error("re-registration changed or duplicated a loaded card adapter")
					return
				}
				r.mu.RLock()
				current := r.behaviorFactories["1021006"]
				r.mu.RUnlock()
				if current != entry {
					t.Error("re-registration replaced a lazy factory entry")
					return
				}
			}
		}()
	}
	close(start)
	workers.Wait()
}

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
