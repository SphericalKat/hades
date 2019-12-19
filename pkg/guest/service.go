package guest

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Service interface {
	SaveGuestEvent(*entities.Guest, uint) error
	GetGuestEvent(string, uint) (*entities.Guest, error)
	GetAllGuestEvent(uint) ([]entities.Guest, error)
	DeleteGuest(string) error
	RemoveGuestEvent(string, uint) error
}

type guestSvc struct {
	repo Repository
}

func NewGuestService(repo Repository) Service {
	return &guestSvc{repo: repo}
}

func (s *guestSvc) SaveGuestEvent(guest *entities.Guest, eventId uint) error {
	return s.repo.SaveGuestEvent(guest, eventId)
}

func (s *guestSvc) GetGuestEvent(email string, eventId uint) (*entities.Guest, error) {
	return s.repo.FindGuestEvent(email, eventId)
}

func (s *guestSvc) GetAllGuestEvent(eventId uint) ([]entities.Guest, error) {
	return s.repo.FindAllGuestEvent(eventId)
}

func (s *guestSvc) DeleteGuest(email string) error {
	return s.repo.Delete(email)
}

func (s *guestSvc) RemoveGuestEvent(email string, eventId uint) error {
	return s.repo.RemoveGuestEvent(email, eventId)
}