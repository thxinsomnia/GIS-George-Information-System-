package middlewares

import (
	"net/http"
	"strings"
	"GIS/config"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTVerif() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Message": "Tidak Terverifikasi!! Harap Login Terlebih dahulu!!"})
			c.Abort()
			return
		}

		claims := &config.JWTClaims{}

		//parsing token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				c.JSON(http.StatusUnauthorized, gin.H{"Message": "Tidak Terverifikasi!! Harap Login Terlebih dahulu!!"})
				c.Abort()
				return

			case jwt.ValidationErrorExpired:
				c.JSON(http.StatusUnauthorized, gin.H{"Message": "Silahkan Login Ulang Sesi Sudah Kadaluarsa!!"})
				c.Abort()
				return

			default:
				c.JSON(http.StatusUnauthorized, gin.H{"Message": "Tidak Terverifikasi!! Harap Login Terlebih dahulu!!"})
				c.Abort()
				return
			}
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"Message": "Tidak Terverifikasi!! Harap Login Terlebih dahulu!!"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// 1. (FIXED) Use the Supabase secret key for validation
		secretKey := os.Getenv("SUPABASE_KEY")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
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
