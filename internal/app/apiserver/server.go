package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"gismart-rest-api/internal/app/model"
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
	ctxKeyUser ctxKey = iota
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
	s.router.HandleFunc("/new_dish", s.handleMenuCreate()).Methods("POST")
	s.router.HandleFunc("/create_order", s.handleOrderCreate()).Methods("POST")
	s.router.HandleFunc("/change_order_status", s.handlerChangeOrderStatus()).Methods("POST")
	s.router.HandleFunc("/get_orders", s.handlerGetOrders())
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
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

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session, err := s.sessionStore.Get(request, sessionName)
		if err != nil {
			s.error(writer, request, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(writer, request, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(writer, request, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyUser, u)) )

	})
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		s.respond(writer, request, http.StatusOK, request.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w,r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email: req.Email,
			Password: req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err )
			return
		}

		u.Sanitaze()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w,r, http.StatusBadRequest, err)
			return
		}
		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil{
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
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