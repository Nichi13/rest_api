package apiserver

import (
	"encoding/json"
	"errors"
	"gismart-rest-api/internal/app/model"
	"net/http"
	"strconv"
	"strings"
)

var (
	errStatusValidateError = errors.New("Incorrect status")
)

func (s *server) handleMenuCreate() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
		Number string `json:"number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w,r, http.StatusBadRequest, err)
			return
		}

		menu := &model.Menu{
			Number: req.Number,
			Name: req.Name,
		}
		if err := s.store.Menu().Create(menu); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err )
			return
		}
		s.respond(w, r, http.StatusCreated, menu)
	}
}

func (s *server) handleOrderCreate() http.HandlerFunc {
	type request struct {
		Dishes string `json:"dishes"`
		Count string `json:"count"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w,r, http.StatusBadRequest, err)
			return
		}
		dishes := strings.Split(string(req.Dishes), `,`)
		counts := strings.Split(string(req.Count), `,`)
		if len(dishes) != len(counts) {
			s.error(w, r, http.StatusUnauthorized, errLenValidateError)
			return
		}
		order_number, err := s.store.Order().Create(dishes, counts);
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errValueValidateError )
			return
		}
		s.respond(w, r, http.StatusCreated, map[string]string{"order": strconv.Itoa(order_number)})
		return
	}
}

func (s *server) handlerChangeOrderStatus() http.HandlerFunc {
	type request struct {
		Number int `json:"number,string"`
		Status string `json:"status"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w,r, http.StatusBadRequest, err)
			return
		}
		if req.Status == "ready" || req.Status == "close" {
			order := &model.Order{
				Number: req.Number,
				Status: req.Status,
			}
			if err := s.store.Order().Update(order); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err )
				return
			}
			s.respond(w, r, http.StatusCreated, order)
		} else {
			s.error(w,r, http.StatusBadRequest, errStatusValidateError)
			return
		}

	}
}

func (s *server) handlerGetOrders() http.HandlerFunc {
	type request struct {
		Status string `json:"status"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if req.Status == "new" || req.Status == "ready" {
			m, err := s.store.Order().Get(req.Status);
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err )
				return
			}
			s.respond(w, r, http.StatusCreated, m)
		} else {
			s.error(w,r, http.StatusBadRequest, errStatusValidateError)
			return
		}

	}
}
