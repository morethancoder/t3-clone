package services

import (
	"morethancoder/t3-clone/utils"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

const expiryKey = "expiry"

var SessionStore sync.Map

type Session struct {
	ID     string
	Data   map[string]interface{}
	Expiry time.Time
}

func NewSession(data map[string]interface{}) *Session {
	session := &Session{
		ID:   uuid.NewString(),
		Data: data,
	}
	if session.Expiry.IsZero() {
		if os.Getenv("ENV") == "dev" {
			session.Expiry = time.Now().Add(time.Minute * 10)
		} else {
			session.Expiry = time.Now().Add(time.Hour * 24)
		}
	}
	session.Save()
	return session
}

func (s *Session) Save() {
	s.Data[expiryKey] = s.Expiry
	SessionStore.Store(s.ID, s.Data)
	utils.Log.Debug("Saved session %s", s.ID)
}

func (s *Session) Load() bool {
	value, ok := SessionStore.Load(s.ID)
	if ok {
		data := value.(map[string]interface{})
		s.Data = data
		s.Expiry = data[expiryKey].(time.Time)
	}
	utils.Log.Debug("Loaded session %s", s.ID)
	return ok
}

func hasSessions() bool {
	has := false
	SessionStore.Range(func(key, value any) bool {
		has = true
		return false
	})
	utils.Log.Debug("Has sessions: %v", has)
	return has
}

func LoopAndCleanSessionStore() {
	// ctx, _ := context.WithCancel(context.Background())
	for {
		// select {
		// case <- ctx.Done():
		// 	utils.Log.Debug("LoopAndCleanSessionStore stopped")
		// 	return
		// default:
		// }

		if os.Getenv("ENV") == "dev" {
			time.Sleep(time.Minute * 1)
		} else {
			time.Sleep(time.Hour * 1)
		}

		if !hasSessions() {
			continue
		}

		now := time.Now()
		utils.Log.Debug("Cleaning session store %v", now.Format(time.RFC3339))
		SessionStore.Range(func(key, value any) bool {
			expiry, ok := value.(map[string]interface{})[expiryKey].(time.Time)
			if ok {
				if expiry.Before(now) {
					SessionStore.Delete(key)
					utils.Log.Debug("Deleted session %s", key)
				}
			}
			return true
		})
		utils.Log.Debug(
			"Cleaned session store %v (%v)",
			time.Now().Format(time.RFC3339),
			time.Duration(time.Now().UnixNano()-now.UnixNano()))
	}
}
