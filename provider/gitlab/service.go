package gitlab

import (
	"github.com/etherlabsio/reviewit"

	gitlab "github.com/xanzy/go-gitlab"
)

type service struct {
	client *gitlab.Client
}

func isPending(mr *gitlab.MergeRequest) bool {
	// idleHours := time.Now().Sub(*mr.UpdatedAt).Hours()
	return !mr.WorkInProgress &&
		mr.Upvotes == 0 &&
		mr.Downvotes == 0 &&
		mr.State == "opened" &&
		mr.MergeStatus == "can_be_merged"
	// idleHours > float64(5)
}

// NewGroupFilter returns an instance of user scope PR filter
func NewGroupFilter(name, token string) reviewit.Filterer {
	s := service{
		client: gitlab.NewClient(nil, token),
	}
	return &Group{service: s, name: name}
}
