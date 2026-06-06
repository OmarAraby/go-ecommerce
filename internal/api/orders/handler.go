package orders

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/OmarAraby/go-ecommerce/internal/api/response"
	"github.com/OmarAraby/go-ecommerce/internal/api/validate"
	orderapp "github.com/OmarAraby/go-ecommerce/internal/application/services/order"
	"github.com/OmarAraby/go-ecommerce/internal/domain"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/auth"
)

type Handler struct {
	svc orderapp.Service
}

func NewHandler(svc orderapp.Service) *Handler {
	return &Handler{svc: svc}
}

// POST /orders
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Items []struct {
			ProductID int64 `json:"product_id" validate:"required,gt=0"`
			Quantity  int   `json:"quantity"   validate:"required,gte=1"`
		} `json:"items" validate:"required,min=1,dive"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validate.Check(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	claims := r.Context().Value(auth.ContextKey).(*auth.Claims)

	items := make([]orderapp.CreateOrderItemDTO, len(req.Items))
	for i, item := range req.Items {
		items[i] = orderapp.CreateOrderItemDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	order, err := h.svc.Create(r.Context(), orderapp.CreateOrderDTO{
		UserID: claims.UserID,
		Items:  items,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			response.NotFound(w, "one or more products not found")
		case errors.Is(err, domain.ErrInsufficientStock):
			response.BadRequest(w, "insufficient stock for one or more items")
		default:
			response.InternalError(w)
		}
		return
	}
	response.JSON(w, http.StatusCreated, order)
}

// GET /orders
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page := queryIntOrder(r, "page", 1)
	limit := queryIntOrder(r, "limit", 20)
	if limit > 100 {
		limit = 100
	}

	claims := r.Context().Value(auth.ContextKey).(*auth.Claims)

	orders, total, err := h.svc.ListByUser(r.Context(), claims.UserID, orderapp.ListOrdersParamsDTO{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Paginated(w, orders, page, limit, total)
}

// GET /orders/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid order id")
		return
	}

	claims := r.Context().Value(auth.ContextKey).(*auth.Claims)

	order, err := h.svc.GetByID(r.Context(), claims.UserID, orderID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "order not found")
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, order)
}

func queryIntOrder(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return def
	}
	return n
}
