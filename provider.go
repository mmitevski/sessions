package sessions

type Provider interface {
	// implements the initialization of a session, and returns new a session if it succeeds.
	SessionInit(sid string) (Session, error)

	// returns a session represented by the corresponding sid. Creates a new session and returns it
	// if it does not already exist.
	SessionRead(sid string) (Session, error)

	// given an sid, deletes the corresponding session.
	SessionDestroy(sid string) error

	// deletes expired session variables according to maxLifeTime.
	DestroyOutdatedSessions(maxLifeTime int64)
}