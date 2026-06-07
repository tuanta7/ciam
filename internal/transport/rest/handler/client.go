package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tuanta7/ciam/internal/oauth2client"
	"github.com/tuanta7/ciam/internal/transport/rest/middleware"
	"github.com/tuanta7/ciam/pkg/httpx"
)

type ClientHandler struct {
	uc *oauth2client.UseCase
}

func NewClientHandler(uc *oauth2client.UseCase) *ClientHandler {
	return &ClientHandler{
		uc: uc,
	}
}

func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	page, pageSize, _ := middleware.GetPaginationParams(r.Context())
	clients, err := h.uc.List(r.Context(), page, pageSize)
	if err != nil {
		_ = httpx.ErrorJSON(w, httpx.NewInternalError(httpx.WithDescription(err.Error())))
		return
	}

	_ = httpx.ResponseJSON(w, http.StatusOK, clients)
}

func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	item, err := h.uc.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, err)
		return
	}

	_ = httpx.ResponseJSON(w, http.StatusOK, item)
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var input oauth2client.CreateInput
	if err := httpx.DecodeAndValidateJSON(r.Body, &input); err != nil {
		_ = httpx.ErrorJSON(w, httpx.NewInvalidArgumentError(httpx.WithDescription(err.Error())))
		return
	}

	item, err := h.uc.Create(r.Context(), input)
	if err != nil {
		h.writeError(w, err)
		return
	}

	_ = httpx.ResponseJSON(w, http.StatusCreated, item)
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	var input oauth2client.UpdateInput
	if err := httpx.DecodeAndValidateJSON(r.Body, &input); err != nil {
		_ = httpx.ErrorJSON(w, httpx.NewInvalidArgumentError(httpx.WithDescription(err.Error())))
		return
	}

	item, err := h.uc.Update(r.Context(), chi.URLParam(r, "id"), input)
	if err != nil {
		h.writeError(w, err)
		return
	}

	_ = httpx.ResponseJSON(w, http.StatusOK, item)
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ClientHandler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, oauth2client.ErrNotFound):
		_ = httpx.ErrorJSON(w, httpx.NewError(http.StatusNotFound, "not found", httpx.WithDescription(err.Error())))
	case errors.Is(err, oauth2client.ErrInvalidClient):
		_ = httpx.ErrorJSON(w, httpx.NewInvalidArgumentError(httpx.WithDescription(err.Error())))
	default:
		_ = httpx.ErrorJSON(w, httpx.NewInternalError(httpx.WithDescription(err.Error())))
	}
}
