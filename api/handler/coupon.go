package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/coupon"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/ATechnoHazard/janus"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func saveCoupon(couponService coupon.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		coup := &entities.Coupon{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		jtk := ctx.Value("janus_context").(*janus.Account)

		if jtk.Role != "admin" {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(coup); err != nil {
			views.Wrap(err, w)
			return
		}

		eve := &entities.Event{}
		eve, err := eventService.ReadEvent(coup.EventId)

		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		if err := couponService.SaveCoupon(coup); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "saved coupon successfully."))
	}
}

func deleteCoupon(couponService coupon.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		coup := &entities.Coupon{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		jtk := ctx.Value("janus_context").(*janus.Account)

		if jtk.Role != "admin" {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(coup); err != nil {
			views.Wrap(err, w)
			return
		}

		eve := &entities.Event{}
		eve, err := eventService.ReadEvent(coup.EventId)

		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		v, err := couponService.VerifyCoupon(coup.EventId, coup.CouponId)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if !v {
			u.Respond(w, u.Message(http.StatusConflict, "The coupon and event are not related."))
			return
		}

		if err := couponService.DeleteCoupon(coup.CouponId); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Coupon deleted successfully."))
	}
}

func redeemCoupon(couponService coupon.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		composite := &views.CouponParticipantComposite{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(composite); err != nil {
			views.Wrap(err, w)
			return
		}

		eve := &entities.Event{}
		eve, err := eventService.ReadEvent(composite.EventID)

		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		v, err := couponService.VerifyCoupon(composite.EventID, composite.CouponID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if !v {
			u.Respond(w, u.Message(http.StatusConflict, "The coupon and event are not related."))
			return
		}

		if err := couponService.RemoveCouponParticipant(composite.CouponID, composite.RegNo); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Successfully redeemed coupon."))
	}
}

func addAllCoupons(couponService coupon.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		coup := &entities.Coupon{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(coup); err != nil {
			views.Wrap(err, w)
			return
		}

		eve := &entities.Event{}
		eve, err := eventService.ReadEvent(coup.EventId)

		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		if err := couponService.AddCouponsToAll(eve.ID); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Added coupons to all participants successfully."))
	}
}

func GetCoupons(couponService coupon.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		coup := &entities.Coupon{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(coup); err != nil {
			views.Wrap(err, w)
			return
		}

		eve := &entities.Event{}
		eve, err := eventService.ReadEvent(coup.EventId)

		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		var coups []entities.Coupon

		coups, err = couponService.GetCoupons(coup.EventId)

		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := u.Message(http.StatusOK, "Retrieved all coupons successfully")
		msg["coupons"] = coups
		u.Respond(w, msg)
	}
}

func MakeCouponHandler(r *httprouter.Router, couponService coupon.Service, eventService event.Service, j *janus.Janus) {
	r.HandlerFunc("POST", "/api/v2/coupon/save-coupon",
		middleware.JwtAuthentication(j.GetHandler(saveCoupon(couponService, eventService))))
	r.HandlerFunc("DELETE", "/api/v2/coupon/delete-coupon",
		middleware.JwtAuthentication(j.GetHandler(deleteCoupon(couponService, eventService))))
	r.HandlerFunc("POST", "/api/v2/coupon/redeem-coupon",
		middleware.JwtAuthentication(redeemCoupon(couponService, eventService)))
	r.HandlerFunc("POST", "/api/v2/coupon/add-all-coupons",
		middleware.JwtAuthentication(j.GetHandler(addAllCoupons(couponService, eventService))))
	r.HandlerFunc("POST", "/api/v2/coupon/get-coupons",
		middleware.JwtAuthentication(GetCoupons(couponService, eventService)))
}
