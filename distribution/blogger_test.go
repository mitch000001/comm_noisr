package distribution

import (
	"strings"
	"testing"
)

func TestBloggerSend(t *testing.T) {
	blogger := NewBloggerClient(BloggerBlogId)
	response, err := blogger.Send(strings.NewReader("foobar"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(response, "http://") {
		t.Fatalf("expected response to start with 'http://', got %s", response)
	}
}
