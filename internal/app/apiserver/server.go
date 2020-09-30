package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"gismart-rest-api/internal/app/store"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const(
	sessionName = "gismart_api"
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("Incorrect email or password")
	errNotAuthenticated = errors.New("Not authenticated")
	errLenValidateError = errors.New("Various number of dishes and quantities")
	errValueValidateError = errors.New("Invalid request, check the entered data ")
)

type ctxKey int8

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store: store,
		sessionStore: sessionStore,
	}
	s.cofigureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	s.router.ServeHTTP(w,r)
}

func (s *server) cofigureRouter()  {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	// создание нового блюда в меню
	s.router.HandleFunc("/new_dish", s.handleMenuCreate()).Methods("POST")
	// создание нового заказа
	s.router.HandleFunc("/create_order", s.handleOrderCreate()).Methods("POST")
	// изменение сатуса заказа
	s.router.HandleFunc("/change_order_status", s.handlerChangeOrderStatus()).Methods("POST")
	// получение списка заказов
	s.router.HandleFunc("/get_orders", s.handlerGetOrders())
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id := uuid.New().String()
		writer.Header().Set("X-Request-ID", id)
		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remoute_addr": request.RemoteAddr,
			"request_id": request.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", request.Method, request.RequestURI)

		start := time.Now()
		rw := &responseWriter{writer, http.StatusOK}
		next.ServeHTTP(rw, request)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
			)
	})
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error)  {
	s.respond(w, r, code, map[string]string{"error": err.Error()})

}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}