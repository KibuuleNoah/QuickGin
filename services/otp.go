package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/KibuuleNoah/QuickGin/db"
	"github.com/KibuuleNoah/QuickGin/models"
	"github.com/KibuuleNoah/QuickGin/models/cache"
	"github.com/KibuuleNoah/QuickGin/utils"
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
	var user models.User
	err = s.DB.Get(&user, "SELECT id, identifier, password, name, updated_at, created_at FROM public.user WHERE identifier=LOWER($1) LIMIT 1", identifier)
	if err != nil {
		return "", "", "", time.Time{}, errors.New("user does not exist")
	}

	otpKey := otpKeyPrefix + user.ID
	dailyKey := otpMaxDailyRequestsPrefix + user.ID

	// Check Max Daily Limit
	val, ok := s.cache.Get(dailyKey)
	dailyCount, _ := val.(int) // Cast to int

	if ok && dailyCount >= otpMaxDailyRequests {
		return "", "", "", time.Time{}, errors.New("daily otp limit reached")
	}

	// Check Cool Down
	if _, exists := s.cache.Get(otpKey); exists {
		return "", "", "", time.Time{}, errors.New("otp cooldown active")
	}

	//Generate OTP
	otp, err = utils.GenerateSecureOTP()
	if err != nil {
		return "", "", "", time.Time{}, err
	}

	expiry = time.Now().Add(otpTTL)
	randStr, _ := utils.RandomString(6)
	otpResendKey = fmt.Sprintf("otp-resend-%s%s", randStr, user.ID)

	//SAVE TO CACHE

	// Set Resend Key
	s.cache.Set(otpResendKey, identifier, 24*time.Hour)

	// Set OTP with cooldown
	s.cache.Set(otpKey, otp, otpTTL)

	// Reset attempts
	s.cache.Set(otpAttPrefix+user.ID, 0, otpTTL)

	// Update Daily Counter
	s.cache.Set(dailyKey, dailyCount+1, 24*time.Hour)

	return otp, otpResendKey, user.ID, expiry, nil
}

// Verify checks the OTP for the given identifier.
// Returns true on success and deletes the stored OTP so it cannot be reused.
// Returns false + a descriptive error on failure.

func (s *OTPService) Verify(ctx context.Context, otp string, userId string) (bool, error) {
	attKey := otpAttPrefix + userId
	otpKey := otpKeyPrefix + userId

	// Handle Attempt Counter
	val, ok := s.cache.Get(attKey)
	attempts, _ := val.(int)
	attempts++

	// Update attempts in cache with the same TTL as OTP
	s.cache.Set(attKey, attempts, otpTTL)

	log.Println("****", attempts, otpMaxAttempts)
	if attempts > otpMaxAttempts {
		return false, ErrTooManyAttempts
	}

	// Check if OTP exists
	storedVal, ok := s.cache.Get(otpKey)
	if !ok {
		return false, ErrOTPNotFound
	}

	stored, ok := storedVal.(string)
	if !ok || stored != otp {
		return false, ErrInvalidOTP
	}

	// On success Clean up single-use keys
	s.cache.Delete(otpKey)
	s.cache.Delete(attKey)

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
