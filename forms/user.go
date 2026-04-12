package forms

// CreateUserForm ...
type CreateUserForm struct {
	Name       string `form:"name" json:"name" binding:"required,min=3,max=20,username"`
	Identifier string `form:"identifier" json:"identifier" binding:"required,identifier"`
	Password   string `form:"password" json:"password" binding:"strong_password"`
}

// RegisterMessages defines validation error messages for register form.
var CreateUserFormMessages = ValidationMessages{
	"Name": {
		"required": "Please enter your name",
		"min":      "Your name should be between 3 to 20 characters",
		"max":      "Your name should be between 3 to 20 characters",
		"username": "Name should not include any special characters or numbers",
	},
	"Identifier": {
		"required":   "Please enter your Identifier email/phone no.",
		"identifier": "Please enter a valid Identifier email/phone no.",
	},
	"Password": {
		"min":             "Your password should be between 8 and 32 characters",
		"max":             "Your password should be between 8 and 32 characters",
		"strong_password": "Password must be 8–32 characters and include at least one uppercase letter, one lowercase letter, one number, and one special character (!@#$ etc.). Avoid sequences like 'abc' or '123', repeated characters like 'aaa', and commonly used passwords.",
	},
}
