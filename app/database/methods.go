package database

import (
	"database/sql"
	"errors"
	"net/mail"
	"postapi/app/models"

	"golang.org/x/crypto/bcrypt"
)

func (d *DB) CreatePost(p *models.Post) error {
	if p.Title == "" || p.Content == "" {
		return errors.New("Invalid Title / content")
	}
	err := d.db.QueryRow(insertPostSchema, p.Title, p.Content, p.Author).Scan(&p.ID)
	return err
}

func (d *DB) GetPostsByUser(author string) ([]*models.Post, error) {
	var posts []*models.Post
	err := d.db.Select(&posts, "SELECT * FROM posts WHERE author = $1", author)

	return posts, err
}

func (d *DB) RegisterUser(p *models.User) error {
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
	_, err = d.db.Exec(insertUserSchema, p.Username, p.Email, hashedPassword)
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

func (d *DB) LoginUser(p *models.User) (*models.User, error) {
	user := &models.User{}
	err := d.db.Get(user, "SELECT * FROM users WHERE username = $1", p.Username)

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

func (d *DB) GetPost(id int64) (*models.Post, error) {
	post := &models.Post{}
	err := d.db.Get(post, "SELECT * FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (d *DB) DeletePost(id int64, username string) error {
	result, err := d.db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (d *DB) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := d.db.Get(user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
