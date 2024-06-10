package store

import (
	"sync"
	"time"
)

type Values struct {
	Value string
	ExpireTo time.Time
}

type Storage struct {
	m sync.Map
}

func NewStorage() *Storage {
	return &Storage{
		m: sync.Map{},
	}
}

func (s *Storage) GetMap() *sync.Map {
	return &s.m
}

func (s *Storage) Initialize(m map[string]Values) {
	s.m = sync.Map{}
	for k, v := range m {
		s.m.Store(k, v)
	}
}

func (s *Storage) Load(key string) (Values, bool) {
	v, ok := s.m.Load(key)
	if !ok {
		return Values{}, ok
	}
	
	str, ok := v.(Values)
	if !ok {
		return Values{}, ok
	}

	return str, ok
}

func (s *Storage) Store(key string, value Values) {
	s.m.Store(key, value)
}

func (s *Storage) Delete(key string) {
	s.m.Delete(key)
}

func (s *Storage) Range(f func(key string, value Values) bool) {
	s.m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(Values))
	})
}
