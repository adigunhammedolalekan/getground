package handler

import (
	"encoding/json"
	"github.com/getground/tech-tasks/backend/service"
	"github.com/getground/tech-tasks/backend/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type (
	Handler struct {
		svc    service.Service
		logger *logrus.Logger
	}
)

func NewHandler(s service.Service) *Handler {
	return &Handler{svc: s, logger: logrus.StandardLogger()}
}

func (h *Handler) AddTableHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Capacity      json.Number `json:"capacity"`
		AllowedExtras json.Number `json:"allowed_extras"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		badRequest(w, r, "request body is not formatted in JSON")
		return
	}

	capacity, cerr := request.Capacity.Int64()
	allowedExtras, _ := request.AllowedExtras.Int64() // ignore error, defaults to 0
	if cerr != nil {
		badRequest(w, r, "value of capacity is not a number")
		return
	}

	newTable, err := h.svc.AddTable(int(capacity), int(allowedExtras))
	if err != nil {
		respond(w, r, err)
		return
	}

	ok(w, r, newTable)
}

func (h *Handler) AddGuestHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Table              json.Number `json:"table"`
		AccompanyingGuests json.Number `json:"accompanying_guests"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		badRequest(w, r, "request body is not formatted in JSON")
		return
	}

	table, cerr := request.Table.Int64()
	accompanyingGuests, err := request.AccompanyingGuests.Int64()
	name := chi.URLParam(r, "name")
	if cerr != nil || err != nil {
		badRequest(w, r, "value of table or accompanying guests is not a number")
		return
	}

	if name == "" {
		badRequest(w, r, "name is empty")
		return
	}

	guest, err := h.svc.AddGuest(name, int(table), int(accompanyingGuests))
	if err != nil {
		respond(w, r, err)
		return
	}

	response := struct {
		Name string `json:"name"`
	}{Name: guest.Name}
	ok(w, r, response)
}

func (h *Handler) GetGuestsListHandler(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetGuestsList()
	if err != nil {
		respond(w, r, err)
		return
	}

	response := struct {
		Guests []*types.Guest `json:"guests"`
	}{Guests: data}
	ok(w, r, response)
}

func (h *Handler) GuestArrivesHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		AccompanyingGuests json.Number `json:"accompanying_guests"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		badRequest(w, r, "request body is not formatted in JSON")
		return
	}

	accompanyingGuests, err := request.AccompanyingGuests.Int64()
	name := chi.URLParam(r, "name")
	if err != nil {
		badRequest(w, r, "value of accompanying guests is not a number")
		return
	}

	if name == "" {
		badRequest(w, r, "name is empty")
		return
	}

	guest, err := h.svc.GuestArrives(name, int(accompanyingGuests))
	if err != nil {
		respond(w, r, err)
		return
	}

	response := struct {
		Name string `json:"name"`
	}{Name: guest.Name}
	ok(w, r, response)
}

func (h *Handler) GuestLeavesHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		badRequest(w, r, "guest name is empty")
		return
	}

	err := h.svc.GuestsLeave(name)
	if err != nil {
		respond(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.Respond(w, r, nil)
}

func (h *Handler) GetArrivedGuestsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetArrivedGuests()
	if err != nil {
		respond(w, r, err)
		return
	}

	formattedData := make([]types.GetArrivedGuestsResponse, 0, len(data))
	for _, next := range data {
		formattedData = append(formattedData, types.GetArrivedGuestsResponse{
			Name:               next.Name,
			AccompanyingGuests: next.AccompanyingGuests,
			TimeArrived:        next.CreatedAt.Format(time.RFC1123),
		})
	}

	response := struct {
		Guests []types.GetArrivedGuestsResponse `json:"guests"`
	}{Guests: formattedData}
	ok(w, r, response)
}

func (h *Handler) AvailableSeatsHandler(w http.ResponseWriter, r *http.Request) {
	availableSeats, err := h.svc.AvailableSeats()
	if err != nil {
		respond(w, r, err)
		return
	}

	response := struct {
		SeatsEmpty int `json:"seats_empty"`
	}{SeatsEmpty: availableSeats}
	ok(w, r, response)
}
