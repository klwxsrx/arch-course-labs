package transport

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/common/app/log"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/common/infra/transport"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/order/app/service"
	"github.com/klwxsrx/arch-course-labs/idempotence/pkg/order/domain"
	"net/http"
)

type createOrderItemData struct {
	ID        uuid.UUID `json:"id"`
	ItemPrice int       `json:"item_price"`
	Quantity  int       `json:"quantity"`
}

type createOrderData struct {
	Items []createOrderItemData `json:"items"`
}

type route struct {
	Name    string
	Method  string
	Pattern string
	Handler func(*service.OrderService, http.ResponseWriter, *http.Request)
}

func getRoutes() []route {
	return []route{
		{
			"createOrder",
			http.MethodPut,
			"/orders",
			createOrderHandler,
		},
		{
			"health",
			http.MethodGet,
			"/healthz",
			healthCheckHandler,
		},
	}
}

func createOrderHandler(srv *service.OrderService, w http.ResponseWriter, r *http.Request) {
	subjectID, err := parseSubjectID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	idempotenceKey, err := parseIdempotenceKey(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var createOrder createOrderData
	err = json.NewDecoder(r.Body).Decode(&createOrder)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderItems := make([]domain.OrderItem, 0, len(createOrder.Items))
	for _, item := range createOrder.Items {
		orderItems = append(orderItems, domain.OrderItem{
			ID:        item.ID,
			ItemPrice: item.ItemPrice,
			Quantity:  item.Quantity,
		})
	}

	err = srv.Create(idempotenceKey, subjectID, orderItems)
	if errors.Is(err, service.ErrOrderAlreadyCreated) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if errors.Is(err, service.ErrEmptyOrder) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func healthCheckHandler(_ *service.OrderService, w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{"OK"})
}

func parseSubjectID(r *http.Request) (uuid.UUID, error) {
	id := r.Header.Get("X-Auth-User-ID")
	return uuid.Parse(id)
}

func parseIdempotenceKey(r *http.Request) (string, error) {
	key := r.Header.Get("X-Idempotence-Key")
	if key == "" {
		return "", errors.New("idempotence key not found")
	}
	return key, nil
}

func getHandlerFunc(
	userProfileService *service.OrderService,
	f func(*service.OrderService, http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(userProfileService, w, r)
	}
}

func NewHTTPHandler(userProfileService *service.OrderService, logger log.Logger) (http.Handler, error) {
	router := mux.NewRouter()

	for _, route := range getRoutes() {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(getHandlerFunc(userProfileService, route.Handler))
	}

	router.Use(transport.NewLoggingMiddleware(logger))
	return router, nil
}
