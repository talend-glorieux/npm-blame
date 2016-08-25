package npmblame

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

var (
	errors = []int{42, 42}
)

// String is a utility function that allocates
// and returns the address of a given string
func String(s string) *string { return &s }

// Int is a utility function that allocates
// and returns the address of a given int
func Int(i int) *int { return &i }

func TestNewReport(t *testing.T) {
	r := NewReport("npm-blame", "test", errors)
	if r.Title != defaultTitle {
		t.Errorf("Wrong title. Expected test got %s", r.Title)
	}
}

func TestSend(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url

	mux.HandleFunc("/repos/npm-blame/test/issues", func(w http.ResponseWriter, r *http.Request) {
		ir := new(github.IssueRequest)
		json.NewDecoder(r.Body).Decode(ir)

		expectedTitle := defaultTitle
		expectedBody := defaultBody

		expected := &github.IssueRequest{
			Title: String(expectedTitle),
			Body:  String(expectedBody),
		}
		if !reflect.DeepEqual(ir, expected) {
			t.Errorf("Request body = %+v, want %+v", ir, expected)
		}
		fmt.Fprint(w, `{"number":1}`)
	})

	r := NewReport("npm-blame", "test", errors)
	issue, err := r.Send(client, "testToken")
	if err != nil {
		t.Error(err)
	}
	want := &github.Issue{Number: Int(1)}
	if !reflect.DeepEqual(issue, want) {
		t.Errorf("Issues.Create returned %+v, want %+v", issue, want)
	}
}
