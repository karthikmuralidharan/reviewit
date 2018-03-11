package gitlab

import (
	"sync"

	"github.com/etherlabsio/reviewit"
	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
)

// Group returns pending PR's for a gitlab project
type Group struct {
	service
	name    string
	err     error
	pending []reviewit.MergeRequest
}

func (f *Group) projects() []*gitlab.Project {
	if f.err != nil {
		return []*gitlab.Project{}
	}
	// List all projects
	projects, _, err := f.client.Groups.ListGroupProjects(f.name, &gitlab.ListGroupProjectsOptions{})
	f.err = errors.WithMessage(err, "failed to get project list")
	return projects
}

func (f *Group) getPending(projects []*gitlab.Project) {
	if f.err != nil {
		return
	}
	var mtx sync.Mutex
	var filterForProject = func(p *gitlab.Project) func() error {
		return func() error {
			pf := &Project{
				service: f.service,
				project: p,
			}
			mrs, err := pf.Filter()
			mtx.Lock()
			f.pending = append(f.pending, mrs...)
			mtx.Unlock()
			return err
		}
	}
	var g errgroup.Group
	for _, project := range projects {
		g.Go(filterForProject(project))
	}
	f.err = errors.WithMessage(g.Wait(), "group: filter err")
}

// Filter returns a list of pending merge requests based on a group's scope
func (f *Group) Filter() (result []reviewit.MergeRequest, err error) {
	projects := f.projects()
	f.getPending(projects)
	if f.err != nil {
		return result, f.err
	}
	return f.pending, nil
}
