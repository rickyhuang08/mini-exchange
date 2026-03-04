package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rickyhuang08/mini-exchange.git/helpers"
	"github.com/rickyhuang08/mini-exchange.git/internal/entity"
)

func (h *Handler) GetMarketSnapshot(c *gin.Context) {
	stock := c.Param("stock")
	response := entity.APIResponse{
		Status: helpers.Success,
		Message: helpers.MarketSnapshotRetrieved,
	}

	snapshot, err := h.MarketUsecase.GetSnapshot(stock)
	if err != nil {
		response.Status = helpers.Error
		response.Message = helpers.MarketSnapshotRetrievalFailed
		response.Error = err.Error()

		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = snapshot
	c.JSON(http.StatusOK, response)
}