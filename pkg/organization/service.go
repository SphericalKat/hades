package organization

type Service interface {
	CreateOrganization(*Organization) error
	GetOrganizations() ([]Organization, error)
	UpdateOrganization(string, *Organization) error
	DeleteOrganization(string) error
}

type organizationSvc struct {
	repo Repository
}

func NewOrganizationService(rp Repository) Service {
	return &organizationSvc{repo: rp}
}

func (sv *organizationSvc) CreateOrganization(organization *Organization) error {
	return sv.repo.Create(organization)
}

func (sv *organizationSvc) GetOrganizations() ([]Organization, error) {
	return sv.repo.FindAll()
}

func (sv *organizationSvc) UpdateOrganization(name string, organization *Organization) error {
	return sv.repo.Update(name, organization)
}

func (sv *organizationSvc) DeleteOrganization(name string) error {
	return sv.repo.Delete(name)
}
