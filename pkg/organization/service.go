package organization

import (
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	SaveOrg(org *entities.Organization) error
	GetAllOrgs() ([]entities.Organization, error)
	DelOrg(orgID uint) error
	GetOrgEvents(orgID uint) ([]entities.Event, error)
	SendJoinRequest(orgID uint, email string) error
	GetAllJoinReqs(orgID uint) ([]entities.JoinRequest, error)
	AcceptJoinReq(orgID uint, email string) error
	LoginOrg(orgID uint, email string) (*jwt.Token, error)
}

type orgSvc struct {
	repo Repository
}

func NewOrganizationService(rp Repository) Service {
	return &orgSvc{repo: rp}
}

func (o *orgSvc) GetAllJoinReqs(orgID uint) ([]entities.JoinRequest, error) {
	return o.repo.FindAllJoinReq(orgID)
}

func (o *orgSvc) AcceptJoinReq(orgID uint, email string) error {
	return o.repo.AcceptJoinReq(&entities.JoinRequest{OrganizationID: orgID, Email: email})
}

func (o *orgSvc) SendJoinRequest(orgID uint, email string) error {
	return o.repo.SaveJoinReq(&entities.JoinRequest{OrganizationID: orgID, Email: email})
}

func (o *orgSvc) SaveOrg(organization *entities.Organization) error {
	return o.repo.Save(organization)
}

func (o *orgSvc) GetAllOrgs() ([]entities.Organization, error) {
	return o.repo.FindAll()
}

func (o *orgSvc) DelOrg(orgID uint) error {
	return o.repo.Delete(orgID)
}

func (o *orgSvc) GetOrgEvents(orgID uint) ([]entities.Event, error) {
	org, err := o.repo.Find(orgID)
	if err != nil {
		return nil, err
	}

	return org.Events, err
}

func (o *orgSvc) LoginOrg(orgID uint, email string) (*jwt.Token, error) {
	org, err := o.repo.Find(orgID)
	if err != nil {
		return nil, err
	}

	for _, user := range org.Users {
		if user.Email == email {
			tk := middleware.Token{Email: email, OrgID: orgID, Role: "admin"}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
			return token, nil
		}
	}

	return nil, pkg.ErrNotFound
}
