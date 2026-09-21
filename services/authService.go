package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"QuickGin/db"
	"QuickGin/forms"

	"github.com/jmoiron/sqlx"
	"QuickGin/models"
	"QuickGin/models/cache"

	jwt "github.com/golang-jwt/jwt/v4"
)

type AuthTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	AtExpires    int64  `json:"atExpires"`
	RtExpires    int64  `json:"rtExpires"`
}

// TokenDetails ...
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// AccessDetails ...
type AccessDetails struct {
	AccessUUID string
	UserID     string
}

type AuthServiceConfig struct {
	DB    *sqlx.DB
	cache cache.Cache
}

type AuthService struct {
	cfg AuthServiceConfig
}

func NewAuthService() *AuthService {
	return &AuthService{cfg: AuthServiceConfig{
		DB:    db.AppDB(),
		cache: db.AppCache(),
	}}
}

// Login With Password...
// func (s *AuthService) AuthWithPassword(form forms.AuthWithPasswordForm) (user models.User, token AuthTokenResponse, err error) {
//
// 	err = s.cfg.DB.Get(&user, "SELECT id, identifier, password, name, updated_at, created_at FROM public.user WHERE identifier=LOWER($1) LIMIT 1", form.Identifier)
// 	if err != nil {
// 		return user, token, err
// 	}
//
// 	//Compare the password form and database if match
// 	err = utils.CompareHashAndPassword(user.Password, form.Password)
// 	if err != nil {
// 		return user, token, err
// 	}
//
// 	//Generate the JWT auth token
// 	tokenDetails, err := s.createToken(user.ID)
// 	if err != nil {
// 		return user, token, err
// 	}
//
// 	if err = s.createAuth(user.ID, tokenDetails); err != nil {
// 		return user, token, err
// 	}
//
// 	token.AccessToken = tokenDetails.AccessToken
// 	token.RefreshToken = tokenDetails.RefreshToken
// 	token.AtExpires = tokenDetails.AtExpires
// 	token.RtExpires = tokenDetails.RtExpires
//
// 	return user, token, nil
// }

func (s *AuthService) AuthWithOTP(form forms.AuthWithOTPForm) (user models.User, token AuthTokenResponse, err error) {
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
	tokenDetails, err := s.createToken(user.ID)
	if err != nil {
		return user, token, err
	}

	if err = s.createAuth(user.ID, tokenDetails); err != nil {
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

// ExtractTokenMetadata ...
func (s *AuthService) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := s.verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	accessUUID, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, errors.New("invalid access_uuid claim")
	}

	userIDRaw, ok := claims["user_id"]
	if !ok {
		return nil, errors.New("user_id not found in claims")
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		return nil, errors.New("user_id is not a string")
	}

	return &AccessDetails{
		AccessUUID: accessUUID,
		UserID:     userID,
	}, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (AuthTokenResponse, error) {

	log.Println("***", refreshToken)
	// Parse and Verify Token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		log.Println(err)
		return AuthTokenResponse{}, fmt.Errorf("invalid or expired token")
	}

	// Extract Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return AuthTokenResponse{}, fmt.Errorf("invalid claims")
	}

	refreshUUID, okUUID := claims["refresh_uuid"].(string)
	userID, okUser := claims["user_id"].(string) // Combined string assertion and check

	if !okUUID || !okUser {
		return AuthTokenResponse{}, fmt.Errorf("missing token metadata")
	}

	// Revoke Old Token
	delErr := s.DeleteAuth(refreshUUID)
	if delErr != nil {
		return AuthTokenResponse{}, fmt.Errorf("token already revoked or expired")
	}

	// Generate New Token Pair
	ts, createErr := s.createToken(userID)
	if createErr != nil {
		return AuthTokenResponse{}, createErr
	}

	// Save New Metadata
	if saveErr := s.createAuth(userID, ts); saveErr != nil {
		return AuthTokenResponse{}, saveErr
	}

	return AuthTokenResponse{
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
		AtExpires:    ts.AtExpires,
		RtExpires:    ts.RtExpires,
	}, nil
}

// FetchAuth ...
func (s *AuthService) FetchAuth(authD *AccessDetails) (string, error) {
	val, ok := s.cfg.cache.Get(authD.AccessUUID)
	if !ok {
		return "", errors.New("Key Not Found")
	}

	str, err := val.String()
	if err != nil {
		return "", errors.New("Failed to convert to string")
	}

	return str, nil
}

// DeleteAuth ...
func (s *AuthService) DeleteAuth(givenUUID string) error {
	return s.cfg.cache.Delete(givenUUID)
}

// CreateAuth ...
func (s *AuthService) createAuth(userid string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := s.cfg.cache.Set(td.AccessUUID, userid, at.Sub(now))
	if errAccess != nil {
		return errAccess
	}
	errRefresh := s.cfg.cache.Set(td.RefreshUUID, userid, rt.Sub(now))
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

/**
*** PRIVATE METHODS
**/

// createToken ...
func (s *AuthService) createToken(userID string) (*TokenDetails, error) {

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.New().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// extractToken ...
func (s *AuthService) extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// verifyToken ...
func (s *AuthService) verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := s.extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
