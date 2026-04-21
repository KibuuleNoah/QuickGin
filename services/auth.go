package services

import (
	"context"
	"errors"

	"github.com/KibuuleNoah/QuickGin/db"
	"github.com/KibuuleNoah/QuickGin/forms"

	"github.com/KibuuleNoah/QuickGin/models"
	"github.com/KibuuleNoah/QuickGin/models/cache"
	"github.com/KibuuleNoah/QuickGin/utils"
	"github.com/jmoiron/sqlx"
)

type AuthServiceConfig struct {
	DB        *sqlx.DB
	cache     cache.Cache
	AuthModel models.AuthModel
}

type AuthService struct {
	cfg AuthServiceConfig
}

func NewAuthService() *AuthService {
	return &AuthService{cfg: AuthServiceConfig{
		DB:        db.AppDB(),
		cache:     db.AppCache(),
		AuthModel: *models.NewAuthModel(),
	}}
}

// Login With Password...
func (s *AuthService) AuthWithPassword(form forms.AuthWithPasswordForm) (user models.User, token models.AuthTokenResponse, err error) {

	err = s.cfg.DB.Get(&user, "SELECT id, identifier, password, name, updated_at, created_at FROM public.user WHERE identifier=LOWER($1) LIMIT 1", form.Identifier)
	if err != nil {
		return user, token, err
	}

	//Compare the password form and database if match
	err = utils.CompareHashAndPassword(user.Password, form.Password)
	if err != nil {
		return user, token, err
	}

	//Generate the JWT auth token
	tokenDetails, err := s.cfg.AuthModel.CreateToken(user.ID)
	if err != nil {
		return user, token, err
	}

	if err = s.cfg.AuthModel.CreateAuth(user.ID, tokenDetails); err != nil {
		return user, token, err
	}

	token.AccessToken = tokenDetails.AccessToken
	token.RefreshToken = tokenDetails.RefreshToken
	token.AtExpires = tokenDetails.AtExpires
	token.RtExpires = tokenDetails.RtExpires

	return user, token, nil
}

func (s *AuthService) AuthWithOTP(form forms.AuthWithOTPForm) (user models.User, token models.AuthTokenResponse, err error) {
	otpSVC := NewOTPService()
	ctx := context.Background()

	ok, err := otpSVC.Verify(ctx, form.OTP, form.UserID)
	if err != nil {
		return user, token, err
	}

	if !ok {
		return user, token, errors.New("Invalid Otp")
	}

	err = s.cfg.DB.Get(&user, "SELECT id, identifier, name, updated_at, created_at FROM public.user WHERE id=$1 LIMIT 1", form.UserID)
	if err != nil {
		return user, token, err
	}

	//Generate the JWT auth token
	tokenDetails, err := s.cfg.AuthModel.CreateToken(user.ID)
	if err != nil {
		return user, token, err
	}

	if err = s.cfg.AuthModel.CreateAuth(user.ID, tokenDetails); err != nil {
		return user, token, err
	}

	token.AccessToken = tokenDetails.AccessToken
	token.RefreshToken = tokenDetails.RefreshToken
	token.AtExpires = tokenDetails.AtExpires
	token.RtExpires = tokenDetails.RtExpires

	return user, token, err
}

func (s *AuthService) QueryOtpResendKeyOwner(otpResendKey string) (string, error) {
	val, ok := s.cfg.cache.Get(otpResendKey)
	if !ok {
		return "", errors.New("Key Not Found")
	}

	str, err := val.String()
	if err != nil {
		return "", errors.New("Failed to convert to string")
	}

	return str, nil
}
