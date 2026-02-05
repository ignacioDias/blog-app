package app

import (
	"postapi/app/database"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     database.PostDB
}

func New() *App {
	a := &App{
		Router: mux.NewRouter(),
	}
	a.initRoutes()
	return a
}

func (a *App) initRoutes() {
	a.Router.HandleFunc("/", a.IndexHandler()).Methods("GET")
	a.Router.HandleFunc("/api/posts", a.AuthMiddleware(a.CreatePostHandler())).Methods("POST")
	a.Router.HandleFunc("/api/posts", a.GetPostsHandler()).Methods("GET")
	a.Router.HandleFunc("/api/{username}/posts", a.GetPostsByUserHanlder()).Methods("GET")
	a.Router.HandleFunc("/api/{username}", a.GetProfileHandler()).Methods("GET")
	a.Router.HandleFunc("/api/register", a.RegisterUserHandler()).Methods("POST")
	a.Router.HandleFunc("/api/login", a.LoginHandler()).Methods("POST")
}
