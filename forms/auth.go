package forms

// Token ...
type RefreshTokenForm struct {
	Token string `form:"token" json:"token" binding:"required"`
}

type RequestOTPForm struct {
	Identifier string `form:"identifier" json:"identifier" binding:"required,identifier"`
}

type AuthWithOTPForm struct {
	UserID string `form:"userId" json:"userId" binding:"required,len=8"`
	OTP    string `form:"otp" json:"otp" binding:"required,len=6"`
}

type AuthWithPasswordForm struct {
	Identifier string `form:"identifier" json:"identifier" binding:"required,identifier"`
	Password   string `form:"password" json:"password" binding:"required,strong_password,min=8,max=32"`
}

var AuthWithPasswordFormMessages = ValidationMessages{
	"Identifier": {
		"required":   "Please enter your Identifier email/phone no.",
		"identifier": "Please enter a valid Identifier email/phone no.",
	},
	"Password": {
		"required":        "Please enter your password",
		"min":             "Your password should be between 8 and 32 characters",
		"max":             "Your password should be between 8 and 32 characters",
		"eqfield":         "Your passwords does not match",
		"strong_password": "Password must be 8–32 characters and include at least one uppercase letter, one lowercase letter, one number, and one special character (!@#$ etc.). Avoid sequences like 'abc' or '123', repeated characters like 'aaa', and commonly used passwords.",
	},
}

var RefreshTokenMessages = ValidationMessages{
	"Token": {
		"required": "Refresh token is required to renew your session.",
	},
}

var RequestOTPFormMessages = ValidationMessages{
	"Identifier": {
		"required":   "Please enter your email or phone number to receive an OTP.",
		"identifier": "Please enter a valid Identifier email/phone no.",
	},
}

var AuthWithOTPFormMessages = ValidationMessages{
	"UserID": {
		"required": "User ID is required.",
		"len":      "User ID must be exactly 8 characters long.",
	},
	"OTP": {
		"required": "Please enter the OTP sent to your device.",
		"len":      "The OTP must be exactly 6 digits.",
	},
}
