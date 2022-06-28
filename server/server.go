package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"rest-ws/database"
	"rest-ws/repository"

	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JWTScret    string
	DatabaseURL string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (b Broker) Config() *Config {
	return b.config
}

func NewConfig(port string, jwtScret string, dataBaseURL string) (*Config, error) {
	if port == "" || jwtScret == "" || dataBaseURL == "" {
		return nil, errors.New("All fields are required for to create a new config")
	}
	return &Config{
		Port:        port,
		JWTScret:    jwtScret,
		DatabaseURL: dataBaseURL,
	}, nil
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTScret == "" {
		return nil, errors.New("secret is required")
	}

	if config.DatabaseURL == "" {
		return nil, errors.New("database url is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}

	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)

	repo, err := database.NewPostgresRepository(b.config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	repository.SetRepository(repo)

	log.Println("Starting server on port: ", b.Config().Port)

	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
