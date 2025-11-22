package scheduling

import (
	"encoding/json"
	"net/http"

	"github.com/Nikitaannusewicz/carwash-crm/internal/middleware"
	"github.com/Nikitaannusewicz/carwash-crm/internal/modules/identity"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleCreateShift(w http.ResponseWriter, r *http.Request) {
	roleStr, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}

	userIDVal := r.Context().Value(middleware.UserIDKey)

	var userID int64
	switch v := userIDVal.(type) {
	case int64:
		userID = v
	case float64:
		userID = int64(v)
	default:
		userID = 0
	}

	var req CreateShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body: ensure dates are ISO8601/RFC33339", http.StatusBadRequest)
		return
	}

	shift, err := h.service.CreateShift(r.Context(), req, userID, identity.Role(roleStr))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shift)
}
