package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/ATechnoHazard/hades-2/pkg/participant"
	"github.com/ATechnoHazard/janus"
	"github.com/julienschmidt/httprouter"
)

func createAttendee(pSvc participant.Service, eSvc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &views.Participant{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		jtk := ctx.Value("janus_context").(*janus.Account)

		if jtk.Role != "admin" {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

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
			return
		}

		if err := pSvc.CreateAttendee(p.Transform(), p.EventId); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Attendee successfully created"))
		return
	}
}

func deleteAttendee(pSvc participant.Service, eSvc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &views.Participant{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		jtk := ctx.Value("janus_context").(*janus.Account)

		if jtk.Role != "admin" {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

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

		if err := pSvc.DeleteAttendee(p.Email); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Attendee successfully deleted"))
		return
	}
}

func readAttendee(pSvc participant.Service, eSvc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := r.URL.Query()["email"]
		if !ok || len(email) < 1 {
			u.Respond(w, u.Message(http.StatusBadRequest, "Invalid email ID"))
			return
		}

		eventID, ok := r.URL.Query()["event_id"]
		if !ok || len(eventID) < 1 {
			u.Respond(w, u.Message(http.StatusBadRequest, "Invalid event ID"))
			return
		}

		eID, err := strconv.Atoi(eventID[0])
		if err != nil {
			views.Wrap(err, w)
			return
		}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		e, err := eSvc.ReadEvent(uint(eID))
		if err != nil {
			views.Wrap(err, w)
			return
		}
		if e.OrganizationID != tk.OrgID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from accessing this resource"))
		}

		a, err := pSvc.ReadAttendee(email[0], e.ID)
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

func rmAttendeeEvent(pSvc participant.Service, eSvc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &views.Participant{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		jtk := ctx.Value("janus_context").(*janus.Account)

		if jtk.Role != "admin" {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			views.Wrap(err, w)
			return
		}

		e, err := eSvc.ReadEvent(p.EventId)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		for _, part := range e.Attendees {
			if part.Email == p.Email {
				if e.OrganizationID != tk.OrgID {
					u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
					return
				}

				if err := pSvc.RemoveAttendeeEvent(p.Email, p.EventId); err != nil {
					views.Wrap(err, w)
					return
				}

				u.Respond(w, u.Message(http.StatusOK, "Attendee successfully deleted"))
				return
			}
		}

		u.Respond(w, u.Message(http.StatusNotFound, "This attendee does not exist for this event"))
	}
}

func MakeParticipantHandler(r *httprouter.Router, partSvc participant.Service, eventSvc event.Service, j *janus.Janus) {
	r.HandlerFunc("POST", "/api/v2/participants/create-attendee",
		middleware.JwtAuthentication(j.GetHandler(createAttendee(partSvc, eventSvc))))
	r.HandlerFunc("DELETE", "/api/v2/participants/delete-attendee",
		middleware.JwtAuthentication(j.GetHandler(deleteAttendee(partSvc, eventSvc))))
	r.HandlerFunc("GET", "/api/v2/participants/read-attendee",
		middleware.JwtAuthentication(readAttendee(partSvc, eventSvc)))
	r.HandlerFunc("DELETE", "/api/v2/participants/rm-attendee",
		middleware.JwtAuthentication(j.GetHandler(rmAttendeeEvent(partSvc, eventSvc))))
}
