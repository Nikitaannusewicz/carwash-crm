package operations

import (
	"encoding/json"
	"net/http"

	"github.com/Nikitaannusewicz/carwash-crm/internal/middleware"
	"github.com/Nikitaannusewicz/carwash-crm/internal/modules/identity"
)

type Handler struct {
	service *OperationsService
}

func NewHandler(service *OperationsService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleCreateBay(w http.ResponseWriter, r *http.Request) {
	roleStr, ok := r.Context().Value(middleware.RoleKey).(string)

	if !ok {
		http.Error(w, "unathorized", http.StatusUnauthorized)
		return
	}

	var req CreateBayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	b, err := h.service.CreateBay(r.Context(), req, identity.Role(roleStr))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

func (h *Handler) HandleCreateLocation(w http.ResponseWriter, r *http.Request) {
	roleStr, ok := r.Context().Value(middleware.RoleKey).(string)

	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	loc, err := h.service.CreateLocation(r.Context(), req, identity.Role(roleStr))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loc)
}

func (h *Handler) HandleCreateService(w http.ResponseWriter, r *http.Request) {
	roleStr, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	ser, err := h.service.CreateService(r.Context(), req, identity.Role(roleStr))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "applicaiton/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ser)
}
