package guest

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Repository interface {
	SaveGuestEvent(*entities.Guest, uint) error
	FindGuestEvent(string, uint) (*entities.Guest, error)
	FindAllGuestEvent(uint) ([]entities.Guest, error)
	Delete(string) error
	RemoveGuestEvent(string, uint) error
}
