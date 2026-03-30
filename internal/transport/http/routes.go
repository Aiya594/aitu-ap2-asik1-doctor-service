package httpdoc

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *DocHandler) {
	r.POST("/doctors", h.CreateDoctor)
	r.GET("/doctors/:id", h.GetDoctor)
	r.GET("/doctors", h.GetAll)
}
