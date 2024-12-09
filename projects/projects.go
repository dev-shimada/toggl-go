package projects

import "fmt"

const (
	getWorkspaceProjectsUsers = "/api/v9/workspaces/%d/project_users"
)

type GetWorkspaceProjectsUsersQuery struct {
	ProjectIds       string
	UserId           string
	WithGroupMembers bool
}
type GetWorkspaceProjectsUsersInput struct {
	WorkspaceID *int
	Query       GetWorkspaceProjectsUsersQuery
}
type GetWorkspaceProjectsUsersOutput struct {
	At                   string `json:"at"`         //	When was last modified
	Gid                  int    `json:"gid"`        //	Group ID, legacy field
	GroupId              int    `json:"group_id"`   //	Group ID
	Id                   int    `json:"id"`         //	Project User ID
	LaborCost            *int   `json:"labor_cost"` //	null
	LaborCostLastUpdated string `json:""`           //	Date for labor cost last updated
	Manager              bool   `json:""`           //	Whether the user is manager of the project
	ProjectId            int    `json:""`           //	Project ID
	Rate                 *int   `json:""`           //	null
	RateLastUpdated      string `json:""`           //	Date for rate last updated
	UserId               int    `json:""`           //	User ID
	WorkspaceId          int    `json:""`           //	Workspace ID
}

func (c Client) GetWorkspaceProjectsUsers(workspaceID int) string {
	return fmt.Sprintf(getWorkspaceProjectsUsers, workspaceID)
}
