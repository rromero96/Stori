package system_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rromero96/stori/cmd/api/system"
	"github.com/rromero96/stori/internal/web"
	"github.com/stretchr/testify/assert"
)

func TestHTTPHandler_GetInfoV1_success(t *testing.T) {
	processTransaction := system.MockProcessTransactions(system.MockEmail(), nil)
	getInfoV1 := system.GetInfoV1(processTransaction)

	ctx, w := context.Background(), httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/test", strings.NewReader(""))

	got := getInfoV1(w, r)

	assert.Nil(t, got)
}

func TestHTTPHandler_GetInfoV1_fails(t *testing.T) {
	processTransaction := system.MockProcessTransactions(system.Email{}, system.ErrCantGetCsvFile)
	getInfoV1 := system.GetInfoV1(processTransaction)

	ctx, w := context.Background(), httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/test", strings.NewReader(""))

	want := web.NewError(http.StatusInternalServerError, system.CantGetInfo)
	got := getInfoV1(w, r)

	assert.Equal(t, got, want)
}

func TestHTTPHandler_GetHTMLInfoV1_success(t *testing.T) {
	processTransaction := system.MockHTMLProcessTransactions([]byte{}, nil)
	getHTMLInfoV1 := system.GetHTMLInfoV1(processTransaction)

	ctx, w := context.Background(), httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/test", strings.NewReader(""))

	got := getHTMLInfoV1(w, r)

	assert.Nil(t, got)
}

func TestHTTPHandler_GetHTMLInfoV1_fails(t *testing.T) {
	processTransaction := system.MockHTMLProcessTransactions([]byte{}, system.ErrCantGetTransactionInfo)
	getHTMLInfoV1 := system.GetHTMLInfoV1(processTransaction)

	ctx, w := context.Background(), httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/test", strings.NewReader(""))

	want := web.NewError(http.StatusInternalServerError, system.CantGetInfo)
	got := getHTMLInfoV1(w, r)

	assert.Equal(t, got, want)
}

func TestHTTPHandler_GetLogoV1_success(t *testing.T) {
	getLogoV1 := system.GetLogoV1()

	ctx, w := context.Background(), httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/test", strings.NewReader(""))

	got := getLogoV1(w, r)

	assert.Nil(t, got)
}
