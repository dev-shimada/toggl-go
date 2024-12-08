package toggl_test

import (
	"os"
	"testing"
	"time"

	"github.com/dev-shimada/toggl-go/timeentries"
	"github.com/dev-shimada/toggl-go/toggl"
	"github.com/google/go-cmp/cmp"
)

func TestNewClient(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in local environment")
	} else if os.Getenv("TOKEN") == "" {
		t.Skip("Skipping test because TOGGL_API_TOKEN is not set")
	}

	token := os.Getenv("TOKEN")
	lastWeek := time.Now().Add(time.Hour).Unix()
	client := toggl.NewClient(token)

	want := []timeentries.GetTimeEntriesOutput{}

	got, err := client.TimeEntriesClient.GetTimeEntries(
		timeentries.GetTimeEntriesInput{
			Query: timeentries.GetTimeEntriesQuery{
				Since: &lastWeek,
			},
		},
	)
	if err != nil {
		t.Fatalf("Failed to get time entries: %v", err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("diff: %v", cmp.Diff(got, want))
	}
}
