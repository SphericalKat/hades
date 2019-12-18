package organization

import (
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	SaveOrg(org *Organization) error
	GetAllOrgs() ([]Organization, error)
	DelOrg(orgID uint) error
	GetOrgEvents(orgID uint) ([]event.Event, error)
	SendJoinRequest(orgID uint, email string) error
	GetAllJoinReqs(orgID uint) ([]JoinRequest, error)
	AcceptJoinReq(orgID uint, email string) error
}

type orgSvc struct {
	repo Repository
}

func NewOrganizationService(rp Repository) Service {
	return &orgSvc{repo: rp}
}

func (o *orgSvc) GetAllJoinReqs(orgID uint) ([]JoinRequest, error) {
	return o.repo.FindAllJoinReq(orgID)
}

func (o *orgSvc) AcceptJoinReq(orgID uint, email string) error {
	return o.repo.AcceptJoinReq(&JoinRequest{OrgID: orgID, Email: email})
}

func (o *orgSvc) SendJoinRequest(orgID uint, email string) error {
	return o.repo.SaveJoinReq(&JoinRequest{OrgID: orgID, Email: email})
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
	tk := middleware.Token{Email: email, OrgID: orgID, Role: "admin"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	return token, nil
}
