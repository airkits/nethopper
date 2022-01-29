package session

import "sync"

var instance *SessionPool
var once sync.Once

//GS session instance
func GS() *SessionPool {
	once.Do(func() {
		instance = NewSessionPool()

	})
	return instance
}

// CreateSession from session pool
func CreateSession(srcID int32, host string, port string) *Session {
	return GS().Alloc(srcID, host, port)
}

// GetSession get Session By sessionID
func GetSession(sessionID string) *Session {
	if v, ok := GS().Objs.Load(sessionID); ok {
		return v.(*Session)
	}
	return nil

}

// RemoveSession remove from pool
func RemoveSession(sess *Session) {
	GS().Free(sess)
}
