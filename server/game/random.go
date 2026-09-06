package game

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"eraofarcane/model"
	"fmt"
	"math/rand"
)

// NewEngineWithSeed supports reproducible engine tests and private replays.
// The seed is secret game data: exposing it to a player reveals deck order.
// All random choices and generated object IDs belong to this engine instance.
func NewEngineWithSeed(gameID string, callback EventCallback, seed int64) *Engine {
	return &Engine{State: NewGameState(gameID), callback: callback,
		log: make([]GameEvent, 0), randomSeed: seed, rng: rand.New(rand.NewSource(seed))}
}

func newRandomSeed() int64 {
	var seed [8]byte
	if _, err := cryptorand.Read(seed[:]); err != nil {
		panic(err)
	}
	return int64(binary.LittleEndian.Uint64(seed[:]))
}

// ReplaySeed is for private server diagnostics, never player/spectator views.
// It is immutable after construction and can be read without taking Engine.mu.
func (e *Engine) ReplaySeed() int64 { return e.randomSeed }

func (e *Engine) nextObjectID() string {
	e.objectSequence++
	return fmt.Sprintf("ci_%s_%d", e.State.GameID, e.objectSequence)
}

func (e *Engine) newCardInstance(card *model.Card, ownerID, turn int) *CardInstance {
	return newCardInstanceWithID(card, ownerID, turn, e.nextObjectID())
}

// randomIntn and shuffleCards run under the same engine ownership/lock as the
// operation using them. They never use the process-global random source.
func (e *Engine) randomIntn(n int) int { return e.rng.Intn(n) }

func (e *Engine) shuffleCards(deck []*CardInstance) {
	e.rng.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}
