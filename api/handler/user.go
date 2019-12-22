package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/ATechnoHazard/hades-2/pkg/user"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"time"
)

func signUp(uSvc user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acc := &entities.User{}
		if err := json.NewDecoder(r.Body).Decode(acc); err != nil {
			views.Wrap(err, w)
			return
		}

		acc.CreatedAt = time.Now()
		tk, err := uSvc.CreateUser(acc)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		tkString, err := tk.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := u.Message(http.StatusOK, "User account successfully saved")
		msg["token"] = tkString
		u.Respond(w, msg)
	}
}

func login(uSvc user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acc := &entities.User{}
		if err := json.NewDecoder(r.Body).Decode(acc); err != nil {
			views.Wrap(err, w)
			return
		}

		tk, err := uSvc.Login(acc.Email, acc.Password)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		tkString, err := tk.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := u.Message(http.StatusOK, "Successfully logged in")
		msg["token"] = tkString
		u.Respond(w, msg)
	}
}

func getUserOrgs(uSvc user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		orgs, err := uSvc.GetUserOrgs(tk.Email)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := u.Message(http.StatusOK, "Successfully retrieved user organizations")
		msg["organizations"] = orgs
		u.Respond(w, msg)
	}
}

func getAllUsers(uSvc user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		users, err := uSvc.GetOrgUsers(tk.OrgID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := u.Message(http.StatusOK, "Successfully retrieved organization users")
		msg["users"] = users
		u.Respond(w, msg)
		return
	}
}

func MakeUserHandler(r *httprouter.Router, uSvc user.Service) {
	r.HandlerFunc("POST", "/api/v2/org/signup", signUp(uSvc))
	r.HandlerFunc("POST", "/api/v2/org/login", login(uSvc))
	r.HandlerFunc("GET", "/api/v2/org/", middleware.JwtAuthentication(getUserOrgs(uSvc)))
	r.HandlerFunc("GET", "/api/v2/org/users", middleware.JwtAuthentication(getAllUsers(uSvc)))
}
