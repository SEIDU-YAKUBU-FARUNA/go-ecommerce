/**

package main

import (
	"fmt"
	"go-ecommerce/config"
	"go-ecommerce/database"
	"go-ecommerce/handlers"
	"go-ecommerce/middleware"
	"net/http"
	"os"
)

func main() {
	config.LoadEnv()

	database.ConnectDB()

	//Auth route/pulic route
	http.HandleFunc("/register", handlers.RegisterUser)
	http.HandleFunc("/login", middleware.RateLimiter(handlers.LoginUser))

	//customers roles/pulic route
	http.HandleFunc("/products", middleware.AuthMiddleware(handlers.GetProducts, false))
	http.HandleFunc("/product", middleware.AuthMiddleware(handlers.GetProduct, false))
	http.HandleFunc("/product/getorders", middleware.AuthMiddleware(handlers.GetOrders, false))

	//Admin roles/private route
	http.HandleFunc("/product/add", middleware.AuthMiddleware(handlers.AddProduct, true))
	http.HandleFunc("/product/update", middleware.AuthMiddleware(handlers.UpdateProduct, true))
	http.HandleFunc("/product/delete", middleware.AuthMiddleware(handlers.DeleteProduct, true))
	http.HandleFunc("/product/creatorder", middleware.AuthMiddleware(handlers.CreateOrder, true))

	//log.Fatal(http.ListenAndServe(":8080", nil))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port:", port)
	http.ListenAndServe(":"+port, nil)

}

**/

package main

import (
	"fmt"
	"go-ecommerce/config"
	"go-ecommerce/database"
	"go-ecommerce/handlers"
	"go-ecommerce/middleware"
	"net/http"
	"os"
)

func main() {

	config.LoadEnv()
	database.ConnectDB()

	// Create router
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/register", handlers.RegisterUser)
	mux.HandleFunc("/login", middleware.RateLimiter(handlers.LoginUser))

	// Customer routes
	mux.HandleFunc("/products", middleware.AuthMiddleware(handlers.GetProducts, false))
	mux.HandleFunc("/product", middleware.AuthMiddleware(handlers.GetProduct, false))
	mux.HandleFunc("/product/getorders", middleware.AuthMiddleware(handlers.GetOrders, false))

	// Admin routes
	mux.HandleFunc("/product/add", middleware.AuthMiddleware(handlers.AddProduct, true))
	mux.HandleFunc("/product/update", middleware.AuthMiddleware(handlers.UpdateProduct, true))
	mux.HandleFunc("/product/delete", middleware.AuthMiddleware(handlers.DeleteProduct, true))
	mux.HandleFunc("/product/creatorder", middleware.AuthMiddleware(handlers.CreateOrder, true))

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port:", port)

	// Apply CORS middleware
	handler := middleware.EnableCORS(mux)

	http.ListenAndServe(":"+port, handler)
}
