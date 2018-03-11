package slack

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/etherlabsio/reviewit"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

type postMessageReq struct {
	slack.PostMessageParameters
	Text string `json:"text,omitempty"`
}

// Slack transport for the PR message pusher
type Slack struct {
	msg     postMessageReq
	pending []reviewit.MergeRequest

	Channel    string
	WebhookURL string
}

// New returns a slack transport instance
func New(channelName, webhookURL string) (*Slack, error) {
	s := &Slack{
		Channel:    channelName,
		WebhookURL: webhookURL,
	}
	if s.Channel == "" {
		s.Channel = "#general"
	}
	if s.WebhookURL == "" {
		return nil, errors.New("slack webhook URL is not valid")
	}
	return s, nil
}

func (s *Slack) build() {
	const defaultText = "Pull Requests To Review"

	s.msg = postMessageReq{slack.NewPostMessageParameters(), defaultText}
	s.msg.Channel = s.Channel
	s.msg.Markdown = true

	for _, mr := range s.pending {
		attachment := slack.Attachment{
			Title:     mr.Title,
			TitleLink: mr.URL,
			Fields: []slack.AttachmentField{
				{
					Title: "Author",
					Value: mr.Author.Name,
				},
				{
					Title: "Last Updated",
					Value: mr.UpdatedAt(),
				},
				{
					Title: "Reviewers",
					Value: mr.ReviewerNames(),
				},
			},
		}
		s.msg.Attachments = append(s.msg.Attachments, attachment)
	}
}

// Send builds and sends the message to slack
func (s *Slack) Send(pending []reviewit.MergeRequest) error {
	s.pending = pending
	s.build()
	buf, err := json.Marshal(&s.msg)
	if err != nil {
		return errors.WithMessage(err, "failure to Marshal")
	}
	req, _ := http.NewRequest("POST", s.WebhookURL, bytes.NewReader(buf))
	if _, err := http.DefaultClient.Do(req); err != nil {
		return errors.WithMessage(err, "slack post message failed")
	}
	return nil
}
