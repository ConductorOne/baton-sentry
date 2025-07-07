package client

const (
	BaseUrl                  = "https://sentry.io/api/0/"
	OrganizationsUrl         = BaseUrl + "organizations/"
	OrganizationMembersUrl   = OrganizationsUrl + "%s/members/"
	OrganizationOneMemberUrl = OrganizationsUrl + "%s/members/%s/"
	OrganizationTeamsUrl     = OrganizationsUrl + "%s/teams/"
	OrganizationProjectsUrl  = OrganizationsUrl + "%s/projects/"

	//https://docs.sentry.io/api/teams/list-a-teams-members/
	//	teams/{organization_id_or_slug}/{team_id_or_slug}/members/
	TeamMembersUrl = BaseUrl + "teams/%s/%s/members/"

	//- grant team member https://docs.sentry.io/api/teams/add-an-organization-member-to-a-team/
	//- revoke team member https://docs.sentry.io/api/teams/delete-an-organization-member-from-a-team/
	//
	//	organizations/{organization_id_or_slug}/members/{member_id}/teams/{team_id_or_slug}/
	ProvisionTeamMemberUrl = OrganizationMembersUrl + "%s/teams/%s/"

	//	projects/{organization_id_or_slug}/{project_id_or_slug}/
	ProjectsUrl = BaseUrl + "projects/%s/%s/"

	//	projects/{organization_id_or_slug}/{project_id_or_slug}/
	ProjectMembersUrl = ProjectsUrl + "members/"

	// provision project members
	//	projects/{organization_id_or_slug}/{project_id_or_slug}/teams/{team_id_or_slug}/
	ProvisionProjectTeamUrl = ProjectsUrl + "teams/%s/"

	// TODO:
	// grant organization https://docs.sentry.io/api/organizations/add-a-member-to-an-organization/
	//revoke organization https://docs.sentry.io/api/organizations/delete-an-organization-member/
)
