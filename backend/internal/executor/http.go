package executor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	Runner *Runner
}

func NewHTTPServer(r *Runner) *HTTPServer {
	return &HTTPServer{Runner: r}
}

func (s *HTTPServer) Register(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.POST("/execute", func(c *gin.Context) {
		var req ExecuteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp := s.Runner.Execute(c.Request.Context(), req)
		code := http.StatusOK
		if !resp.OK && resp.Error == "timeout" {
			code = http.StatusRequestTimeout
		}
		c.JSON(code, resp)
	})
}
