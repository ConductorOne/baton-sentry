package client

import "fmt"

// Path patterns (relative to base URL).
const (
	organizationsPath         = "organizations/"
	organizationMembersPath   = "organizations/%s/members/"
	organizationOneMemberPath = "organizations/%s/members/%s/"
	organizationTeamsPath     = "organizations/%s/teams/"
	organizationProjectsPath  = "organizations/%s/projects/"

	// https://docs.sentry.io/api/teams/retrieve-a-team/
	//	teams/{organization_id_or_slug}/{team_id_or_slug}/
	TeamUrl = BaseUrl + "teams/%s/%s/"

	//https://docs.sentry.io/api/teams/list-a-teams-members/
	//	teams/{organization_id_or_slug}/{team_id_or_slug}/members/
	teamMembersPath = "teams/%s/%s/members/"

	//- grant team member https://docs.sentry.io/api/teams/add-an-organization-member-to-a-team/
	//- revoke team member https://docs.sentry.io/api/teams/delete-an-organization-member-from-a-team/
	//
	//	organizations/{organization_id_or_slug}/members/{member_id}/teams/{team_id_or_slug}/
	provisionTeamMemberPath = "organizations/%s/members/%s/teams/%s/"

	//	projects/{organization_id_or_slug}/{project_id_or_slug}/
	projectsPath = "projects/%s/%s/"

	//	projects/{organization_id_or_slug}/{project_id_or_slug}/members/
	projectMembersPath = "projects/%s/%s/members/"

	// provision project members
	//	projects/{organization_id_or_slug}/{project_id_or_slug}/teams/{team_id_or_slug}/
	provisionProjectTeamPath = "projects/%s/%s/teams/%s/"

	// teams/{organization_id_or_slug}/{team_id_or_slug}/projects/.
	teamProjectsPath = "teams/%s/%s/projects/"
)

// URL builder methods.

func (c *Client) OrganizationsURL() string {
	return c.baseURL + organizationsPath
}

func (c *Client) OrganizationMembersURL(orgID string) string {
	return c.baseURL + fmt.Sprintf(organizationMembersPath, orgID)
}

func (c *Client) OrganizationOneMemberURL(orgID, memberID string) string {
	return c.baseURL + fmt.Sprintf(organizationOneMemberPath, orgID, memberID)
}

func (c *Client) OrganizationTeamsURL(orgID string) string {
	return c.baseURL + fmt.Sprintf(organizationTeamsPath, orgID)
}

func (c *Client) OrganizationProjectsURL(orgID string) string {
	return c.baseURL + fmt.Sprintf(organizationProjectsPath, orgID)
}

func (c *Client) TeamMembersURL(orgID, teamID string) string {
	return c.baseURL + fmt.Sprintf(teamMembersPath, orgID, teamID)
}

func (c *Client) ProvisionTeamMemberURL(orgID, memberID, teamID string) string {
	return c.baseURL + fmt.Sprintf(provisionTeamMemberPath, orgID, memberID, teamID)
}

func (c *Client) ProjectsURL(orgID, projectID string) string {
	return c.baseURL + fmt.Sprintf(projectsPath, orgID, projectID)
}

func (c *Client) ProjectMembersURL(orgID, projectID string) string {
	return c.baseURL + fmt.Sprintf(projectMembersPath, orgID, projectID)
}

func (c *Client) ProvisionProjectTeamURL(orgID, projectID, teamID string) string {
	return c.baseURL + fmt.Sprintf(provisionProjectTeamPath, orgID, projectID, teamID)
}

func (c *Client) TeamProjectsURL(orgID, teamID string) string {
	return c.baseURL + fmt.Sprintf(teamProjectsPath, orgID, teamID)
}
