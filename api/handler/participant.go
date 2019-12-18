package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/middleware"
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
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Attendee successfully deleted"))
		return
	}
}

func ReadAttendee(svc participant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		regNo, ok := r.URL.Query()["reg_no"]
		if !ok || len(regNo) < 1 {
			u.Respond(w, u.Message(http.StatusBadRequest, "Invalid Registration number"))
		}

		a, err := svc.ReadAttendee(regNo[0])
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := u.Message(http.StatusOK, "Attendee found")
		msg["attendee"] = a

		u.Respond(w, msg)
		return
	}
}

func RmAttendeeEvent(svc participant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &views.Participant{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			views.Wrap(err, w)
			return
		}

		if err := svc.RemoveAttendeeEvent(p.RegNo, p.EventId); err != nil {
			views.Wrap(err, w)
			return
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

	r.Handle("/api/v1/participants/create-attendee",
		middleware.JwtAuthentication(CreateAttendee(svc))).Methods("POST")
	r.Handle("/api/v1/participants/delete-attendee",
		middleware.JwtAuthentication(DeleteAttendee(svc))).Methods("POST")
	r.Handle("/api/v1/participants/read-attendee",
		middleware.JwtAuthentication(ReadAttendee(svc))).Methods("GET")
	r.Handle("/api/v1/participants/rm-attendee",
		middleware.JwtAuthentication(RmAttendeeEvent(svc))).Methods("POST")
}
