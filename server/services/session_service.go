package services

//go:generate mockgen -destination=../mocks/mock_session_service.go -package=mocks github.com/aliparlakci/country-roads/services SessionRepository,SessionFetcher,SessionUpdater,SessionCreator

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type SessionStore struct {
	Store *redis.Client
}

type SessionService interface {
	SessionFetcher
	SessionCreator
	SessionRevoker
}

type SessionFetcher interface {
	FetchSession(c context.Context, sessionId string) (string, error)
}

type SessionCreator interface {
	CreateSession(c context.Context, userId string) (string, error)
}

type SessionRevoker interface {
	RevokeSession(c context.Context, sessionId string) error
}

func (s *SessionStore) FetchSession(c context.Context, sessionId string) (string, error) {
	session, err := s.Store.Get(sessionId).Result()
	if err != nil {
		return session, fmt.Errorf("session does not exist [sessionId=%v]", sessionId)
	}
	return session, nil
}

func (s *SessionStore) CreateSession(c context.Context, userId string) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("cannot create new uuid: %v",err.Error())
	}
	// TODO: Set expiration
	if err := s.Store.Set(id.String(), userId, 0).Err(); err != nil {
		return "", fmt.Errorf("cannot create a new session: %v", err.Error())
	}

	return id.String(), nil
}

func (s *SessionStore) RevokeSession(c context.Context, sessionId string) error {
	return s.Store.Del(sessionId).Err()
}