package cache

import (
	"sync"

	"github.com/paveloborin/fetch-task/pkg/model"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	data map[uuid.UUID]*model.Task
	mx   *sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[uuid.UUID]*model.Task),
		mx:   &sync.RWMutex{},
	}
}

func (s *Storage) Get(id uuid.UUID) (*model.Task, bool) {
	s.mx.RLock()
	val, ok := s.data[id]
	s.mx.RUnlock()
	if !ok {
		return nil, false
	}

	return val, true
}

func (s *Storage) Add(task *model.Task) {
	s.mx.Lock()
	s.data[task.ID] = task
	s.mx.Unlock()
}

func (s *Storage) Delete(id uuid.UUID) bool {
	s.mx.Lock()
	defer s.mx.Unlock()

	_, ok := s.data[id]
	if !ok {

		return false
	}
	delete(s.data, id)

	return true
}

func (s *Storage) GetAll(pageNum, count int) []*model.Task {
	var flat []*model.Task

	s.mx.RLock()
	for _, value := range s.data {
		flat = append(flat, value)
	}
	s.mx.RUnlock()

	if pageNum != 0 && count != 0 {
		start := (pageNum - 1) * count
		if len(flat) < start {
			return []*model.Task{}
		}

		end := start + count
		if len(flat) < end {
			end = len(flat)
		}

		return flat[start:end]
	}
	return flat
}
