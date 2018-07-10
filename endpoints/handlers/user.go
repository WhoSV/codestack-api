package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strconv"

	"github.com/WhoSV/codestack-api/database"
	"github.com/WhoSV/codestack-api/errors"
	"github.com/WhoSV/codestack-api/model"
	"github.com/gorilla/mux"
)

// GetPeople from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []model.Person

	var db = database.DB()

	if err := db.Find(&people).Error; err != nil {
		fmt.Printf("can not get all people from db: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"can not get all people from db"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&people)
}

// GetPerson display's a single User
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person model.Person

	// Fetch user from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
			return
		}

		var db = database.DB()

		q := db.First(&person, id)
		if q.RecordNotFound() {
			fmt.Printf("record not found: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"record not found"})
			return
		} else if q.Error != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
			return
		}

		// Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&person)
	}

}

// CreatePerson ...
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person model.Person

	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		fmt.Printf("json decode failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
		return
	}

	var db = database.DB()

	// Check if Email provided by user is unique
	q := db.Where(model.Person{Email: person.Email}).First(&person)
	if !q.RecordNotFound() {
		fmt.Printf("email must be unique\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"email must be unique"})
		return
	}

	// Create new person in DB.
	if err := db.Create(&person).Error; err != nil {
		fmt.Printf("person creation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"person creation failed"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeletePerson ...
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person model.Person

	// Fetch user from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
			return
		}

		var db = database.DB()

		q := db.First(&person, id)
		if q.RecordNotFound() {
			fmt.Printf("record not found: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"record not found"})
			return
		} else if q.Error != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
			return
		}

		if err := db.Delete(&person).Error; err != nil {
			fmt.Printf("can not delete person: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"can not delete person"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// UpdatePerson ...
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID       uint   `json:"id"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}

	params := mux.Vars(r)

	var person model.Person

	// Fetch user from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
			return
		}

		var db = database.DB()

		q := db.First(&person, id)
		if q.RecordNotFound() {
			fmt.Printf("record not found: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"record not found"})
			return
		} else if q.Error != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
			return
		}

		// Decode request body into data.
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			fmt.Printf("json decode failed: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
			return
		}

		// Check if Email provided by user is unique
		if person.Email == data.Email {
			fmt.Printf("email must be unique\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"email must be unique"})
			return
		}

		// Update account info in DB.
		person.FullName = data.FullName
		person.Email = data.Email

		if err := db.Save(&person).Error; err != nil {
			fmt.Printf("unknown database error: %v\n", q.Error)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"unknown database error"})
			return
		}

		// Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&person)
	}
}

// UpdatePersonPassword ...
func UpdatePersonPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID          uint   `json:"id"`
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}

	params := mux.Vars(r)

	var person model.Person

	// Fetch user from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
			return
		}

		var db = database.DB()

		q := db.First(&person, id)
		if q.RecordNotFound() {
			fmt.Printf("record not found: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"record not found"})
			return
		} else if q.Error != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
			return
		}

		// Decode request body into data.
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			fmt.Printf("json decode failed: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
			return
		}

		// Check if the password is correct.
		if data.Password != person.Password {
			fmt.Printf("invalid password\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"invalid password"})
			return
		}

		person.Password = data.NewPassword

		if err := db.Save(&person).Error; err != nil {
			fmt.Printf("unknown database error: %v\n", q.Error)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"unknown database error"})
			return
		}

		// Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&person)
	}
}

// ResetPassword send's password by email to User
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string `json:"email"`
	}

	// Decode request body into data.
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Printf("json decode failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
		return
	}

	// Validate user provided email.
	if len(data.Email) == 0 {
		err := fmt.Errorf("email can not be empty")
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

	from := "codeestacks@gmail.com"
	pass := "Neareastuniversity"
	to := data.Email
	body := "Your password is " + person.Password
	msg := "From: " + from + "\n" + "To: " + to + "\n" + "Subject: Reset Password from CodeStack\n\n" + body

	if err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{to}, []byte(msg)); err != nil {
		fmt.Printf("smtp error: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"smtp error"})
		return
	}

	// Success
	w.WriteHeader(http.StatusNoContent)
}
