package client

const (
	BaseUrl                = "https://sentry.io/api/0/"
	OrganizationsUrl       = BaseUrl + "organizations/"
	OrganizationMembersUrl = OrganizationsUrl + "%s/members/"
	OrganizationTeamsUrl   = OrganizationsUrl + "%s/teams/"

	// grant organization https://docs.sentry.io/api/organizations/add-a-member-to-an-organization/
	//revoke organization https://docs.sentry.io/api/organizations/delete-an-organization-member/

	//	https://docs.sentry.io/api/teams/list-a-teams-members/
	//	first %s is the org ID, second %s is the team slug
	TeamMembersUrl = BaseUrl + "teams/%s/%s/members/"

	// grant team https://docs.sentry.io/api/teams/add-an-organization-member-to-a-team/
	// revoke team https://docs.sentry.io/api/teams/delete-an-organization-member-from-a-team/

	// PROJECTS
	// list projects https://docs.sentry.io/api/projects/list-your-projects/

	// list project members  https://docs.sentry.io/api/projects/list-a-projects-users/
	// or also https://docs.sentry.io/api/projects/list-a-projects-organization-members/

	// grant teams to projects https://docs.sentry.io/api/projects/add-a-team-to-a-project/
	// revoke team from project https://docs.sentry.io/api/projects/delete-a-team-from-a-project/

)
