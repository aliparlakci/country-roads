package services

//go:generate mockgen -destination=../mocks/mock_session_service.go -package=mocks github.com/aliparlakci/country-roads/services SessionRepository,SessionFetcher,SessionUpdater,SessionCreator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"sync"
)

type Session struct {
	UserId string `json:"userId,omitempty"`

}

type SessionStore struct {
	Store *redis.Client
	mutex sync.Mutex
}

type SessionService interface {
	SessionFetcher
	SessionCreator
	SessionUpdater
}

type SessionFetcher interface {
	FetchSession(c context.Context, sessionId string) (Session, error)
	Lock()
	Unlock()
}

type SessionCreator interface {
	CreateSession(c context.Context, value Session) (string, error)
	Lock()
	Unlock()
}

type SessionUpdater interface {
	UpdateSession(c context.Context, sessionId string, value Session) error
	Lock()
	Unlock()
}

func (s *SessionStore) FetchSession(c context.Context, sessionId string) (Session, error) {
	var session Session
	raw, err := s.Store.Get(sessionId).Result()
	if err != nil {
		return session, errors.New("session does not exist")
	}
	if err := json.Unmarshal([]byte(raw), &session); err != nil {
		return session, fmt.Errorf("cannot parse session value: %v", raw)
	}
	return session, nil
}

func (s *SessionStore) CreateSession(c context.Context, value Session) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", errors.New("uuid cannot get created")
	}
	rawSession, err := json.Marshal(value)
	if err != nil {
		return "", errors.New("cannot marshal given session value")
	}
	// TODO: Set expiration
	if err := s.Store.Set(id.String(), string(rawSession), 0).Err(); err != nil {
		return "", fmt.Errorf("cannot create a new session: %v", err)
	}

	return id.String(), nil
}

func (s *SessionStore) UpdateSession(c context.Context, sessionId string, value Session) error {
	if _, err := s.FetchSession(c, sessionId); err != nil {
		return fmt.Errorf("sessionId does not exist [sessionId=%v]", sessionId)
	}
	rawSession, err := json.Marshal(value)
	if err != nil {
		return errors.New("cannot marshal given session value")
	}
	if err := s.Store.Do("SET", sessionId, string(rawSession), "keepttl").Err(); err != nil {
		return fmt.Errorf("session could not get update: %v", err)
	}
	return nil
}

func (s *SessionStore) Lock() {
	s.mutex.Lock()
}

func (s *SessionStore) Unlock() {
	s.mutex.Unlock()
}
