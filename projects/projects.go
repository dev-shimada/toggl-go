package projects

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

const (
	getWorkspaceProjectsUsers                    = "/api/v9/workspaces/%d/project_users"
	postAddAnUserIntoWorkspaceProjectsUsers      = "/api/v9/workspaces/%d/project_users"
	patchProjectUsersFromWorkspace               = "/api/v9/workspaces/%d/project_users/%s"
	putUpdateAnUserIntoWorkspaceProjectsUsers    = "/api/v9/workspaces/%d/project_users/%s"
	deleteAProjectUserFromWorkspaceProjectsUsers = "/api/v9/workspaces/%d/project_users/%s"
	getWorkspaceProjects                         = "/api/v9/workspaces/%d/projects"
	postWorkspaceProjects                        = "/api/v9/workspaces/%d/projects"
	patchWorkspaceProjects                       = "/api/v9/workspaces/%d/projects/%s"
	getWorkspaceProject                          = "/api/v9/workspaces/%d/projects/%d"
	putWorkspaceProject                          = "/api/v9/workspaces/%d/projects/%d"
	deleteWorkspaceProject                       = "/api/v9/workspaces/%d/projects/%d"
)

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
	ProjectId            int    `json:"project_id"`              //	Project ID
	Rate                 int    `json:"rate"`                    //	null
	RateLastUpdated      string `json:"rate_last_updated"`       //	Date for rate last updated
	UserId               int    `json:"user_id"`                 //	User ID
	WorkspaceId          int    `json:"workspace_id"`            //	Workspace ID
}

func (c Client) GetWorkspaceProjectsUsers(input GetWorkspaceProjectsUsersInput) ([]GetWorkspaceProjectsUsersOutput, error) {
	iq := input.Query
	q := url.Values{}
	q.Add("project_ids", iq.ProjectIds)
	q.Add("user_id", iq.UserId)
	q.Add("with_group_members", fmt.Sprintf("%v", iq.WithGroupMembers))

	toggl := c.Get(url.URL{RawQuery: q.Encode()})
	toggl.URL.Path = fmt.Sprintf(getWorkspaceProjectsUsers, input.WorkspaceID)

	resp, err := c.HttpClient.Do(&toggl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
	default:
		slog.Error(fmt.Sprintf("Error response status code: %v, message: %v", resp.Status, string(body)))
		return nil, ErrorStatusNotOK
	}
	out := make([]GetWorkspaceProjectsUsersOutput, 0)
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}

	return out, nil
}
