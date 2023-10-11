package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHTMLInfoV1 show the information about the csv balance file in html format
func GetHTMLInfoV1(htmlProcessTransactions HTMLProcessTransactions) gin.HandlerFunc {
	return func(c *gin.Context) {
		html, err := htmlProcessTransactions(c)
		if err != nil {
			WebError(c, http.StatusInternalServerError, CantGetInfo)
		}

		c.Writer.WriteHeader(http.StatusOK)
		_, err = c.Writer.Write(html)
		if err != nil {
			WebError(c, http.StatusInternalServerError, CantWriteHtml)
		}
	}
}
