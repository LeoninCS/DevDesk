package handler

import (
	"net/http"

	"moment/internal/service"

	"github.com/gin-gonic/gin"
)

func NewHandler() http.Handler {
	r := gin.Default()

	// ===== CORS 中间件（前后端联调必须要有）=====
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// 你也可以写死成 http://localhost:5173
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}

		// 预检请求直接返回
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	// =======================================

	// CodeShare 分组
	svc := service.NewCodeShareService()
	codeHandler := NewCodeShareHandler(svc)
	cg := r.Group("/codeshare")
	{
		cg.POST("/upload", codeHandler.Upload)
		cg.GET("/code/:hash", codeHandler.Get)
	}

	// WorkPlan 分组
	workPlanHandler := NewWorkPlanHandler()
	wg := r.Group("/workplan")
	{
		wg.GET("/new", workPlanHandler.NewPersonlPlan)
		wg.POST("/add", workPlanHandler.AddTODO)
		wg.POST("/delete", workPlanHandler.DeleteTODO)
		wg.POST("/edit", workPlanHandler.EditTODO)
		wg.POST("/done", workPlanHandler.SetDone)
		wg.GET("/:hash", workPlanHandler.GetTODOs)
	}

	return r
}
