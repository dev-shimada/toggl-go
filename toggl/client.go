// Package toggl provides a client for interacting with the Toggl API,
// allowing access to time entries and other Toggl resources.
package toggl

import (
	"github.com/dev-shimada/toggl-go/timeentries"
)

// Client represents a Toggl client with access to time entries.
type Client struct {
	TimeEntriesClient timeentries.Client
}

// NewClient creates a new Toggl client with the provided API token.
// It returns a Client struct with a TimeEntriesClient initialized.
func NewClient(token string) Client {
	return Client{
		TimeEntriesClient: timeentries.NewClient(token),
	}
}
