package handler

import (
	"net/http"

	"DevDesk/internal/service"

	"github.com/gin-gonic/gin"
)

type HttpTestHandler struct {
	svc *service.HttpTestService
}

func NewHttpTestHandler(svc *service.HttpTestService) *HttpTestHandler {
	return &HttpTestHandler{svc: svc}
}

func (h *HttpTestHandler) Do(c *gin.Context) {
	var req service.HttpTestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	resp, err := h.svc.Do(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
