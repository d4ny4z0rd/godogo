package router

import (
	"github.com/d4ny4z0rd/godogo/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/todos",middleware.GetAllTodos).Methods("GET","OPTIONS")
	router.HandleFunc("/api/todos",middleware.CreateTodo).Methods("POST","OPTIONS")
	router.HandleFunc("/api/completeTodo/{id}",middleware.CompleteTodo).Methods("PUT","OPTIONS")
	router.HandleFunc("/api/undoTodo/{id}",middleware.UndoTodo).Methods("PUT","OPTIONS")
	router.HandleFunc("/api/todos/{id}",middleware.DeleteTodo).Methods("DELETE","OPTIONS")
	router.HandleFunc("/api/todos",middleware.DeleteAll).Methods("DELETE","OPTIONS")

	return router
}