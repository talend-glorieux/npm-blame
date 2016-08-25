package npmblame

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	defaultTitle = `Errors from npm-blame`
	defaultBody  = ``
)

// Report represents a npm package issue report
type Report struct {
	Title      string
	Body       string
	Owner      string
	Repository string
	Errors     []int
	Solutions  []int
}

// NewReport returns a new issue report
// based on the errors types
func NewReport(owner string, repo string, Errors []int) *Report {
	return &Report{
		Title:      defaultTitle,
		Body:       defaultBody,
		Owner:      owner,
		Repository: repo,
	}
}

func getDefaultClient(authToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return client
}

// Send sends a report to the appropriate npm package issue tracker
func (r *Report) Send(client *github.Client, authToken string) (issue *github.Issue, err error) {
	if client == nil {
		client = getDefaultClient(authToken)
	}

	issue, _, err = client.Issues.Create(r.Owner, r.Repository, &github.IssueRequest{
		Title: &r.Title,
		Body:  &r.Body,
	})
	return issue, err
}
