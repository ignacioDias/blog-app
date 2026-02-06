package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"postapi/internal/application"
	"postapi/internal/infrastructure/handlers"
	httpserver "postapi/internal/infrastructure/httpserver"
	"postapi/internal/infrastructure/persistence"
	"postapi/internal/middleware"
	"syscall"
	"time"
)

func main() {
	// Configuración de la base de datos
	database := &persistence.DB{}

	err := database.Open()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	userRepo := database.UserRepository
	postRepo := database.PostRepository
	profileRepo := database.ProfileRepository
	followRepo := database.UserFollowRepository

	jwtService := application.NewJWTService("secret-key")

	postUseCase := application.PostUseCase{PostRepo: postRepo}
	userUseCase := application.UserUseCase{UserRepo: userRepo, FollowRepo: followRepo}
	profileUseCase := application.ProfileUseCase{ProfileRepository: profileRepo}

	postHandler := &handlers.PostHandler{PostUseCase: postUseCase}
	followHandler := &handlers.FollowHandler{UserUseCase: userUseCase}
	userHandler := &handlers.UserHandler{UserUseCase: userUseCase, JWTService: jwtService}
	profileHandler := &handlers.ProfileHandler{ProfileUseCase: profileUseCase}

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	router := httpserver.NewRouter(
		postHandler,
		followHandler,
		userHandler,
		profileHandler,
		authMiddleware,
	)

	server := httpserver.NewServer("8080", router)

	// Canal para manejar señales de interrupción
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Println("Server started successfully")

	// Esperar señal de interrupción
	<-done
	log.Println("Server stopping...")

	// Shutdown gracefully con timeout de 30 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
