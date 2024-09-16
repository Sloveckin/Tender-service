package repositories

type OrgRepository interface {
	OrganizationExistsById(organisationId string) (bool, error)
}
