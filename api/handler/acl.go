package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/janus"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func setRights(j *janus.Janus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		jtk := ctx.Value("janus_context").(*janus.Account)
		acc := &janus.Account{}
		if err := json.NewDecoder(r.Body).Decode(acc); err != nil {
			views.Wrap(err, w)
			return
		}

		if jtk.OrganizationID != acc.OrganizationID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		if jtk.Role != "admin" {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		err := j.SetRights(acc)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Rights successfully set"))
		return
	}
}

func MakeAclHandler(r *httprouter.Router, j *janus.Janus) {
	r.HandlerFunc("POST", "/api/v2/acl/rights", middleware.JwtAuthentication(j.GetHandler(setRights(j))))
}
