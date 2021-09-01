package services

import (
	"crypto/rand"
	"fmt"
	"math"
	"sync"

	"github.com/go-redis/redis"
)

type AuthStore struct {
	Store *redis.Client
	mutex sync.Mutex
}

type AuthService interface {
	CreateOTP(owner string) error
	RevokeOTP(owner string) error
	VerifyOTP(owner, otp string) (bool, error)
}

func (a *AuthStore) CreateOTP(owner string) error {
	if err := a.Store.Get(owner).Err(); err == nil {
		return nil
	}

	bytes := make([]byte, 3)
	rand.Read(bytes)

	otp := uint32(0)
	for i, b := range bytes {
		otp += uint32(b) * uint32(math.Pow(10.0, float64(i*3)))
	}

	otp = otp % 999999

	// TODO: Set expiration
	return a.Store.Set(owner, fmt.Sprintf("%06d",otp), 0).Err()
}

func (a *AuthStore) VerifyOTP(owner, givenOTP string) (bool, error) {
	otpInDb, err := a.Store.Get(owner).Result()
	if err != nil {
		return false, fmt.Errorf("owner does not exist in DB [owner=%v]", owner)
	}
	return givenOTP != "" && otpInDb == givenOTP, err
}

func (a *AuthStore) RevokeOTP(owner string) error {
	return a.Store.Del(owner).Err()
}
