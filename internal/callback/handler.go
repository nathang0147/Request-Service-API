package callback

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	chi "github.com/go-chi/chi/v5"

	"github.com/nathang0147/Request-Service-API/internal/shared/apierror"
)

const callbackSecretHeader = "X-Callback-Secret"

type CallbackService interface {
	HandleCallback(context.Context, []byte) (Result, error)
}

type Authenticator interface {
	Authenticate(*http.Request) error
}

type Handler struct {
	service       CallbackService
	authenticator Authenticator
}

func NewHandler(service CallbackService, authenticator Authenticator) *Handler {
	return &Handler{
		service:       service,
		authenticator: authenticator,
	}
}

func (handler *Handler) RegisterRoutes(router chi.Router) {
	router.Post("/api/v1/callbacks/walt", handler.handleWaltCallback)
}

func (handler *Handler) handleWaltCallback(w http.ResponseWriter, r *http.Request) {
	if handler.authenticator != nil {
		if err := handler.authenticator.Authenticate(r); err != nil {
			apierror.Write(w, http.StatusUnauthorized, "UNAUTHORIZED", "callback authentication failed")
			return
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		apierror.Write(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid callback body")
		return
	}

	result, err := handler.service.HandleCallback(r.Context(), body)
	if err != nil {
		apierror.Write(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusAccepted, result)
}

type staticAuthenticator struct {
	secret string
}

func NewStaticAuthenticator(secret string) Authenticator {
	return staticAuthenticator{secret: secret}
}

func (auth staticAuthenticator) Authenticate(request *http.Request) error {
	if auth.secret == "" {
		return nil
	}

	if request.Header.Get(callbackSecretHeader) == auth.secret {
		return nil
	}

	if token, ok := strings.CutPrefix(request.Header.Get("Authorization"), "Bearer "); ok && strings.TrimSpace(token) == auth.secret {
		return nil
	}

	return errors.New("unauthorized")
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
