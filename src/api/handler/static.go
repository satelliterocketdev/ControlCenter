package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
)

func (ctrl *Controller) static(c *gin.Context) {
	rootDir, _ := os.Getwd()
	filePath := path.Join(rootDir, "src/api/static", "index.html")
	// check if file exist
	if _, errFileExist := os.Stat(filePath); os.IsNotExist(errFileExist) {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	c.Writer.Header().Set("Pragma", "no-cache") // HTTP 1.0.
	c.Writer.Header().Set("Expires", "0") // Proxies.

	http.ServeFile(c.Writer, c.Request, filePath)
	return
}
