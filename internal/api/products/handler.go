package products

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/OmarAraby/go-ecommerce/internal/api/response"
	"github.com/OmarAraby/go-ecommerce/internal/api/validate"
	productapp "github.com/OmarAraby/go-ecommerce/internal/application/services/product"
	"github.com/OmarAraby/go-ecommerce/internal/domain"
)

func queryInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return def
	}
	return n
}

func queryFloat(r *http.Request, key string) float64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return 0
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || f < 0 {
		return 0
	}
	return f
}

type Handler struct {
	svc productapp.Service
}

func NewHandler(svc productapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	limit := queryInt(r, "limit", 20)
	if limit > 100 {
		limit = 100
	}

	params := productapp.ListParams{
		Page:     page,
		Limit:    limit,
		Sort:     r.URL.Query().Get("sort"),
		Order:    r.URL.Query().Get("order"),
		Name:     r.URL.Query().Get("name"),
		MinPrice: queryFloat(r, "min_price"),
		MaxPrice: queryFloat(r, "max_price"),
	}

	prods, total, err := h.svc.List(r.Context(), params)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Paginated(w, prods, page, limit, total)
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
		Name        string  `json:"name"        validate:"required,min=1,max=200"`
		Description string  `json:"description" validate:"max=1000"`
		Price       float64 `json:"price"       validate:"required,gt=0"`
		Stock       int     `json:"stock"       validate:"gte=0"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validate.Check(req); errs != nil {
		response.ValidationError(w, errs)
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
		Name        string  `json:"name"        validate:"required,min=1,max=200"`
		Description string  `json:"description" validate:"max=1000"`
		Price       float64 `json:"price"       validate:"required,gt=0"`
		Stock       int     `json:"stock"       validate:"gte=0"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validate.Check(req); errs != nil {
		response.ValidationError(w, errs)
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
