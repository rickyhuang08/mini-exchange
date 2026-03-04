package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rickyhuang08/mini-exchange.git/helpers"
	"github.com/rickyhuang08/mini-exchange.git/internal/entity"
)

func (h *Handler) LoginHandler(c *gin.Context) {
	var request entity.LoginRequest
	response := entity.APIResponse{
		Status: helpers.Success,
		Message: helpers.LoginSuccess,
	}

	if err := c.ShouldBind(&request); err != nil {
		response.Status = helpers.Error
		response.Message = helpers.InvalidRequest
		response.Error = err.Error()

		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	loginResponse, err := h.AuthUsecase.Login(request)
	if err != nil {
		response.Status = helpers.Error
		response.Message = helpers.InvalidCredentials
		response.Error = err.Error()

		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response.Data = loginResponse

	c.JSON(http.StatusOK, response)
}