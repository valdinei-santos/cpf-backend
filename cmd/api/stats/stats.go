package stats

import (
	"sync"
	"time"
)

type AccessStats struct {
	mu        sync.RWMutex
	Accesses  map[string]int
	StartTime time.Time
}

var GlobalStats = NewAccessStats()

func NewAccessStats() *AccessStats {
	return &AccessStats{
		Accesses:  make(map[string]int),
		StartTime: time.Now(),
	}
}

func (s *AccessStats) Increment(path string) {
	s.mu.Lock()         // Bloqueia para escrita
	defer s.mu.Unlock() // Libera após a execução
	s.Accesses[path]++
}

func (s *AccessStats) GetStats() map[string]int {
	s.mu.RLock()         // Bloqueia apenas para leitura
	defer s.mu.RUnlock() // Libera a leitura

	// Cria e retorna uma cópia para evitar modificações externas
	statsCopy := make(map[string]int)
	for k, v := range s.Accesses {
		statsCopy[k] = v
	}
	return statsCopy
}
