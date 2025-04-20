package middleware

import (
	// "net/http"
	// "strings"

	"context"
	"goth-todo/internal/auth"
	"net/http"

	"github.com/gin-gonic/gin"

	"goth-todo/internal/auth"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {
	// return func(c *gin.Context) {
	// 	token := auth.GetUserToken(c)
	// 	if token == nil {
	// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
	// 		return
	// 	}

	// 	c.Next()
	// }
	return func(c *gin.Context) {
		cookie, err := c.Cookie("jwt_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims := &auth.CustomClaims{} // ✅ use your actual claims type

		token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "account_id", claims.AccountId) // ✅ matches your struct
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
