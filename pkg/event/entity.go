package event

import (
	"hades-2.0/pkg/participant"
	"time"
)

type Event struct {
	EventID               string                    `json:"event_id" gorm:"primary_key;AUTO_INCREMENT"`
	ClubName              string                    `json:"club_name"`
	Name                  string                    `json:"name"`
	Budget                string                    `json:"budget"`
	Description           string                    `json:"description"`
	Category              string                    `json:"category"`
	Venue                 string                    `json:"venue"`
	Attendance            string                    `json:"attendance"`
	ExpectedParticipants  string                    `json:"expected_participants"`
	PROrequest            string                    `json:"PROrequest"`
	CampusEngineerRequest string                    `json:"campus_engineer_request"`
	Duration              string                    `json:"duration"`
	Status                string                    `json:"status"`
	ToDate                time.Time                 `json:"to_date"`
	FromDate              time.Time                 `json:"from_date"`
	ToTime                time.Time                 `json:"to_time"`
	FromTime              time.Time                 `json:"from_time"`
	Attendees             []participant.Participant `json:"attendees" gorm:"many2many:participant_events;"`
	// MainSponsor           Participant `json:"mainSponsor"`
	//StudentCoordinator    Participant `json:"studentCoordinator"`
	//FacultyCoordinator    Participant `json:"facultyCoordinator"`
}
