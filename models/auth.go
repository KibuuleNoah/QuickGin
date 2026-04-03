package models

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KibuuleNoah/QuickGin/db"
	redis "github.com/go-redis/redis/v7"
	jwt "github.com/golang-jwt/jwt/v4"
	uuid "github.com/google/uuid"
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

// AuthModel ...
type AuthModel struct {
	redisDB *redis.Client
}

func NewAuthModel() *AuthModel {
	return &AuthModel{
		redisDB: db.GetRedis(),
	}
}

// CreateToken ...
func (m *AuthModel) CreateToken(userID string) (*TokenDetails, error) {

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

// CreateAuth ...
func (m *AuthModel) CreateAuth(userid string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := m.redisDB.Set(td.AccessUUID, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := m.redisDB.Set(td.RefreshUUID, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

// ExtractToken ...
func (m *AuthModel) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// VerifyToken ...
func (m *AuthModel) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := m.ExtractToken(r)
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

// ExtractTokenMetadata ...
func (m *AuthModel) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := m.VerifyToken(r)
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

func (m *AuthModel) RefreshTokens(refreshToken string) (AuthTokenResponse, error) {

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
	deleted, delErr := m.DeleteAuth(refreshUUID)
	if delErr != nil || deleted == 0 {
		return AuthTokenResponse{}, fmt.Errorf("token already revoked or expired")
	}

	// Generate New Token Pair
	ts, createErr := m.CreateToken(userID)
	if createErr != nil {
		return AuthTokenResponse{}, createErr
	}

	// Save New Metadata
	if saveErr := m.CreateAuth(userID, ts); saveErr != nil {
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
func (m *AuthModel) FetchAuth(authD *AccessDetails) (string, error) {
	log.Println("*****", authD, m.redisDB)
	return m.redisDB.Get(authD.AccessUUID).Result()
}

// DeleteAuth ...
func (m *AuthModel) DeleteAuth(givenUUID string) (int64, error) {
	return m.redisDB.Del(givenUUID).Result()
}
