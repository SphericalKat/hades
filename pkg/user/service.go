package user

import (
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/ATechnoHazard/hades-2/pkg/organization"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(user *User) (*jwt.Token, error)
	Login(email string, password string) (*jwt.Token, error)
	GetUserOrgs(email string) ([]organization.Organization, error)
}

func NewUserService(r Repository) Service {
	return &userSvc{repo: r}
}

type userSvc struct {
	repo Repository
}

func (u *userSvc) CreateUser(user *User) (*jwt.Token, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	user.Password = string(hashedPass)
	err = u.repo.Create(user)
	if err != nil {
		return nil, err
	}
	tk := middleware.Token{Email: user.Email}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	return token, err
}

func (u *userSvc) Login(email string, password string) (*jwt.Token, error) {
	user, err := u.repo.Find(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	switch err {
	case nil:
		tk := middleware.Token{Email: user.Email}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
		return token, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, pkg.ErrUnauthorized
	default:
		return nil, pkg.ErrDatabase
	}
}

func (u *userSvc) GetUserOrgs(email string) ([]organization.Organization, error) {
	user, err := u.repo.Find(email)
	if err != nil {
		return nil, err
	}
	return user.Organizations, err
}
