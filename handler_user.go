package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prmzk/go-backend-1/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdateAt:  user.UpdatedAt,
		ApiKey:    user.ApiKey,
	}
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, "Invalid request payload")
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:   uuid.New(),
		Name: params.Name,
	})

	if err != nil {
		respondWithError(w, 500, "Internal server error")
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := apiCfg.DB.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "Internal server error")
		return
	}

	var userList []User
	for _, u := range users {
		userList = append(userList, databaseUserToUser(u))
	}

	respondWithJSON(w, 200, userList)
}

func (apiCfg *apiConfig) handlerGetUsersByApiKey(w http.ResponseWriter, r *http.Request) {
	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		respondWithError(w, 500, "Internal server error")
		return
	}

	// if !user == nil {
	// 	respondWithError(w, 404, "User not found")
	// 	return
	// }

	respondWithJSON(w, 200, databaseUserToUser(user))
}
