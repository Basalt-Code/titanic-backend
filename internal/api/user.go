package api

import (
	"cmd/app/main.go/internal/model"
	"cmd/app/main.go/internal/store/sqlstore"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateUserRequest struct {
	Username string     `json:"username"`
	Password string     `json:"password"`
	Email    string     `json:"email"`
	Role     model.Role `json:"role"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, store *sqlstore.Store) {
	var req CreateUserRequest

	// Декодируем JSON из тела запроса в структуру
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	u := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := store.User().Create(u); err != nil {
		fmt.Printf("Create method called incorrect: %v\n", err)
		message := err.Error()
		if message == "user with this username or email already exists" {
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
