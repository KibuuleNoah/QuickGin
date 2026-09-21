package forms

// CreateUserForm ...
type CreateUserForm struct {
	Username  string `form:"username" json:"username" binding:"required,min=5,max=20,username"`
	PhoneNo   string `form:"phone_no" json:"phoneNo" binding:"required,identifier"`
	AvatarUrl string `form:"avatarurl" json:"avatarUrl" binding:"required"`
}

type UpdateUserForm struct {
	Username  string `form:"username" json:"username" binding:"required,min=5,max=20,username"`
	PhoneNo   string `form:"phone_no" json:"phoneNo" binding:"required,identifier"`
	AvatarUrl string `form:"avatarurl" json:"avatarUrl"`
}

// RegisterMessages defines validation error messages for register form.
var CreateUserFormMessages = ValidationMessages{
	"Username": {
		"required": "Please enter your name",
		"min":      "Your name should be between 5 to 20 characters",
		"max":      "Your name should be between 5 to 20 characters",
		"username": "Name should not include any special characters or numbers",
	},
	"PhoneNo": {
		"required":   "Please enter your PhoneNo phone no.",
		"identifier": "Please enter a valid PhoneNo phone no.",
	},
}

var UpdateUserFormMessages = ValidationMessages{
	"Username": {
		"required": "Please enter your name",
		"min":      "Your name should be between 3 to 20 characters",
		"max":      "Your name should be between 3 to 20 characters",
		"username": "Name should not include any special characters or numbers",
	},
	"PhoneNo": {
		"required":   "Please enter your PhoneNo email/phone no.",
		"identifier": "Please enter a valid PhoneNo email/phone no.",
	},
}
