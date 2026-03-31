package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/Massad/gin-boilerplate/services"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

type AuthControllerConfig struct {
	AuthModel models.AuthModel
	Asvc      services.AuthService // Auth service
}

// AuthController ...
type AuthController struct {
	cfg AuthControllerConfig
}

func NewAuthController() *AuthController {
	return &AuthController{cfg: AuthControllerConfig{
		AuthModel: models.AuthModel{},
		Asvc:      *services.NewAuthService(),
	}}
}

func (ctl AuthController) AuthWithPassword(c *gin.Context) {
	var authForm forms.AuthWithPasswordForm

	if validationErr := c.ShouldBindJSON(&authForm); validationErr != nil {
		message := forms.Translate(validationErr, forms.AuthWithPasswordFormMessages)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	fmt.Println(authForm)

	user, token, err := ctl.cfg.Asvc.AuthWithPassword(authForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid login details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "user": user, "token": token})
}

func (ctl AuthController) AuthWithOTP(c *gin.Context) {}

// Refresh Token godoc
// @Summary Refresh Token example
// @Schemes
// @Description Refresh Token example
// @Tags Auth
// @Accept json
// @Produce json
// @Param auth body forms.Token true "Auth"
// @Success 	 200  {object}  models.AuthResponse
// @Failure      406  {object}  models.MessageResponse
// @Router /auth/token/refresh [POST]
func (ctl AuthController) Refresh(c *gin.Context) {
	var tokenForm forms.Token

	if c.ShouldBindJSON(&tokenForm) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid form", "form": tokenForm})
		c.Abort()
		return
	}

	//verify the token
	token, err := jwt.Parse(tokenForm.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
			return
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := ctl.cfg.AuthModel.DeleteAuth(refreshUUID)
		if delErr != nil || deleted == 0 { //if any goes wrong
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
			return
		}

		//Create new pairs of refresh and access tokens
		ts, createErr := ctl.cfg.AuthModel.CreateToken(userID)
		if createErr != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid authorization, please login again"})
			return
		}
		//save the tokens metadata to redis
		saveErr := ctl.cfg.AuthModel.CreateAuth(userID, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid authorization, please login again"})
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusOK, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
	}
}

// Logout User godoc
// @Summary Logout User example
// @Schemes
// @Description Logout User example
// @Tags User
// @Accept json
// @Produce json
// @Success 	 200  {object}  models.MessageResponse
// @Failure      406  {object}  models.MessageResponse
// @Router /auth/logout [GET]
func (ctl AuthController) AuthLogout(c *gin.Context) {
	au, err := ctl.cfg.AuthModel.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not logged in"})
		return
	}

	deleted, delErr := ctl.cfg.AuthModel.DeleteAuth(au.AccessUUID)
	if delErr != nil || deleted == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// RequestOTP godoc
// @Summary      Request a one-time password
// @Description  Sends a 6-digit OTP to the provided email or phone number.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      forms.RequestOTPForm  true  "email or phone"
// @Success      200   {object}  map[string]string
// @Router       /auth/request-otp [post]
// func (ctl AuthController) RequestOTP(c *gin.Context) {
// 	var form forms.RequestOTPForm
// 	if err := c.ShouldBindJSON(&form); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	if form.Email == nil && form.Phone == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone is required"})
// 		return
// 	}
// 	if form.Email != nil && form.Phone != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "provide either email or phone, not both"})
// 		return
// 	}
//
// 	ctx := c.Request.Context()
//
// 	if form.Email != nil {
// 		identifier := strings.ToLower(strings.TrimSpace(*form.Email))
// 		otp, err := a.otp.Generate(ctx, identifier)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate OTP"})
// 			return
// 		}
// 		if err = a.mail.SendOTP(identifier, otp); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not send OTP email"})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email"})
// 		return
// 	}
//
// 	// Phone path
// 	identifier := strings.TrimSpace(*form.Phone)
// 	otp, err := a.otp.Generate(ctx, identifier)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate OTP"})
// 		return
// 	}
// 	if err = a.sms.SendOTP(identifier, otp); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not send OTP SMS"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to phone"})
// }

// VerifyOTP godoc
// @Summary      Verify OTP and authenticate
// @Description  Validates the OTP. On success, creates the user if they don't exist
//
//	and returns a JWT access token + refresh token (PocketBase-style).
//
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      forms.VerifyOTPForm  true  "identifier + otp"
// @Success      200   {object}  map[string]interface{}
// @Router       /auth/verify-otp [post]

// func (a *AuthController) VerifyOTP(c *gin.Context) {
// 	var form forms.VerifyOTPForm
// 	if err := c.ShouldBindJSON(&form); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	if form.Email == nil && form.Phone == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone is required"})
// 		return
// 	}
//
// 	ctx := c.Request.Context()
// 	var identifier string
// 	isEmail := form.Email != nil
//
// 	if isEmail {
// 		identifier = strings.ToLower(strings.TrimSpace(*form.Email))
// 	} else {
// 		identifier = strings.TrimSpace(*form.Phone)
// 	}
//
// 	ok, err := a.otp.Verify(ctx, identifier, form.OTP)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, services.ErrTooManyAttempts):
// 			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many attempts, request a new OTP"})
// 		case errors.Is(err, services.ErrOTPNotFound):
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired or not found"})
// 		case errors.Is(err, services.ErrInvalidOTP):
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid OTP"})
// 		default:
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "verification failed"})
// 		}
// 		return
// 	}
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid OTP"})
// 		return
// 	}
//
// 	// Upsert user — same as PocketBase: first verify = auto-register
// 	var user interface{ GetID() string }
// 	if isEmail {
// 		u, err := a.userRepo.UpsertByEmail(ctx, identifier)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
// 			return
// 		}
// 		accessToken, err := a.tokens.IssueAccessToken(u.ID)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue token"})
// 			return
// 		}
// 		refreshToken, err := a.tokens.IssueRefreshToken(ctx, u.ID)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue refresh token"})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{
// 			"token":         accessToken,
// 			"refresh_token": refreshToken,
// 			"record":        u,
// 		})
// 		return
// 	}
//
// 	// Phone path
// 	_ = user
// 	u, err := a.userRepo.UpsertByPhone(ctx, identifier)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
// 		return
// 	}
// 	accessToken, err := a.tokens.IssueAccessToken(u.ID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue token"})
// 		return
// 	}
// 	refreshToken, err := a.tokens.IssueRefreshToken(ctx, u.ID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue refresh token"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"token":         accessToken,
// 		"refresh_token": refreshToken,
// 		"record":        u,
// 	})
// }
