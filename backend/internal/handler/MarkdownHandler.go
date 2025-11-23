package handler

import (
	"DevDesk/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MarkdownHandler struct {
	md *service.Markdown
}

func NewMarkdownHandler(md *service.Markdown) *MarkdownHandler {
	return &MarkdownHandler{md: md}
}

func (h *MarkdownHandler) NewDocument(c *gin.Context) {
	doc, err := h.md.NewDocument()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"hash": doc.Hash})
}

func (h *MarkdownHandler) GetDocument(c *gin.Context) {
	hash := c.Param("hash")
	doc, ok := h.md.GetDocument(hash)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": service.ErrDocNotFound.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": doc.Content})
}

func (h *MarkdownHandler) UpdateDocument(c *gin.Context) {
	var req struct {
		Hash    string `json:"hash"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hash is required"})
		return
	}

	_, err := h.md.UpdateDocument(req.Hash, req.Content)
	if err != nil {
		if err == service.ErrDocNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *MarkdownHandler) StreamDocument(c *gin.Context) {
	hash := c.Param("hash")
	doc, ok := h.md.GetDocument(hash)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": service.ErrDocNotFound.Error()})
		return
	}
	// 获取底层 ResponseWriter
	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming unsupported"})
		return
	}

	// 注册客户端通道
	ch := doc.AddClient()
	defer doc.RemoveClient(ch)

	// 设置 SSE 相关头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher.Flush()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case content, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(map[string]string{"content": content})
			_, _ = w.Write([]byte("data: " + string(data) + "\n\n"))
			flusher.Flush()
		}
	}

}
