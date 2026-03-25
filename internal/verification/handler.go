package verification

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	chi "github.com/go-chi/chi/v5"

	"github.com/nathang0147/Request-Service-API/internal/shared/apierror"
)

const errCodeInternal = "INTERNAL_ERROR"

type RequestService interface {
	CreateRequest(context.Context, CreateRequestInput) (CreateRequestResponse, error)
	GetRequest(context.Context, string) (GetRequestResponse, error)
}

type Handler struct {
	service RequestService
}

func NewHandler(service RequestService) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/api/v1/verification-requests", func(router chi.Router) {
		router.Post("/", handler.handleCreateRequest)
		router.Get("/{requestID}", handler.handleGetRequest)
	})
}

func (handler *Handler) handleCreateRequest(w http.ResponseWriter, r *http.Request) {
	var body CreateRequestInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if strings.TrimSpace(body.BusinessRef) == "" || strings.TrimSpace(body.CandidateRef) == "" {
		apierror.Write(w, http.StatusBadRequest, "INVALID_REQUEST", "businessRef and candidateRef are required")
		return
	}

	response, err := handler.service.CreateRequest(r.Context(), body)
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

func (handler *Handler) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	requestID := strings.TrimSpace(chi.URLParam(r, "requestID"))
	if requestID == "" {
		apierror.Write(w, http.StatusBadRequest, "INVALID_REQUEST", "request id is required")
		return
	}

	response, err := handler.service.GetRequest(r.Context(), requestID)
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (handler *Handler) writeServiceError(w http.ResponseWriter, err error) {
	var serviceError *Error
	if errors.As(err, &serviceError) && serviceError.Code == ErrCodeRequestNotFound {
		apierror.Write(w, http.StatusNotFound, ErrCodeRequestNotFound, "verification request not found")
		return
	}

	apierror.Write(w, http.StatusInternalServerError, errCodeInternal, "internal server error")
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
