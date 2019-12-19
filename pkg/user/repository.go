package user

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Repository interface {
	Create(user *entities.User) error
	Find(email string) (*entities.User, error)
	Delete(email string) error
}
