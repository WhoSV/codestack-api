package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/WhoSV/codestack-api/database"
	"github.com/WhoSV/codestack-api/errors"
	model "github.com/WhoSV/codestack-api/model"
	"github.com/google/uuid"
)

// TokenAuthMiddleware : Generate auth token
func TokenAuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := strings.Split(r.Header.Get("Authorization"), " ")

		if len(s) != 2 {
			fmt.Printf("malformed authorization token\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"Authorization token should be in the form of Authorization: Bearer <token>"})
			return
		}

		if s[0] != "Bearer" {
			fmt.Printf("authorization token is not bearer kind\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"Authorization token should be in the form of Authorization: Bearer <token>"})
			return
		}

		token := s[1]

		// Check token
		var person model.Person

		var db = database.DB()

		q := db.Where(model.Person{Token: token}).First(&person)
		if q.RecordNotFound() {
			fmt.Printf("not authorized\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"Not Authorized"})
			return
		} else if q.Error != nil {
			fmt.Printf("unknown database error: %v\n", q.Error)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"unknown database error"})
			return
		}

		// Permitted
		h.ServeHTTP(w, r)
	}
}

// Authorize : Authorization Handler
func Authorize(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode request body into data.
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Printf("json decode failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
		return
	}

	// Validate user provided email and password.
	if len(data.Email) == 0 {
		err := fmt.Errorf("email can not be empty")
		fmt.Printf("validation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"validation failed"})
		return
	}
	if len(data.Password) == 0 {
		err := fmt.Errorf("password can not be empty")
		fmt.Printf("validation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"validation failed"})
		return
	}

	// Fetch user by email address provided.
	var person model.Person

	var db = database.DB()

	q := db.Where(model.Person{Email: data.Email}).First(&person)

	if q.RecordNotFound() {
		fmt.Printf("record not found\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"invalid email or password"})
		return
	} else if q.Error != nil {
		fmt.Printf("unknown database error: %v\n", q.Error)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"unknown database error"})
		return
	}

	// Check if the password is correct.
	if data.Password != person.Password {
		fmt.Printf("invalid password\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"invalid email or password"})
		return
	}

	// Create an authorization token.
	var token string

	randUUID, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("unknown uuid generation error: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"unknown uuid generation error"})
		return
	}
	token = randUUID.String()

	// Store creted token in DB.
	person.Token = token
	if err := db.Save(&person).Error; err != nil {
		fmt.Printf("unknown database error: %v\n", q.Error)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"unknown database error"})
		return
	}

	// Success
	var resp struct {
		Token string `json:"token"`
		ID    uint   `json:"id"`
		Role  string `json:"user_role"`
	}

	resp.Token = token
	resp.ID = person.ID
	resp.Role = person.Role

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&resp)
}
