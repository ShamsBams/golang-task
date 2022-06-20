package router

import (
	"tz/middleware"
	"tz/parser"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/post/{id}", middleware.GetPost).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/post", middleware.GetAllPost).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/post/{id}", middleware.UpdatePost).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletepost/{id}", middleware.DeletePost).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/parse", parser.Parser).Methods("POST", "OPTIONS")

	return router
}
