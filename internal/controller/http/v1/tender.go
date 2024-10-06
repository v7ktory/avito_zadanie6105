package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.codenrock.com/tender/internal/model"
	"git.codenrock.com/tender/internal/service"
	"github.com/google/uuid"
)

type tenderRouter struct {
	tenderService service.Tender
}

func newTenderRouter(tenderService service.Tender, services *service.Services) http.Handler {
	r := &tenderRouter{
		tenderService: tenderService,
	}

	mux := http.NewServeMux()

	createTenderMiddleware := createTenderMiddleware(services)
	getTenderStatusMiddleware := getTenderStatusMiddleware(services)
	updateTenderStatusMiddleware := updateTenderStatusMiddleware(services)
	UpdateTenderMiddleware := UpdateTenderMiddleware(services)

	mux.Handle("POST /new", createTenderMiddleware(http.HandlerFunc(r.createTenderHandler())))
	mux.Handle("GET /", r.getTendersHandler())
	mux.Handle("GET /my", r.getUserTendersHandler(services.Employee))
	mux.Handle("GET /{tenderId}/status", getTenderStatusMiddleware(http.HandlerFunc(r.getTenderStatusHandler())))
	mux.Handle("PUT /{tenderId}/status", updateTenderStatusMiddleware(http.HandlerFunc(r.updateTenderStatusHandler())))
	mux.Handle("PATCH /{tenderId}/edit", UpdateTenderMiddleware(http.HandlerFunc(r.updateTenderHandler())))

	return http.StripPrefix("/api/tenders", mux)
}

type ResponseTender struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ServiceType string    `json:"serviceType"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (tr *tenderRouter) createTenderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tender, problems, err := decodeValid[model.Tender](r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
			return
		}

		if len(problems) > 0 {
			respondWithValidationErrors(w, problems)
			return
		}

		tID, err := uuid.Parse(tender.OrganizationID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		createdTender, err := tr.tenderService.CreateTender(r.Context(), &service.CreateTenderInput{
			Name:            tender.Name,
			Description:     tender.Description,
			ServiceType:     tender.ServiceType,
			OrganizationID:  tID,
			CreatorUsername: tender.CreatorUsername,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := ResponseTender{
			ID:          createdTender.ID,
			Name:        createdTender.Name,
			Description: createdTender.Description,
			Status:      createdTender.Status,
			ServiceType: createdTender.ServiceType,
			Version:     createdTender.Version,
			CreatedAt:   createdTender.CreatedAt,
		}

		respondWithJSON(w, http.StatusOK, response)
	}
}

func (tr *tenderRouter) getTendersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 5
		offset := 0
		serviceTypes := []string{}

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

		if serviceTypesParam := r.URL.Query().Get("service_type"); serviceTypesParam != "" {
			serviceTypes = strings.Split(serviceTypesParam, ",")
			for i, s := range serviceTypes {
				serviceTypes[i] = strings.TrimSpace(s)
			}
		}

		tenders, err := tr.tenderService.GetTenders(r.Context(), &service.GetTendersInput{
			Limit:        limit,
			Offset:       offset,
			ServiceTypes: serviceTypes,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to get tenders: "+err.Error())
			return
		}

		var tendersResponse []ResponseTender
		for _, tender := range tenders {
			tendersResponse = append(tendersResponse, ResponseTender{
				ID:          tender.ID,
				Name:        tender.Name,
				Description: tender.Description,
				Status:      tender.Status,
				ServiceType: tender.ServiceType,
				Version:     tender.Version,
				CreatedAt:   tender.CreatedAt,
			})
			respondWithJSON(w, http.StatusOK, tendersResponse)
		}
	}
}
func (tr *tenderRouter) getUserTendersHandler(es service.Employee) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 5
		offset := 0
		username := r.URL.Query().Get("username")

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
		_, err := es.GetByUsername(r.Context(), username)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		tenders, err := tr.tenderService.GetUserTenders(r.Context(), &service.GetUserTendersInput{
			Limit:    limit,
			Offset:   offset,
			Username: username,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var tendersResponse []ResponseTender
		for _, tender := range tenders {
			tendersResponse = append(tendersResponse, ResponseTender{
				ID:          tender.ID,
				Name:        tender.Name,
				Description: tender.Description,
				Status:      tender.Status,
				ServiceType: tender.ServiceType,
				Version:     tender.Version,
				CreatedAt:   tender.CreatedAt,
			})
		}
		respondWithJSON(w, http.StatusOK, tendersResponse)
	}
}

func (tr *tenderRouter) getTenderStatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenderID := r.PathValue("tenderId")
		tID, err := uuid.Parse(tenderID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid tender ID format: "+err.Error())
			return
		}

		tender, err := tr.tenderService.GetTenderByID(r.Context(), tID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, tender.Status)
	}
}
func (tr *tenderRouter) updateTenderStatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		tenderID := r.PathValue("tenderId")
		tID, _ := uuid.Parse(tenderID)

		if status != "Opened" && status != "Closed" && status != "Published" {
			respondWithError(w, http.StatusBadRequest, "Invalid status")
			return
		}

		tender, err := tr.tenderService.UpdateTenderStatus(r.Context(), &service.UpdateTenderStatusInput{
			TenderID: tID,
			Status:   status,
		})
		if err != nil {
			if errors.Is(err, service.ErrTenderNotFound) {
				respondWithError(w, http.StatusNotFound, "Tender not found "+err.Error())
			} else {
				respondWithError(w, http.StatusInternalServerError, "Failed to update tender status "+err.Error())
			}
			return
		}

		response := ResponseTender{
			ID:          tender.ID,
			Name:        tender.Name,
			Description: tender.Description,
			Status:      tender.Status,
			ServiceType: tender.ServiceType,
			Version:     tender.Version,
			CreatedAt:   tender.CreatedAt,
		}
		respondWithJSON(w, http.StatusOK, response)
	}
}
func (tr *tenderRouter) updateTenderHandler() http.HandlerFunc {
	type Request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ServiceType string `json:"serviceType"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := decode[Request](r)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
			return
		}

		tenderID := r.PathValue("tenderId")
		tID, err := uuid.Parse(tenderID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid tender ID format: "+err.Error())
			return
		}

		tender, err := tr.tenderService.UpdateTender(r.Context(), &service.UpdateTenderInput{
			TenderID:    tID,
			Name:        data.Name,
			Description: data.Description,
			ServiceType: data.ServiceType,
		})
		if err != nil {
			if errors.Is(err, service.ErrTenderNotFound) {
				respondWithError(w, http.StatusNotFound, "Tender not found: "+err.Error())
			} else {
				respondWithError(w, http.StatusInternalServerError, "Failed to update tender: "+err.Error())
			}
			return
		}
		response := ResponseTender{
			ID:          tender.ID,
			Name:        tender.Name,
			Description: tender.Description,
			Status:      tender.Status,
			ServiceType: tender.ServiceType,
			Version:     tender.Version,
			CreatedAt:   tender.CreatedAt,
		}
		respondWithJSON(w, http.StatusOK, response)
	}
}
