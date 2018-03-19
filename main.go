package main

import (
	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	Mail     string `json:"mail,omitempty"`
	Role     string `json:"role,omitempty"`
	Password string `json:"password,omitempty"`
}

type ErrorMsg struct {
	Message string `json:"message"`
}

var db *gorm.DB

// Display all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person

	if err := db.Find(&people).Error; err != nil {
		fmt.Printf("can not get all people from db: %v", err)
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
			fmt.Printf("can not convert from string to int: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&person, id)
		if q.RecordNotFound() {
			fmt.Printf("record not found: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorMsg{"record not found"})
			return
		} else if q.Error != nil {
			fmt.Printf("can not convert from string to int: %v", err)
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
		fmt.Printf("json decode failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
		return
	}

	// Create new person in DB.
	if err := db.Create(&person).Error; err != nil {
		fmt.Printf("person creation failed: %v", err)
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
			fmt.Printf("can not convert from string to int: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&person, id)
		if q.RecordNotFound() {
			fmt.Printf("record not found: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorMsg{"record not found"})
			return
		} else if q.Error != nil {
			fmt.Printf("can not convert from string to int: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		if err := db.Delete(&person).Error; err != nil {
			fmt.Printf("can not delete person: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"can not delete person"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// HandleCORS is a CORS handler.
func HandleCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Stop if the request is OPTIONS.
	if r.Method == "OPTIONS" {
		return
	}
}

// CORSMiddleware sets up CORS headers.
func CORSMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		HandleCORS(w, r)
		if r.Method == "OPTIONS" {
			return
		}
		// Call the next handler.
		handler.ServeHTTP(w, r)
	}
}

// main function to boot up everything
func main() {
	defer db.Close()
	router := mux.NewRouter()

	router.HandleFunc("/people", CORSMiddleware(GetPeople)).Methods("OPTIONS", "GET")
	router.HandleFunc("/people/{id}", CORSMiddleware(GetPerson)).Methods("OPTIONS", "GET")
	router.HandleFunc("/people", CORSMiddleware(CreatePerson)).Methods("OPTIONS", "POST")
	router.HandleFunc("/people/{id}", CORSMiddleware(DeletePerson)).Methods("OPTIONS", "DELETE")
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
