package gitlab

import (
	"github.com/etherlabsio/reviewit"
	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
)

// User returns pending PR's for a gitlab project
type User struct {
	service
	err     error
	pending []reviewit.MergeRequest
}

func (f *User) projects() []*gitlab.Project {
	if f.err != nil {
		return []*gitlab.Project{}
	}
	// List all projects
	projects, _, err := f.client.Projects.ListProjects(&gitlab.ListProjectsOptions{
		OrderBy:    gitlab.String("updated_at"),
		Archived:   gitlab.Bool(false),
		Sort:       gitlab.String("desc"),
		Visibility: gitlab.Visibility(gitlab.PrivateVisibility),
	})
	f.err = errors.WithMessage(err, "failed to get project list")
	return projects
}

func (f *User) getPending(projects []*gitlab.Project) {
	if f.err != nil {
		return
	}
	var g errgroup.Group
	for _, project := range projects {
		g.Go(func() error {
			pf := &Project{
				service: f.service,
				project: project,
			}
			mrs, err := pf.Filter()
			f.pending = append(f.pending, mrs...)
			return err
		})
	}
	f.err = errors.WithMessage(g.Wait(), "user: filter err")
}

// Filter returns a list of pending merge requests based on a User's scope
func (f *User) Filter() (result []reviewit.MergeRequest, err error) {
	projects := f.projects()
	f.getPending(projects)
	if f.err != nil {
		return result, f.err
	}
	return f.pending, nil
}
