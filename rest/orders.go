package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *API) GetOrder(c *gin.Context) {
	uid := c.Param("order_uid")
	order, err := h.db.GetOrder(uid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, order)
}
