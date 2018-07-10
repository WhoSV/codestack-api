package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/WhoSV/codestack-api/database"
	"github.com/WhoSV/codestack-api/errors"
	"github.com/WhoSV/codestack-api/model"
	"github.com/gorilla/mux"
)

// GetCourses display's all from the course var
func GetCourses(w http.ResponseWriter, r *http.Request) {
	var courses []model.Course

	var db = database.DB()

	if err := db.Find(&courses).Error; err != nil {
		fmt.Printf("can not get all courses from db: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"can not get all courses from db"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&courses)
}

// GetCourse display's a single Course
func GetCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var course model.Course

	var db = database.DB()

	// Fetch course from db.
	if id := params["id"]; len(id) > 0 {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("can not convert from string to int: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
			return
		}

		q := db.First(&course, id)
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
			json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
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
		model.Course
		FileBody string `json:"file_body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		fmt.Printf("json decode failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"json decode failed"})
		return
	}

	// Acquire the data
	b64data := course.FileBody[strings.IndexByte(course.FileBody, ',')+1:]

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		fmt.Println("error:", err)
	}

	var db = database.DB()

	// Create new course in DB.
	if err := db.Create(&course.Course).Error; err != nil {
		fmt.Printf("course creation failed: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ErrorMsg{"course creation failed"})
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

	var course model.Course

	// Fetch course from db.
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

		q := db.First(&course, id)
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
			json.NewEncoder(w).Encode(errors.ErrorMsg{"unknown database error"})
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

	var course model.Course

	// Fetch course from db.
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

		q := db.First(&course, id)
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

		// Update course status in DB.
		course.Status = data.Status

		if err := db.Save(&course).Error; err != nil {
			fmt.Printf("unknown database error: %v\n", q.Error)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ErrorMsg{"unknown database error"})
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

	var course model.Course

	// Fetch course from db.
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

		q := db.First(&course, id)
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
			json.NewEncoder(w).Encode(errors.ErrorMsg{"can not delete course"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
