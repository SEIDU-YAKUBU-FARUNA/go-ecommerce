package middleware

import "net/http"

// CORSMiddleware allows browsers (React, Vue, HTML, Next.js)
// to access your backend API.
func CORSMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow all origins (for development)
		w.Header().Set("Access-Control-Allow-Origin", "https://mini-s.netlify.app/")

		// Allow specific request methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow headers sent by frontend
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle browser preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
