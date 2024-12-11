package projects

// const (
// 	getWorkspaceProjectsUsers = "/api/v9/workspaces/%d/project_users"
// )

type GetWorkspaceProjectsUsersQuery struct {
	ProjectIds       string
	UserId           string
	WithGroupMembers bool
}
type GetWorkspaceProjectsUsersInput struct {
	WorkspaceID int
	Query       GetWorkspaceProjectsUsersQuery
}
type GetWorkspaceProjectsUsersOutput struct {
	At                   string `json:"at"`                      //	When was last modified
	Gid                  int    `json:"gid"`                     //	Group ID, legacy field
	GroupId              int    `json:"group_id"`                //	Group ID
	Id                   int    `json:"id"`                      //	Project User ID
	LaborCost            int    `json:"labor_cost"`              //	null
	LaborCostLastUpdated string `json:"labor_cost_last_updated"` //	Date for labor cost last updated
	Manager              bool   `json:"manager"`                 //	Whether the user is manager of the project
	ProjectId            int    `json:"prpoject_id"`             //	Project ID
	Rate                 int    `json:"rate"`                    //	null
	RateLastUpdated      string `json:"rate_last_updated"`       //	Date for rate last updated
	UserId               int    `json:"user_id"`                 //	User ID
	WorkspaceId          int    `json:"workspace_id"`            //	Workspace ID
}

func (c Client) GetWorkspaceProjectsUsers(input GetWorkspaceProjectsUsersInput) (GetWorkspaceProjectsUsersOutput, error) {
	return GetWorkspaceProjectsUsersOutput{}, nil
}
