package organization

type Repository interface {
	Save(organization *Organization) error
	FindAll() ([]Organization, error)
	Find(orgID uint) (*Organization, error)
	Delete(orgID uint) error
}
