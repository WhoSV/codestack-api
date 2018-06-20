package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Favorite Type
type Favorite struct {
	ID       uint `json:"id,omitempty" gorm:"primary_key"`
	UserID   int  `json:"user_id,omitempty"`
	CourseID int  `json:"course_id,omitempty"`
}

// // ErrorMsg Type
// type ErrorMsg struct {
// 	Message string `json:"message"`
// }

// AddFavorite course to user
func AddFavorite(w http.ResponseWriter, r *http.Request) {
	var favorite Favorite

	if err := json.NewDecoder(r.Body).Decode(&favorite); err != nil {
		fmt.Printf("json decode failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
		return
	}

	// Check if course already favorite
	q := db.Where(Favorite{UserID: favorite.UserID, CourseID: favorite.CourseID}).First(&favorite)
	if !q.RecordNotFound() {
		fmt.Printf("course is already favorite\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(ErrorMsg{"course is already favorite"})
		return
	}

	// Create favorite in DB.
	if err := db.Create(&favorite).Error; err != nil {
		fmt.Printf("favorite creation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"favorite creation failed"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetFavorites display's all from Favorite db
func GetFavorites(w http.ResponseWriter, r *http.Request) {
	var favorite []Favorite

	if err := db.Find(&favorite).Error; err != nil {
		fmt.Printf("can not get all favorite from db: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"can not get all favorite from db"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&favorite)
}

// DeleteFavorite course
func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var favorite Favorite

	// Fetch favorite from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&favorite, id)
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

		if err := db.Delete(&favorite).Error; err != nil {
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
