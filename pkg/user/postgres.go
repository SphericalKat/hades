package user

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB: db}
}

func (r *repo) Create(user *entities.User) error {
	err := r.DB.Save(user).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) Find(email string) (*entities.User, error) {
	user := &entities.User{Email: email}
	err := r.DB.Find(user).Association("Organizations").Find(&user.Organizations).Error

	switch err {
	case nil:
		return user, nil
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) Delete(email string) error {
	user := &entities.User{Email: email}
	err := r.DB.Delete(user).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) GetAllUsers(orgID uint) ([]entities.User, error) {
	if orgID == 0 {
		return nil, pkg.ErrNotFound
	}
	o := &entities.Organization{ID: orgID}
	err := r.DB.Find(o).Association("Users").Find(&o.Users).Error
	switch err {
	case nil:
		return o.Users, nil
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}
