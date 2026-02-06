package handlers

import (
	"fmt"
	"log"
	"net/http"
	"postapi/internal/application"
	models "postapi/internal/domain"
	"postapi/internal/middleware"
	"strconv"

	"github.com/gorilla/mux"
)

type PostHandler struct {
	PostUseCase application.PostUseCase
}

func (p *PostHandler) CreatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(middleware.UsernameKey).(string)
		if !ok {
			middleware.SendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		req := models.PostRequest{}
		err := middleware.Parse(w, r, &req)
		if err != nil {
			log.Printf("Cannot middleware.Parse body. err = %v \n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}

		post := &models.Post{
			ID:      0,
			Title:   req.Title,
			Author:  username,
			Content: req.Content,
		}

		err = p.PostUseCase.PostRepo.Create(post)
		if err != nil {
			log.Printf("Cannot save post in DB. err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to create post"}, http.StatusInternalServerError)
			return
		}

		resp := application.MapPostToJson(post)
		middleware.SendResponse(w, r, resp, http.StatusOK)
	}
}
func (p *PostHandler) UpdatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(middleware.UsernameKey).(string)
		if !ok {
			middleware.SendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		req := models.PostRequest{}
		err := middleware.Parse(w, r, &req)

		if err != nil {
			log.Printf("Cannot middleware.Parse body. err = %v \n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
			return
		}
		vars := mux.Vars(r)
		id := vars["post_id"]
		idAsNumber, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			middleware.SendResponse(w, r, map[string]string{"error": fmt.Sprintf("Invalid ID %s", id)}, http.StatusBadRequest)
			return
		}

		content := req.Content
		title := req.Title

		oldPost, err := p.PostUseCase.PostRepo.FindByID(idAsNumber)
		if err != nil {
			log.Printf("Cannot get post, err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to get post"}, http.StatusInternalServerError)
			return
		}
		if req.Content == "" {
			content = oldPost.Content
		}
		if req.Title == "" {
			title = oldPost.Title
		}
		post := &models.Post{
			ID:      idAsNumber,
			Title:   title,
			Author:  username,
			Content: content,
		}

		err = p.PostUseCase.PostRepo.Update(post)
		if err != nil {
			log.Printf("Cannot save post in postRepo. err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to update post"}, http.StatusInternalServerError)
			return
		}

		resp := application.MapPostToJson(post)
		middleware.SendResponse(w, r, resp, http.StatusOK)
	}
}
func (p *PostHandler) DeletePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(middleware.UsernameKey).(string)
		if !ok {
			middleware.SendResponse(w, r, map[string]string{"error": "Unauthorized"}, http.StatusUnauthorized)
			return
		}
		vars := mux.Vars(r)
		id := vars["post_id"]
		idAsNumber, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			middleware.SendResponse(w, r, map[string]string{"error": fmt.Sprintf("Invalid ID %s", id)}, http.StatusBadRequest)
			return
		}

		err = p.PostUseCase.PostRepo.Delete(idAsNumber, username)
		if err != nil {
			log.Printf("Cannot delete post in postRepo. err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to delete post"}, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (p *PostHandler) GetPostsByUserHandler() http.HandlerFunc { //TODO: PAGINATION
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		if username == "" {
			middleware.SendResponse(w, r, map[string]string{"error": "Username required"}, http.StatusBadRequest)
			return
		}
		posts, err := p.PostUseCase.PostRepo.FindByAuthor(username)
		if err != nil {
			log.Printf("Cannot get posts, err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to get posts"}, http.StatusInternalServerError)
			return
		}
		var resp = make([]models.JsonPost, len(posts))
		for idx, post := range posts {
			resp[idx] = application.MapPostToJson(post)
		}
		middleware.SendResponse(w, r, resp, http.StatusOK)
	}
}

func (p *PostHandler) GetPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["post_id"]

		idAsNumber, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			middleware.SendResponse(w, r, map[string]string{"error": fmt.Sprintf("Invalid ID %s", id)}, http.StatusBadRequest)
			return
		}

		post, err := p.PostUseCase.PostRepo.FindByID(idAsNumber)
		if err != nil {
			log.Printf("Cannot get post, err = %v\n", err)
			middleware.SendResponse(w, r, map[string]string{"error": "Failed to get post"}, http.StatusInternalServerError)
			return
		}

		resp := application.MapPostToJson(post)
		middleware.SendResponse(w, r, resp, http.StatusOK)
	}
}
