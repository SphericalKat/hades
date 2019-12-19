package handler

import (
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/organization"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func acceptJoinRequest(oSvc organization.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		if err := oSvc.AcceptJoinReq(tk.OrgID, tk.Email); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Join request accepted"))
	}
}

func sendJoinRequest(oSvc organization.Service) http.HandlerFunc {

}

func MakeOrgHandler(r *httprouter.Router, oSvc organization.Service) {
	r.HandlerFunc("POST","/api/v1/org/accept", middleware.JwtAuthentication(acceptJoinRequest(oSvc)))
}
