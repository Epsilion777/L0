package rest

import (
	"L0/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllOrders(c *gin.Context) {
	data := gin.H{
		"OrderCache": model.OrderCache,
	}
	// Rendering a template with data
	c.HTML(http.StatusOK, "index.tmpl", data)
}
