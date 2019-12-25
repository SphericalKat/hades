package segment

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Repository interface {
	GetSegments(uint) ([]entities.EventSegment, error)
	GetParticipantsInSegment(uint) ([]entities.Participant, error)
	DeleteSegment(uint) error
	AddParticipantToSegment(string, uint, uint) error
	Find(day uint, eventID uint) (*entities.EventSegment, error)
}
