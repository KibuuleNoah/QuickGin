package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/KibuuleNoah/QuickGin/db"
	"github.com/KibuuleNoah/QuickGin/forms"
	"github.com/KibuuleNoah/QuickGin/models"
	"github.com/KibuuleNoah/QuickGin/services"
	"github.com/gin-gonic/gin"
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
// @Router       /auth/with-password [post]
func (ctl *AuthController) AuthWithPassword(c *gin.Context) {
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
func (ctl *AuthController) AuthRequestOtp(c *gin.Context) {
	var form forms.RequestOTPForm

	var json map[string]interface{}
	if err := c.BindJSON(&json); err != nil {
		// handle error
	}

	if json["Identifier"] != nil {
		key, ok := json["Identifier"].(string)
		if ok && strings.HasPrefix(key, "otp-resend-") {
			identifier, err := ctl.cfg.Asvc.QueryOtpResendKeyOwner(key)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid or Expired otp resend key"})
				return
			}

			form.Identifier = identifier
		}
	}

	if form.Identifier == "" {
		if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
			message := forms.Translate(validationErr, forms.RequestOTPFormMessages)
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
			return
		}
	}

	log.Println(form, form.Identifier)

	otpSVC := services.NewOTPService(db.GetRedis())
	ctx := c.Request.Context()

	identifier := strings.ToLower(strings.TrimSpace(form.Identifier))
	otp, otpResendKey, userId, expiry, err := otpSVC.Generate(ctx, identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate OTP: " + err.Error()})
		return
	}

	fmt.Printf("*******OTP %s For User %s", otp, userId)

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully ", "userId": userId, "otpResendKey": otpResendKey, "expiry": expiry})
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
// @Router       /auth/with-otp [post]
func (ctl *AuthController) AuthWithOTP(c *gin.Context) {
	var form forms.AuthWithOTPForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := forms.Translate(validationErr, forms.AuthWithOTPFormMessages)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	user, token, err := ctl.cfg.Asvc.AuthWithOTP(form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid login details: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
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

func (ctl *AuthController) RefreshToken(c *gin.Context) {
	var form forms.RefreshTokenForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := forms.Translate(validationErr, forms.RefreshTokenMessages)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	token, err := ctl.cfg.AuthModel.RefreshTokens(form.Token)
	if err != nil {
		log.Println("****", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Please login again"})
		return
	}

	c.JSON(http.StatusOK, token)
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
func (ctl *AuthController) AuthLogout(c *gin.Context) {
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
