package client

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tuanta7/ciam/internal/domain"
	"github.com/tuanta7/ciam/internal/handler/rest"
	"github.com/tuanta7/ciam/internal/handler/rest/middleware"
	"github.com/tuanta7/ciam/internal/usecase/client"
)

type Handler struct {
	uc *client.UseCase
}

func NewHandler(uc *client.UseCase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) ListClients(w http.ResponseWriter, r *http.Request) {
	page, pageSize, _ := middleware.GetPaginationParams(r.Context())
	clients, err := h.uc.List(r.Context(), page, pageSize)
	if err != nil {
		_ = rest.ErrorJSON(w, rest.InternalError().WithDescription(err.Error()))
		return
	}

	_ = rest.WriteJSON(w, http.StatusOK, clients)
}

func (h *Handler) GetClient(w http.ResponseWriter, r *http.Request) {
	item, err := h.uc.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, err)
		return
	}

	_ = rest.WriteJSON(w, http.StatusOK, item)
}

func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var input client.CreateInput
	if err := rest.ParseJSON(r.Body, &input); err != nil {
		_ = rest.ErrorJSON(w, rest.InvalidArgumentError().WithDescription(err.Error()))
		return
	}

	item, secret, err := h.uc.Create(r.Context(), input)
	if err != nil {
		h.writeError(w, err)
		return
	}

	_ = rest.WriteJSON(w, http.StatusCreated, rest.JSON{
		"client": item,
		"secret": secret,
	})
}

func (h *Handler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	var input client.UpdateInput
	if err := rest.ParseJSON(r.Body, &input); err != nil {
		_ = rest.ErrorJSON(w, rest.InvalidArgumentError().WithDescription(err.Error()))
		return
	}

	item, err := h.uc.Update(r.Context(), chi.URLParam(r, "id"), input)
	if err != nil {
		h.writeError(w, err)
		return
	}

	_ = rest.WriteJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrClientNotFound):
		_ = rest.ErrorJSON(w, rest.Error(http.StatusNotFound, "not found").WithDescription(err.Error()))
	case errors.Is(err, domain.ErrInvalidClient):
		_ = rest.ErrorJSON(w, rest.InvalidArgumentError().WithDescription(err.Error()))
	default:
		_ = rest.ErrorJSON(w, rest.InternalError().WithDescription(err.Error()))
	}
}
