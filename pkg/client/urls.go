package client

const (
	BaseUrl                = "https://sentry.io/api/0"
	OrganizationsUrl       = BaseUrl + "/organizations/"
	OrganizationMembersUrl = OrganizationsUrl + "%s/members/"
	OrganizationTeamsUrl   = OrganizationsUrl + "%s/teams/"
)
