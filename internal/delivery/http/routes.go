package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rickyhuang08/mini-exchange.git/middleware"
)

func RegisterRoutes(router *gin.Engine, handler *Handler, mw *middleware.MiddlewareModule) {
	// Public routes
	mw.RegisterGlobalMiddleware(router)
	public := router.Group("/api/public")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
		public.POST("/login", handler.LoginHandler)

		// This endpoint is intentionally left public for simplicity, but in a real application, it should be protected with authentication and authorization checks.
		public.POST("/orders", handler.CreateOrder)
		public.GET("/orders", handler.GetOrders)
		public.GET("/trades", handler.GetTrades)
		public.GET("/market/:stock/snapshot", handler.GetMarketSnapshot)
	}

	// WebSocket route
	realtime := router.Group("/api/realtime")
	{
		realtime.GET("/ws", handler.WebSocketHandler)
	}

	// Protected routes
	protected := router.Group("/api/protected")
	{
		err := mw.RegisterAuthMiddleware(protected)
		if err != nil {
			panic("Failed to register auth middleware: " + err.Error())
		}

		// Add protected routes here (e.g., user profile, order management)
		protected.POST("/orders", handler.CreateOrder)
		protected.GET("/orders", handler.GetOrders)
		protected.GET("/trades", handler.GetTrades)
		protected.GET("/market/:stock/snapshot", handler.GetMarketSnapshot)
	}
}