package guest

type Service interface {
	CreateGuest(*Guest) error
	GetGuest(string) (*Guest, error)
	GetAllGuests() ([]Guest, error)
	DeleteGuest(string) error
}

type guestSvc struct {
	repo Repository
}

func NewGuestService(repo Repository) Service {
	return &guestSvc{repo: repo}
}

func (s *guestSvc) CreateGuest(guest *Guest) error {
	return s.repo.Save(guest)
}

func (s *guestSvc) GetGuest(email string) (*Guest, error) {
	return s.repo.Find(email)
}

func (s *guestSvc) GetAllGuests() ([]Guest, error) {
	return s.repo.FindAll()
}

func (s *guestSvc) DeleteGuest(email string) error {
	return s.repo.Delete(email)
}