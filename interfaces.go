package reviewit

// Sender sends the formatted message to a specific transport such as slack.
type Sender interface {
	Send() error
}

// Filterer filters for pending merge requests
type Filterer interface {
	Filter() ([]MergeRequest, error)
}
