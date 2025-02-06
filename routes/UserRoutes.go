package routes

import (
	"github.com/gorilla/mux"
	"github.com/subhendu/go-auth/controllers"
	"github.com/subhendu/go-auth/middlewares"
)

func UserRoutes(route *mux.Router) {
	protectedRoutes := route.PathPrefix("/user").Subrouter()
	protectedRoutes.Use(middlewares.IsAuthenticated)

	protectedRoutes.HandleFunc("", controllers.GetUsers).Methods("GET")
	protectedRoutes.HandleFunc("/{user_id}", controllers.GetUser).Methods("GET")
}
