package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"git.codenrock.com/tender/internal/model"
	"git.codenrock.com/tender/internal/repo/repoerrs"
	"git.codenrock.com/tender/internal/service"
	"github.com/google/uuid"
)

type bidRouter struct {
	bidService service.Bid
}

func newbidRouter(bidService service.Bid, services *service.Services) http.Handler {
	r := &bidRouter{
		bidService: bidService,
	}

	mux := http.NewServeMux()

	createBidMiddleware := createBidMiddleware(services)
	authorOrResponsibleMiddleware := authorOrResponsibleMiddleware(services)
	accessMiddleware := accessMiddleware(services)

	mux.Handle("POST /new", createBidMiddleware(http.HandlerFunc(r.createBidHandler())))
	mux.Handle("GET /my", r.GetUserBidsHandler(services.Employee))
	mux.Handle("GET /{tenderId}/list", r.getBidsByTenderHandler(services.Employee))
	mux.Handle("GET /{bidId}/status", authorOrResponsibleMiddleware(http.HandlerFunc(r.getBidStatusHandler())))
	mux.Handle("PUT /{bidId}/status", authorOrResponsibleMiddleware(http.HandlerFunc(r.updateBidStatusHandler())))
	mux.Handle("PATCH /{bidId}/edit", authorOrResponsibleMiddleware(http.HandlerFunc(r.updateBidHandler())))

	mux.Handle("PUT /{bidId}/submit_decision", accessMiddleware(http.HandlerFunc(r.updateBidDecisionHandler(services.Tender))))
	mux.Handle("PUT /{bidId}/feedback", accessMiddleware(http.HandlerFunc(r.updateBidFeedbackHandler())))

	return http.StripPrefix("/api/bids", mux)
}

type ResponseBid struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	AuthorType string    `json:"authorType"`
	AuthorID   uuid.UUID `json:"authorId"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (br *bidRouter) createBidHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bid, problems, err := decodeValid[*model.Bid](r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(problems) > 0 {
			respondWithValidationErrors(w, problems)
			return
		}

		createdBid, err := br.bidService.CreateBid(r.Context(), &service.CreateBidInput{
			Name:        bid.Name,
			Description: bid.Description,
			TenderID:    bid.TenderID,
			AuthorType:  bid.AuthorType,
			AuthorID:    bid.AuthorID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := ResponseBid{
			ID:         createdBid.ID,
			Name:       createdBid.Name,
			Status:     createdBid.Status,
			AuthorType: createdBid.AuthorType,
			AuthorID:   createdBid.AuthorID,
			Version:    createdBid.Version,
			CreatedAt:  createdBid.CreatedAt,
		}
		respondWithJSON(w, http.StatusOK, response)
	}
}
func (br *bidRouter) GetUserBidsHandler(es service.Employee) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		if username == "" {
			respondWithError(w, http.StatusBadRequest, "Username is required")
			return
		}

		limit := 5
		offset := 0
		if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
			if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
			if o, err := strconv.Atoi(offsetParam); err == nil && o >= 0 {
				offset = o
			}
		}

		user, err := es.GetByUsername(r.Context(), username)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Failed to retrieve user "+err.Error())
			return
		}

		bids, err := br.bidService.GetUserBids(r.Context(), &service.GetBidsByUsernameInput{
			Limit:    limit,
			Offset:   offset,
			AuthorID: user.ID,
		})
		if err != nil {
			if errors.Is(err, repoerrs.ErrNotFound) {
				respondWithError(w, http.StatusNotFound, "No bids found for the given user "+err.Error())
			} else {
				respondWithError(w, http.StatusInternalServerError, "Failed to retrieve bids "+err.Error())
			}
			return
		}

		var bidsResponse []ResponseBid
		for _, bid := range bids {
			bidsResponse = append(bidsResponse, ResponseBid{
				ID:         bid.ID,
				Name:       bid.Name,
				Status:     bid.Status,
				AuthorType: bid.AuthorType,
				AuthorID:   bid.AuthorID,
				Version:    bid.Version,
				CreatedAt:  bid.CreatedAt,
			})
		}
		respondWithJSON(w, http.StatusOK, bidsResponse)
	}
}

func (br *bidRouter) getBidsByTenderHandler(es service.Employee) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		if username == "" {
			respondWithError(w, http.StatusBadRequest, "Username is required")
			return
		}

		limit := 5
		offset := 0
		if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
			if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
			if o, err := strconv.Atoi(offsetParam); err == nil && o >= 0 {
				offset = o
			}
		}

		tenderID := r.PathValue("tenderId")
		tID, err := uuid.Parse(tenderID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid tender ID format")
			return
		}
		user, err := es.GetByUsername(r.Context(), username)
		if err != nil {
			if errors.Is(err, service.ErrEmployeeNotFound) {
				respondWithError(w, http.StatusUnauthorized, "User does not exist "+err.Error())
			} else {
				respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user "+err.Error())
			}
			return
		}
		bids, err := br.bidService.GetBidsByTender(r.Context(), &service.GetBidsByTenderInput{
			Limit:    limit,
			Offset:   offset,
			TenderID: tID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve bids "+err.Error())
			return
		}

		for _, bid := range bids {
			if bid.AuthorID != user.ID {
				respondWithError(w, http.StatusForbidden, "User is not authorized to view these bids")
				return
			}
		}

		bidsResponse := make([]ResponseBid, 0, len(bids))
		for _, bid := range bids {
			bidsResponse = append(bidsResponse, ResponseBid{
				ID:         bid.ID,
				Name:       bid.Name,
				Status:     bid.Status,
				AuthorType: bid.AuthorType,
				AuthorID:   bid.AuthorID,
				Version:    bid.Version,
				CreatedAt:  bid.CreatedAt,
			})
		}

		respondWithJSON(w, http.StatusOK, bidsResponse)
	}
}

func (br *bidRouter) getBidStatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bidID := r.PathValue("bidId")
		bID, err := uuid.Parse(bidID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid bid ID format")
			return
		}

		bid, err := br.bidService.GetBidByID(r.Context(), bID)
		if err != nil {
			if errors.Is(err, service.ErrBidNotFound) {
				respondWithError(w, http.StatusNotFound, "Bid not found "+err.Error())
			} else {
				respondWithError(w, http.StatusInternalServerError, "Failed to retrieve bid "+err.Error())
			}
			return
		}

		respondWithJSON(w, http.StatusOK, bid.Status)
	}
}

func (br *bidRouter) updateBidStatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bidID := r.PathValue("bidId")
		bID, err := uuid.Parse(bidID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid bid ID format")
			return
		}

		status := r.URL.Query().Get("status")
		if status == "" {
			respondWithError(w, http.StatusBadRequest, "Status is required")
			return
		}
		if status != "Opened" && status != "Canceled" && status != "Published" {
			respondWithError(w, http.StatusBadRequest, "Invalid status")
			return
		}
		updatedBid, err := br.bidService.UpdateBidStatus(r.Context(), &service.UpdateBidStatusInput{
			BidID:  bID,
			Status: status,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update bid status "+err.Error())
			return
		}

		response := ResponseBid{
			ID:         updatedBid.ID,
			Name:       updatedBid.Name,
			Status:     updatedBid.Status,
			AuthorType: updatedBid.AuthorType,
			AuthorID:   updatedBid.AuthorID,
			Version:    updatedBid.Version,
			CreatedAt:  updatedBid.CreatedAt,
		}
		respondWithJSON(w, http.StatusOK, response)
	}
}

func (br *bidRouter) updateBidHandler() http.HandlerFunc {
	type Request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := decode[Request](r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
			return
		}
		bidID := r.PathValue("bidId")
		bID, err := uuid.Parse(bidID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid bid ID format")
			return
		}

		updatedBid, err := br.bidService.UpdateBid(r.Context(), &service.UpdateBidInput{
			BidID:       bID,
			Name:        data.Name,
			Description: data.Description,
		})
		if err != nil {
			if errors.Is(err, service.ErrBidNotFound) {
				respondWithError(w, http.StatusNotFound, "Bid not found "+err.Error())
			} else {
				respondWithError(w, http.StatusInternalServerError, "Failed to update bid "+err.Error())
			}
			return
		}

		response := ResponseBid{
			ID:         updatedBid.ID,
			Name:       updatedBid.Name,
			Status:     updatedBid.Status,
			AuthorType: updatedBid.AuthorType,
			AuthorID:   updatedBid.AuthorID,
			Version:    updatedBid.Version,
			CreatedAt:  updatedBid.CreatedAt,
		}
		respondWithJSON(w, http.StatusOK, response)
	}
}

func (br *bidRouter) updateBidDecisionHandler(ts service.Tender) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		bidID := r.PathValue("bidId")
		bID, err := uuid.Parse(bidID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid bid ID format")
			return
		}
		decision := r.URL.Query().Get("decision")
		if decision != "Approved" && decision != "Rejected" {
			respondWithError(w, http.StatusBadRequest, "Invalid decision value. Available values: 'Approved', 'Rejected'")
			return
		}

		if decision != "Approved" && decision != "Rejected" {
			respondWithError(w, http.StatusBadRequest, "Invalid decision value. Available values: 'Approved', 'Rejected'")
			return
		}

		updatedBid, err := br.bidService.UpdateBidDecision(r.Context(), &service.UpdateBidDecisionInput{
			BidID:    bID,
			Decision: decision,
		})
		if err != nil {
			if errors.Is(err, service.ErrBidNotFound) {
				respondWithError(w, http.StatusNotFound, "Bid not found "+err.Error())
			} else {
				respondWithError(w, http.StatusInternalServerError, "Failed to update bid decision "+err.Error())
			}
			return
		}

		if decision == "Approved" {
			_, err := ts.UpdateTenderStatus(r.Context(), &service.UpdateTenderStatusInput{
				TenderID: updatedBid.TenderID,
				Status:   "Closed",
			})
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Failed to update tender status "+err.Error())
				return
			}
		}

		response := ResponseBid{
			ID:         updatedBid.ID,
			Name:       updatedBid.Name,
			Status:     updatedBid.Status,
			AuthorType: updatedBid.AuthorType,
			AuthorID:   updatedBid.AuthorID,
			Version:    updatedBid.Version,
			CreatedAt:  updatedBid.CreatedAt,
		}
		respondWithJSON(w, http.StatusOK, response)
	}
}

func (br *bidRouter) updateBidFeedbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bidID := r.PathValue("bidId")
		bID, err := uuid.Parse(bidID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid bid ID format")
			return
		}
		feedback := r.URL.Query().Get("feedback")
		updatedBid, err := br.bidService.UpdateBidFeedback(r.Context(), &service.UpdateBidFeedbackInput{
			BidID:    bID,
			Feedback: feedback,
		})
		if err != nil {
			if errors.Is(err, service.ErrBidNotFound) {
				respondWithError(w, http.StatusNotFound, "Bid not found "+err.Error())
			} else {
				respondWithError(w, http.StatusInternalServerError, "Failed to update bid feedback "+err.Error())
			}
			return
		}

		response := ResponseBid{
			ID:         updatedBid.ID,
			Name:       updatedBid.Name,
			Status:     updatedBid.Status,
			AuthorType: updatedBid.AuthorType,
			AuthorID:   updatedBid.AuthorID,
			Version:    updatedBid.Version,
			CreatedAt:  updatedBid.CreatedAt,
		}
		respondWithJSON(w, http.StatusOK, response)
	}
}
