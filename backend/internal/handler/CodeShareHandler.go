package handler

import (
	"net/http"

	"moment/internal/service"

	"github.com/gin-gonic/gin"
)

type CodeShareHandler struct {
	cs *service.CodeShare
}

func NewCodeShareHandler(cs *service.CodeShare) *CodeShareHandler {
	return &CodeShareHandler{cs: cs}
}

// POST /codeshare/upload
func (h *CodeShareHandler) Upload(c *gin.Context) {
	var req uploadRequest {
		Author   string `json:"author"   binding:"required"`
		Language string `json:"language" binding:"required"`
		Content  string `json:"content"  binding:"required"`
		TTL      int64  `json:"ttl"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	if req.TTL == 0 {
		req.TTL = 3600
	}

	code := h.cs.Upload(req.Author, req.Language, req.Content, req.TTL)
	if code == nil {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "Content too large",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"hash":    code.Hash,
		"url":     "/codeshare/code/" + code.Hash,
	})
}

// GET /codeshare/code/:hash
func (h *CodeShareHandler) Get(c *gin.Context) {
	hash := c.Param("hash")

	code, ok := h.cs.Get(hash)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Code not found or expired",
		})
		return
	}

	c.JSON(http.StatusOK, code)
}
