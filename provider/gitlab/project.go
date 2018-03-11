package gitlab

import (
	"fmt"

	"github.com/etherlabsio/reviewit"
	"github.com/pkg/errors"
	client "github.com/xanzy/go-gitlab"
)

// Project returns pending PR's for a gitlab project
type Project struct {
	service
	project *client.Project
	all     []*client.MergeRequest
	pending []reviewit.MergeRequest
	err     error
}

func (f *Project) allMergeRequests() {
	if f.err != nil {
		return
	}
	fmt.Println("finding merge requests for project: ", f.project.Name)

	var err error
	f.all, _, err = f.client.MergeRequests.ListProjectMergeRequests(
		f.project.ID, &client.ListProjectMergeRequestsOptions{
			State:   client.String("opened"),
			OrderBy: client.String("updated_at"),
			Sort:    client.String("desc"),
			Scope:   client.String("all"),
		},
	)
	fmt.Println(f.all)
	f.err = errors.WithMessage(err, "find merge requests failed")
}

func (f *Project) getPending() {
	if f.err != nil {
		return
	}
	for _, mr := range f.all {
		if isPending(mr) {
			mergeReq := reviewit.MergeRequest{
				Author: reviewit.Author{
					Name: mr.Author.Name,
				},
				Reviewers: []reviewit.Reviewer{
					{
						Name: mr.Assignee.Name,
					},
				},
				Title:       mr.Title,
				URL:         mr.WebURL,
				CreatedAt:   *mr.CreatedAt,
				LastUpdated: *mr.UpdatedAt,
			}
			f.pending = append(f.pending, mergeReq)
		}
	}
}

// Filter performs the filter operation for a project
func (f *Project) Filter() ([]reviewit.MergeRequest, error) {
	f.allMergeRequests()
	f.getPending()
	if f.err != nil {
		return nil, errors.WithMessage(f.err, "gitlab: failed for project: "+f.project.NameWithNamespace)
	}
	return f.pending, nil
}
