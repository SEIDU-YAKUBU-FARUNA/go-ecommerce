/**
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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		//http.Error(w, "invalid method", http.StatusBadRequest)
		utils.RespondWithError(w, http.StatusBadRequest, "invalid method")
		return

	}

	var request struct {
		User     string   `json:"user"`
		Products []string `json:"products"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		//http.Error(w, "invalid request body ", http.StatusBadRequest)
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if request.User == "" || len(request.Products) == 0 {
		//http.Error(w, "user and products required ", http.StatusBadRequest)
		utils.RespondWithError(w, http.StatusBadRequest, "failed to create order")
		return
	}

	userID, err := primitive.ObjectIDFromHex(request.User)
	if err != nil {
		//http.Error(w, " invalid ID", http.StatusBadRequest)
		utils.RespondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}

	productsCollection := database.DB.Collection("products")
	oderCollection := database.DB.Collection("oders")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var productIDs []primitive.ObjectID
	var total float64 = 0

	for _, pid := range request.Products {

		//find the product id
		objID, err := primitive.ObjectIDFromHex(pid)
		if err != nil {
			//http.Error(w, "invalid product id ", http.StatusBadGateway)
			utils.RespondWithError(w, http.StatusBadRequest, "invalid product id")
			return
		}

		var product models.Product
		err = productsCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&product)
		if err != nil {
			//http.Error(w, "product not found", http.StatusInternalServerError)
			utils.RespondWithError(w, http.StatusInternalServerError, "product not found")
			return
		}

		//add price to total

		total += product.Price
		productIDs = append(productIDs, objID)

	}

	order := models.Order{
		UserID:      userID,
		Product:     productIDs,
		TotalAmount: total,
	}

	_, err = oderCollection.InsertOne(ctx, order)
	if err != nil {
		//http.Error(w, " fialed to create order", http.StatusInternalServerError)
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{

		"message": "oder created successfully",
		"total":   total,
	})

}

func GetOrders(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {

		//http.Error(w, " method not allowed", http.StatusMethodNotAllowed)
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if r.Header.Get("X-Admin") == "true" {
		//http.Error(w, " admin access required", http.StatusForbidden)
		utils.RespondWithError(w, http.StatusForbidden, "admin access required")
		return
	}

	collection := database.DB.Collection("oders")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		//http.Error(w, "failed to fetch odrders", http.StatusInternalServerError)
		utils.RespondWithError(w, http.StatusInternalServerError, "all field is required")

		return
	}
	defer cursor.Close(ctx)

	var orders []models.Order

	for cursor.Next(ctx) {
		var order models.Order
		cursor.Decode(&order)
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "applications/json")
	json.NewEncoder(w).Encode(orders)

}

**/

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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST is allowed")
		return
	}

	var request struct {
		User     string   `json:"user"`
		Products []string `json:"products"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if request.User == "" || len(request.Products) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "user ID and products are required")
		return
	}

	userID, err := primitive.ObjectIDFromHex(request.User)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	productsCollection := database.DB.Collection("products")
	orderCollection := database.DB.Collection("orders") // Fixed typo: oders -> orders

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var productIDs []primitive.ObjectID
	var total float64 = 0

	for _, pid := range request.Products {
		objID, err := primitive.ObjectIDFromHex(pid)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid product id: "+pid)
			return
		}

		var product models.Product
		// 🟢 FIX: MongoDB uses "_id", not "id"
		err = productsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "product not found: "+pid)
			return
		}

		total += product.Price
		productIDs = append(productIDs, objID)
	}

	// 🟢 FIX: Matching your models.Order struct field exactly
	order := models.Order{
		UserID:      userID,
		product:     productIDs, // Use 'Product' to match your models.go exactly
		TotalAmount: total,
	}

	_, err = orderCollection.InsertOne(ctx, order)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "order created successfully",
		"total":   total,
	})
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Note: You should ideally use your AuthMiddleware here instead of manual header checks
	isAdmin := r.Header.Get("X-Admin")
	if isAdmin != "true" {
		utils.RespondWithError(w, http.StatusForbidden, "admin access required")
		return
	}

	collection := database.DB.Collection("orders") // Fixed typo: oders -> orders
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch orders")
		return
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to decode orders")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
