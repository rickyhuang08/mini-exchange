package middleware

import (
	"context"
	"crypto/rsa"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rickyhuang08/mini-exchange.git/helpers"
	"github.com/rickyhuang08/mini-exchange.git/internal/entity"
	pkg_jwt "github.com/rickyhuang08/mini-exchange.git/pkg/jwt"
)


type AuthModule struct{}

func NewAuthModule() *AuthModule {
	return &AuthModule{}
}

// AuthMiddleware verifies JWT tokens
func (m *AuthModule) AuthMiddleware(publicKey *rsa.PublicKey) gin.HandlerFunc {
	response := entity.APIResponse {
		Status: helpers.Success,
	}
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Status = helpers.Error
			response.Message = helpers.AuthenticationFailed
			response.Error = helpers.MissingToken

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.ValidateJWT(tokenString, publicKey)
		if err != nil {
			response.Status = helpers.Error
			response.Message = helpers.ValidateJWTFailed
			response.Error = err.Error()

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Store user info in context
		ctx := context.WithValue(c.Request.Context(), pkg_jwt.UserKey, claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// ValidateJWT validates a JWT token
func (m *AuthModule) ValidateJWT(tokenString string, publicKey *rsa.PublicKey) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, jwt.ErrTokenExpired
		}
	}

	return claims, nil
}