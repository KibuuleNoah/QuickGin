package forms

// Token ...
type Token struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}

// RequestOTPForm is the payload for /auth/request-otp.
// Exactly one of Email or Phone must be provided.
type RequestOTPForm struct {
	Identifier string `form:"identifier" json:"identifier" binding:"required, identifier"`
}

// VerifyOTPForm is the payload for /auth/verify-otp.
type VerifyOTPForm struct {
	OTPID *string `json:"otpId" binding:"-"`
	OTP   string  `json:"otp"   binding:"required,len=6"`
}
