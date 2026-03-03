package handlers

import (
	"context"
	"encoding/json"
	"go-ecommerce/database"
	"go-ecommerce/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}

	var request struct {
		User     string   `json:"user"`
		Products []string `json:"products"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body ", http.StatusBadRequest)
		return
	}

	if request.User == "" || len(request.Products) == 0 {
		http.Error(w, "user and products required ", http.StatusBadRequest)
		return
	}

	userID, err := primitive.ObjectIDFromHex(request.User)
	if err != nil {
		http.Error(w, " invalid ID", http.StatusBadRequest)
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
			http.Error(w, "invalid product id ", http.StatusBadGateway)
		}

		var product models.Product
		err = productsCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&product)
		if err != nil {
			http.Error(w, "product not found", http.StatusInternalServerError)
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
		http.Error(w, " fialed to create order", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{

		"message": "oder created successfully",
		"total":   total,
	})

}

func GetOrders(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {

		http.Error(w, " method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("X-Admin") == "true" {
		http.Error(w, " admin access required", http.StatusForbidden)
		return
	}

	collection := database.DB.Collection("oders")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "failed to fetch odrders", http.StatusInternalServerError)
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
