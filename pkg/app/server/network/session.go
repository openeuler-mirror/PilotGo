package network

import (
	"errors"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	sconfig "openeuler.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeuler.org/PilotGo/PilotGo/sdk/logger"
)

const (
	default_age           = 60 * 30
	default_session_count = 100
)

type SessionManage struct {
	sessionMap map[string]*SessionInfo
	mutex      sync.RWMutex
	maxAge     int
	maxLen     int
}

type SessionInfo struct {
	sessionTime time.Time
}

func SessionManagerInit(conf *sconfig.HttpServer) error {
	var sessionManage SessionManage
	sessionManage.Init(conf.SessionMaxAge, conf.SessionCount)
	return nil
}

func (s *SessionManage) Init(maxAge, maxLen int) bool {
	logger.Debug("the MaxAge is %d,the SessionCount is %d", maxAge, maxLen)
	s.maxAge = maxAge
	if s.maxAge <= 0 {
		s.maxAge = default_age
	}

	s.maxLen = maxLen
	if s.maxLen <= 0 {
		s.maxLen = default_session_count
	}
	s.sessionMap = make(map[string]*SessionInfo, maxLen)
	go checkOutSessionId(s)
	return true
}

func (s *SessionManage) Set(key string, value *SessionInfo) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	mapLen := len(s.sessionMap)
	if mapLen > s.maxLen {
		return errors.New("out of len")
	}

	s.sessionMap[key] = value
	logger.Debug("set the session id:%s", key)
	return nil
}

func (s *SessionManage) FindAndFlush(key string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	session, ok := s.sessionMap[key]
	if ok {
		session.sessionTime = time.Now()
	}
	return ok
}

func CreateSessionId() string {
	return uuid.NewV4().String()
}

func checkOutSessionId(s *SessionManage) {
	for {
		func() {
			s.mutex.RLock()
			defer s.mutex.RUnlock()
			for k, v := range s.sessionMap {
				now := time.Now()
				sub := now.Sub(v.sessionTime)
				if sub > time.Duration(s.maxAge*1000000000) || sub < 0 {
					logger.Debug("rm the session %s\n", k)
					delete(s.sessionMap, k)
				}
			}
		}()
		time.Sleep(time.Second * 10)
	}
}
