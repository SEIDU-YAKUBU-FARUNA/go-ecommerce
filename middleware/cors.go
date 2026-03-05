package middleware

import "net/http"

// EnableCORS is a middleware function that allows your frontend
// application (like React, Vue, Next.js, etc.) to communicate with
// your Go backend API when they are hosted on different domains.
func EnableCORS(next http.Handler) http.Handler {

	// http.HandlerFunc converts a function into an HTTP handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow requests only from your frontend domain
		// This prevents unknown websites from calling your API
		w.Header().Set("Access-Control-Allow-Origin", "https://mini-s.netlify.app")

		// Allow these HTTP methods when the frontend makes requests
		// Example:
		// GET    -> Fetch products
		// POST   -> Create product
		// PUT    -> Update product
		// DELETE -> Delete product
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow these headers in requests from the frontend
		// Content-Type -> for JSON data
		// Authorization -> for JWT tokens
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Browsers send a "preflight" request using OPTIONS
		// before making the real request to check permissions
		// If the request method is OPTIONS, we stop here
		// and return the headers above.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)

			return
		}

		// If it is not an OPTIONS request,
		// continue to the actual API handler
		next.ServeHTTP(w, r)
	})
}
