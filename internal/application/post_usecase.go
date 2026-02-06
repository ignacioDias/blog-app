package application

import (
	repo "postapi/internal/domain"
)

type PostUseCase struct {
	PostRepo repo.PostRepository
}

func MapPostToJson(p *repo.Post) repo.JsonPost {
	return repo.JsonPost{
		ID:      p.ID,
		Author:  p.Author,
		Content: p.Content,
		Title:   p.Title,
	}
}
