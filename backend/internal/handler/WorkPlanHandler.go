package handler

import (
	"net/http"

	"DevDesk/internal/service"

	"github.com/gin-gonic/gin"
)

type WorkPlanHandler struct {
	wp *service.WorkPlan
}

func NewWorkPlanHandler(srv *service.WorkPlan) *WorkPlanHandler {
	return &WorkPlanHandler{
		wp: srv,
	}
}

// 创建新的个人计划
func (h *WorkPlanHandler) NewPersonlPlan(c *gin.Context) {
	pp := h.wp.NewPersonalPlan()
	c.JSON(http.StatusOK, gin.H{
		"hash": pp.Hash,
	})
}

// 添加 TODO
func (h *WorkPlanHandler) AddTODO(c *gin.Context) {
	var req struct {
		Hash    string `json:"hash"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pp := h.wp.GetPlan(req.Hash)
	if pp == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan not found"})
		return
	}

	if err := pp.AddTODO(req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// 删除 TODO
func (h *WorkPlanHandler) DeleteTODO(c *gin.Context) {
	var req struct {
		Hash string `json:"hash"`
		Id   int    `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pp := h.wp.GetPlan(req.Hash)
	if pp == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan not found"})
		return
	}

	if err := pp.DeleteTODO(req.Id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// 修改 TODO 内容
func (h *WorkPlanHandler) EditTODO(c *gin.Context) {
	var req struct {
		Hash    string `json:"hash"`
		Id      int    `json:"id"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pp := h.wp.GetPlan(req.Hash)
	if pp == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan not found"})
		return
	}

	if err := pp.EditTODO(req.Id, req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// 设置 TODO 完成 / 未完成
func (h *WorkPlanHandler) SetDone(c *gin.Context) {
	var req struct {
		Hash string `json:"hash"`
		Id   int    `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pp := h.wp.GetPlan(req.Hash)
	if pp == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan not found"})
		return
	}

	if err := pp.SetTODODone(req.Id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// 获取 TODO 列表
func (h *WorkPlanHandler) GetTODOs(c *gin.Context) {
	hash := c.Param("hash")

	pp := h.wp.GetPlan(hash)
	if pp == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"todos": pp.GetTODOs(),
	})
}
