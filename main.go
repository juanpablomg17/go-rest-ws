package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"rest-ws/handlers"
	"rest-ws/middleware"
	"rest-ws/server"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	fmt.Println(DATABASE_URL)

	ctx := context.Background()

	mainConfig, errConfig := server.NewConfig(PORT, JWT_SECRET, DATABASE_URL)
	if errConfig != nil {
		log.Fatal(errConfig)
	}

	mainServer, errServer := server.NewServer(ctx, mainConfig)
	if errServer != nil {
		log.Fatal(errServer)
	}

	mainServer.Start(BindRoutes)

}

func BindRoutes(s server.Server, r *mux.Router) {
	r.Use(middleware.CheckAuthMiddleware(s))
	r.HandleFunc("/", handlers.HomeHandlers(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SingUpHanlder(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
}
