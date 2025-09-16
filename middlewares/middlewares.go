package middlewares

import (
	"net/http"
	"strings"
	"os"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)



func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secretKey := os.Getenv("SUPABASE_JWT_KEY")

		// The validation function is now more flexible
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// This dynamically checks the algorithm from the token itself
			// instead of forcing a single method.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			// Add a log to see the actual validation error
			log.Println("JWT validation error:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}


		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to process token claims"})
			return
		}

		// 2. (FIXED) User ID from Supabase is a string (UUID)
		userID, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		// 3. (SIMPLIFIED) Set only the user ID in the context
		c.Set("userID", userID)
		c.Next()
	}
}
