package participant

import "hades-2.0/pkg"

type Service interface {
	CreateAttendee(participant *Participant, eventID string) error
	DeleteAttendee(regNo string) error
	ReadAttendee(regNo string) (*Participant, error)
	RemoveAttendeeEvent(regNo string, eventName string, clubName string) error
}

type participantSvc struct {
	repo Repository
}

func NewParticipantService(r Repository) Service {
	return &participantSvc{repo: r}
}

func (s *participantSvc) CreateAttendee(participant *Participant, eventID string) error {
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

func (s *participantSvc) RemoveAttendeeEvent(regNo string, eventName string, clubName string) error {
	return s.repo.RemoveAttendeeEvent(regNo, eventName)
}
