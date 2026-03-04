package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)

type JWTHelper struct {
	Logger *logger.Logger
	PrivateKey string
	Epxiration int
}

func NewJWTHelper(logger *logger.Logger, privateKey string, exp int) *JWTHelper {
	return &JWTHelper{
		Logger:     logger,
		PrivateKey: privateKey,
		Epxiration: exp,
	}
}

// GenerateJWT creates a signed JWT
func (m *JWTHelper) GenerateJWT(userID, role int, email string) (string, error) {
	m.Logger.LogLevel(logger.LogLevelInfo, "Generating JWT for user: "+email)
	// Initialize usecases
	rsaPrivateKey, err := LoadPrivateKey(m.PrivateKey)

	if err != nil {
		m.Logger.LogLevel(logger.LogLevelError, fmt.Sprintf("Failed to load private key: %v", err))
		return "", err
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"email": email,
		"role": role,
		"exp": time.Now().Add(time.Hour * time.Duration(m.Epxiration)).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(rsaPrivateKey)
}