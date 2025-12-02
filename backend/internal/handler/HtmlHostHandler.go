package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"DevDesk/internal/service"

	"github.com/gin-gonic/gin"
)

type HtmlHostHandler struct {
	svc          *service.HTMLHostService
	publicPrefix string
}

func NewHtmlHostHandler(svc *service.HTMLHostService, publicPrefix string) *HtmlHostHandler {
	if publicPrefix == "" {
		publicPrefix = "/hosted"
	}
	return &HtmlHostHandler{
		svc:          svc,
		publicPrefix: publicPrefix,
	}
}

// POST /html/upload
func (h *HtmlHostHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	filename, err := h.svc.SaveHTML(fileHeader)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrHTMLTooLarge) {
			status = http.StatusRequestEntityTooLarge
		}
		c.JSON(status, gin.H{
			"error":       err.Error(),
			"limit_bytes": h.svc.MaxSizeBytes(),
		})
		return
	}

	relative := strings.TrimRight(h.publicPrefix, "/") + "/" + filename
	sharePath := "/api" + relative

	scheme := "http"
	if c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	fullURL := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, sharePath)

	c.JSON(http.StatusOK, gin.H{
		"url":         sharePath,
		"full_url":    fullURL,
		"filename":    filename,
		"limit_bytes": h.svc.MaxSizeBytes(),
	})
}
