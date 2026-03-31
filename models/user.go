package models

import (
	"errors"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"

	"golang.org/x/crypto/bcrypt"
)

type UserAuthResponse struct {
	User    User   `json:"user"`
	Token   Token  `json:"token"`
	Message string `json:"message"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// User ...
type User struct {
	ID          int64  `db:"id" json:"id"`
	Identifier  string `db:"identifier" json:"identifier"`
	Password    string `db:"password" json:"-"`
	Verified    bool   `db:"verified"      json:"verified"`
	Name        string `db:"name" json:"name"`
	UpdatedAt   int64  `db:"updated_at" json:"-"`
	CreatedAt   int64  `db:"created_at" json:"-"`
	LastLoginAt int64  `db:"last_login_at" json:"last_login_at,omitempty"`
}

// UserModel ...
type UserModel struct{}

var authModel = new(AuthModel)

func HashPassword(password string) ([]byte, error) {

	bytePassword := []byte(password)
	return bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
}

// Create New User ...
func (m UserModel) Create(form forms.CreateUserForm) (user User, err error) {
	getDb := db.GetDB()

	//Check if the user exists in database
	var checkUser int64
	err = getDb.Get(&checkUser, "SELECT count(id) FROM public.user WHERE identifier=LOWER($1) LIMIT 1", form.Identifier)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	if checkUser > 0 {
		return user, errors.New("email/phone no. already exists")
	}

	if len(form.Password) > 0 {
		password, err := HashPassword(form.Password)
		if err != nil {
			return user, errors.New("something went wrong, please try again later")
		}
		user.Password = string(password)
	}

	//Create the user and return back the user ID
	err = getDb.QueryRow("INSERT INTO public.user(identifier, password, name) VALUES($1, $2, $3) RETURNING id", form.Identifier, user.Password, form.Name).Scan(&user.ID)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	user.Password = ""
	user.Name = form.Name
	user.Identifier = form.Identifier

	return user, err
}

// One ...
func (m UserModel) One(userID int64) (user User, err error) {
	err = db.GetDB().Get(&user, "SELECT id, email, name FROM public.user WHERE id=$1 LIMIT 1", userID)
	return user, err
}
