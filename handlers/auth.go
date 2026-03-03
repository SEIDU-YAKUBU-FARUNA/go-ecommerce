package handlers

import (
	"context"
	"encoding/json"
	"go-ecommerce/database"
	"go-ecommerce/models"
	"go-ecommerce/utils"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {

		http.Error(w, "only post is alllowed", http.StatusMethodNotAllowed)
		return

	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid body request", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "all field is required", http.StatusBadRequest)
		return
	}

	collection := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	/**var existingUser models.User
		err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)

		if err == nil {
			http.Error(w, "email already exist", http.StatusConflict)
			return
		}
	**/

	var existingUser models.User
	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {

		http.Error(w, "user email already exist", http.StatusConflict)
		return
	}

	hashedpassword, err := utils.Hashpassword(user.Password)
	if err != nil {
		http.Error(w, "password not hashed ", http.StatusInternalServerError)
		return
	}

	user.Password = hashedpassword

	user.IsAdmin = false

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "apllicaton/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "user successfully created"})

}

//login function

func LoginUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "invalid reques ", http.StatusBadRequest)
		return
	}

	if loginData.Email == "" || loginData.Password == "" {
		http.Error(w, "all field are required", http.StatusBadRequest)
		return
	}

	collection := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	err = collection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "inavalid email or password", http.StatusUnauthorized)
		return
	}

	err = utils.Checkpassword(user.Password, loginData.Password)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{

		"message": "Login Successful",
		"user": map[string]interface{}{
			"id":       user.ID.Hex(),
			"name":     user.Name,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
		},
	})

}
