package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"git.codenrock.com/tender/internal/model"
	"git.codenrock.com/tender/internal/service"
	"github.com/google/uuid"
)

func createTenderMiddleware(services *service.Services) func(http.Handler) http.Handler {
	type Request struct {
		OrganizationID  string `json:"organizationId"`
		CreatorUsername string `json:"creatorUsername"`
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var data Request

			bodyBuf := new(bytes.Buffer)
			teeReader := io.TeeReader(r.Body, bodyBuf)

			if err := json.NewDecoder(teeReader).Decode(&data); err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
				return
			}

			r.Body = io.NopCloser(bodyBuf)

			orgID, err := uuid.Parse(data.OrganizationID)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid organization ID format")
				return
			}

			employee, err := services.Employee.GetByUsername(r.Context(), data.CreatorUsername)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Invalid creator username")
				return
			}

			_, err = services.Organization.GetOrganizationResponsible(r.Context(), &service.OrganizationResponsibleInput{
				OrganizationID: orgID,
				EmployeeID:     employee.ID,
			})
			if err != nil {
				respondWithError(w, http.StatusForbidden, "User is not organization responsible")
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
func getTenderStatusMiddleware(services *service.Services) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tenderID := r.PathValue("tenderId")
			tID, err := uuid.Parse(tenderID)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid tender ID format: "+err.Error())
				return
			}

			tender, err := services.Tender.GetTenderByID(r.Context(), tID)
			if err != nil {
				if errors.Is(err, service.ErrTenderNotFound) {
					respondWithError(w, http.StatusNotFound, "Tender not found: "+err.Error())
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve tender: "+err.Error())
				}
				return
			}

			if tender.Status == "PUBLISHED" {
				h.ServeHTTP(w, r)
				return
			}

			username := r.URL.Query().Get("username")
			user, err := services.Employee.GetByUsername(r.Context(), username)
			if err != nil {
				if errors.Is(err, service.ErrEmployeeNotFound) {
					respondWithError(w, http.StatusUnauthorized, "Employee does not exist: "+err.Error())
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user: "+err.Error())
				}
				return
			}

			_, err = services.Organization.GetOrganizationResponsible(r.Context(), &service.OrganizationResponsibleInput{
				OrganizationID: tender.OrganizationID,
				EmployeeID:     user.ID,
			})
			if err != nil {
				respondWithError(w, http.StatusForbidden, "User is not responsible for this organization: "+err.Error())
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func updateTenderStatusMiddleware(services *service.Services) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tenderID := r.PathValue("tenderId")
			tID, err := uuid.Parse(tenderID)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid tender ID format "+err.Error())
				return
			}

			tender, err := services.Tender.GetTenderByID(r.Context(), tID)
			if err != nil {
				if errors.Is(err, service.ErrTenderNotFound) {
					respondWithError(w, http.StatusNotFound, "Tender not found "+err.Error())
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve tender "+err.Error())
				}
				return
			}
			username := r.URL.Query().Get("username")
			user, err := services.Employee.GetByUsername(r.Context(), username)
			if err != nil {
				if errors.Is(err, service.ErrEmployeeNotFound) {
					respondWithError(w, http.StatusUnauthorized, "Employee does not exist: "+err.Error())
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user: "+err.Error())
				}
				return
			}
			if tender.CreatorUsername != user.Username {
				respondWithError(w, http.StatusForbidden, "Access denied: only the creator can update this tender")
				return
			}

			_, err = services.Organization.GetOrganizationResponsible(r.Context(), &service.OrganizationResponsibleInput{
				OrganizationID: tender.OrganizationID,
				EmployeeID:     user.ID,
			})
			if err != nil {
				respondWithError(w, http.StatusForbidden, "User is not responsible for this organization: "+err.Error())
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
func UpdateTenderMiddleware(services *service.Services) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenderID := r.PathValue("tenderId")
			tID, err := uuid.Parse(tenderID)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid tender ID format: "+err.Error())
				return
			}

			tender, err := services.Tender.GetTenderByID(r.Context(), tID)
			if err != nil {
				if errors.Is(err, service.ErrTenderNotFound) {
					respondWithError(w, http.StatusNotFound, "Tender not found: "+err.Error())
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve tender: "+err.Error())
				}
				return
			}

			username := r.URL.Query().Get("username")
			if username != tender.CreatorUsername {
				respondWithError(w, http.StatusForbidden, "Access denied: only the creator can update this tender")
				return
			}

			user, err := services.Employee.GetByUsername(r.Context(), username)
			if err != nil {
				if errors.Is(err, service.ErrEmployeeNotFound) {
					respondWithError(w, http.StatusUnauthorized, "Employee does not exist: "+err.Error())
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user: "+err.Error())
				}
				return
			}

			_, err = services.Organization.GetOrganizationResponsible(r.Context(), &service.OrganizationResponsibleInput{
				OrganizationID: tender.OrganizationID,
				EmployeeID:     user.ID,
			})
			if err != nil {
				respondWithError(w, http.StatusForbidden, "User is not responsible for this organization: "+err.Error())
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func createBidMiddleware(services *service.Services) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var data model.Bid

			bodyBuf := new(bytes.Buffer)
			teeReader := io.TeeReader(r.Body, bodyBuf)

			if err := json.NewDecoder(teeReader).Decode(&data); err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
				return
			}

			r.Body = io.NopCloser(bodyBuf)

			tender, err := services.Tender.GetTenderByID(r.Context(), data.TenderID)
			if err != nil {
				if errors.Is(err, service.ErrTenderNotFound) {
					respondWithError(w, http.StatusNotFound, err.Error())
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve tender: "+err.Error())
				}
				return
			}

			author, err := services.Employee.GetByID(r.Context(), data.AuthorID)
			if err != nil {
				if errors.Is(err, service.ErrEmployeeNotFound) {
					respondWithError(w, http.StatusUnauthorized, "Author does not exist "+err.Error())
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve author: "+err.Error())
				}
				return
			}

			orgResponsible, err := services.Organization.GetOrganizationResponsible(r.Context(), &service.OrganizationResponsibleInput{
				OrganizationID: tender.OrganizationID,
				EmployeeID:     author.ID,
			})
			if err != nil || orgResponsible.EmployeeID != author.ID {
				respondWithError(w, http.StatusForbidden, "User is not responsible for this organization's tender")
				return
			}

			existingBid, err := services.Bid.GetBidByTenderAndAuthor(r.Context(), tender.ID, author.ID)
			if err == nil && existingBid != nil {
				respondWithError(w, http.StatusConflict, "A bid for this tender already exists from this author")
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func authorOrResponsibleMiddleware(services *service.Services) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bidID := r.PathValue("bidId")
			bID, err := uuid.Parse(bidID)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid bid ID format")
				return
			}

			bid, err := services.GetBidByID(r.Context(), bID)
			if err != nil {
				if errors.Is(err, service.ErrBidNotFound) {
					respondWithError(w, http.StatusNotFound, "Bid not found")
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve bid: "+err.Error())
				}
				return
			}

			username := r.URL.Query().Get("username")
			if username == "" {
				respondWithError(w, http.StatusBadRequest, "Username is required")
				return
			}

			user, err := services.Employee.GetByUsername(r.Context(), username)
			if err != nil {
				if errors.Is(err, service.ErrEmployeeNotFound) {
					respondWithError(w, http.StatusUnauthorized, "User does not exist")
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user: "+err.Error())
				}
				return
			}

			if bid.AuthorID == user.ID {
				h.ServeHTTP(w, r)
				return
			}

			isResponsible, err := services.IsResponsible(r.Context(), user.ID)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Failed to check if user is responsible: "+err.Error())
				return
			}

			if !isResponsible {
				respondWithError(w, http.StatusForbidden, "User is neither the author nor responsible for this bid")
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
func accessMiddleware(services *service.Services) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bidID := r.PathValue("bidId")
			bID, err := uuid.Parse(bidID)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid bid ID format")
				return
			}

			bid, err := services.Bid.GetBidByID(r.Context(), bID)
			if err != nil {
				if errors.Is(err, service.ErrBidNotFound) {
					respondWithError(w, http.StatusNotFound, "Bid not found")
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve bid: "+err.Error())
				}
				return
			}

			username := r.URL.Query().Get("username")
			if username == "" {
				respondWithError(w, http.StatusBadRequest, "Username is required")
				return
			}

			user, err := services.Employee.GetByUsername(r.Context(), username)
			if err != nil {
				if errors.Is(err, service.ErrEmployeeNotFound) {
					respondWithError(w, http.StatusUnauthorized, "User does not exist")
				} else {
					respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user: "+err.Error())
				}
				return
			}

			isResponsible, err := services.Organization.IsResponsibleForTender(r.Context(), user.ID, bid.TenderID)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Failed to check responsibility: "+err.Error())
				return
			}

			if !isResponsible {
				respondWithError(w, http.StatusForbidden, "User is not responsible for this tender")
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
