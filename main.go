package main

import (
	"fmt"
	"log"
	"net/smtp"

	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	e "../codestack-api/endpoints"
)

// User roles
const (
	RoleUndefined = "UNDEFINED"

	RoleAdmin   = "ADMIN"
	RoleTeacher = "TEACHER"
	RoleStudent = "STUDENT"
)

// The person Type
type Person struct {
	ID       uint   `json:"id,omitempty" gorm:"primary_key"`
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty" gorm:"unique,not null"`
	Role     string `json:"role,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token"`
}

type ErrorMsg struct {
	Message string `json:"message"`
}

var db *gorm.DB

// Display all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person

	if err := db.Find(&people).Error; err != nil {
		fmt.Printf("can not get all people from db: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"can not get all people from db"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&people)
}

// Display a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person

	// Fetch user from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&person, id)
		if q.RecordNotFound() {
			fmt.Printf("record not found: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorMsg{"record not found"})
			return
		} else if q.Error != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		// Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&person)
	}

}

// Create a new user
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person Person

	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		fmt.Printf("json decode failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
		return
	}

	// Check if Email provided by user is unique
	q := db.Where(Person{Email: person.Email}).First(&person)
	if !q.RecordNotFound() {
		fmt.Printf("email must be unique\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(ErrorMsg{"email must be unique"})
		return
	}

	// Create new person in DB.
	if err := db.Create(&person).Error; err != nil {
		fmt.Printf("person creation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"person creation failed"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Delete an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person

	// Fetch user from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&person, id)
		if q.RecordNotFound() {
			fmt.Printf("record not found: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorMsg{"record not found"})
			return
		} else if q.Error != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		if err := db.Delete(&person).Error; err != nil {
			fmt.Printf("can not delete person: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"can not delete person"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// Update user
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Id       uint   `json:"id"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}

	params := mux.Vars(r)

	var person Person

	// Fetch user from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&person, id)
		if q.RecordNotFound() {
			fmt.Printf("record not found: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorMsg{"record not found"})
			return
		} else if q.Error != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		// Decode request body into data.
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			fmt.Printf("json decode failed: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		// Check if Email provided by user is unique
		if person.Email == data.Email {
			fmt.Printf("email must be unique\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorMsg{"email must be unique"})
			return
		}

		// Update account info in DB.
		person.FullName = data.FullName
		person.Email = data.Email

		if err := db.Save(&person).Error; err != nil {
			fmt.Printf("unknown database error: %v\n", q.Error)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"unknown database error"})
			return
		}

		// Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&person)
	}
}

// Update Person Password
func UpdatePersonPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Id          uint   `json:"id"`
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}

	params := mux.Vars(r)

	var person Person

	// Fetch user from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&person, id)
		if q.RecordNotFound() {
			fmt.Printf("record not found: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorMsg{"record not found"})
			return
		} else if q.Error != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		// Decode request body into data.
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			fmt.Printf("json decode failed: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		// Check if the password is correct.
		if data.Password != person.Password {
			fmt.Printf("invalid password\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorMsg{"invalid password"})
			return
		}

		person.Password = data.NewPassword

		if err := db.Save(&person).Error; err != nil {
			fmt.Printf("unknown database error: %v\n", q.Error)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"unknown database error"})
			return
		}

		// Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&person)
	}
}

// Reset password
func ResertPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string `json:"email"`
	}

	// Decode request body into data.
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Printf("json decode failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
		return
	}

	// Validate user provided email.
	if len(data.Email) == 0 {
		err := fmt.Errorf("email can not be empty")
		fmt.Printf("validation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{"validation failed"})
		return
	}

	// Fetch user by email address provided.
	var person Person

	q := db.Where(Person{Email: data.Email}).First(&person)

	if q.RecordNotFound() {
		fmt.Printf("record not found\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorMsg{"invalid email or password"})
		return
	} else if q.Error != nil {
		fmt.Printf("unknown database error: %v\n", q.Error)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"unknown database error"})
		return
	}

	from := "codeestacks@gmail.com"
	pass := "Neareastuniversity"
	to := "vladislavsiumbeli@gmail.com" // data.Email
	body := "Your password is " + person.Password
	msg := "From: " + from + "\n" + "To: " + to + "\n" + "Subject: Reset Password from CodeStack\n\n" + body

	if err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{to}, []byte(msg)); err != nil {
		fmt.Printf("smtp error: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"smtp error"})
		return
	}

	// Success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&person)
}

// Authorization Handler
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
		json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
		return
	}

	// Validate user provided email and password.
	if len(data.Email) == 0 {
		err := fmt.Errorf("email can not be empty")
		fmt.Printf("validation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{"validation failed"})
		return
	}
	if len(data.Password) == 0 {
		err := fmt.Errorf("password can not be empty")
		fmt.Printf("validation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{"validation failed"})
		return
	}

	// Fetch user by email address provided.
	var person Person

	q := db.Where(Person{Email: data.Email}).First(&person)

	if q.RecordNotFound() {
		fmt.Printf("record not found\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorMsg{"invalid email or password"})
		return
	} else if q.Error != nil {
		fmt.Printf("unknown database error: %v\n", q.Error)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"unknown database error"})
		return
	}

	// Check if the password is correct.
	if data.Password != person.Password {
		fmt.Printf("invalid password\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorMsg{"invalid email or password"})
		return
	}

	// Create an authorization token.
	var token string

	randUUID, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("unknown uuid generation error: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"unknown uuid generation error"})
		return
	}
	token = randUUID.String()

	// Store creted token in DB.
	person.Token = token
	if err := db.Save(&person).Error; err != nil {
		fmt.Printf("unknown database error: %v\n", q.Error)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"unknown database error"})
		return
	}

	// Success
	var resp struct {
		Token string `json:"token"`
		Id    uint   `json:"id"`
	}

	resp.Token = token
	resp.Id = person.ID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&resp)
}

func TokenAuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := strings.Split(r.Header.Get("Authorization"), " ")

		if len(s) != 2 {
			fmt.Printf("malformed authorization token\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorMsg{"Authorization token should be in the form of Authorization: Bearer <token>"})
			return
		}

		if s[0] != "Bearer" {
			fmt.Printf("authorization token is not bearer kind\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorMsg{"Authorization token should be in the form of Authorization: Bearer <token>"})
			return
		}

		token := s[1]

		// Check token
		var person Person

		q := db.Where(Person{Token: token}).First(&person)
		if q.RecordNotFound() {
			fmt.Printf("not authorized\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorMsg{"Not Authorized"})
			return
		} else if q.Error != nil {
			fmt.Printf("unknown database error: %v\n", q.Error)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"unknown database error"})
			return
		}

		// Permitted
		h.ServeHTTP(w, r)
	}
}

// main function to boot up everything
func main() {
	defer db.Close()
	router := mux.NewRouter()

	// Authorization
	router.HandleFunc("/auth", e.CORSMiddleware(Authorize)).Methods("OPTIONS", "POST")

	// Users
	router.HandleFunc("/people", e.CORSMiddleware(TokenAuthMiddleware(GetPeople))).Methods("OPTIONS", "GET")
	router.HandleFunc("/people/{id}", e.CORSMiddleware(TokenAuthMiddleware(GetPerson))).Methods("OPTIONS", "GET")
	router.HandleFunc("/people", e.CORSMiddleware(CreatePerson)).Methods("OPTIONS", "POST")
	router.HandleFunc("/people/{id}", e.CORSMiddleware(UpdatePerson)).Methods("OPTIONS", "PUT", "PATCH")
	router.HandleFunc("/people/{id}/update", e.CORSMiddleware(UpdatePersonPassword)).Methods("OPTIONS", "PUT", "PATCH")
	router.HandleFunc("/people/{id}", e.CORSMiddleware(TokenAuthMiddleware(DeletePerson))).Methods("OPTIONS", "DELETE")

	// Reset Password
	router.HandleFunc("/people/reset", e.CORSMiddleware(ResertPassword)).Methods("OPTIONS", "POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func init() {
	// Initialize the database.
	var err error
	db, err = gorm.Open("sqlite3", "codestack.db")
	if err != nil {
		fmt.Printf("database connection failed: %v", err)
		return
	}
	db.AutoMigrate(&Person{})
}
