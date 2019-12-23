package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	"github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/coupon"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func SaveCoupon(couponService coupon.Service, eventService event.Service) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
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
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		if err := couponService.SaveCoupon(coup); err != nil {
			views.Wrap(err, w)
			return
		}

		utils.Respond(w, utils.Message(http.StatusOK, "saved coupon successfully."))
	}
}

func DeleteCoupon(couponService coupon.Service, eventService event.Service) http.HandlerFunc {
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
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		v, err := couponService.VerifyCoupon(coup.CouponId, coup.EventId)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if !v {
			utils.Respond(w, utils.Message(http.StatusConflict, "The coupon and event are not related."))
			return
		}


		if err := couponService.DeleteCoupon(coup.CouponId); err != nil {
			views.Wrap(err, w)
			return
		}

		utils.Respond(w, utils.Message(http.StatusOK, "Coupon deleted successfully."))
	}
}

func RedeemCoupon(couponService coupon.Service, eventService event.Service) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
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
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		v, err := couponService.VerifyCoupon(composite.EventID, composite.CouponID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if !v {
			utils.Respond(w, utils.Message(http.StatusConflict, "The coupon and event are not related."))
			return
		}

		if err := couponService.RemoveCouponParticipant(composite.CouponID, composite.RegNo); err != nil {
			views.Wrap(err, w)
			return
		}

		utils.Respond(w, utils.Message(http.StatusOK, "Successfully redeemed coupon."))
	}
}

func AddAllCoupons(couponService coupon.Service, eventService event.Service) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
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
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		if err := couponService.AddCouponsToAll(eve.ID); err != nil {
			views.Wrap(err, w)
			return
		}

		utils.Respond(w, utils.Message(http.StatusOK, "Added coupons to all participants successfully."))
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
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		var coups []entities.Coupon

		coups, err = couponService.GetCoupons(coup.EventId)

		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := utils.Message(http.StatusOK, "Retrieved all coupons successfully")
		msg["coupons"] = coups
		utils.Respond(w, msg)
	}
}

func MakeCouponHandler(r *httprouter.Router, couponService coupon.Service, eventService event.Service) {
	r.HandlerFunc("POST", "/api/v2/coupon/save-coupon",
		middleware.JwtAuthentication(SaveCoupon(couponService, eventService)))
	r.HandlerFunc("DELETE", "/api/v2/coupon/delete-coupon",
		middleware.JwtAuthentication(DeleteCoupon(couponService, eventService)))
	r.HandlerFunc("POST", "/api/v2/coupon/redeem-coupon",
		middleware.JwtAuthentication(RedeemCoupon(couponService, eventService)))
	r.HandlerFunc("POST", "/api/v2/coupon/add-all-coupons",
		middleware.JwtAuthentication(AddAllCoupons(couponService, eventService)))
	r.HandlerFunc("POST", "/api/v2/coupon/get-coupons",
		middleware.JwtAuthentication(GetCoupons(couponService, eventService)))
}