package guest

type Repository interface {
	Save(*Guest) error
	Find(string) (*Guest, error)
	FindAll() ([]Guest, error)
	Delete(string) error
}
