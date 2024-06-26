package function

import (
	"fmt"
	"function/entity"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	gzb "github.com/ifo/gozulipbot"
	"os"
	"time"
)

// Handle an HTTP Request.
func Handle(event cloudevents.Event) error {
	var cloudEventData entity.CloudEventData
	if err := event.DataAs(&cloudEventData); err != nil {
		fmt.Println(err)
		return nil
	}

	bot := gzb.Bot{
		APIKey:  os.Getenv("API_KEY"),
		APIURL:  os.Getenv("API_URL"),
		Email:   os.Getenv("EMAIL"),
		Backoff: 1 * time.Second,
	}

	bot.Init()

	m := gzb.Message{
		Stream:  os.Getenv("STREAM"),
		Topic:   os.Getenv("TOPIC"),
		Content: fmt.Sprintf("```bash\n$ %s\n%s\n$ %s\n", cloudEventData.IODocument.Input, cloudEventData.IODocument.Output, cloudEventData.IODocument.PS1),
	}

	_, err := bot.Message(m)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}
