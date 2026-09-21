package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"QuickGin/db"
	"QuickGin/forms"
	"QuickGin/models"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("phone number already exists")
	ErrInternal     = errors.New("something went wrong, please try again later")
)

type UserServiceConfig struct {
	DB *sqlx.DB
}

type UserService struct {
	cfg UserServiceConfig
}

func NewUserService() *UserService {
	return &UserService{cfg: UserServiceConfig{
		DB: db.AppDB(),
	}}
}

func (s *UserService) GetProfile(id string) (*models.User, error) {
	user := &models.User{}

	err := s.cfg.DB.Get(user, "SELECT id, username, phone_no, avatar_url, created_at FROM users WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		log.Printf("GetProfile(%s): %v", id, err)
		return nil, fmt.Errorf("%w: fetch profile", ErrInternal)
	}

	return user, nil
}

func (s *UserService) Create(form forms.CreateUserForm) (*models.User, error) {
	var count int64
	err := s.cfg.DB.Get(&count, "SELECT count(id) FROM users WHERE phone_no = LOWER($1) LIMIT 1", form.PhoneNo)
	if err != nil {
		log.Printf("Create: check existing user: %v", err)
		return nil, ErrInternal
	}
	if count > 0 {
		return nil, ErrUserExists
	}

	user := &models.User{
		Username: form.Username,
		PhoneNo:  form.PhoneNo,
	}

	err = s.cfg.DB.QueryRow(
		"INSERT INTO users(phone_no, username) VALUES($1, $2) RETURNING id",
		form.PhoneNo, form.Username,
	).Scan(&user.ID)
	if err != nil {
		log.Printf("Create: insert user: %v", err)
		return nil, ErrInternal
	}

	return user, nil
}

func (s *UserService) UpdateUser(id string, form forms.UpdateUserForm) (*models.User, error) {
	query := `
		UPDATE users
		SET
			username = $1,
			phone_no = $2,
			avatar_url = $3
		WHERE id = $4
		RETURNING id, username, phone_no, avatar_url
	`

	user := &models.User{}
	err := s.cfg.DB.QueryRowx(
		query,
		form.Username,
		form.PhoneNo,
		form.AvatarUrl,
		id,
	).StructScan(user)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		log.Printf("UpdateUser(%s): %v", id, err)
		return nil, fmt.Errorf("%w: update user", ErrInternal)
	}

	return user, nil
}

func (s *UserService) DeleteUser(id string) error {
	result, err := s.cfg.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Printf("DeleteUser(%s): %v", id, err)
		return fmt.Errorf("%w: delete user", ErrInternal)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("DeleteUser(%s): rows affected: %v", id, err)
		return fmt.Errorf("%w: delete user", ErrInternal)
	}
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}
