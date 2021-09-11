package services

import (
	"crypto/rand"
	"fmt"
	"math"
	"sync"

	"github.com/go-redis/redis"
)

type OTPStore struct {
	Store *redis.Client
	mutex sync.Mutex
}

type OTPService interface {
	CreateOTP(owner string) (string, error)
	RevokeOTP(owner string) error
	VerifyOTP(owner, otp string) (bool, error)
}

func (a *OTPStore) CreateOTP(owner string) (string, error) {
	if existingOtp, err := a.Store.Get(owner).Result(); err == nil {
		return existingOtp, nil
	}

	bytes := make([]byte, 3)
	rand.Read(bytes)

	randomNumber := uint32(0)
	for i, b := range bytes {
		randomNumber += uint32(b) * uint32(math.Pow(10.0, float64(i*3)))
	}
	randomNumber = randomNumber%999999

	otp := fmt.Sprintf("%06d", randomNumber)

	// TODO: Set expiration
	return otp, a.Store.Set(owner, otp, 0).Err()
}

func (a *OTPStore) VerifyOTP(owner, givenOTP string) (bool, error) {
	otpInDb, err := a.Store.Get(owner).Result()
	if err != nil {
		return false, fmt.Errorf("owner does not exist in DB [owner=%v]", owner)
	}
	return givenOTP != "" && otpInDb == givenOTP, err
}

func (a *OTPStore) RevokeOTP(owner string) error {
	return a.Store.Del(owner).Err()
}
