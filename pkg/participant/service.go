package participant

import "github.com/ATechnoHazard/hades-2/pkg"

type Service interface {
	CreateAttendee(participant *Participant, eventID uint) error
	DeleteAttendee(regNo string) error
	ReadAttendee(regNo string) (*Participant, error)
	RemoveAttendeeEvent(regNo string, eventID uint) error
}

type participantSvc struct {
	repo Repository
}

func NewParticipantService(r Repository) Service {
	return &participantSvc{repo: r}
}

func (s *participantSvc) CreateAttendee(participant *Participant, eventID uint) error {
	a, err := s.repo.FindByRegNo(participant.RegNo)
	if err != nil && a != nil {
		return pkg.ErrAlreadyExists
	}
	return s.repo.Save(participant, eventID)
}

func (s *participantSvc) DeleteAttendee(regNo string) error {
	return s.repo.Delete(regNo)
}

func (s *participantSvc) ReadAttendee(regNo string) (*Participant, error) {
	return s.repo.FindByRegNo(regNo)
}

func (s *participantSvc) RemoveAttendeeEvent(regNo string, eventID uint) error {
	return s.repo.RemoveAttendeeEvent(regNo, eventID)
}
