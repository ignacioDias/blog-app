package persistence

import (
	"database/sql"
	"errors"
	models "postapi/internal/domain"

	"github.com/jmoiron/sqlx"
)

type PostRepositoryImpl struct {
	db *sqlx.DB
}

func (p *PostRepositoryImpl) Create(post *models.Post) error {
	if post.Title == "" || post.Content == "" {
		return errors.New("Invalid Title / content")
	}
	err := p.db.QueryRow(insertPostSchema, post.Title, post.Content, post.Author).Scan(&post.ID)
	return err
}

func (p *PostRepositoryImpl) Update(post *models.Post) error {
	result, err := p.db.Exec(
		`UPDATE posts
		 SET title = $1, content = $2
		 WHERE id = $3 AND author = $4`,
		post.Title,
		post.Content,
		post.ID,
		post.Author,
	)
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

func (p *PostRepositoryImpl) Delete(id int64, author string) error {
	result, err := p.db.Exec("DELETE FROM posts WHERE id = $1 AND author = $2", id, author)
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

func (p *PostRepositoryImpl) FindByID(id int64) (*models.Post, error) {
	post := &models.Post{}
	err := p.db.Get(post, "SELECT * FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *PostRepositoryImpl) FindByAuthor(author string) ([]*models.Post, error) {
	var posts []*models.Post
	err := p.db.Select(&posts, "SELECT * FROM posts WHERE author = $1", author)

	return posts, err
}
