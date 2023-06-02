package system

import (
	"net/http"

	"github.com/rromero96/stori/internal/web"
)

func GetInfoV1(processTransactions ProcessTransactions) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		info, err := processTransactions(r.Context())
		if err != nil {
			return web.NewError(http.StatusInternalServerError, CantGetInfo)
		}

		return web.EncodeJSON(w, info.toDTO(), http.StatusOK)
	}
}

func GetHTMLInfoV1() web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return web.EncodeJSON(w, "Hello World!", http.StatusOK)
	}
}
