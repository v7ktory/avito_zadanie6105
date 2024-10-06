package v1

import (
	"net/http"

	"git.codenrock.com/tender/internal/service"
)

func NewRouter(mux *http.ServeMux, services *service.Services) {
	mux.Handle("GET /api/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))

	tenderRouter := newTenderRouter(services.Tender, services)
	bidRouter := newbidRouter(services.Bid, services)

	mux.Handle("/api/tenders/", tenderRouter)
	mux.Handle("/api/bids/", bidRouter)

}
