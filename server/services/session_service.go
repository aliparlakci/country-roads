package services

//go:generate mockgen -destination=../mocks/mock_session_service.go -package=mocks github.com/aliparlakci/country-roads/services SessionRepository,SessionFetcher,SessionUpdater,SessionCreator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type Session struct {
	UserId string `json:"userId"`
	
}

type SessionService struct {
	Client *redis.Client
}

type SessionRepository interface {
	SessionFetcher
	SessionCreator
	SessionUpdater
}

type SessionFetcher interface {
	FetchSession(c context.Context, sessionId string) (Session, error)
}

type SessionCreator interface {
	CreateSession(c context.Context, value Session) (string, error)
}

type SessionUpdater interface {
	UpdateSession(c context.Context, sessionId string, value Session) error
}

func (s SessionService) FetchSession(c context.Context, sessionId string) (Session, error) {
	var session Session
	raw, err := s.Client.Get(sessionId).Result()
	if err != nil {
		return session, errors.New("session does not exist")
	}
	if err := json.Unmarshal([]byte(raw), &session); err != nil {
		return session, fmt.Errorf("cannot parse session value: %v", raw)
	}
	return session, nil
}

func (s SessionService) CreateSession(c context.Context, value Session) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", errors.New("uuid cannot get created")
	}
	rawSession, err := json.Marshal(value)
	if err != nil {
		return "", errors.New("cannot marshal given session value")
	}
	s.Client.Set(id.String(), string(rawSession), 0)

	return id.String(), nil
}

func (s SessionService) UpdateSession(c context.Context, sessionId string, value Session) error {
	return nil
}
