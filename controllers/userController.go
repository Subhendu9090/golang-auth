package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/subhendu/go-auth/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.GetCollection("users")

func HashPassword() {

}
func VerifyPassword() {

}
func Login(w http.ResponseWriter, r *http.Request) {

}
func SignUp(w http.ResponseWriter, r *http.Request) {

}
func GetUsers(w http.ResponseWriter, r *http.Request) {

}
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["user_id"]
	fmt.Println("Searching for user with ID:", userId)

	// Convert userId string to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var user bson.M
	err = userCollection.FindOne(r.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	// Return the found user as a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
