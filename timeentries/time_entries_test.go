package timeentries_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/dev-shimada/toggl-go/timeentries"
	"github.com/google/go-cmp/cmp"
)

type MockHttpClient struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (m MockHttpClient) Do(r *http.Request) (*http.Response, error) {
	return m.DoFunc(r)
}

func fakeClient(res *http.Response) timeentries.Client {
	return timeentries.Client{
		HttpClient: MockHttpClient{
			DoFunc: func(r *http.Request) (*http.Response, error) {
				return res, nil
			},
		},
	}
}

func TestGetTimeEntries(t *testing.T) {
	testFile, err := os.ReadFile("testdata/time_entries/time_entries.json")
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
		arg      timeentries.GetTimeEntriesInput
		client   timeentries.Client
		wantJson []byte
		wantErr  error
	}{
		{"success", timeentries.GetTimeEntriesInput{}, successClient, testFile, nil},
		{"error", timeentries.GetTimeEntriesInput{}, errorClient, errorWant, timeentries.ErrorStatusNotOK},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.GetTimeEntries(timeentries.GetTimeEntriesInput{})
			if tt.wantErr != nil && err == nil {
				t.Errorf("Expected error, got nil")
			} else if tt.wantErr != nil && err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.wantErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			jgot, err := json.MarshalIndent(got, "", "  ")
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(tt.wantJson, jgot) {
				t.Errorf("diff: %v", cmp.Diff(tt.wantJson, jgot))
			}
		})
	}
}

func TestGetCurrentTimeEntry(t *testing.T) {
	testFile, err := os.ReadFile("testdata/time_entries/time_entry.json")
	if err != nil {
		t.Fatal(err)
	}
	testFile = bytes.TrimSpace(testFile)
	okBody := io.NopCloser(bytes.NewBuffer(testFile))

	successClient := fakeClient(&http.Response{StatusCode: http.StatusOK, Body: okBody})
	errorClient := fakeClient(&http.Response{StatusCode: http.StatusBadRequest, Body: http.NoBody})
	errorWant, err := json.MarshalIndent(timeentries.GetCurrentTimeEntry{}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	test := []struct {
		name     string
		client   timeentries.Client
		wantJson []byte
		wantErr  error
	}{
		{"success", successClient, testFile, nil},
		{"error", errorClient, errorWant, timeentries.ErrorStatusNotOK},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.GetCurrentTimeEntry()
			if tt.wantErr != nil && err == nil {
				t.Errorf("Expected error, got nil")
			} else if tt.wantErr != nil && err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.wantErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			jgot, err := json.MarshalIndent(got, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			jgot = bytes.TrimSpace(jgot)
			if !cmp.Equal(tt.wantJson, jgot) {
				t.Errorf("diff: %v", cmp.Diff(tt.wantJson, jgot))
			}
		})
	}
}

func TestGetATimeEntryById(t *testing.T) {
	testFile, err := os.ReadFile("testdata/time_entries/time_entry.json")
	if err != nil {
		t.Fatal(err)
	}
	testFile = bytes.TrimSpace(testFile)
	okBody := io.NopCloser(bytes.NewBuffer(testFile))

	successClient := fakeClient(&http.Response{StatusCode: http.StatusOK, Body: okBody})
	errorClient := fakeClient(&http.Response{StatusCode: http.StatusBadRequest, Body: http.NoBody})
	errorWant, err := json.MarshalIndent(timeentries.GetATimeEntryByIdOutput{}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	timeEntryId := 123456789
	test := []struct {
		name     string
		client   timeentries.Client
		arg      timeentries.GetATimeEntryByIdInput
		wantJson []byte
		wantErr  error
	}{
		{
			name:     "success",
			client:   successClient,
			arg:      timeentries.GetATimeEntryByIdInput{TimeEntryId: timeEntryId},
			wantJson: testFile,
			wantErr:  nil,
		},
		{
			name:     "parameter error",
			client:   successClient,
			arg:      timeentries.GetATimeEntryByIdInput{},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorRequiredParameter,
		},
		{
			name:     "http error",
			client:   errorClient,
			arg:      timeentries.GetATimeEntryByIdInput{TimeEntryId: timeEntryId},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorStatusNotOK,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.GetATimeEntryById(tt.arg)
			if tt.wantErr != nil && err == nil {
				t.Errorf("Expected error, got nil")
			} else if tt.wantErr != nil && err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.wantErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			jgot, err := json.MarshalIndent(got, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			jgot = bytes.TrimSpace(jgot)
			if !cmp.Equal(tt.wantJson, jgot) {
				t.Errorf("diff: %v", cmp.Diff(tt.wantJson, jgot))
			}
		})
	}
}

func TestPostTimeEntries(t *testing.T) {
	testFile, err := os.ReadFile("testdata/time_entries/time_entry.json")
	if err != nil {
		t.Fatal(err)
	}
	testFile = bytes.TrimSpace(testFile)
	okBody := io.NopCloser(bytes.NewBuffer(testFile))

	successClient := fakeClient(&http.Response{StatusCode: http.StatusOK, Body: okBody})
	errorClient := fakeClient(&http.Response{StatusCode: http.StatusBadRequest, Body: http.NoBody})
	errorWant, err := json.MarshalIndent(timeentries.PostTimeEntriesOutput{}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	workspaceId := 123456789
	test := []struct {
		name     string
		client   timeentries.Client
		arg      timeentries.PostTimeEntriesInput
		wantJson []byte
		wantErr  error
	}{
		{
			name:     "success",
			client:   successClient,
			arg:      timeentries.PostTimeEntriesInput{WorkspaceId: workspaceId, Body: timeentries.PostTimeEntriesBody{}},
			wantJson: testFile,
			wantErr:  nil,
		},
		{
			name:     "parameter error",
			client:   successClient,
			arg:      timeentries.PostTimeEntriesInput{},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorRequiredParameter,
		},
		{
			name:     "http error",
			client:   errorClient,
			arg:      timeentries.PostTimeEntriesInput{WorkspaceId: workspaceId},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorStatusNotOK,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.PostTimeEntries(tt.arg)
			if tt.wantErr != nil && err == nil {
				t.Errorf("Expected error, got nil")
			} else if tt.wantErr != nil && err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.wantErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			jgot, err := json.MarshalIndent(got, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			jgot = bytes.TrimSpace(jgot)
			if !cmp.Equal(tt.wantJson, jgot) {
				t.Errorf("diff: %v", cmp.Diff(tt.wantJson, jgot))
			}
		})
	}
}

func TestPatchBulkEditingTimeEntries(t *testing.T) {
	testFile, err := os.ReadFile("testdata/time_entries/patch_bulk_editing_time_entries.json")
	if err != nil {
		t.Fatal(err)
	}
	testFile = bytes.TrimSpace(testFile)
	okBody := io.NopCloser(bytes.NewBuffer(testFile))

	successClient := fakeClient(&http.Response{StatusCode: http.StatusOK, Body: okBody})
	errorClient := fakeClient(&http.Response{StatusCode: http.StatusBadRequest, Body: http.NoBody})
	errorWant, err := json.MarshalIndent(timeentries.PatchBulkEditingTimeEntriesOutput{}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	workspaceId := 123456789
	timeEntryIds := "1234567890,1234567891"
	test := []struct {
		name     string
		client   timeentries.Client
		arg      timeentries.PatchBulkEditingTimeEntriesInput
		wantJson []byte
		wantErr  error
	}{
		{
			name:     "success",
			client:   successClient,
			arg:      timeentries.PatchBulkEditingTimeEntriesInput{WorkspaceId: workspaceId, TimeEntryIds: timeEntryIds},
			wantJson: testFile,
			wantErr:  nil,
		},
		{
			name:     "WorkspaceId parameter error",
			client:   successClient,
			arg:      timeentries.PatchBulkEditingTimeEntriesInput{TimeEntryIds: timeEntryIds},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorRequiredParameter,
		},
		{
			name:     "TimeEntryIds parameter error",
			client:   successClient,
			arg:      timeentries.PatchBulkEditingTimeEntriesInput{WorkspaceId: workspaceId},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorRequiredParameter,
		},
		{
			name:     "http error",
			client:   errorClient,
			arg:      timeentries.PatchBulkEditingTimeEntriesInput{WorkspaceId: workspaceId, TimeEntryIds: timeEntryIds},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorStatusNotOK,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.PatchBulkEditingTimeEntries(tt.arg)
			if tt.wantErr != nil && err == nil {
				t.Errorf("Expected error, got nil")
			} else if tt.wantErr != nil && err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.wantErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			jgot, err := json.MarshalIndent(got, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			jgot = bytes.TrimSpace(jgot)
			if !cmp.Equal(tt.wantJson, jgot) {
				t.Errorf("diff: %v", cmp.Diff(tt.wantJson, jgot))
			}
		})
	}
}

func TestPutTimeEntries(t *testing.T) {
	testFile, err := os.ReadFile("testdata/time_entries/time_entry.json")
	if err != nil {
		t.Fatal(err)
	}
	testFile = bytes.TrimSpace(testFile)
	okBody := io.NopCloser(bytes.NewBuffer(testFile))

	successClient := fakeClient(&http.Response{StatusCode: http.StatusOK, Body: okBody})
	errorClient := fakeClient(&http.Response{StatusCode: http.StatusBadRequest, Body: http.NoBody})
	errorWant, err := json.MarshalIndent(timeentries.PutTimeEntriesOutput{}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	workspaceId := 123456789
	timeEntryId := 1234567890
	test := []struct {
		name     string
		client   timeentries.Client
		arg      timeentries.PutTimeEntriesInput
		wantJson []byte
		wantErr  error
	}{
		{
			name:     "success",
			client:   successClient,
			arg:      timeentries.PutTimeEntriesInput{WorkspaceId: workspaceId, TimeEntryId: timeEntryId},
			wantJson: testFile,
			wantErr:  nil,
		},
		{
			name:     "WorkspaceId parameter error",
			client:   successClient,
			arg:      timeentries.PutTimeEntriesInput{TimeEntryId: timeEntryId},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorRequiredParameter,
		},
		{
			name:     "TimeEntryIds parameter error",
			client:   successClient,
			arg:      timeentries.PutTimeEntriesInput{WorkspaceId: workspaceId},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorRequiredParameter,
		},
		{
			name:     "http error",
			client:   errorClient,
			arg:      timeentries.PutTimeEntriesInput{WorkspaceId: workspaceId, TimeEntryId: timeEntryId},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorStatusNotOK,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.PutTimeEntries(tt.arg)
			if tt.wantErr != nil && err == nil {
				t.Errorf("Expected error, got nil")
			} else if tt.wantErr != nil && err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.wantErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			jgot, err := json.MarshalIndent(got, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			jgot = bytes.TrimSpace(jgot)
			if !cmp.Equal(tt.wantJson, jgot) {
				t.Errorf("diff: %v", cmp.Diff(tt.wantJson, jgot))
			}
		})
	}
}

func TestDeleteTimeEntries(t *testing.T) {
	successClient := fakeClient(&http.Response{StatusCode: http.StatusOK, Body: http.NoBody})
	errorClient := fakeClient(&http.Response{StatusCode: http.StatusBadRequest, Body: http.NoBody})

	workspaceId := 123456789
	timeEntryId := 1234567890
	test := []struct {
		name    string
		client  timeentries.Client
		arg     timeentries.DeleteTimeEntriesInput
		wantErr error
	}{
		{
			name:    "success",
			client:  successClient,
			arg:     timeentries.DeleteTimeEntriesInput{WorkspaceId: workspaceId, TimeEntryId: timeEntryId},
			wantErr: nil,
		},
		{
			name:    "WorkspaceId parameter error",
			client:  successClient,
			arg:     timeentries.DeleteTimeEntriesInput{TimeEntryId: timeEntryId},
			wantErr: timeentries.ErrorRequiredParameter,
		},
		{
			name:    "TimeEntryIds parameter error",
			client:  successClient,
			arg:     timeentries.DeleteTimeEntriesInput{WorkspaceId: workspaceId},
			wantErr: timeentries.ErrorRequiredParameter,
		},
		{
			name:    "http error",
			client:  errorClient,
			arg:     timeentries.DeleteTimeEntriesInput{WorkspaceId: workspaceId, TimeEntryId: timeEntryId},
			wantErr: timeentries.ErrorStatusNotOK,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.DeleteTimeEntries(tt.arg)
			if tt.wantErr != nil && err == nil {
				t.Errorf("Expected error, got nil")
			} else if tt.wantErr != nil && err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.wantErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestPatchStopTimeEntry(t *testing.T) {
	testFile, err := os.ReadFile("testdata/time_entries/time_entry.json")
	if err != nil {
		t.Fatal(err)
	}
	testFile = bytes.TrimSpace(testFile)
	okBody := io.NopCloser(bytes.NewBuffer(testFile))

	successClient := fakeClient(&http.Response{StatusCode: http.StatusOK, Body: okBody})
	errorClient := fakeClient(&http.Response{StatusCode: http.StatusBadRequest, Body: http.NoBody})
	errorWant, err := json.MarshalIndent(timeentries.PatchStopTimeEntryOutput{}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	workspaceId := 123456789
	timeEntryId := 1234567890
	test := []struct {
		name     string
		client   timeentries.Client
		arg      timeentries.PatchStopTimeEntryInput
		wantJson []byte
		wantErr  error
	}{
		{
			name:     "success",
			client:   successClient,
			arg:      timeentries.PatchStopTimeEntryInput{WorkspaceId: workspaceId, TimeEntryId: timeEntryId},
			wantJson: testFile,
			wantErr:  nil,
		},
		{
			name:     "WorkspaceId parameter error",
			client:   successClient,
			arg:      timeentries.PatchStopTimeEntryInput{TimeEntryId: timeEntryId},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorRequiredParameter,
		},
		{
			name:     "TimeEntryIds parameter error",
			client:   successClient,
			arg:      timeentries.PatchStopTimeEntryInput{WorkspaceId: workspaceId},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorRequiredParameter,
		},
		{
			name:     "http error",
			client:   errorClient,
			arg:      timeentries.PatchStopTimeEntryInput{WorkspaceId: workspaceId, TimeEntryId: timeEntryId},
			wantJson: errorWant,
			wantErr:  timeentries.ErrorStatusNotOK,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.PatchStopTimeEntry(tt.arg)
			if tt.wantErr != nil && err == nil {
				t.Errorf("Expected error, got nil")
			} else if tt.wantErr != nil && err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.wantErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			jgot, err := json.MarshalIndent(got, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			jgot = bytes.TrimSpace(jgot)
			if !cmp.Equal(tt.wantJson, jgot) {
				t.Errorf("diff: %v", cmp.Diff(tt.wantJson, jgot))
			}
		})
	}
}

func TestGetTimeEntriesRemoteAccess(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in local environment")
	} else if os.Getenv("WORKSPACE_ID") == "" {
		t.Skip("Skipping test because WORKSPACE_ID is not set")
	} else if os.Getenv("TOKEN") == "" {
		t.Skip("Skipping test because TOKEN is not set")
	}

	now := time.Now()
	nowTime := now.Format(time.RFC3339)
	base := now.Add(-time.Hour * 24)
	baseTime := base.Format(time.RFC3339)
	workspace, err := strconv.Atoi(os.Getenv("WORKSPACE_ID"))
	if err != nil {
		t.Fatal(err)
	}
	client := timeentries.NewClient(os.Getenv("TOKEN"))

	// cleanup
	cleanup := func() {
		ct, err := client.GetCurrentTimeEntry()
		if err != nil {
			t.Log(err)
		} else {
			if _, err := client.PatchStopTimeEntry(timeentries.PatchStopTimeEntryInput{WorkspaceId: workspace, TimeEntryId: ct.Id}); err != nil {
				t.Log(err)
			}
		}

		timeEntry, err := client.GetTimeEntries(timeentries.GetTimeEntriesInput{Query: timeentries.GetTimeEntriesQuery{StartDate: &baseTime, EndDate: &nowTime}})
		if err != nil {
			t.Fatal(err)
		}

		wg := sync.WaitGroup{}
		for _, te := range timeEntry {
			if te.WorkspaceId == workspace {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					if err := client.DeleteTimeEntries(timeentries.DeleteTimeEntriesInput{WorkspaceId: workspace, TimeEntryId: id}); err != nil {
						t.Log(err)
					}
				}(te.Id)
			}
		}
		wg.Wait()

		afTimeEntry, err := client.GetTimeEntries(timeentries.GetTimeEntriesInput{Query: timeentries.GetTimeEntriesQuery{StartDate: &baseTime, EndDate: &nowTime}})
		if err != nil {
			t.Fatal(err)
		}
		if len(afTimeEntry) != 0 {
			t.Fatal("Failed to delete all time entries")
		}
	}
	defer cleanup()

	// PostTimeEntries
	t.Run("PostTimeEntries", func(t *testing.T) {
		cleanup()
		time.Sleep(1 * time.Second)
		description := "test PostTimeEntries"
		start := now.Add(-time.Hour).Format("2006-01-02T15:04:05Z")
		stop := now.Format("2006-01-02T15:04:05Z")
		_, err := client.PostTimeEntries(timeentries.PostTimeEntriesInput{
			WorkspaceId: workspace,
			Body: timeentries.PostTimeEntriesBody{
				Description: description,
				Start:       &start,
				Stop:        stop,
				WorkspaceId: workspace,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		got, err := client.GetTimeEntries(timeentries.GetTimeEntriesInput{Query: timeentries.GetTimeEntriesQuery{StartDate: &start, EndDate: &stop}})
		if err != nil {
			t.Fatal(err)
		}
		if description != got[0].Description {
			t.Errorf("Expected %s, got %s", description, got[0].Description)
		}
	})

	// GetCurrentTimeEntry
	t.Run("GetCurrentTimeEntry", func(t *testing.T) {
		cleanup()
		time.Sleep(1 * time.Second)
		description := "test GetCurrentTimeEntry"
		start := now.Add(-time.Hour).Format("2006-01-02T15:04:05Z")
		duration := -1
		_, err := client.PostTimeEntries(timeentries.PostTimeEntriesInput{
			WorkspaceId: workspace,
			Body: timeentries.PostTimeEntriesBody{
				Description: description,
				Start:       &start,
				Duration:    duration,
				WorkspaceId: workspace,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		got, err := client.GetCurrentTimeEntry()
		if err != nil {
			t.Fatal(err)
		}
		if description != got.Description {
			t.Errorf("Expected %s, got %s", description, got.Description)
		}
	})

	// GetATimeEntryById
	t.Run("GetATimeEntryById", func(t *testing.T) {
		cleanup()
		time.Sleep(1 * time.Second)
		description := "test GetATimeEntryById"
		start := now.Add(-time.Hour).Format("2006-01-02T15:04:05Z")
		stop := now.Format("2006-01-02T15:04:05Z")
		postTimeEntries, err := client.PostTimeEntries(timeentries.PostTimeEntriesInput{
			WorkspaceId: workspace,
			Body: timeentries.PostTimeEntriesBody{
				Description: description,
				Start:       &start,
				Stop:        stop,
				WorkspaceId: workspace,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		got, err := client.GetATimeEntryById(timeentries.GetATimeEntryByIdInput{TimeEntryId: postTimeEntries.Id})
		if err != nil {
			t.Fatal(err)
		}
		if description != got.Description {
			t.Errorf("Expected %s, got %s", description, got.Description)
		}
	})

	// PutTimeEntriesBody
	t.Run("PutTimeEntriesBody", func(t *testing.T) {
		cleanup()
		time.Sleep(1 * time.Second)
		description := "test PutTimeEntriesBody"
		start := now.Add(-time.Hour).Format("2006-01-02T15:04:05Z")
		stop := now.Format("2006-01-02T15:04:05Z")
		postTimeEntries, err := client.PostTimeEntries(timeentries.PostTimeEntriesInput{
			WorkspaceId: workspace,
			Body: timeentries.PostTimeEntriesBody{
				Description: description,
				Start:       &start,
				Stop:        stop,
				WorkspaceId: workspace,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		want := "test updated PutTimeEntriesBody"
		got, err := client.PutTimeEntries(timeentries.PutTimeEntriesInput{
			WorkspaceId: workspace,
			// TimeEntryId: GetTimeEntries[0].Id,
			TimeEntryId: postTimeEntries.Id,
			Body:        timeentries.PutTimeEntriesBody{Description: want, WorkspaceId: workspace},
		})
		if err != nil {
			t.Fatal(err)
		}
		if want != got.Description {
			t.Errorf("Expected %s, got %s", want, got.Description)
		}
	})

	// PatchBulkEditingTimeEntries
	t.Run("PatchBulkEditingTimeEntries", func(t *testing.T) {
		cleanup()
		time.Sleep(1 * time.Second)
		description := "test PatchBulkEditingTimeEntries"
		start1 := now.Add(-time.Hour).Format("2006-01-02T15:04:05Z")
		stop1 := now.Format("2006-01-02T15:04:05Z")
		postTimeEntries1, err := client.PostTimeEntries(timeentries.PostTimeEntriesInput{
			WorkspaceId: workspace,
			Body: timeentries.PostTimeEntriesBody{
				Description: description,
				Start:       &start1,
				Stop:        stop1,
				WorkspaceId: workspace,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		start2 := now.Add(-time.Hour * 2).Format("2006-01-02T15:04:05Z")
		stop2 := now.Add(-time.Hour * 1).Format("2006-01-02T15:04:05Z")
		postTimeEntries2, err := client.PostTimeEntries(timeentries.PostTimeEntriesInput{
			WorkspaceId: workspace,
			Body: timeentries.PostTimeEntriesBody{
				Description: description,
				Start:       &start2,
				Stop:        stop2,
				WorkspaceId: workspace,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		timeEntryIds := fmt.Sprintf("%d,%d", postTimeEntries1.Id, postTimeEntries2.Id)
		body := `[{"op": "replace", "path": "/description", "value":"test updated PatchBulkEditingTimeEntries"}]`
		patchBulkEditingTimeEntries, err := client.PatchBulkEditingTimeEntries(
			timeentries.PatchBulkEditingTimeEntriesInput{
				WorkspaceId:  workspace,
				TimeEntryIds: timeEntryIds,
				Body:         []byte(body),
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		got, err := client.GetATimeEntryById(timeentries.GetATimeEntryByIdInput{TimeEntryId: patchBulkEditingTimeEntries.Success[0]})
		if err != nil {
			t.Fatal(err)
		}
		if got.Description != "test updated PatchBulkEditingTimeEntries" {
			t.Errorf("Expected %s, got %s", "test updated PatchBulkEditingTimeEntries", got.Description)
		}
	})
}

// remote
/*
	- GET TimeEntries
	- POST TimeEntries
	- GET Get current time entry
	- PATCH Stop TimeEntry
	PATCH Bulk editing time entries
	- GET Get a time entry by ID.
	- PUT TimeEntries
	- DELETE TimeEntries
*/
