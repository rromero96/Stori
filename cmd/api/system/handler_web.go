package system

import (
	"net/http"

	"github.com/rromero96/roro-lib/cmd/web"
)

// GetInfoV1 show the information about the csv balance file
func GetInfoV1(processTransactions ProcessTransactions) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		info, err := processTransactions(r.Context())
		if err != nil {
			return web.NewError(http.StatusInternalServerError, CantGetInfo)
		}

		return web.EncodeJSON(w, info.toDTO(), http.StatusOK)
	}
}

// GetHTMLInfoV1 show the information about the csv balance file in html format
func GetHTMLInfoV1(htmlProcessTransactions HTMLProcessTransactions) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "text/html")

		html, err := htmlProcessTransactions(r.Context())
		if err != nil {
			return web.NewError(http.StatusInternalServerError, CantGetInfo)

		}

		_, err = w.Write(html)
		if err != nil {
			return web.NewError(http.StatusInternalServerError, CantWriteHtml)
		}

		return web.EncodeJSON(w, nil, http.StatusOK)
	}
}
