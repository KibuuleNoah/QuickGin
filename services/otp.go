package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/KibuuleNoah/QuickGin/db"
	"github.com/KibuuleNoah/QuickGin/models"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
)

const (
	otpTTL                    = 2 * time.Minute
	otpMaxAttempts            = 5
	otpMaxDailyRequests       = 9
	otpKeyPrefix              = "otp:"
	otpAttPrefix              = "otp_att:"
	otpMaxDailyRequestsPrefix = "otp_daily_limit:"
)

// OTPService handles generation, storage, and verification of OTPs via Redis.
type OTPService struct {
	rdb *redis.Client
	DB  *sqlx.DB
}

func NewOTPService(rdb *redis.Client) *OTPService {
	return &OTPService{
		rdb: rdb,
		DB:  db.GetDB(),
	}
}

// Generate creates a 6-digit OTP, stores it in Redis, and returns it.
// The key is "otp:<identifier>" where identifier is an e-mail or E.164 phone number.

func (s *OTPService) Generate(ctx context.Context, identifier string) (otp string, userId string, err error) {
	var user models.User
	err = s.DB.Get(&user, "SELECT id, identifier, password, name, updated_at, created_at FROM public.user WHERE identifier=LOWER($1) LIMIT 1", identifier)
	if err != nil {
		return "", "", errors.New("user does not exist")
	}

	otpKey := otpKeyPrefix + user.ID
	dailyKey := otpMaxDailyRequestsPrefix + user.ID

	// Check Max Daily Limit
	dailyCount, err := s.rdb.Get(dailyKey).Int()
	if err != nil && err != redis.Nil {
		return "", "", ErrOTPMaxDailyRequestsHit
	}

	log.Println("dailyCount**** ", dailyCount)
	if dailyCount >= otpMaxDailyRequests {
		return "", "", errors.New("daily otp limit reached")
	}

	// Check Cool Down
	_, err = s.rdb.Get(otpKey).Result()
	if err == nil {
		return "", "", ErrOTPCoolDownActive
	} else if err != redis.Nil {
		return "", "", err
	}

	// Generate and Save
	otp, err = generateSecureOTP()
	if err != nil {
		return "", "", err
	}

	pipe := s.rdb.Pipeline()
	// Set the OTP with cooldown TTL
	pipe.Set(otpKey, otp, otpTTL)

	// ReSet the OTP attempts counter
	pipe.Set(otpAttPrefix+user.ID, 0, otpTTL)

	// Increment daily counter and set 24h expiry
	pipe.Incr(dailyKey)

	pipe.Expire(dailyKey, 24*time.Hour)

	_, err = pipe.Exec()

	return otp, user.ID, err
}

// Verify checks the OTP for the given identifier.
// Returns true on success and deletes the stored OTP so it cannot be reused.
// Returns false + a descriptive error on failure.
func (s *OTPService) Verify(ctx context.Context, otp string, userId string) (bool, error) {
	attKey := otpAttPrefix + userId
	otpKey := otpKeyPrefix + userId

	// Increment attempt counter before checking — prevents timing-based brute force
	attempts, err := s.rdb.Incr(attKey).Result()
	if err != nil {
		return false, fmt.Errorf("track attempts: %w", err)
	}
	// Give the attempt key the same TTL as the OTP itself so it auto-cleans
	s.rdb.Expire(attKey, otpTTL)

	log.Println("****", attempts, otpMaxAttempts)
	if attempts > otpMaxAttempts {
		return false, ErrTooManyAttempts
	}

	stored, err := s.rdb.Get(otpKey).Result()
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
	s.rdb.Del(otpKey, attKey)
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

// DetectType checks if a string is a valid email or mobile number (E.164 format)
func detectIdentifierType(input string) string {
	input = strings.TrimSpace(input)

	// Check for Email using the standard library
	if _, err := mail.ParseAddress(input); err == nil {
		// Standard ParseAddress can accept "Name <email@domain.com>",
		// so we check if the input contains '@' to be safe for simple strings.
		if strings.Contains(input, "@") && !strings.Contains(input, " ") {
			return "email"
		}
	}

	// Check for Mobile Number (E.164 format: +1234567890)
	// This regex matches an optional '+' followed by 10 to 15 digits.
	mobileRegex := regexp.MustCompile(`^\+?[1-9]\d{9,14}$`)
	if mobileRegex.MatchString(input) {
		return "mobile"
	}

	return "invalid"
}

// Sentinel errors
var (
	ErrOTPMaxDailyRequestsHit = fmt.Errorf("otp daily request limit hit, request new OTP within 24")
	ErrOTPCoolDownActive      = fmt.Errorf("otp request denied; please wait for the cooldown period to expire")
	ErrOTPNotFound            = fmt.Errorf("OTP not found or expired")
	ErrInvalidOTP             = fmt.Errorf("invalid OTP")
	ErrTooManyAttempts        = fmt.Errorf("too many verification attempts")
)

// func (s *OTPService) Generate(ctx context.Context, identifier string) (otp string, userId string, err error) {
//
// 	var user models.User
//
// 	err = s.DB.Get(&user, "SELECT id, identifier, password, name, updated_at, created_at FROM public.user WHERE identifier=LOWER($1) LIMIT 1", identifier)
//
// 	if err != nil {
// 		return "", "", errors.New("User Doesn't Exists")
// 	}
//
// 	otpKey := otpAttPrefix + user.ID
//
//
// 	_, err = s.rdb.Get(otpKey).Result()
// 	if err == redis.Nil {
// 		return "", "", ErrOTPCoolDownActive
// 	}
//
// 	otp, err = generateSecureOTP()
//
// 	pipe := s.rdb.Pipeline()
//
// 	pipe.Set(otpKey, otp, otpTTL)
// 	// Reset attempt counter whenever a fresh OTP is issued
// 	pipe.Del(otpKey)
// 	_, err = pipe.Exec()
//
// 	userId = user.ID
// 	return otp, userId, err
// }
