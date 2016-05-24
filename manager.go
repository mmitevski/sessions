package sessions
import (
	"sync"
	"net/http"
	"net/url"
	"crypto/rand"
	"encoding/base64"
	"time"
	"log"
	"strings"
	"errors"
)

const _SESSION_ID_LENGTH = 32
const _SESSION_COOKIE_NAME_DEFAULT = "gosession"

type Manager struct {
	cookieName  string     //private cookie name
	lock        sync.Mutex // protects session
	provider    Provider
	maxLifeTime int64
	secure      bool
	closed      chan bool
}

func (manager *Manager) sessionId() string {
	b := make([]byte, _SESSION_ID_LENGTH)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func NewManager(provider Provider, cookieName string, maxLifeTime int64, secure bool) (*Manager, error) {
	if provider == nil {
		return nil, errors.New("Attemt to create sessions.Manager instance without session.Provider.")
	}
	cookieName = strings.TrimSpace(cookieName)
	if len(cookieName) <= 0 {
		cookieName = _SESSION_COOKIE_NAME_DEFAULT
	}
	manager := &Manager{
		provider: provider,
		cookieName: cookieName,
		maxLifeTime: maxLifeTime,
		secure: secure,
	}
	defer manager.init()
	return manager, nil
}

func (manager *Manager) Start(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{
			Name: manager.cookieName,
			Value: url.QueryEscape(sid),
			Path: "/", HttpOnly: true,
			MaxAge: int(manager.maxLifeTime),
			Secure:manager.secure,
		}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

//Destroy session
func (manager *Manager) Destroy(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) destroyOutdatedSessions() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if manager.provider != nil {
		manager.provider.DestroyOutdatedSessions(manager.maxLifeTime)
	}
}

func (manager *Manager) init() {
	manager.closed = make(chan bool)
	go func() {
		log.Println("Start monitoring sessions...")
		for {
			log.Println("Checking sessions...")
			select {
			case <- manager.closed:
				log.Println("Stop monitoring sessions...")
				return
			default:
				log.Println("Destroying outdated sessions...")
				manager.destroyOutdatedSessions()
				time.Sleep(10 * time.Minute)
			}
		}
	}()
}