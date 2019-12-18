package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/ATechnoHazard/hades-2/pkg/participant"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateAttendee(pSvc participant.Service, eSvc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &views.Participant{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			views.Wrap(err, w)
			return
		}

		e, err := eSvc.ReadEvent(p.EventId)
		if err != nil {
			views.Wrap(err, w)
			return
		}
		if e.OrganizationID != tk.OrgID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
		}

		if err := pSvc.CreateAttendee(p.Transform(), p.EventId); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Attendee successfully created"))
		return
	}
}

func DeleteAttendee(pSvc participant.Service, eSvc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &views.Participant{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(middleware.Token)
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			views.Wrap(err, w)
			return
		}

		e, err := eSvc.ReadEvent(p.EventId)
		if err != nil {
			views.Wrap(err, w)
			return
		}
		if e.OrganizationID != tk.OrgID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
		}

		if err := pSvc.DeleteAttendee(p.RegNo); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Attendee successfully deleted"))
		return
	}
}

func ReadAttendee(pSvc participant.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		regNo, ok := r.URL.Query()["reg_no"]
		if !ok || len(regNo) < 1 {
			u.Respond(w, u.Message(http.StatusBadRequest, "Invalid Registration number"))
		}

		a, err := pSvc.ReadAttendee(regNo[0])
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

func RmAttendeeEvent(pSvc participant.Service, eSvc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &views.Participant{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(middleware.Token)
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			views.Wrap(err, w)
			return
		}

		e, err := eSvc.ReadEvent(p.EventId)
		if err != nil {
			views.Wrap(err, w)
			return
		}
		if e.OrganizationID != tk.OrgID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
		}

		if err := pSvc.RemoveAttendeeEvent(p.RegNo, p.EventId); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Attendee successfully deleted"))
		return
	}
}

func MakeParticipantHandler(r *mux.Router, partSvc participant.Service, eventSvc event.Service) {
	r.Handle("/api/v1/admin/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})).Methods("GET")

	r.Handle("/api/v1/participants/create-attendee",
		middleware.JwtAuthentication(CreateAttendee(partSvc, eventSvc))).Methods("POST")
	r.Handle("/api/v1/participants/delete-attendee",
		middleware.JwtAuthentication(DeleteAttendee(partSvc, eventSvc))).Methods("POST")
	r.Handle("/api/v1/participants/read-attendee",
		middleware.JwtAuthentication(ReadAttendee(partSvc))).Methods("GET")
	r.Handle("/api/v1/participants/rm-attendee",
		middleware.JwtAuthentication(RmAttendeeEvent(partSvc, eventSvc))).Methods("POST")
}
