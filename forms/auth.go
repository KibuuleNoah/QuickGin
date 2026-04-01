package forms

// Token ...
type Token struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}

// RequestOTPForm is the payload for /auth/request-otp.
// Exactly one of Email or Phone must be provided.
type RequestOTPForm struct {
	Identifier string `form:"identifier" json:"identifier" binding:"required"`
}

// VerifyOTPForm is the payload for /auth/verify-otp.
type AuthWithOTPForm struct {
	UserID string `form:"userId" json:"userId" binding:"required,len=8"`
	OTP    string `form:"otp" json:"otp" binding:"required,len=6"`
}

type AuthWithPasswordForm struct {
	Identifier string `form:"identifier" json:"identifier" binding:"required"`
	Password   string `form:"password" json:"password" binding:"required,min=3,max=50"`
}

var AuthWithPasswordFormMessages = ValidationMessages{
	"Identifier": {
		"required":   "Please enter your Identifier email/phone no.",
		"identifier": "Please enter a valid Identifier email/phone no.",
	},
	"Password": {
		"required": "Please enter your password",
		"min":      "Your password should be between 3 and 50 characters",
		"max":      "Your password should be between 3 and 50 characters",
		"eqfield":  "Your passwords does not match",
	},
}
