package services

import (
	"context"
	"errors"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceConfig struct {
	DB        *sqlx.DB
	rDB       *redis.Client
	AuthModel models.AuthModel
}

type AuthService struct {
	cfg AuthServiceConfig
}

func NewAuthService() *AuthService {
	return &AuthService{cfg: AuthServiceConfig{
		DB:        db.GetDB(),
		rDB:       db.GetRedis(),
		AuthModel: *models.NewAuthModel(),
	}}
}

// Login With Password...
func (s *AuthService) AuthWithPassword(form forms.AuthWithPasswordForm) (user models.User, token models.Token, err error) {

	err = s.cfg.DB.Get(&user, "SELECT id, identifier, password, name, updated_at, created_at FROM public.user WHERE identifier=LOWER($1) LIMIT 1", form.Identifier)

	if err != nil {
		return user, token, err
	}

	//Compare the password form and database if match
	bytePassword := []byte(form.Password)
	byteHashedPassword := []byte(user.Password)

	err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)

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

	return user, token, nil
}

func (s *AuthService) AuthWithOTP(form forms.AuthWithOTPForm) (user models.User, token models.Token, err error) {
	otpSVC := NewOTPService(s.cfg.rDB)
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

	return user, token, err
}
