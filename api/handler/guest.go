package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	"github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/ATechnoHazard/hades-2/pkg/guest"
	"github.com/ATechnoHazard/janus"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func createGuest(guestService guest.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gus := &views.Guest{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		jtk := ctx.Value("janus_context").(*janus.Account)

		if jtk.Role != "admin" {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(gus); err != nil {
			views.Wrap(err, w)
			return
		}

		eve, err := eventService.ReadEvent(gus.EventId)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if eve.OrganizationID != tk.OrgID {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		if err := guestService.SaveGuestEvent(gus.Transform(), gus.EventId); err != nil {
			views.Wrap(err, w)
			return
		}

		utils.Respond(w, utils.Message(http.StatusOK, "Guest successfully created."))
		return
	}
}

func removeGuestEvent(guestService guest.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gus := &views.Guest{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		jtk := ctx.Value("janus_context").(*janus.Account)

		if jtk.Role != "admin" {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(gus); err != nil {
			views.Wrap(err, w)
			return
		}

		eve, err := eventService.ReadEvent(gus.EventId)

		if err != nil {
			views.Wrap(err, w)
			return
		}

		if eve.OrganizationID != tk.OrgID {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		if _, err := guestService.GetGuestEvent(gus.Email, gus.EventId); err != nil {
			views.Wrap(err, w)
			return
		}

		if err := guestService.RemoveGuestEvent(gus.Email, gus.EventId); err != nil {
			views.Wrap(err, w)
			return
		}

		utils.Respond(w, utils.Message(http.StatusOK, "Removed guest from event successfully."))
		return
	}
}

func getAllGuests(guestService guest.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gus := &views.Guest{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(gus); err != nil {
			views.Wrap(err, w)
			return
		}

		eve, err := eventService.ReadEvent(gus.EventId)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		guests, err := guestService.GetAllGuestEvent(gus.EventId)

		if err != nil {
			views.Wrap(err, w)
			return
		}
		msg := utils.Message(http.StatusOK, "Guests successfully retrieved")
		msg["guests"] = guests
		utils.Respond(w, msg)
		return
	}
}

func getGuest(guestService guest.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gus := &views.Guest{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(gus); err != nil {
			views.Wrap(err, w)
			return
		}

		eve, err := eventService.ReadEvent(gus.EventId)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		g, err := guestService.GetGuestEvent(gus.Email, gus.EventId)

		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := utils.Message(http.StatusOK, "Guest successfully retrieved")
		msg["guest"] = g
		utils.Respond(w, msg)
		return
	}
}

func MakeGuestHandler(r *httprouter.Router, guestService guest.Service, eventService event.Service, j *janus.Janus) {
	r.HandlerFunc("POST", "/api/v2/guests/create-guest",
		middleware.JwtAuthentication(j.GetHandler(createGuest(guestService, eventService))))
	r.HandlerFunc("POST", "/api/v2/guests/get-guest",
		middleware.JwtAuthentication(getGuest(guestService, eventService)))
	r.HandlerFunc("POST", "/api/v2/guests/all-guests",
		middleware.JwtAuthentication(getAllGuests(guestService, eventService)))
	r.HandlerFunc("POST", "/api/v2/guests/remove-guest",
		middleware.JwtAuthentication(j.GetHandler(removeGuestEvent(guestService, eventService))))
}
