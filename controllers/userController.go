package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/subhendu/go-auth/database"
	"github.com/subhendu/go-auth/helper"
	"github.com/subhendu/go-auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.GetCollection("users")
var validate = validator.New()

func HashPassword() {

}
func VerifyPassword() {

}
func Login(w http.ResponseWriter, r *http.Request) {

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Sign up function works")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Error", err)
		http.Error(w, "Error in Sending data", 404)
	}
	// validate the user data
	//  after getting data check user is already exist or not
	// hash the password using bcrypt
	// generate access token and refresh token
	// create a unique user_id (primitive)
	// save the data in mongodb

	err = validate.Struct(user)
	if err != nil {
		http.Error(w, "Validation Error"+err.Error(), 400)
		return
	}

	var existingUser models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"$or": []bson.M{{"email": user.Emil}, {"phone": user.Phone}}}).Decode(&existingUser)

	if err == nil {
		http.Error(w, "User already exist", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	*user.Password = string(hashedPassword)

	accessToken, err := helper.GenerateToken(user, 15*time.Minute)
	refreshToken, err := helper.GenerateToken(user, 7*24*time.Minute)

	user.Token = &accessToken
	user.Refresh_token = &refreshToken
	user.Id = primitive.NewObjectID()
	user.User_id = user.Id.Hex()
	user.Created_at = time.Now()
	user.Updated_at = time.Now()

	_, err = userCollection.InsertOne(context.TODO(), user)

	if err != nil {
		http.Error(w, "User Creation Failed", http.StatusInternalServerError)
		return
	}
	user.Password = nil
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["user_id"]
	if err := helper.MatchUserTypeToUid(r, userId); err != nil {
		json.NewEncoder(w).Encode(err)
	}
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var user models.User
	err = userCollection.FindOne(r.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
