package products

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/OmarAraby/go-ecommerce/internal/api/response"
	productapp "github.com/OmarAraby/go-ecommerce/internal/application/services/product"
	"github.com/OmarAraby/go-ecommerce/internal/domain"
)

type Handler struct {
	svc productapp.Service
}

func NewHandler(svc productapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	prods, err := h.svc.List(r.Context())
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, prods)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid product id")
		return
	}
	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "product not found")
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, p)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	p, err := h.svc.Create(r.Context(), productapp.CreateProductDTO{
		Name: req.Name, Description: req.Description, Price: req.Price, Stock: req.Stock,
	})
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusCreated, p)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid product id")
		return
	}
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	p, err := h.svc.Update(r.Context(), productapp.UpdateProductDTO{
		ID: id, Name: req.Name, Description: req.Description, Price: req.Price, Stock: req.Stock,
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "product not found")
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, p)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid product id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		response.InternalError(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
