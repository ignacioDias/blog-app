package persistence

import (
	"errors"
	"net/mail"
	models "postapi/internal/domain"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func (u *UserRepositoryImpl) Create(p *models.User) error {
	if p.Username == "" {
		return errors.New("username required")
	}
	if len(p.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if !checkValidEmail(p.Email) {
		return errors.New("invalid email format")
	}

	hashedPassword, err := hashPassword(p.Password)
	if err != nil {
		return err
	}
	_, err = u.db.Exec(insertUserSchema, p.Username, p.Email, hashedPassword)
	return err
}

func checkValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (u *UserRepositoryImpl) LoginUser(p *models.User) (*models.User, error) {
	user := &models.User{}
	err := u.db.Get(user, "SELECT * FROM users WHERE username = $1", p.Username)

	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password))
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (u *UserRepositoryImpl) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := u.db.Get(user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
