package participant

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
)

type Service interface {
	CreateAttendee(participant *entities.Participant, eventID uint) error
	DeleteAttendee(email string) error
	SaveAttendee(participant *entities.Participant, eventID uint) error
	ReadAttendee(email string, eventID uint) (*entities.Participant, error)
	RemoveAttendeeEvent(email string, eventID uint) error
}

type participantSvc struct {
	repo Repository
}

func (s *participantSvc) SaveAttendee(participant *entities.Participant, eventID uint) error {
	return s.repo.Save(participant, eventID)
}

func NewParticipantService(r Repository) Service {
	return &participantSvc{repo: r}
}

func (s *participantSvc) CreateAttendee(participant *entities.Participant, eventID uint) error {
	a, err := s.repo.FindByEmail(participant.Email, eventID)
	if err != nil && a != nil {
		return pkg.ErrAlreadyExists
	}
	return s.repo.Save(participant, eventID)
}

func (s *participantSvc) DeleteAttendee(email string) error {
	return s.repo.Delete(email)
}

func (s *participantSvc) ReadAttendee(email string, eventID uint) (*entities.Participant, error) {
	return s.repo.FindByEmail(email, eventID)
}

func (s *participantSvc) RemoveAttendeeEvent(email string, eventID uint) error {
	return s.repo.RemoveAttendeeEvent(email, eventID)
}
