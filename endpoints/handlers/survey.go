package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WhoSV/codestack-api/database"
	"github.com/WhoSV/codestack-api/errors"
	"github.com/WhoSV/codestack-api/model"
)

// CreateSurvey ...
func CreateSurvey(w http.ResponseWriter, r *http.Request) {
	var survey model.Survey

	if err := json.NewDecoder(r.Body).Decode(&survey); err != nil {
		fmt.Printf("json decode failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
		return
	}

	var db = database.DB()

	// Create new survey in DB.
	if err := db.Create(&survey).Error; err != nil {
		fmt.Printf("survey creation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"survey creation failed"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetSurvey from the survey var
func GetSurvey(w http.ResponseWriter, r *http.Request) {
	var surveys []model.Survey

	var db = database.DB()

	if err := db.Find(&surveys).Error; err != nil {
		fmt.Printf("can not get all surveys from db: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"can not get all surveys from db"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&surveys)
}
