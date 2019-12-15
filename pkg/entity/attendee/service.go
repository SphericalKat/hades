package attendee

import "hades-2.0/pkg/entity"

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{r}
}

func (s *Service) Save(p *entity.Attendee) error {
	return s.repo.Save(p)
}

func (s *Service) FindByRegNo(regNo string) (*entity.Attendee, error) {
	return s.repo.FindByRegNo(regNo)
}

func (s *Service) Delete(regNo string) error {
	return s.repo.Delete(regNo)
}
