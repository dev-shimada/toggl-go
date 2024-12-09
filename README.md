[![Go Report Card](https://goreportcard.com/badge/github.com/dev-shimada/toggl-go)](https://goreportcard.com/report/github.com/dev-shimada/toggl-go)
[![Coverage Status](https://coveralls.io/repos/github/dev-shimada/toggl-go/badge.svg?branch=main)](https://coveralls.io/github/dev-shimada/toggl-go?branch=main)
[![CI](https://github.com/dev-shimada/toggl-go/actions/workflows/CI.yaml/badge.svg)](https://github.com/dev-shimada/toggl-go/actions/workflows/CI.yaml)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/dev-shimada/toggl-go/blob/master/LICENSE)

# Toggl Time Entries Go Client

This project provides a Go client for interacting with the [Toggl Track API](https://developers.track.toggl.com/docs/). It allows developers to integrate Toggl time tracking functionalities into their Go applications.

## Features

- Retrieve time entries with various query parameters
- Fetch the current running time entry
- Get a specific time entry by ID
- Create new time entries
- Bulk edit time entries
- Update existing time entries
- Delete time entries
- Stop a running time entry

## Installation

Use `go get` to install the package:
```plaintext
go get github.com/dev-shimada/toggl-go
```

## Usage

Here is an example of how to use the package:

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dev-shimada/toggl-go/timeentries"
	"github.com/dev-shimada/toggl-go/toggl"
)

func main() {
	// Get your API token from https://toggl.com/app/profile
	token := os.Getenv("TOKEN")
	client := toggl.NewClient(token)

	// Get time entries for the last 7 days
	now := time.Now()
	start := now.Add(-24 * 7 * time.Hour).Format(time.RFC3339)
	end := now.Format(time.RFC3339)
	result, err := client.TimeEntriesClient.GetTimeEntries(
		timeentries.GetTimeEntriesInput{
			Query: timeentries.GetTimeEntriesQuery{
				StartDate: &start,
				EndDate:   &end,
			},
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range result {
		fmt.Printf(
			"ID: %d, Description: %s, Start: %s, Stop: %s, Duration: %d\n",
			*v.Id, *v.Description, *v.Start, *v.Stop, *v.Duration,
		)
	}
}
```
