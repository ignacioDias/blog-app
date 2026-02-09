package infrastructure

import (
	"net/http"
	"postapi/internal/infrastructure/handlers"
	"postapi/internal/middleware"

	"github.com/gorilla/mux"
)

type Router struct {
	router         *mux.Router
	postHandler    *handlers.PostHandler
	followHandler  *handlers.FollowHandler
	userHandler    *handlers.UserHandler
	profileHandler *handlers.ProfileHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(
	postHandler *handlers.PostHandler,
	followHandler *handlers.FollowHandler,
	userHandler *handlers.UserHandler,
	profileHandler *handlers.ProfileHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		router:         mux.NewRouter(),
		postHandler:    postHandler,
		followHandler:  followHandler,
		userHandler:    userHandler,
		profileHandler: profileHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) SetupRoutes() *mux.Router {

	// Rutas de autenticación
	r.router.HandleFunc("/api/register", r.userHandler.RegisterUserHandler()).Methods("POST")
	r.router.HandleFunc("/api/login", r.userHandler.LoginHandler()).Methods("POST")

	// Rutas de posts
	r.router.HandleFunc("/api/posts", r.authMiddleware.AuthMiddleware(r.postHandler.CreatePostHandler())).Methods("POST")
	r.router.HandleFunc("/api/posts/{post_id}", r.postHandler.GetPostHandler()).Methods("GET")
	r.router.HandleFunc("/api/posts/{post_id}", r.authMiddleware.AuthMiddleware(r.postHandler.UpdatePostHandler())).Methods("PATCH")
	r.router.HandleFunc("/api/posts/{post_id}", r.authMiddleware.AuthMiddleware(r.postHandler.DeletePostHandler())).Methods("DELETE")

	// Rutas de usuarios
	r.router.HandleFunc("/api/users/{username}", r.userHandler.GetUserByUsernameHandler()).Methods("GET")
	r.router.HandleFunc("/api/users/{username}/posts", r.postHandler.GetPostsByUserHandler()).Methods("GET")
	r.router.HandleFunc("/api/follow/{username}", r.authMiddleware.AuthMiddleware(r.followHandler.FollowHandler())).Methods("POST")
	r.router.HandleFunc("/api/unfollow/{username}", r.authMiddleware.AuthMiddleware(r.followHandler.UnfollowHandler())).Methods("DELETE")
	r.router.HandleFunc("/api/users/{username}/followers", r.followHandler.GetFollowersHandler()).Methods("GET")
	r.router.HandleFunc("/api/users/{username}/following", r.followHandler.GetFollowingHandler()).Methods("GET")

	// Rutas de perfiles
	r.router.HandleFunc("/api/profiles/{username}", r.profileHandler.GetProfileHandler()).Methods("GET")
	r.router.HandleFunc("/api/profiles/me", r.authMiddleware.AuthMiddleware(r.profileHandler.CreateProfileHandler())).Methods("POST")
	r.router.HandleFunc("/api/profiles/me", r.authMiddleware.AuthMiddleware(r.profileHandler.UpdateProfileHandler())).Methods("PATCH")

	fs := http.FileServer(http.Dir("./web"))
	r.router.PathPrefix("/").Handler(fs)

	return r.router
}
