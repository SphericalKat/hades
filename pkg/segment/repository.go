package segment

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Repository interface {
	GetSegments(uint) ([]entities.EventSegment, error)
	GetParticipantsInSegment(uint) ([]entities.Participant, error)
	DeleteSegment(uint) error
	AddPartipantToSegment(string, uint) error
	Find(uint) (*entities.EventSegment, error)
}
