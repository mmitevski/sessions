package memory

import (
	"time"
)

type session struct {
	id           string                 // unique session id
	timeAccessed time.Time              // last access time
	started      time.Time              // session start time
	values       map[string]interface{} // session value stored inside
	provider     *provider
}

func (s *session) update() {
	s.provider.SessionUpdate(s.id)
}

func (s *session) Set(key string, value interface{}) error {
	s.values[key] = value
	s.update()
	return nil
}

func (s *session) Get(key string) interface{} {
	if v, ok := s.values[key]; ok {
		s.update()
		return v
	}
	return nil
}

func (s *session) Delete(key string) error {
	delete(s.values, key)
	s.update()
	return nil
}
