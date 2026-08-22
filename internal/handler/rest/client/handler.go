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
	page, pageSize := middleware.GetPaginationParams(r.Context())

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	clients, total, err := h.uc.List(r.Context(), page, pageSize)
	if err != nil {
		_ = rest.ErrorJSON(w, rest.InternalError().WithDescription(err.Error()))
		return
	}

	_ = rest.WriteJSON(w, http.StatusOK, rest.JSON{
		"status": "sucess",
		"data": rest.JSON{
			"clients": clients,
			"pagination": rest.JSON{
				"page":     page,
				"pageSize": pageSize,
				"total":    total,
			},
		},
	})
}

func (h *Handler) GetClient(w http.ResponseWriter, r *http.Request) {
	client, err := h.uc.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, err)
		return
	}

	_ = rest.WriteJSON(w, http.StatusOK, rest.JSON{
		"status": "sucess",
		"data": rest.JSON{
			"client": client,
		},
	})
}

func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var input ClientInput
	if err := rest.ParseJSON(r.Body, &input); err != nil {
		_ = rest.ErrorJSON(w, rest.InvalidArgumentError().WithDescription(err.Error()))
		return
	}

	if err := input.validate(); err != nil {
		h.writeError(w, err)
		return
	}

	client, secret, err := h.uc.Create(r.Context(), input.toDomain())
	if err != nil {
		h.writeError(w, err)
		return
	}

	_ = rest.WriteJSON(w, http.StatusCreated, rest.JSON{
		"status": "sucess",
		"data": rest.JSON{
			"client": client,
			"secret": secret,
		},
	})
}

func (h *Handler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	var input ClientInput
	if err := rest.ParseJSON(r.Body, &input); err != nil {
		_ = rest.ErrorJSON(w, rest.InvalidArgumentError().WithDescription(err.Error()))
		return
	}

	if err := input.validate(); err != nil {
		h.writeError(w, err)
		return
	}

	dc := input.toDomain()
	dc.ID = chi.URLParam(r, "id")

	client, err := h.uc.Update(r.Context(), dc)
	if err != nil {
		h.writeError(w, err)
		return
	}

	_ = rest.WriteJSON(w, http.StatusOK, rest.JSON{
		"status": "sucess",
		"data": rest.JSON{
			"client": client,
		},
	})
}

func (h *Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		h.writeError(w, err)
		return
	}

	_ = rest.WriteJSON(w, http.StatusNoContent, rest.JSON{
		"status": "sucess",
		"data":   nil,
	})

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
