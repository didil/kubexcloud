package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// Cors middleware
func Cors(next http.Handler) http.Handler {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})

	return c.Handler(next)
}
