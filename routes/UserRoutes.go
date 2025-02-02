package routes

import (
	"github.com/gorilla/mux"
	"github.com/subhendu/go-auth/controllers"
)

func UserRoutes(route *mux.Router) {
	// route.Use(middlewares.Authenticate)
	route.HandleFunc("/user", controllers.GetUsers).Methods("GET")
	route.HandleFunc("/user/{user_id}", controllers.GetUser).Methods("GET")
}
