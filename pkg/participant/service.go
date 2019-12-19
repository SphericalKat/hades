package participant

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
)

type Service interface {
	CreateAttendee(participant *entities.Participant, eventID uint) error
	DeleteAttendee(regNo string) error
	ReadAttendee(regNo string, eventID uint) (*entities.Participant, error)
	RemoveAttendeeEvent(regNo string, eventID uint) error
}

type participantSvc struct {
	repo Repository
}

func NewParticipantService(r Repository) Service {
	return &participantSvc{repo: r}
}

func (s *participantSvc) CreateAttendee(participant *entities.Participant, eventID uint) error {
	a, err := s.repo.FindByRegNo(participant.RegNo, eventID)
	if err != nil && a != nil {
		return pkg.ErrAlreadyExists
	}
	return s.repo.Save(participant, eventID)
}

func (s *participantSvc) DeleteAttendee(regNo string) error {
	return s.repo.Delete(regNo)
}

func (s *participantSvc) ReadAttendee(regNo string, eventID uint) (*entities.Participant, error) {
	return s.repo.FindByRegNo(regNo, eventID)
}

func (s *participantSvc) RemoveAttendeeEvent(regNo string, eventID uint) error {
	return s.repo.RemoveAttendeeEvent(regNo, eventID)
}
