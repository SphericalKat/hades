package organization

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Repository interface {
	Save(organization *entities.Organization) (*entities.Organization, error)
	FindAll() ([]entities.Organization, error)
	Find(orgID uint) (*entities.Organization, error)
	Delete(orgID uint) error
	SaveJoinReq(request *entities.JoinRequest) error
	FindAllJoinReq(orgID uint) ([]entities.JoinRequest, error)
	AcceptJoinReq(request *entities.JoinRequest) error
	AddUserToOrg(orgID uint, email string) error
}
