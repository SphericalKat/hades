package organization

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	SaveOrg(org *Organization) error
	GetAllOrgs() ([]Organization, error)
	DelOrg(orgID uint) error
	GetOrgEvents(orgID uint) ([]event.Event, error)
}

type orgSvc struct {
	repo Repository
}

func NewOrganizationService(rp Repository) Service {
	return &orgSvc{repo: rp}
}

func (o *orgSvc) SaveOrg(organization *Organization) error {
	return o.repo.Save(organization)
}

func (o *orgSvc) GetAllOrgs() ([]Organization, error) {
	return o.repo.FindAll()
}

func (o *orgSvc) DelOrg(orgID uint) error {
	return o.repo.Delete(orgID)
}

func (o *orgSvc) GetOrgEvents(orgID uint) ([]event.Event, error) {
	org, err := o.repo.Find(orgID)
	if err != nil {
		return nil, err
	}

	return org.Events, err
}

func (o *orgSvc) LoginOrg(orgID uint, email string) (*jwt.Token, error) {

	return nil, pkg.ErrNotFound

}
