package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/ATechnoHazard/hades-2/pkg/organization"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func acceptJoinRequest(oSvc organization.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		j := &entities.JoinRequest{}
		if err := json.NewDecoder(r.Body).Decode(j); err != nil {
			views.Wrap(err, w)
			return
		}

		if err := oSvc.AcceptJoinReq(j.OrganizationID, tk.Email); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Join request accepted"))
	}
}

func sendJoinRequest(oSvc organization.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		j := &entities.JoinRequest{}
		if err := json.NewDecoder(r.Body).Decode(j); err != nil {
			views.Wrap(err, w)
			return
		}
		if err := oSvc.SendJoinRequest(j.OrganizationID, tk.Email); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Join request created successfully"))
	}
}

func MakeOrgHandler(r *httprouter.Router, oSvc organization.Service) {
	r.HandlerFunc("POST", "/api/v1/org/accept", middleware.JwtAuthentication(acceptJoinRequest(oSvc)))
	r.HandlerFunc("POST", "/api/v1/org/join", middleware.JwtAuthentication(sendJoinRequest(oSvc)))
}
