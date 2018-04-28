package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"

	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	e "codestack-api/endpoints"
)

// User roles
const (
	RoleUndefined = "UNDEFINED"

	RoleAdmin   = "ADMIN"
	RoleTeacher = "TEACHER"
	RoleStudent = "STUDENT"
)

// Course status
const (
	StatusUndefined = "UNDEFINED"

	StatusActive  = "ACTIVE"
	StatusBlocked = "BLOCKED"
)

// Person Type
type Person struct {
	ID       uint   `json:"id,omitempty" gorm:"primary_key"`
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty" gorm:"unique,not null"`
	Role     string `json:"role,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token"`
}

// Course Type
type Course struct {
	ID          uint   `json:"id,omitempty" gorm:"primary_key"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Teacher     string `json:"teacher,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	FileName    string `json:"file_name,omitempty"`
}

// Favorite Type
type Favorite struct {
	ID       uint `json:"id,omitempty" gorm:"primary_key"`
	UserID   int  `json:"user_id,omitempty"`
	CourseID int  `json:"course_id,omitempty"`
}

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

// ErrorMsg Type
type ErrorMsg struct {
	Message string `json:"message"`
}

var db *gorm.DB

// GetPeople from the people var
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

// GetPerson display's a single User
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

// CreatePerson ...
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

// DeletePerson ...
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

// UpdatePerson ...
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID       uint   `json:"id"`
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

// UpdatePersonPassword ...
func UpdatePersonPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID          uint   `json:"id"`
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
	to := data.Email
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
	w.WriteHeader(http.StatusNoContent)
}

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

// GetCourses display's all from the course var
func GetCourses(w http.ResponseWriter, r *http.Request) {
	var courses []Course

	if err := db.Find(&courses).Error; err != nil {
		fmt.Printf("can not get all courses from db: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"can not get all courses from db"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&courses)
}

// GetCourse display's a single Course
func GetCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var course Course

	// Fetch course from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&course, id)
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
		json.NewEncoder(w).Encode(&course)
	}

}

// OpenCourse sends pdf string to browser
func OpenCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Fetch course from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		// Read data
		iid := fmt.Sprint(id)
		data, err := ioutil.ReadFile("data/" + iid + ".pdf")
		if err != nil {
			fmt.Println("error:", err)
		}

		// Encode base64 data
		sEnc := base64.StdEncoding.EncodeToString([]byte(data))

		// Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&sEnc)
	}
}

// CreateCourse ...
func CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course struct {
		Course
		FileBody string `json:"file_body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		fmt.Printf("json decode failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
		return
	}

	// Acquire the data
	b64data := course.FileBody[strings.IndexByte(course.FileBody, ',')+1:]

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Create new course in DB.
	if err := db.Create(&course.Course).Error; err != nil {
		fmt.Printf("course creation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMsg{"course creation failed"})
		return
	}

	// Create data
	id := fmt.Sprint(course.Course.ID)
	f, err := os.Create("data/" + id + ".pdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateCourse ...
func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		FileName    string `json:"file_name"`
		FileBody    string `json:"file_body"`
	}

	params := mux.Vars(r)

	var course Course

	// Fetch course from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&course, id)
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

		// Acquire the data
		b64data := data.FileBody[strings.IndexByte(data.FileBody, ',')+1:]

		// Decode base64 data
		newData, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			fmt.Println("error:", err)
		}

		iid := fmt.Sprint(course.ID)
		f, err := os.Create("data/" + iid + ".pdf")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := f.Write(newData); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}

		// Update course info in DB.
		course.Name = data.Name
		course.Description = data.Description
		course.FileName = data.FileName

		if err := db.Save(&course).Error; err != nil {
			fmt.Printf("unknown database error: %v\n", q.Error)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"unknown database error"})
			return
		}

		// Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&course)
	}
}

// UpdateCourseStatus ...
func UpdateCourseStatus(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID     uint   `json:"id"`
		Status string `json:"status"`
	}

	params := mux.Vars(r)

	var course Course

	// Fetch course from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&course, id)
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

		// Update course status in DB.
		course.Status = data.Status

		if err := db.Save(&course).Error; err != nil {
			fmt.Printf("unknown database error: %v\n", q.Error)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"unknown database error"})
			return
		}

		// Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&course)
	}
}

// DeleteCourse ...
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var course Course

	// Fetch course from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&course, id)
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

		// Delete data
		iid := fmt.Sprint(course.ID)
		errr := os.Remove("data/" + iid + ".pdf")

		if errr != nil {
			fmt.Println(err)
		}

		if err := db.Delete(&course).Error; err != nil {
			fmt.Printf("can not delete course: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorMsg{"can not delete course"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

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

// TokenAuthMiddleware : Generate auth token
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
	router.HandleFunc("/people/reset", e.CORSMiddleware(ResetPassword)).Methods("OPTIONS", "POST")

	// Favorite
	router.HandleFunc("/favorite", e.CORSMiddleware(TokenAuthMiddleware(AddFavorite))).Methods("OPTIONS", "POST")
	router.HandleFunc("/favorite", e.CORSMiddleware(TokenAuthMiddleware(GetFavorites))).Methods("OPTIONS", "GET")
	router.HandleFunc("/favorite/{id}", e.CORSMiddleware(TokenAuthMiddleware(DeleteFavorite))).Methods("OPTIONS", "DELETE")

	// Courses
	router.HandleFunc("/courses", e.CORSMiddleware(TokenAuthMiddleware(GetCourses))).Methods("OPTIONS", "GET")
	router.HandleFunc("/courses", e.CORSMiddleware(TokenAuthMiddleware(CreateCourse))).Methods("OPTIONS", "POST")
	router.HandleFunc("/courses/{id}", e.CORSMiddleware(TokenAuthMiddleware(DeleteCourse))).Methods("OPTIONS", "DELETE")
	router.HandleFunc("/courses/{id}/status", e.CORSMiddleware(TokenAuthMiddleware(UpdateCourseStatus))).Methods("OPTIONS", "PUT", "PATCH")
	router.HandleFunc("/courses/{id}", e.CORSMiddleware(TokenAuthMiddleware(GetCourse))).Methods("OPTIONS", "GET")
	router.HandleFunc("/courses/{id}", e.CORSMiddleware(TokenAuthMiddleware(UpdateCourse))).Methods("OPTIONS", "PUT", "PATCH")
	router.HandleFunc("/courses/{id}/open", e.CORSMiddleware(TokenAuthMiddleware(OpenCourse))).Methods("OPTIONS", "GET")

	// Surveys
	router.HandleFunc("/survey", e.CORSMiddleware(TokenAuthMiddleware(CreateSurvey))).Methods("OPTIONS", "POST")
	router.HandleFunc("/survey", e.CORSMiddleware(TokenAuthMiddleware(GetSurvey))).Methods("OPTIONS", "GET")

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
	db.AutoMigrate(&Person{}, &Course{}, &Favorite{}, &Survey{})
}
