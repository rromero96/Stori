package system_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/rromero96/stori/cmd/api/system"
)

func TestHTTPHandler_GetHTMLInfoV1_success(t *testing.T) {
	processTransaction := system.MockHTMLProcessTransactions([]byte{}, nil)
	getHTMLInfoV1 := system.GetHTMLInfoV1(processTransaction)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getHTMLInfoV1(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHTTPHandler_GetHTMLInfoV1_fails(t *testing.T) {
	processTransaction := system.MockHTMLProcessTransactions([]byte{}, system.ErrCantGetTransactionInfo)
	getHTMLInfoV1 := system.GetHTMLInfoV1(processTransaction)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getHTMLInfoV1(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
