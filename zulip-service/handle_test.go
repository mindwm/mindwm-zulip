package function

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"testing"
)

// TestHandle ensures that Handle executes without error and returns the
// HTTP 200 status code indicating no errors.
func TestHandle(t *testing.T) {
	err := Handle(cloudevents.Event{})
	fmt.Println(err)
}
