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

func deleteSegment(segmentService segment.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seg := &entities.EventSegment{}
		eve := &entities.Event{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(seg); err != nil {
			views.Wrap(err, w)
			return
		}

		seg, err := segmentService.ReadEventSegment(seg.Day, seg.EventID)
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

		if err := segmentService.DeleteSegment(seg.Day); err != nil {
			views.Wrap(err, w)
			return
		}

		utils.Respond(w, utils.Message(http.StatusOK, "Deleted segment successfully."))
	}
}

func getSegments(segmentService segment.Service, eventService event.Service) http.HandlerFunc {
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

func getParticipantsInSegment(segmentService segment.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seg := &entities.EventSegment{}
		eve := &entities.Event{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(seg); err != nil {
			views.Wrap(err, w)
			return
		}

		seg, err := segmentService.ReadEventSegment(seg.Day, seg.EventID)
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

		msg := utils.Message(http.StatusOK, "Retrieved all participants in event segment successfully")
		msg["participants"] = seg.PresentParticipants
		utils.Respond(w, msg)
	}
}

func markPresent(segmentService segment.Service, eventService event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		composite := &views.SegmentParticipantComposite{}
		eve := &entities.Event{}

		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		if err := json.NewDecoder(r.Body).Decode(composite); err != nil {
			views.Wrap(err, w)
			return
		}

		seg, err := segmentService.ReadEventSegment(composite.Day, composite.EventID)
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
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from modifying this resource."))
			return
		}

		for _, part := range eve.Attendees {
			if part.RegNo == composite.RegNo {
				if err := segmentService.AddParticipantToSegment(composite.RegNo, composite.Day, composite.EventID); err != nil {
					views.Wrap(err, w)
					return
				}

				utils.Respond(w, utils.Message(http.StatusOK, "Participant Attendance recorded successfully."))
				return
			}
		}

		utils.Respond(w, utils.Message(http.StatusNotFound, "Participant doesn't exist in this event"))
	}
}

func MakeEventSegmentHandler(r *httprouter.Router, segmentService segment.Service, eventService event.Service) {
	r.HandlerFunc("POST", "/api/v2/participant/get-days",
		middleware.JwtAuthentication(getSegments(segmentService, eventService)))
	r.HandlerFunc("POST", "/api/v2/participant/get-present",
		middleware.JwtAuthentication(getParticipantsInSegment(segmentService, eventService)))
	r.HandlerFunc("POST", "/api/v2/participant/mark-present",
		middleware.JwtAuthentication(markPresent(segmentService, eventService)))
}
