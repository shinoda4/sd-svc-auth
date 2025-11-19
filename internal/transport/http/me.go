package http

import (
	stdhttp "net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HandleMe(c *gin.Context) {
	claims, _ := c.Get("claims")
	c.JSON(stdhttp.StatusOK, gin.H{"claims": claims})
}
