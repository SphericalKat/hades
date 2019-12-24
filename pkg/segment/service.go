package segment

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Service interface {
	GetSegments(uint) ([]entities.EventSegment, error)
	GetParticipantsInSegment(uint) ([]entities.Participant, error)
	DeleteSegment(uint) error
	AddParticipantToSegment(regNo string, day uint, eventID uint) error
	ReadEventSegment(day uint, eventID uint) (*entities.EventSegment, error)
}

type eventSegSvc struct {
	repo Repository
}

func NewEventSegmentService(repo Repository) Service {
	return &eventSegSvc{repo:repo}
}


func (s *eventSegSvc) GetSegments(eventId uint) ([]entities.EventSegment, error) {
	return s.repo.GetSegments(eventId)
}

func (s *eventSegSvc) GetParticipantsInSegment(segmentId uint) ([]entities.Participant, error) {
	return s.repo.GetParticipantsInSegment(segmentId)
}

func (s *eventSegSvc) DeleteSegment(segmentId uint) error {
	return s.repo.DeleteSegment(segmentId)
}

func (s *eventSegSvc) AddParticipantToSegment(regNo string, day uint, eventID uint) error {
	return s.repo.AddParticipantToSegment(regNo, day, eventID)
}

func (s *eventSegSvc) ReadEventSegment(day uint, eventID uint) (*entities.EventSegment, error) {
	return s.repo.Find(day, eventID)
}