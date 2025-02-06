package projects_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/dev-shimada/toggl-go/projects"
	"github.com/google/go-cmp/cmp"
)

type MockHttpClient struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (m MockHttpClient) Do(r *http.Request) (*http.Response, error) {
	return m.DoFunc(r)
}

func fakeClient(res *http.Response) projects.Client {
	return projects.Client{
		HttpClient: MockHttpClient{
			DoFunc: func(r *http.Request) (*http.Response, error) {
				return res, nil
			},
		},
	}
}

func TestGetWorkspaceProjectsUsers(t *testing.T) {
	testFile, err := os.ReadFile("testdata/projects/project_users.json")
	if err != nil {
		t.Fatal(err)
	}
	testFile = bytes.TrimSpace(testFile)
	okBody := io.NopCloser(bytes.NewBuffer(testFile))
	successClient := fakeClient(&http.Response{StatusCode: http.StatusOK, Body: okBody})
	errorClient := fakeClient(&http.Response{StatusCode: http.StatusBadRequest, Body: http.NoBody})
	errorWant, err := json.MarshalIndent(nil, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	test := []struct {
		name     string
		client   projects.Client
		arg      projects.GetWorkspaceProjectsUsersInput
		wantJson []byte
		wantErr  error
	}{
		{
			name:   "success",
			client: successClient,
			arg: projects.GetWorkspaceProjectsUsersInput{
				WorkspaceId: 1,
				Query: projects.GetWorkspaceProjectsUsersQuery{
					ProjectIds: "1", UserId: "1", WithGroupMembers: true},
			},
			wantJson: testFile,
			wantErr:  nil,
		},
		{
			name:   "http error",
			client: errorClient,
			arg: projects.GetWorkspaceProjectsUsersInput{
				WorkspaceId: 1,
				Query: projects.GetWorkspaceProjectsUsersQuery{
					ProjectIds: "1", UserId: "1", WithGroupMembers: true},
			},
			wantJson: errorWant,
			wantErr:  projects.ErrorStatusNotOK,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.GetWorkspaceProjectsUsers(tt.arg)
			if err != tt.wantErr {
				t.Errorf("want %v, got %v", tt.wantErr, err)
			}
			gotJson, err := json.MarshalIndent(got, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(gotJson, tt.wantJson) {
				t.Errorf("diff: %v", cmp.Diff(gotJson, tt.wantJson))
			}
		})
	}
}

func TestPostAddAnUserIntoWorkspaceProjectsUsers(t *testing.T) {
	testFile, err := os.ReadFile("testdata/projects/project_users.json")
	if err != nil {
		t.Fatal(err)
	}
	testFile = bytes.TrimSpace(testFile)
	okBody := io.NopCloser(bytes.NewBuffer(testFile))
	successClient := fakeClient(&http.Response{StatusCode: http.StatusOK, Body: okBody})
	errorClient := fakeClient(&http.Response{StatusCode: http.StatusBadRequest, Body: http.NoBody})
	errorWant, err := json.MarshalIndent(nil, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	test := []struct {
		name     string
		client   projects.Client
		arg      projects.PostAddAnUserIntoWorkspaceProjectsUsersInput
		wantJson []byte
		wantErr  error
	}{
		{
			name:   "success",
			client: successClient,
			arg: projects.PostAddAnUserIntoWorkspaceProjectsUsersInput{
				WorkspaceId: 1,
				Body:        projects.PostAddAnUserIntoWorkspaceProjectsUsersBody{},
			},
			wantJson: testFile,
			wantErr:  nil,
		},
		{
			name:   "http error",
			client: errorClient,
			arg: projects.PostAddAnUserIntoWorkspaceProjectsUsersInput{
				WorkspaceId: 1,
				Body:        projects.PostAddAnUserIntoWorkspaceProjectsUsersBody{},
			},
			wantJson: errorWant,
			wantErr:  projects.ErrorStatusNotOK,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.PostAddAnUserIntoWorkspaceProjectsUsers(tt.arg)
			if err != tt.wantErr {
				t.Errorf("want %v, got %v", tt.wantErr, err)
			}
			gotJson, err := json.MarshalIndent(got, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(gotJson, tt.wantJson) {
				t.Errorf("diff: %v", cmp.Diff(gotJson, tt.wantJson))
			}
		})
	}
}
