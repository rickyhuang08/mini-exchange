package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rickyhuang08/mini-exchange.git/pkg/jwt"
	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)

type MiddlewareModule struct {
	PublicKeyPath string
	Logger        *logger.Logger
}

func NewMiddlewareModule(publicKeyPath string, logger *logger.Logger) *MiddlewareModule {
	return &MiddlewareModule{
		PublicKeyPath: publicKeyPath,
		Logger:        logger,
	}
}

func (m *MiddlewareModule) InitAuthMiddleware() *AuthModule {
	return NewAuthModule()
}

func (m *MiddlewareModule) RegisterGlobalMiddleware(r *gin.Engine) {
	r.Use(m.LoggerMiddleware())
	r.Use(CORSMiddleware())
}

// RegisterAuthMiddleware applies to protected routes
func (m *MiddlewareModule) RegisterAuthMiddleware(r *gin.RouterGroup) error {
	loadPublicKey, err := jwt.LoadPublicKey(m.PublicKeyPath)
	if err != nil {
		return err
	}
	r.Use(m.LoggerMiddleware())
	r.Use(m.InitAuthMiddleware().AuthMiddleware(loadPublicKey)) // Requires JWT authentication

	return nil
}

func (m *MiddlewareModule) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		latency := end.Sub(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		msg := fmt.Sprintf("%s - %s %s [%d] %v", clientIP, method, path, statusCode, latency)
		m.Logger.LogAccess(msg)
	}
}