package models

import (
	"errors"
	"github.com/KibuuleNoah/QuickGin/db"
	"github.com/KibuuleNoah/QuickGin/forms"
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
	ID          string `db:"id" json:"id"`
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

// Fetch One User of the currect struct by id
// same as  err := DB.Get(&user, "SELECT id, identifier, name, updated_at, created_at FROM public.user WHERE id=$1 LIMIT 1", form.UserID)
// func (u *User) () (err error) {
// 	err = db.GetDB().Get(u, "SELECT id, identifier, name, updated_at, created_at FROM public.user WHERE id=$1 LIMIT 1", u.ID)
// d	return err
// }
