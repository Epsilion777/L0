package rest

import (
	"L0/internal/usecase/interfaces"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	orderUsecase interfaces.OrderUsecase
}

func NewHandler(orderUsecase interfaces.OrderUsecase, router *gin.Engine) *Handler {
	h := &Handler{
		orderUsecase: orderUsecase,
	}
	// Loading templates
	router.LoadHTMLGlob("internal/templates/*")

	router.GET("/", h.GetAllOrders)
	router.GET("/orders/:id", h.GetOrderByID)
	return h
}
