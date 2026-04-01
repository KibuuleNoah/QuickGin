package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Massad/gin-boilerplate/db"
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
		AuthModel: *models.NewAuthModel(),
		Asvc:      *services.NewAuthService(),
	}}
}

// AuthWithPassword godoc
// @Summary      Authenticate with password
// @Description  Logs in a user using their identifier and password, returning a JWT access token and refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      forms.AuthWithPasswordForm  true  "Login credentials"
// @Success      200   {object}  map[string]interface{} "Successfully logged in with user and token data"
// @Failure      406   {object}  map[string]string      "Invalid login details or validation error"
// @Router       /auth/login-password [post]
func (ctl AuthController) AuthWithPassword(c *gin.Context) {
	var authForm forms.AuthWithPasswordForm

	if validationErr := c.ShouldBindJSON(&authForm); validationErr != nil {
		message := forms.Translate(validationErr, forms.AuthWithPasswordFormMessages)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	user, token, err := ctl.cfg.Asvc.AuthWithPassword(authForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid login details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "user": user, "token": token})
}

// AuthRequestOtp godoc
// @Summary      Request an OTP
// @Description  Generates and sends a One-Time Password (OTP) to the user's email or phone number.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      forms.RequestOTPForm  true  "User identifier (email or phone)"
// @Success      200   {object}  map[string]interface{} "OTP and OTPID generated successfully"
// @Failure      400   {object}  map[string]string      "Invalid identifier or request body"
// @Failure      500   {object}  map[string]string      "Internal server error generating OTP"
// @Router       /auth/request-otp [post]
func (ctl AuthController) AuthRequestOtp(c *gin.Context) {
	var form forms.RequestOTPForm

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(form.Identifier) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identifier (email/phone) is required"})
		return
	}

	otpSVC := services.NewOTPService(db.GetRedis())
	ctx := c.Request.Context()

	identifier := strings.ToLower(strings.TrimSpace(form.Identifier))
	otp, userId, err := otpSVC.Generate(ctx, identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate OTP: " + err.Error()})
		return
	}

	fmt.Printf("*******OTP %s For User %s", otp, userId)

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to " + form.Identifier, "userId": userId})
}

// AuthWithOTP godoc
// @Summary      Verify OTP and authenticate
// @Description  Validates the OTP. On success, creates the user if they don't exist and returns a JWT access token + refresh token (PocketBase-style).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      forms.AuthWithOTPForm  true  "OTP verification details"
// @Success      200   {object}  map[string]interface{} "Successfully logged in with user and token data"
// @Failure      406   {object}  map[string]string      "Invalid login details or validation error"
// @Router       /auth/verify-otp [post]
func (ctl *AuthController) AuthWithOTP(c *gin.Context) {
	var form forms.AuthWithOTPForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := forms.Translate(validationErr, forms.AuthWithPasswordFormMessages)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	user, token, err := ctl.cfg.Asvc.AuthWithOTP(form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid login details: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "user": user, "token": token})
}

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
		userID := fmt.Sprintf("%.f", claims["user_id"])
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
