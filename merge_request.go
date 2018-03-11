package reviewit

import (
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
)

// MergeRequest for a version control service such as github, gitlab etc..
type MergeRequest struct {
	Title       string
	URL         string
	Author      Author
	CreatedAt   time.Time
	LastUpdated time.Time
	Reviewers   []Reviewer
}

// UpdatedAt returns a friendly updated time
func (mr MergeRequest) UpdatedAt() string {
	return humanize.Time(mr.LastUpdated)
}

// ReviewerNames returns a comma separated list of all reviewer names
func (mr MergeRequest) ReviewerNames() string {
	var reviewerNames []string
	for _, reviewer := range mr.Reviewers {
		reviewerNames = append(reviewerNames, reviewer.Name)
	}
	return strings.Join(reviewerNames, ", ")
}

// Author is the person who raised the merge request for review
type Author struct {
	Name string
}

// Reviewer represents a person assigned to review a merge request
type Reviewer struct {
	Name string
}
