package memory
import (
	"time"
)

type session struct {
	id           string                      // unique session id
	timeAccessed time.Time                   // last access time
	values       map[interface{}]interface{} // session value stored inside
	provider     *provider
}

func (s *session) update() {
	s.provider.SessionUpdate(s.id)
}

func (s *session) Set(key, value interface{}) error {
	s.values[key] = value
	s.update()
	return nil
}

func (s *session) Get(key interface{}) interface{} {
	s.update()
	if v, ok := s.values[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (s *session) Delete(key interface{}) error {
	delete(s.values, key)
	s.update()
	return nil
}