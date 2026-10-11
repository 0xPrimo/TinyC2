package listener

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateHTTPHandler(listener *HTTPListener) func(*gin.Context) {
	return func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
			return
		}

		response, err := listener.Engine.ImplantProcess(listener.Name, body)
		if err != nil {
			c.JSON(200, gin.H{"message": "processing error"})
			return
		}

		c.Data(http.StatusOK, "application/json", response)
	}
}
