package forms

type AuthWithPasswordForm struct {
	Identifier string `form:"identifier" json:"identifier" binding:"required"`
	Password   string `form:"password" json:"password" binding:"required,min=3,max=50"`
}

// CreateUserForm ...
type CreateUserForm struct {
	Name       string `form:"name" json:"name" binding:"required,min=3,max=20,fullName"`
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

// RegisterMessages defines validation error messages for register form.
var CreateUserFormMessages = ValidationMessages{
	"Name": {
		"required": "Please enter your name",
		"min":      "Your name should be between 3 to 20 characters",
		"max":      "Your name should be between 3 to 20 characters",
		"fullName": "Name should not include any special characters or numbers",
	},
	"Identifier": {
		"required":   "Please enter your Identifier email/phone no.",
		"identifier": "Please enter a valid Identifier email/phone no.",
	},
	"Password": {
		"required": "Please enter your password",
		"min":      "Your password should be between 3 and 50 characters",
		"max":      "Your password should be between 3 and 50 characters",
	},
}
