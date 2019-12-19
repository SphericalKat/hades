package event

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Repository interface {
	Find(eventID uint) (*entities.Event, error)
	Save(event *entities.Event) error
	Delete(eventID uint) error
}
