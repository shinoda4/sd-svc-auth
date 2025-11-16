package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HandleMe(c *gin.Context) {
	claims, _ := c.Get("claims")
	c.JSON(http.StatusOK, gin.H{"claims": claims})
}
