package rest

import (
	"L0/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetOrderByID(c *gin.Context) {
	orderUID := c.Param("id")
	data := gin.H{
		"Order": model.OrderCache[orderUID],
	}
	// Rendering a template with data
	c.HTML(http.StatusOK, "order.tmpl", data)
}
