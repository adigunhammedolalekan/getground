package handler

import (
	"github.com/getground/tech-tasks/backend/errors"
	"github.com/go-chi/render"
	"net/http"
)

type errorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func badRequest(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusBadRequest)
	render.Respond(w, r, &errorResponse{Error: true, Message: message})
}

func notFound(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusNotFound)
	render.Respond(w, r, &errorResponse{Error: true, Message: message})
}

func serverError(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusForbidden)
	render.Respond(w, r, &errorResponse{Error: true, Message: message})
}

func ok(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusOK)
	render.Respond(w, r, data)
}

func respond(w http.ResponseWriter, r *http.Request, err error) {
	if e, ok := err.(*errors.Error); ok {
		render.Status(r, e.Code)
		render.Respond(w, r, &errorResponse{Error: true, Message: e.Message})
	} else {
		serverError(w, r, err.Error())
	}
}

func H(w http.ResponseWriter, r *http.Request) {
	ok(w, r, "haha! I'm up!")
}
