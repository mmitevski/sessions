package memory

import (
	"sync"
	"time"
	"github.com/mmitevski/sessions"
	"log"
)

type provider struct {
	lock     sync.Mutex          // lock
	sessions map[string]*session // save in memory
}

func (p *provider) SessionInit(sid string) (sessions.Session, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	s := &session{
		id: sid,
		timeAccessed: time.Now(),
		started: time.Now(),
		values: make(map[string]interface{}, 0),
		provider:p,
	}
	p.sessions[sid] = s
	log.Printf("Started session %s.", sid)
	return s, nil
}

func (p *provider) SessionRead(sid string) (sessions.Session, error) {
	if session, ok := p.sessions[sid]; ok {
		return session, nil
	} else {
		return p.SessionInit(sid)
	}
}

func (p *provider) SessionDestroy(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if session, ok := p.sessions[sid]; ok {
		delete(p.sessions, sid)
		session.provider = nil
	}
	return nil
}

func (p *provider) DestroyOutdatedSessions(maxLifeTime int64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	now := time.Now().Unix() - maxLifeTime
	for sid, session := range p.sessions {
		if session.timeAccessed.Unix() < now {
			delete(p.sessions, sid)
			session.provider = nil
			log.Printf("Deleted session %s started %v.", sid, session.started)
		}
	}
}

func (p *provider) SessionUpdate(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if session, ok := p.sessions[sid]; ok {
		session.timeAccessed = time.Now()
	}
	return nil
}

func New() sessions.Provider {
	p := &provider{
		sessions: make(map[string]*session, 0),
	}
	return p
}