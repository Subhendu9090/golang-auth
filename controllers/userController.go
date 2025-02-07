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
func VerifyPassword(inComingPassword string, password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(inComingPassword))
	if err != nil {
		return false, "Incorrect Password"
	}
	return true, "Correct Password"
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	var existedUser models.User
	defer cancel()
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "pass the email and password", http.StatusNotFound)
		return
	}

	err = userCollection.FindOne(ctx, bson.M{"email": user.Emil}).Decode(&existedUser)

	if err != nil {
		http.Error(w, "User Not found", http.StatusNotFound)
		return
	}
	isPasswordCorrect, msg := VerifyPassword(*user.Password, *existedUser.Password)

	if !isPasswordCorrect {
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	accessToken, _ := helper.GenerateToken(existedUser, 24*time.Hour)
	refreshToken, _ := helper.GenerateToken(existedUser, 7*24*time.Hour)

	existedUser.Refresh_token = &refreshToken
	existedUser.Token = &accessToken
	existedUser.Password = nil
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existedUser)
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

	userCountEmail, err := userCollection.CountDocuments(context.TODO(), bson.M{"email": user.Emil})
	if err != nil {
		http.Error(w, "Internal server error", 404)
		return
	}
	if userCountEmail > 0 {
		http.Error(w, "User Already exists with this email", http.StatusConflict)
		return
	}

	userCountPhone, err := userCollection.CountDocuments(context.TODO(), bson.M{"phone": user.Phone})

	if err != nil {
		http.Error(w, "Internal server error", 404)
		return
	}

	if userCountPhone > 0 {
		http.Error(w, "User Already exists with this phone Number", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	*user.Password = string(hashedPassword)

	accessToken, _ := helper.GenerateToken(user, 15*time.Minute)
	refreshToken, _ := helper.GenerateToken(user, 7*24*time.Minute)

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
