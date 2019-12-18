package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/views"
	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/participant"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateAttendee(svc participant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &views.Participant{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			views.Wrap(err, w)
			return
		}

		if err := svc.CreateAttendee(p.Transform(), p.EventId); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Attendee successfully created"))
		return
	}
}

func DeleteAttendee(svc participant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &views.Participant{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			views.Wrap(err, w)
			return
		}

		if err := svc.DeleteAttendee(p.RegNo); err != nil {
			views.Wrap(err, w)
		}

		u.Respond(w, u.Message(http.StatusOK, "Attendee successfully deleted"))
		return
	}
}

func MakeParticipantHandler(r *mux.Router, svc participant.Service) {
	r.Handle("/api/v1/admin/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})).Methods("GET")

	r.Handle("/api/v1/participants/create-attendee", CreateAttendee(svc)).Methods("POST")
	r.Handle("/api/v1/participants/delete-attendee", DeleteAttendee(svc)).Methods("POST")
}
