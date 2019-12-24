package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	"github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/ATechnoHazard/hades-2/pkg/segment"
	"github.com/gocarina/gocsv"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func getParticipantsCSVData(eventService event.Service, segmentService segment.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seg := &entities.EventSegment{}
		eve := &entities.Event{}

		if err := json.NewDecoder(r.Body).Decode(seg); err != nil {
			views.Wrap(err, w)
			return
		}

		tk := r.Context().Value(middleware.JwtContextKey("token")).(*middleware.Token)

		eve, err := eventService.ReadEvent(seg.EventID)
		if err != nil {
			views.Wrap(err, w)
			return
		}
		if tk.OrgID != eve.OrganizationID {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from accessing this resource."))
			return
		}

		if _, err := segmentService.ReadEventSegment(seg.Day, seg.EventID); err != nil {
			views.Wrap(err, w)
			return
		}

		var peeps []entities.Participant
		peepMap := make(map[string]*entities.CSVParticipant)

		// Populate all peeps in the event and hash them into peepMap
		for _, peep := range eve.Attendees {
			peepMap[peep.RegNo] = entities.P2CSVPTransform(&peep)
		}

		// populate all present peeps.
		peeps, err = segmentService.GetParticipantsInSegment(seg.Day)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		// mark present peeps as present.
		for _, peep := range peeps {
			peepMap[peep.RegNo].IsPresent = true
		}

		csvPeeps := make([]*entities.CSVParticipant, 0)

		// retrieve all values from peepMap
		for _, csvPeep := range peepMap {
			csvPeeps = append(csvPeeps, csvPeep)
		}

		csvContents, err := gocsv.MarshalString(&csvPeeps)
		if err != nil {
			utils.Respond(w, utils.Message(http.StatusInternalServerError,
				"Error occured during CSV generation."))
			return
		}

		msg := utils.Message(http.StatusOK, "Retrieved CSV data successfully.")
		msg["csv-data"] = csvContents
		utils.Respond(w, msg)
	}
}

func getParticipantsJSONData(eventService event.Service, segmentService segment.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seg := &entities.EventSegment{}
		eve := &entities.Event{}

		if err := json.NewDecoder(r.Body).Decode(seg); err != nil {
			views.Wrap(err, w)
			return
		}

		tk := r.Context().Value(middleware.JwtContextKey("token")).(*middleware.Token)

		eve, err := eventService.ReadEvent(seg.EventID)
		if err != nil {
			views.Wrap(err, w)
			return
		}
		if tk.OrgID != eve.OrganizationID {
			utils.Respond(w, utils.Message(http.StatusForbidden, "You are forbidden from accessing this resource."))
			return
		}

		if _, err := segmentService.ReadEventSegment(seg.Day, seg.EventID); err != nil {
			views.Wrap(err, w)
			return
		}

		var peeps []entities.Participant
		peepMap := make(map[string]*entities.CSVParticipant)

		// Populate all peeps in the event and hash them into peepMap
		for _, peep := range eve.Attendees {
			peepMap[peep.RegNo] = entities.P2CSVPTransform(&peep)
		}

		// populate all present peeps.
		peeps, err = segmentService.GetParticipantsInSegment(seg.Day)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		// mark present peeps as present.
		for _, peep := range peeps {
			peepMap[peep.RegNo].IsPresent = true
		}

		jsonPeeps := make([]*entities.CSVParticipant, 0)

		// retrieve all values from peepMap
		for _, jsonPeep := range peepMap {
			jsonPeeps = append(jsonPeeps, jsonPeep)
		}

		msg := utils.Message(http.StatusOK, "Retrieved JSON data successfully.")
		msg["data"] = jsonPeeps
		utils.Respond(w, msg)
	}
}

func MakeExporterHandlers(r *httprouter.Router, eventService event.Service, segmentService segment.Service) {
	r.HandlerFunc("POST", "/api/v2/export/csv",
		middleware.JwtAuthentication(getParticipantsCSVData(eventService, segmentService)))
	r.HandlerFunc("POST", "/api/v2/export/json",
		middleware.JwtAuthentication(getParticipantsJSONData(eventService, segmentService)))
}
