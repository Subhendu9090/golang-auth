package routes

import (
	"github.com/gorilla/mux"
	"github.com/subhendu/go-auth/controllers"
)

func AuthRoutes(router *mux.Router) {
	router.HandleFunc("/user/sign_up", controllers.SignUp).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
}
