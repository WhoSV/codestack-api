package router

import (
	e "github.com/WhoSV/codestack-api/endpoints"
	"github.com/WhoSV/codestack-api/endpoints/handlers"
	"github.com/gorilla/mux"
)

// GetRouter ...
func GetRouter() *mux.Router {
	router := mux.NewRouter()

	// Authorization
	router.HandleFunc("/auth", e.CORSMiddleware(e.Authorize)).Methods("OPTIONS", "POST")

	// Users
	router.HandleFunc("/people", e.CORSMiddleware(handlers.GetPeople)).Methods("OPTIONS", "GET") // add e.TokenAuthMiddleware
	router.HandleFunc("/people/{id}", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.GetPerson))).Methods("OPTIONS", "GET")
	router.HandleFunc("/people", e.CORSMiddleware(handlers.CreatePerson)).Methods("OPTIONS", "POST")
	router.HandleFunc("/people/{id}", e.CORSMiddleware(handlers.UpdatePerson)).Methods("OPTIONS", "PUT", "PATCH")
	router.HandleFunc("/people/{id}/update", e.CORSMiddleware(handlers.UpdatePersonPassword)).Methods("OPTIONS", "PUT", "PATCH")
	router.HandleFunc("/people/{id}", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.DeletePerson))).Methods("OPTIONS", "DELETE")
	router.HandleFunc("/people/reset", e.CORSMiddleware(handlers.ResetPassword)).Methods("OPTIONS", "POST")

	// Favorite
	router.HandleFunc("/favorite", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.AddFavorite))).Methods("OPTIONS", "POST")
	router.HandleFunc("/favorite", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.GetFavorites))).Methods("OPTIONS", "GET")
	router.HandleFunc("/favorite/{id}", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.DeleteFavorite))).Methods("OPTIONS", "DELETE")

	// Courses
	router.HandleFunc("/courses", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.GetCourses))).Methods("OPTIONS", "GET")
	router.HandleFunc("/courses", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.CreateCourse))).Methods("OPTIONS", "POST")
	router.HandleFunc("/courses/{id}", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.DeleteCourse))).Methods("OPTIONS", "DELETE")
	router.HandleFunc("/courses/{id}/status", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.UpdateCourseStatus))).Methods("OPTIONS", "PUT", "PATCH")
	router.HandleFunc("/courses/{id}", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.GetCourse))).Methods("OPTIONS", "GET")
	router.HandleFunc("/courses/{id}", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.UpdateCourse))).Methods("OPTIONS", "PUT", "PATCH")
	router.HandleFunc("/courses/{id}/open", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.OpenCourse))).Methods("OPTIONS", "GET")

	// Surveys
	router.HandleFunc("/survey", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.CreateSurvey))).Methods("OPTIONS", "POST")
	router.HandleFunc("/survey", e.CORSMiddleware(e.TokenAuthMiddleware(handlers.GetSurvey))).Methods("OPTIONS", "GET")

	return router
}
