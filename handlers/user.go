package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"rest-ws/models"
	"rest-ws/repository"
	"rest-ws/server"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

type SingUpLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SingUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func SingUpHanlder(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request SingUpLoginRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var user = models.User{
			Id:       id.String(),
			Email:    request.Email,
			Password: string(hashedPassword),
		}

		errInsert := repository.InsertUser(r.Context(), &user)
		if errInsert != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SingUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})

	}
}

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request SingUpLoginRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			log.Println(err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(s.Config().JWTScret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})

	}
}
