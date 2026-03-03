package main

import (
	"fmt"
	"go-ecommerce/database"
	"go-ecommerce/handlers"
	"net/http"
	"os"
)

func main() {

	database.ConnectDB()

	//Auth route/pulic route
	http.HandleFunc("/register", handlers.RegisterUser)
	http.HandleFunc("/login", handlers.LoginUser)

	//customers roles/pulic route
	http.HandleFunc("/products", handlers.GetProducts)
	http.HandleFunc("/product", handlers.GetProduct)
	http.HandleFunc("/product/getorders", handlers.GetOrders)

	//Admin roles/private route
	http.HandleFunc("/product/add", handlers.AddProduct)
	http.HandleFunc("/product/update", handlers.UpdateProduct)
	http.HandleFunc("/product/delete", handlers.DeleteProduct)
	http.HandleFunc("/product/creatorder", handlers.CreateOrder)

	//log.Fatal(http.ListenAndServe(":8080", nil))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port:", port)
	http.ListenAndServe(":"+port, nil)

}
