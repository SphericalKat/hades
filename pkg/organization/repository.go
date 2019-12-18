package organization

type Repository interface {
	Create(*Organization) error
	FindAll() ([]Organization, error)
	Find(string) (*Organization, error)
	Delete(string) error
	Update(string, *Organization) error
}
