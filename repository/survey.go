package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Survey Type
type Survey struct {
	ID       uint `json:"id,omitempty" gorm:"primary_key"`
	CourseID int  `json:"course_id,omitempty"`
	First    int  `json:"first,omitempty"`
	Second   int  `json:"second,omitempty"`
	Third    int  `json:"third,omitempty"`
	Fourth   int  `json:"fourth,omitempty"`
	Fifth    int  `json:"fifth,omitempty"`
}

// // ErrorMsg Type
// type ErrorMsg struct {
// 	Message string `json:"message"`
// }

// CreateSurvey ...
func CreateSurvey(w http.ResponseWriter, r *http.Request) {
	var survey Survey

	if err := json.NewDecoder(r.Body).Decode(&survey); err != nil {
		fmt.Printf("json decode failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
		return
	}

	// Create new survey in DB.
	if err := db.Create(&survey).Error; err != nil {
		fmt.Printf("survey creation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"survey creation failed"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetSurvey from the survey var
func GetSurvey(w http.ResponseWriter, r *http.Request) {
	var surveys []Survey

	if err := db.Find(&surveys).Error; err != nil {
		fmt.Printf("can not get all surveys from db: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"can not get all surveys from db"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&surveys)
}
