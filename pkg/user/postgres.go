package user

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB: db}
}

func (r *repo) Create(user *User) error {
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

func (r *repo) Find(email string) (*User, error) {
	user := &User{Email: email}
	err := r.DB.Find(user).Association("Organizations").Find(user.Organizations).Error

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
	user := &User{Email: email}
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
