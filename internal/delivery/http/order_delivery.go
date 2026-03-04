package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rickyhuang08/mini-exchange.git/helpers"
	"github.com/rickyhuang08/mini-exchange.git/internal/domain"
	"github.com/rickyhuang08/mini-exchange.git/internal/entity"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	var req entity.OrderRequest
	response := entity.APIResponse{
		Status: helpers.Success,
		Message: helpers.OrderPlaced,
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Status = helpers.Error
		response.Message = helpers.InvalidRequest
		response.Error = err.Error()

		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	order := &domain.Order{
		ID: uuid.New().String(),
		StockCode: req.StockCode,
		Side:      domain.Side(req.Side),
		Price:     req.Price,
		Quantity:  req.Quantity,
		Status:    domain.StatusNew,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.OrderUsecase.PlaceOrder(order); err != nil {
		response.Status = helpers.Error
		response.Message = helpers.OrderPlacementFailed
		response.Error = err.Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = gin.H{"order_id": order.ID}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetOrders(c *gin.Context) {
	stockCode := c.Query("stock_code")
	status := c.Query("status")

	filter := make(map[string]interface{})
	if stockCode != "" {
		filter["stock_code"] = stockCode
	}
	if status != "" {
		filter["status"] = domain.OrderStatus(status)
	}

	response := entity.APIResponse{
		Status: helpers.Success,
		Message: helpers.OrderListRetrieved,
	}

	orders, err := h.OrderUsecase.ListOrders(filter)
	if err != nil {
		response.Status = helpers.Error
		response.Message = helpers.OrderListRetrievalFailed
		response.Error = err.Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = orders
	c.JSON(http.StatusOK, response)
}


func (h *Handler) GetTrades(c *gin.Context) {
	stockCode := c.Query("stock_code")
	response := entity.APIResponse{
		Status: helpers.Success,
		Message: helpers.TradeListRetrieved,
	}

	if stockCode == "" {
		response.Status = helpers.Error
		response.Message = helpers.InvalidRequest
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	trades, err := h.OrderUsecase.GetTradeHistory(stockCode)
	if err != nil {
		response.Status = helpers.Error
		response.Message = helpers.TradeListRetrievalFailed
		response.Error = err.Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	
	response.Data = trades
	c.JSON(http.StatusOK, response)
}