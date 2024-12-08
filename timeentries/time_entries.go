// Package timeentries provides functions to interact with time entries in Toggl.
package timeentries

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

const (
	timeEntriesPath             = "/api/v9/me/time_entries"
	currentTimeEntriesPath      = "/api/v9/me/time_entries/current"
	getATimeEntryByIdPath       = "/api/v9/me/time_entries/%d"
	postTimeEntries             = "/api/v9/workspaces/%d/time_entries"
	patchBulkEditingTimeEntries = "/api/v9/workspaces/%d/time_entries/%s"
	putTimeEntries              = "/api/v9/workspaces/%d/time_entries/%d"
	deleteTimeEntries           = "/api/v9/workspaces/%d/time_entries/%d"
	patchStopTimeEntry          = "/api/v9/workspaces/%d/time_entries/%d/stop"
)

// GetTimeEntriesQuery represents the query parameters for fetching time entries.
type GetTimeEntriesQuery struct {
	Meta           bool    // Should the response contain data for meta entities
	IncludeSharing bool    // Include sharing details in the response
	Since          *int64  // Get entries modified since this date using UNIX timestamp, including deleted ones.
	Before         *string // Get entries with start time, before given date (YYYY-MM-DD) or with time in RFC3339 format.
	StartDate      *string // Get entries with start time, from start_date YYYY-MM-DD or with time in RFC3339 format. To be used with end_date.
	EndDate        *string // Get entries with start time, until end_date YYYY-MM-DD or with time in RFC3339 format. To be used with start_date.
}

// GetTimeEntriesInput contains the input data for GetTimeEntries.
type GetTimeEntriesInput struct {
	Query GetTimeEntriesQuery
}

// SharedWith represents sharing details of a time entry.
type SharedWith struct {
	Accepted bool    `json:"accepted"`
	UserId   *int    `json:"user_id"`
	UserName *string `json:"user_name"`
}

// GetTimeEntriesOutput represents a time entry fetched from Toggl.
type GetTimeEntriesOutput struct {
	At              *string      `json:"at"`
	Billable        *bool        `json:"billable"`
	ClientName      *string      `json:"client_name"`
	Description     *string      `json:"description"`
	Duration        *int         `json:"duration"`
	Duronly         *bool        `json:"duronly"`
	Id              *int         `json:"id"`
	Permissions     []string     `json:"permissions"`
	Pid             *int         `json:"pid"`
	ProjectActive   *bool        `json:"project_active"`
	ProjectBillable *bool        `json:"project_billable"`
	ProjectColor    *string      `json:"project_color"`
	ProjectId       *int         `json:"project_id"`
	ProjectName     *string      `json:"project_name"`
	SharedWith      []SharedWith `json:"shared_with"`
	Start           *string      `json:"start"`
	Stop            *string      `json:"stop"`
	TagIds          []int        `json:"tag_ids"`
	Tags            []string     `json:"tags"`
	TaskId          *int         `json:"task_id"`
	TaskName        *string      `json:"task_name"`
	Tid             *int         `json:"tid"`
	Uid             *int         `json:"uid"`
	UserAvatarUrl   *string      `json:"user_avatar_url"`
	UserId          *int         `json:"user_id"`
	UserName        *string      `json:"user_name"`
	Wid             *int         `json:"wid"`
	WorkspaceId     *int         `json:"workspace_id"`
}

// GetTimeEntries retrieves time entries based on the provided input.
func (c Client) GetTimeEntries(tei GetTimeEntriesInput) ([]GetTimeEntriesOutput, error) {
	teq := tei.Query
	q := url.Values{}
	q.Add("meta", fmt.Sprintf("%v", teq.Meta))
	q.Add("include_sharing", fmt.Sprintf("%v", teq.IncludeSharing))
	if teq.Since != nil {
		q.Add("since", fmt.Sprintf("%v", *teq.Since))
	}
	if teq.Before != nil {
		q.Add("before", *teq.Before)
	}
	if teq.StartDate != nil {
		q.Add("start_date", *teq.StartDate)
	}
	if teq.EndDate != nil {
		q.Add("end_date", *teq.EndDate)
	}
	toggl := c.Get(url.URL{RawQuery: q.Encode()})
	toggl.URL.Path = timeEntriesPath

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
	case http.StatusNotFound:
		return []GetTimeEntriesOutput{}, nil
	default:
		slog.Error(fmt.Sprintf("Error response status code: %v, message: %v", resp.Status, string(body)))
		return nil, ErrorStatusNotOK
	}

	gteo := make([]GetTimeEntriesOutput, 0)
	if err := json.Unmarshal(body, &gteo); err != nil {
		return nil, err
	}

	return gteo, nil
}

// GetCurrentTimeEntry represents the current running time entry.
type GetCurrentTimeEntry = GetTimeEntriesOutput

// GetCurrentTimeEntry retrieves the current running time entry.
func (c Client) GetCurrentTimeEntry() (GetCurrentTimeEntry, error) {
	toggl := c.Get(url.URL{})
	toggl.URL.Path = currentTimeEntriesPath

	resp, err := c.HttpClient.Do(&toggl)
	if err != nil {
		return GetCurrentTimeEntry{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetCurrentTimeEntry{}, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return GetCurrentTimeEntry{}, nil
	default:
		slog.Error(fmt.Sprintf("Error response status code: %v, message: %v", resp.Status, string(body)))
		return GetCurrentTimeEntry{}, ErrorStatusNotOK
	}

	gcte := GetCurrentTimeEntry{}
	if err := json.Unmarshal(body, &gcte); err != nil {
		return GetCurrentTimeEntry{}, err
	}

	return gcte, nil
}

// GetATimeEntryByIdQuery represents the query parameters for fetching a specific time entry.
type GetATimeEntryByIdQuery struct {
	Meta           bool
	IncludeSharing bool
}

// GetATimeEntryByIdInput contains the input data for GetATimeEntryById.
type GetATimeEntryByIdInput struct {
	TimeEntryId *int // required
	Query       GetATimeEntryByIdQuery
}

// GetATimeEntryByIdOutput represents a single time entry retrieved by ID.
type GetATimeEntryByIdOutput = GetTimeEntriesOutput

// GetATimeEntryById retrieves a time entry by its ID.
func (c Client) GetATimeEntryById(input GetATimeEntryByIdInput) (GetATimeEntryByIdOutput, error) {
	if input.TimeEntryId == nil {
		slog.Error("TimeEntryId is required")
		return GetATimeEntryByIdOutput{}, ErrorRequiredParameter
	}
	q := url.Values{}
	q.Add("meta", fmt.Sprintf("%v", input.Query.Meta))
	q.Add("include_sharing", fmt.Sprintf("%v", input.Query.IncludeSharing))
	toggl := c.Get(url.URL{RawQuery: q.Encode()})
	toggl.URL.Path = fmt.Sprintf(getATimeEntryByIdPath, *input.TimeEntryId)

	resp, err := c.HttpClient.Do(&toggl)
	if err != nil {
		return GetATimeEntryByIdOutput{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetATimeEntryByIdOutput{}, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return GetATimeEntryByIdOutput{}, nil
	default:
		slog.Error(fmt.Sprintf("Error response status code: %v, message: %v", resp.Status, string(body)))
		return GetATimeEntryByIdOutput{}, ErrorStatusNotOK
	}

	gatebio := GetATimeEntryByIdOutput{}
	if err := json.Unmarshal(body, &gatebio); err != nil {
		return GetATimeEntryByIdOutput{}, err
	}

	return gatebio, nil
}

// PostTimeEntriesQuery represents the query parameters for creating a time entry.
type PostTimeEntriesQuery struct {
	Meta bool
}

// EventMetadata contains metadata for events related to time entries.
type EventMetadata struct {
	OriginFeature     string `json:"origin_feature,omitempty"`
	VisibleGoalsCount int    `json:"visible_goals_count,omitempty"`
}

// PostTimeEntriesBody represents the body of the request to create a time entry.
type PostTimeEntriesBody struct {
	Billable           bool          `json:"billable,omitempty"`             // Whether the time entry is marked as billable, optional, default false
	CreatedWith        string        `json:"created_with"`                   // Must be provided when creating a time entry and should identify the service/application used to create it
	Description        string        `json:"description,omitempty"`          // Time entry description, optional
	Duration           int           `json:"duration,omitempty"`             // Time entry duration. For running entries should be negative, preferable -1
	Duronly            bool          `json:"duronly,omitempty"`              // Deprecated: Used to create a time entry with a duration but without a stop time. This parameter can be ignored.
	EventMetadata      EventMetadata `json:"event_metadata,omitempty"`       // -
	Pid                int           `json:"pid,omitempty"`                  // Project ID, legacy field
	ProjectId          int           `json:"project_id,omitempty"`           // Project ID, optional
	SharedWith_userIds []int         `json:"shared_with_user_ids,omitempty"` // List of user IDs to share this time entry with
	Start              *string       `json:"start,omitempty"`                // Start time in UTC, required for creation. Format: 2006-01-02T15:04:05Z
	Start_date         *string       `json:"start_date,omitempty"`           // If provided during creation, the date part will take precedence over the date part of "start". Format: 2006-11-07
	Stop               string        `json:"stop,omitempty"`                 // Stop time in UTC, can be omitted if it's still running or created with "duration". If "stop" and "duration" are provided, values must be consistent (start + duration == stop)
	TagAction          string        `json:"tag_action,omitempty"`           // Can be "add" or "delete". Used when updating an existing time entry
	TagIds             []int         `json:"tag_ids,omitempty"`              // IDs of tags to add/remove
	Tags               []string      `json:"tags,omitempty"`                 // Names of tags to add/remove. If name does not exist as tag, one will be created automatically
	TaskId             int           `json:"task_id,omitempty"`              // Task ID, optional
	Tid                int           `json:"tid,omitempty"`                  // Task ID, legacy field
	Uid                int           `json:"uid,omitempty"`                  // Time Entry creator ID, legacy field
	UserId             int           `json:"user_id,omitempty"`              // Time Entry creator ID, if omitted will use the requester user ID
	Wid                int           `json:"wid,omitempty"`                  // Workspace ID, legacy field
	WorkspaceId        int           `json:"workspace_id"`                   // Workspace ID, required
}

// PostTimeEntriesInput contains the input data for PostTimeEntries.
type PostTimeEntriesInput struct {
	WorkspaceId *int // required
	Query       PostTimeEntriesQuery
	Body        PostTimeEntriesBody
}

// PostTimeEntriesOutput represents the response after creating a time entry.
type PostTimeEntriesOutput = GetTimeEntriesOutput

// PostTimeEntries creates a new time entry in Toggl.
func (c Client) PostTimeEntries(input PostTimeEntriesInput) (PostTimeEntriesOutput, error) {
	if input.WorkspaceId == nil {
		slog.Error("WorkspaceId is required")
		return PostTimeEntriesOutput{}, ErrorRequiredParameter
	}
	q := url.Values{}
	q.Add("meta", fmt.Sprintf("%v", input.Query.Meta))
	j, err := json.Marshal(input.Body)
	if err != nil {
		return PostTimeEntriesOutput{}, err
	}
	toggl := c.Post(url.URL{RawQuery: q.Encode()}, j)
	toggl.URL.Path = fmt.Sprintf(postTimeEntries, *input.WorkspaceId)

	resp, err := c.HttpClient.Do(&toggl)
	if err != nil {
		return PostTimeEntriesOutput{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PostTimeEntriesOutput{}, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return PostTimeEntriesOutput{}, nil
	default:
		slog.Error(fmt.Sprintf("Error response status code: %v, message: %v", resp.Status, string(body)))
		return PostTimeEntriesOutput{}, ErrorStatusNotOK
	}

	pteo := PostTimeEntriesOutput{}
	if err := json.Unmarshal(body, &pteo); err != nil {
		return PostTimeEntriesOutput{}, err
	}

	return pteo, nil
}

// PatchBulkEditingTimeEntriesQuery represents the query parameters for bulk editing time entries.
type PatchBulkEditingTimeEntriesQuery struct {
	Meta bool
}

// PatchBulkEditingTimeEntriesInput contains the input data for bulk editing time entries.
type PatchBulkEditingTimeEntriesInput struct {
	/*
		Numeric ID of the workspace
	*/
	WorkspaceId *int // required
	/*
		Numeric IDs of time_entries, separated by comma.
		E.g.: 204301830,202700150,202687559. The limit is 100 IDs per request.
	*/
	TimeEntryIds *string // required
	Query        PatchBulkEditingTimeEntriesQuery
	/*
		Body
		* items	Array of object	-

		items
		* op	string	Operation (add/remove/replace)
		* path	string	The path to the entity to patch (e.g. /description)
		* value	object	The new value for the entity in path.
	*/
	Body []byte
}

// Failure represents a failure in bulk editing time entries.
type Failure struct {
	Id      int    `json:"id"`      // The ID for which the patch operation failed.
	Message string `json:"message"` // The operation failure reason
}

// PatchBulkEditingTimeEntriesOutput represents the response from bulk editing time entries.
type PatchBulkEditingTimeEntriesOutput struct {
	Failure []Failure `json:"failure"`
	Success []int     `json:"success"` // The IDs for which the patch was succesful.
}

// PatchBulkEditingTimeEntries performs bulk edit operations on time entries.
func (c Client) PatchBulkEditingTimeEntries(input PatchBulkEditingTimeEntriesInput) (PatchBulkEditingTimeEntriesOutput, error) {
	if input.WorkspaceId == nil {
		slog.Error("WorkspaceId is required")
		return PatchBulkEditingTimeEntriesOutput{}, ErrorRequiredParameter
	}
	if input.TimeEntryIds == nil {
		slog.Error("TimeEntryIds is required")
		return PatchBulkEditingTimeEntriesOutput{}, ErrorRequiredParameter
	}
	q := url.Values{}
	q.Add("meta", fmt.Sprintf("%v", input.Query.Meta))
	toggl := c.Patch(url.URL{RawQuery: q.Encode()}, input.Body)
	toggl.URL.Path = fmt.Sprintf(patchBulkEditingTimeEntries, *input.WorkspaceId, *input.TimeEntryIds)

	resp, err := c.HttpClient.Do(&toggl)
	if err != nil {
		return PatchBulkEditingTimeEntriesOutput{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PatchBulkEditingTimeEntriesOutput{}, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return PatchBulkEditingTimeEntriesOutput{}, nil
	default:
		slog.Error(fmt.Sprintf("Error response status code: %v, message: %v", resp.Status, string(body)))
		return PatchBulkEditingTimeEntriesOutput{}, ErrorStatusNotOK
	}

	pbeto := PatchBulkEditingTimeEntriesOutput{}
	if err := json.Unmarshal(body, &pbeto); err != nil {
		return PatchBulkEditingTimeEntriesOutput{}, err
	}

	return pbeto, nil
}

// PutTimeEntriesQuery represents the query parameters for updating a time entry.
type PutTimeEntriesQuery struct {
	Meta           bool // Should the response contain data for meta entities
	IncludeSharing bool // Should the response contain time entry sharing details
}

// PutTimeEntriesBody represents the body of the request to update a time entry.
type PutTimeEntriesBody = PostTimeEntriesBody

// PutTimeEntriesInput contains the input data for updating a time entry.
type PutTimeEntriesInput struct {
	WorkspaceId *int // required
	TimeEntryId *int // required
	Query       PutTimeEntriesQuery
	Body        PutTimeEntriesBody
}

// PutTimeEntriesOutput represents the response after updating a time entry.
type PutTimeEntriesOutput = GetTimeEntriesOutput

// PutTimeEntries updates an existing time entry.
func (c Client) PutTimeEntries(input PutTimeEntriesInput) (PutTimeEntriesOutput, error) {
	if input.WorkspaceId == nil {
		slog.Error("WorkspaceId is required")
		return PutTimeEntriesOutput{}, ErrorRequiredParameter
	}
	if input.TimeEntryId == nil {
		slog.Error("TimeEntryId is required")
		return PutTimeEntriesOutput{}, ErrorRequiredParameter
	}
	q := url.Values{}
	q.Add("meta", fmt.Sprintf("%v", input.Query.Meta))
	j, err := json.Marshal(input.Body)
	if err != nil {
		return PutTimeEntriesOutput{}, err
	}
	toggl := c.Put(url.URL{RawQuery: q.Encode()}, j)
	toggl.URL.Path = fmt.Sprintf(putTimeEntries, *input.WorkspaceId, *input.TimeEntryId)

	resp, err := c.HttpClient.Do(&toggl)
	if err != nil {
		return PutTimeEntriesOutput{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PutTimeEntriesOutput{}, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return PutTimeEntriesOutput{}, nil
	default:
		slog.Error(fmt.Sprintf("Error response status code: %v, message: %v", resp.Status, string(body)))
		return PutTimeEntriesOutput{}, ErrorStatusNotOK
	}

	pbeto := PutTimeEntriesOutput{}
	if err := json.Unmarshal(body, &pbeto); err != nil {
		return PutTimeEntriesOutput{}, err
	}

	return pbeto, nil
}

// DeleteTimeEntriesInput contains the input data for deleting a time entry.
type DeleteTimeEntriesInput struct {
	WorkspaceId *int // required
	TimeEntryId *int // required
}

// DeleteTimeEntries deletes a time entry from Toggl.
func (c Client) DeleteTimeEntries(input DeleteTimeEntriesInput) error {
	if input.WorkspaceId == nil {
		slog.Error("WorkspaceId is required")
		return ErrorRequiredParameter
	}
	if input.TimeEntryId == nil {
		slog.Error("TimeEntryId is required")
		return ErrorRequiredParameter
	}
	toggl := c.Delete(url.URL{})
	toggl.URL.Path = fmt.Sprintf(deleteTimeEntries, *input.WorkspaceId, *input.TimeEntryId)
	resp, err := c.HttpClient.Do(&toggl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil
	default:
		slog.Error(fmt.Sprintf("Error response status code: %v, message: %v", resp.Status, string(body)))
		return ErrorStatusNotOK
	}

	return nil
}

// PatchStopTimeEntryInput contains the input data for stopping a running time entry.
type PatchStopTimeEntryInput struct {
	WorkspaceId *int // required
	TimeEntryId *int // required
}

// PatchStopTimeEntryOutput represents the response after stopping a time entry.
type PatchStopTimeEntryOutput = GetTimeEntriesOutput

// PatchStopTimeEntry stops a running time entry.
func (c Client) PatchStopTimeEntry(input PatchStopTimeEntryInput) (PatchStopTimeEntryOutput, error) {
	if input.WorkspaceId == nil {
		slog.Error("WorkspaceId is required")
		return PatchStopTimeEntryOutput{}, ErrorRequiredParameter
	}
	if input.TimeEntryId == nil {
		slog.Error("TimeEntryId is required")
		return PatchStopTimeEntryOutput{}, ErrorRequiredParameter
	}
	toggl := c.Patch(url.URL{}, nil)
	toggl.URL.Path = fmt.Sprintf(patchStopTimeEntry, *input.WorkspaceId, *input.TimeEntryId)
	resp, err := c.HttpClient.Do(&toggl)
	if err != nil {
		return PatchStopTimeEntryOutput{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PatchStopTimeEntryOutput{}, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
	default:
		slog.Error(fmt.Sprintf("Error response status code: %v, message: %v", resp.Status, string(body)))
		return PatchStopTimeEntryOutput{}, ErrorStatusNotOK
	}
	psteo := PatchStopTimeEntryOutput{}
	if err := json.Unmarshal(body, &psteo); err != nil {
		return PutTimeEntriesOutput{}, err
	}

	return psteo, nil
}
