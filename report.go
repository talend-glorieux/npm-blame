package main

// Report represents a npm package issue report
type Report struct {
	Title     string
	Errors    []int
	Solutions []int
}

// NewReport returns a new issue report
func NewReport(title string) *Report {
	return &Report{}
}

// Send sends a report to the appropriate npm package issue tracker
// TODO: Implement
func (r *Report) Send() error {
	return nil
}
