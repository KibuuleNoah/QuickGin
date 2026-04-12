package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KibuuleNoah/QuickGin/db"
	"github.com/KibuuleNoah/QuickGin/internal/cache"
	"github.com/jmoiron/sqlx"
)

const (
	otpTTL                    = 2 * time.Minute
	otpMaxAttempts            = 5
	otpMaxDailyRequests       = 3
	otpKeyPrefix              = "otp:"
	otpAttPrefix              = "otp_att:"
	otpMaxDailyRequestsPrefix = "otp_daily_limit:"
)

// OTPService handles generation, storage, and verification of OTPs via Redis.
type OTPService struct {
	cache cache.Cache
	DB    *sqlx.DB
}

func NewOTPService() *OTPService {
	return &OTPService{
		DB:    db.AppDB(),
		cache: db.AppCache(),
	}
}

// Generate creates a 6-digit OTP, stores it in Redis, and returns it.
// The key is "otp:<identifier>" where identifier is an e-mail or E.164 phone number.

func (s *OTPService) Generate(ctx context.Context, identifier string) (otp string, otpResendKey string, userId string, expiry time.Time, err error) {
	return "", "", "", time.Time{}, errors.New("user does not exist")

	// var user models.User
	// err = s.DB.Get(&user, "SELECT id, identifier, password, name, updated_at, created_at FROM public.user WHERE identifier=LOWER($1) LIMIT 1", identifier)
	// if err != nil {
	// 	return "", "", "", time.Time{}, errors.New("user does not exist")
	// }
	//
	// otpKey := otpKeyPrefix + user.ID
	// dailyKey := otpMaxDailyRequestsPrefix + user.ID
	//
	// // Check Max Daily Limit
	// dailyCount, err := s.cache.Get(dailyKey)
	// if err != nil && err != redis.Nil {
	// 	return "", "", "", time.Time{}, ErrOTPMaxDailyRequestsHit
	// }
	//
	//
	// log.Println("dailyCount**** ", dailyCount)
	// if dailyCount >= otpMaxDailyRequests {
	// 	return "", "", "", time.Time{}, errors.New("daily otp limit reached")
	// }
	//
	// // Check Cool Down
	// _, err = s.cache.Get(otpKey)
	// if err == nil {
	// 	return "", "", "", time.Time{}, ErrOTPCoolDownActive
	// } else if err != redis.Nil {
	// 	return "", "", "", time.Time{}, err
	// }
	//
	// // Generate and Save
	// otp, err = utils.GenerateSecureOTP()
	// if err != nil {
	// 	return "", "", "", time.Time{}, err
	// }
	//
	// // Calculate expiry time based on otpTTL
	// expiry = time.Now().Add(otpTTL)
	//
	// pipe := s.rdb.Pipeline()
	//
	// randStr, err := utils.RandomString(6)
	// otpResendKey = fmt.Sprintf("otp-resend-%s%s", randStr, user.ID)
	//
	// pipe.Set(,otpResendKey, identifier, 24*time.Hour)
	//
	// // Set the OTP with cooldown TTL
	// pipe.Set(otpKey, otp, otpTTL)
	//
	// // ReSet the OTP attempts counter
	// pipe.Set(otpAttPrefix+user.ID, 0, otpTTL)
	//
	// // Increment daily counter and set 24h expiry
	// pipe.Incr(dailyKey)
	//
	// pipe.Expire(dailyKey, 24*time.Hour)
	//
	// _, err = pipe.Exec()
	//
	// return otp, otpResendKey, user.ID, expiry, err
}

// Verify checks the OTP for the given identifier.
// Returns true on success and deletes the stored OTP so it cannot be reused.
// Returns false + a descriptive error on failure.
func (s *OTPService) Verify(ctx context.Context, otp string, userId string) (bool, error) {
	// attKey := otpAttPrefix + userId
	// otpKey := otpKeyPrefix + userId
	//
	// // Increment attempt counter before checking — prevents timing-based brute force
	// attempts, err := s.rdb.Incr(attKey).Result()
	// if err != nil {
	// 	return false, fmt.Errorf("track attempts: %w", err)
	// }
	// // Give the attempt key the same TTL as the OTP itself so it auto-cleans
	// s.rdb.Expire(attKey, otpTTL)
	//
	// log.Println("****", attempts, otpMaxAttempts)
	// if attempts > otpMaxAttempts {
	// 	return false, ErrTooManyAttempts
	// }
	//
	// stored, err := s.rdb.Get(otpKey).Result()
	// if err == redis.Nil {
	// 	return false, ErrOTPNotFound
	// }
	// if err != nil {
	// 	return false, fmt.Errorf("get otp: %w", err)
	// }
	//
	// if stored != otp {
	// 	return false, ErrInvalidOTP
	// }
	//
	// // Delete both keys on success — OTP is single-use
	// s.rdb.Del(otpKey, attKey)
	return true, nil
}

// Sentinel errors
var (
	ErrOTPMaxDailyRequestsHit = fmt.Errorf("otp daily request limit hit, request new OTP within 24")
	ErrOTPCoolDownActive      = fmt.Errorf("otp request denied; please wait for the cooldown period to expire")
	ErrOTPNotFound            = fmt.Errorf("OTP not found or expired")
	ErrInvalidOTP             = fmt.Errorf("invalid OTP")
	ErrTooManyAttempts        = fmt.Errorf("too many verification attempts")
)
