package server

// import (
// 	"fmt"
// 	"sync"
// )

// type InMemorySubjectStore struct {
// 	mu    sync.Mutex
// 	hours map[string]int
// }

// func NewInMemorySubjectStore() *InMemorySubjectStore {
// 	return &InMemorySubjectStore{
// 		hours: make(map[string]int),
// 	}
// }

// func (i *InMemorySubjectStore) GetHours(subject string) (int, error) {
// 	i.mu.Lock()
// 	defer i.mu.Unlock()
// 	h, ok := i.hours[subject]
// 	if !ok {
// 		return 0, fmt.Errorf("failed to find subject %s", subject)
// 	}
// 	return h, nil
// }

// func (i *InMemorySubjectStore) RecordHour(subject string, numHours int) error {
// 	i.mu.Lock()
// 	defer i.mu.Unlock()
// 	i.hours[subject]++
// 	return nil
// }
