package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	"github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/ATechnoHazard/hades-2/pkg/segment"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func AddSegment(segmentService segment.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seg := &entities.EventSegment{}
		eve := &entities.Event{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(seg); err != nil {
			views.Wrap(err, w)
			return
		}

		eve, err := eventService.ReadEvent(seg.EventID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		if err := segmentService.AddSegment(seg); err != nil {
			views.Wrap(err, w)
			return
		}

		utils.Respond(w, utils.Message(http.StatusOK, "Added event segment successfully."))
	}
}

func DeleteSegment(segmentService segment.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seg := &entities.EventSegment{}
		eve := &entities.Event{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(seg); err != nil {
			views.Wrap(err, w)
			return
		}

		seg, err := segmentService.ReadEventSegment(seg.SegmentID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		eve, err = eventService.ReadEvent(seg.EventID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resouce."))
			return
		}

		if err := segmentService.DeleteSegment(seg.SegmentID); err != nil {
			views.Wrap(err, w)
			return
		}

		utils.Respond(w, utils.Message(http.StatusOK, "Deleted segment successfully."))
	}
}

func GetSegments(segmentService segment.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eve := &entities.Event{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(eve); err != nil {
			views.Wrap(err, w)
			return
		}

		eve, err := eventService.ReadEvent(eve.ID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		segments, err := segmentService.GetSegments(eve.ID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := utils.Message(http.StatusOK, "Retrieved all event segments successfully")
		msg["segments"] = segments
		utils.Respond(w, msg)
	}
}

func GetParticipantsInSegment(segmentService segment.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seg := &entities.EventSegment{}
		eve := &entities.Event{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(seg); err != nil {
			views.Wrap(err, w)
			return
		}

		seg, err := segmentService.ReadEventSegment(seg.SegmentID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		eve, err = eventService.ReadEvent(seg.EventID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resouce."))
			return
		}

		peeps, err := segmentService.GetParticipantsInSegment(seg.SegmentID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := utils.Message(http.StatusOK, "Retrieved all participants in event segment successfully")
		msg["participants"] = peeps
		utils.Respond(w, msg)
	}
}

func AttendEventSegment(segmentService segment.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		composite := &views.SegmentParticipantComposite{}
		eve := &entities.Event{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(composite); err != nil {
			views.Wrap(err, w)
			return
		}

		seg, err := segmentService.ReadEventSegment(composite.SegmentID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		eve, err = eventService.ReadEvent(seg.EventID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if tk.OrgID != eve.OrganizationID {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resouce."))
			return
		}

		if err := segmentService.AddParticipantToSegment(composite.RegNo, composite.SegmentID); err != nil {
			views.Wrap(err, w)
			return
		}

		utils.Respond(w, utils.Message(http.StatusOK, "Participant Attendance recorded successfully."))
	}
}

func MakeEventSegmentHandlers(r *httprouter.Router, segmentService segment.Service, eventService event.Service) {
	r.HandlerFunc("POST", "/api/v2/event-segment/add-segment",
		middleware.JwtAuthentication(AddSegment(segmentService, eventService)))
	r.HandlerFunc("POST", "/api/v2/event-segment/rm-segment",
		middleware.JwtAuthentication(DeleteSegment(segmentService, eventService)))
	r.HandlerFunc("POST", "/api/v2/event-segment/get-segments",
		middleware.JwtAuthentication(GetSegments(segmentService, eventService)))
	r.HandlerFunc("POST", "/api/v2/event-segment/get-present-participants",
		middleware.JwtAuthentication(GetParticipantsInSegment(segmentService, eventService)))
	r.HandlerFunc("POST", "/api/v2/event-segment/mark-participant-present",
		middleware.JwtAuthentication(AttendEventSegment(segmentService, eventService)))
}

