package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	otpTTL        = 5 * time.Minute
	otpMaxAttempts = 5
	otpKeyPrefix  = "otp:"
	otpAttPrefix  = "otp_att:"
)

// OTPService handles generation, storage, and verification of OTPs via Redis.
type OTPService struct {
	rdb *redis.Client
}

func NewOTPService(rdb *redis.Client) *OTPService {
	return &OTPService{rdb: rdb}
}

// Generate creates a 6-digit OTP, stores it in Redis, and returns it.
// The key is "otp:<identifier>" where identifier is an e-mail or E.164 phone number.
func (s *OTPService) Generate(ctx context.Context, identifier string) (string, error) {
	otp, err := generateSecureOTP()
	if err != nil {
		return "", fmt.Errorf("generate otp: %w", err)
	}

	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, otpKeyPrefix+identifier, otp, otpTTL)
	// Reset attempt counter whenever a fresh OTP is issued
	pipe.Del(ctx, otpAttPrefix+identifier)
	if _, err = pipe.Exec(ctx); err != nil {
		return "", fmt.Errorf("store otp: %w", err)
	}
	return otp, nil
}

// Verify checks the OTP for the given identifier.
// Returns true on success and deletes the stored OTP so it cannot be reused.
// Returns false + a descriptive error on failure.
func (s *OTPService) Verify(ctx context.Context, identifier, otp string) (bool, error) {
	attKey := otpAttPrefix + identifier
	otpKey := otpKeyPrefix + identifier

	// Increment attempt counter before checking — prevents timing-based brute force
	attempts, err := s.rdb.Incr(ctx, attKey).Result()
	if err != nil {
		return false, fmt.Errorf("track attempts: %w", err)
	}
	// Give the attempt key the same TTL as the OTP itself so it auto-cleans
	s.rdb.Expire(ctx, attKey, otpTTL)

	if attempts > otpMaxAttempts {
		return false, ErrTooManyAttempts
	}

	stored, err := s.rdb.Get(ctx, otpKey).Result()
	if err == redis.Nil {
		return false, ErrOTPNotFound
	}
	if err != nil {
		return false, fmt.Errorf("get otp: %w", err)
	}

	if stored != otp {
		return false, ErrInvalidOTP
	}

	// Delete both keys on success — OTP is single-use
	s.rdb.Del(ctx, otpKey, attKey)
	return true, nil
}

// generateSecureOTP returns a cryptographically secure 6-digit string.
func generateSecureOTP() (string, error) {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// Sentinel errors
var (
	ErrOTPNotFound     = fmt.Errorf("OTP not found or expired")
	ErrInvalidOTP      = fmt.Errorf("invalid OTP")
	ErrTooManyAttempts = fmt.Errorf("too many verification attempts")
)
