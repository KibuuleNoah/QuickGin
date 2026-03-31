package services

import (
	"log"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceConfig struct {
	DB        *sqlx.DB
	AuthModel models.AuthModel
}

type AuthService struct {
	cfg AuthServiceConfig
}

func NewAuthService() *AuthService {
	return &AuthService{cfg: AuthServiceConfig{
		DB:        db.GetDB(),
		AuthModel: models.AuthModel{},
	}}
}

// Login ...
func (r *AuthService) AuthWithPassword(form forms.AuthWithPasswordForm) (user models.User, token models.Token, err error) {

	err = r.cfg.DB.Get(&user, "SELECT id, identifier, password, name, updated_at, created_at FROM public.user WHERE identifier=LOWER($1) LIMIT 1", form.Identifier)

	log.Println(err, "(((())))")
	if err != nil {
		return user, token, err
	}

	log.Println(user, "(((())))")
	//Compare the password form and database if match
	bytePassword := []byte(form.Password)
	byteHashedPassword := []byte(user.Password)

	err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)

	if err != nil {
		return user, token, err
	}

	//Generate the JWT auth token
	tokenDetails, err := r.cfg.AuthModel.CreateToken(user.ID)
	if err != nil {
		return user, token, err
	}

	if err = r.cfg.AuthModel.CreateAuth(user.ID, tokenDetails); err != nil {
		return user, token, err
	}

	token.AccessToken = tokenDetails.AccessToken
	token.RefreshToken = tokenDetails.RefreshToken

	return user, token, nil
}
