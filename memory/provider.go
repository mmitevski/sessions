package memory
import (
	"sv/sessions"
	"sync"
	"container/list"
	"time"
)


type provider struct {
	lock     sync.Mutex               // lock
	sessions map[string]*list.Element // save in memory
	list     *list.List               // gc
}

func (p *provider) SessionInit(sid string) (sessions.Session, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	s := &session{
		id: sid,
		timeAccessed: time.Now(),
		values: v,
		provider:p,
	}
	element := p.list.PushBack(s)
	p.sessions[sid] = element
	return s, nil
}

func (p *provider) SessionRead(sid string) (sessions.Session, error) {
	if element, ok := p.sessions[sid]; ok {
		return element.Value.(*session), nil
	} else {
		sess, err := p.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}

func (p *provider) SessionDestroy(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if element, ok := p.sessions[sid]; ok {
		delete(p.sessions, sid)
		p.list.Remove(element)
		element.Value.(*session).provider = nil
	}
	return nil
}

func (p *provider) DestroyOutdatedSessions(maxLifeTime int64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for element := p.list.Front(); element != nil; element = element.Next() {
		if (element.Value.(*session).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
			p.list.Remove(element)
			delete(p.sessions, element.Value.(*session).id)
		}
	}
}

func (p *provider) SessionUpdate(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if element, ok := p.sessions[sid]; ok {
		element.Value.(*session).timeAccessed = time.Now()
		p.list.MoveToFront(element)
		return nil
	}
	return nil
}

func New() sessions.Provider {
	p := &provider{
		sessions: make(map[string]*list.Element, 0),
		list: list.New(),
	}
	return p
}