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
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST is allowed")
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid body request")
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {

		utils.RespondWithError(w, http.StatusBadRequest, "all field is required")
		return
	}

	collection := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	/**

	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "user email already exist")
		return
	}

	HashPassword, err := utils.HashPassword(user.Password)
	if err != nil {

		utils.RespondWithError(w, http.StatusInternalServerError, "password not hashed")
		return
	}

	user.Password = HashPassword
	user.IsAdmin = true
	_, err = collection.InsertOne(ctx, user)
	if err != nil {

		utils.RespondWithError(w, http.StatusInternalServerError, "password not hashed")
		return
	}

	w.Header().Set("Content-Type", "apllicaton/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "user successfully created"})

	**/
	// 1. Check if user exists
	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil { // If err is NIL, it means a user WAS found!
		utils.RespondWithError(w, http.StatusBadRequest, "user email already exists")
		return
	}

	// 2. Hash the password (MAKE SURE THE 'P' IS CAPITALIZED)
	hashedpassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// 3. Assign and Save
	user.Password = hashedpassword
	user.IsAdmin = false // Default to non-admin. You can change this logic as needed.
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	// 4. Success
	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "user successfully created"})

}

//login function

func LoginUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "all field is required")
		return
	}

	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		//http.Error(w, "invalid reques ", http.StatusBadRequest)
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if loginData.Email == "" || loginData.Password == "" {
		//http.Error(w, "all field are required", http.StatusBadRequest)
		utils.RespondWithError(w, http.StatusBadRequest, "all field is required")
		return
	}

	collection := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	err = collection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil {
		//http.Error(w, "inavalid email or password", http.StatusUnauthorized)
		utils.RespondWithError(w, http.StatusBadRequest, "invalid email or password")
		return

	}

	err = utils.Checkpassword(user.Password, loginData.Password)
	if err != nil {
		//http.Error(w, "invalid email or password", http.StatusUnauthorized)
		utils.RespondWithError(w, http.StatusBadRequest, "invalid email or password")

		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.Hex(), user.IsAdmin)
	if err != nil {
		//http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// ✅ Return token to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})

	/**

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

	**/

}
