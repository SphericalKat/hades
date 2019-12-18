package user

type Repository interface {
	Create(user *User) error
	Find(email string) (*User, error)
	Delete(email string) error
}