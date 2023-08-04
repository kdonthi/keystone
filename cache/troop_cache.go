package cache

import (
	"fmt"
	"github.com/curio-research/go-backend/engine"
	"sync"
)

// TODO want a struct that returns the value of the largest number before it if there is no key for it
type TroopCache struct {
	mu    *sync.Mutex
	cache map[int]*latestUpdateFinder[engine.Pos] // map of tick #: map of troop id: location
}

func NewTroopCache() TroopCache {
	return TroopCache{
		mu:    &sync.Mutex{},
		cache: map[int]*latestUpdateFinder[engine.Pos]{},
	}
}

func (t *TroopCache) AddTroopData(tickNumber, troopID, x, y int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.cache[troopID]; !ok {
		t.cache[troopID] = newLatestUpdateFinder[engine.Pos]()
	}

	t.cache[troopID].addUpdate(tickNumber, engine.Pos{X: x, Y: y})
}

func (t *TroopCache) GetTroopPosition(tickNumber, troopID, secondsBefore int) (engine.Pos, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if troopPositions, ok := t.cache[troopID]; ok {
		if pos, ok := troopPositions.getLatestUpdate(tickNumber); ok {
			return pos, nil
		} else {
			return engine.Pos{}, fmt.Errorf("troop id %v for tick number %v not found", tickNumber, troopID)
		}
	} else {
		return engine.Pos{}, fmt.Errorf("invalid tick number %v", troopID)
	}
}

// keep a store of the last update tick
// have an array as well as a map?
type latestUpdateFinder[K any] struct {
	ticksSeen []int
	cache     map[int]K
}

func newLatestUpdateFinder[K any]() *latestUpdateFinder[K] {
	return &latestUpdateFinder[K]{
		ticksSeen: []int{},
		cache:     map[int]K{},
	}
}

func (l *latestUpdateFinder[K]) addUpdate(tickNumber int, value K) {
	if tickNumber == l.ticksSeen[len(l.ticksSeen)-1] {
		return
	}

	l.ticksSeen = append(l.ticksSeen, tickNumber)
	l.cache[tickNumber] = value
}

func (l *latestUpdateFinder[K]) getLatestUpdate(tickNumber int) (K, bool) {
	res := binarySearchUpperBound(l.ticksSeen, tickNumber)
	if res == -1 {
		return nil, false
	}

	return l.cache[res], true
}

func binarySearchUpperBound(arr []int, target int) int {
	low := 0
	high := len(arr) - 1
	result := -1 // Initialize the result to -1 (indicating not found)

	for low <= high {
		mid := (low + high) / 2

		if arr[mid] < target {
			result = mid  // Update result to the current index
			low = mid + 1 // Target is in the right half
		} else {
			high = mid - 1 // Target is in the left half
		}
	}

	return result
}
