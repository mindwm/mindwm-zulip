package function

import (
	"context"
	"fmt"
	"function/db"
	neo4jRepo "function/db/repository/neo4j"
	"function/entity"
	"os"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	gzb "github.com/ifo/gozulipbot"
)

type Repository interface {
	GetZulipBotState(ctx context.Context) (entity.ZulipBotState, error)
	CreateZulipBotState(ctx context.Context, state string) error
	UpdateZulipBotState(ctx context.Context, state string) error
}

// Handle an HTTP Request.
func Handle(ctx context.Context, event cloudevents.Event) error {
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

	driver := db.InitNeo4j(ctx)
	defer driver.Close(ctx)
	var repository Repository = neo4jRepo.NewRepository(driver)

	if IsZulipStateCmd(cloudEventData.IODocument.Input) {
		state := ConvertState(cloudEventData.IODocument.Input)

		zulipBotState, err := repository.GetZulipBotState(ctx)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if zulipBotState.State == "" {
			if err = repository.CreateZulipBotState(ctx, state); err != nil {
				fmt.Println(err)
				return nil
			}
		} else {
			if err = repository.UpdateZulipBotState(ctx, state); err != nil {
				fmt.Println(err)
				return nil
			}
		}

	} else {
		zulipBotState, err := repository.GetZulipBotState(ctx)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if zulipBotState.State == "" {
			err = repository.CreateZulipBotState(ctx, entity.StateZulipOff)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			zulipBotState.State = entity.StateZulipOff
		}

		if zulipBotState.State == entity.StateZulipOn {
			fmt.Println(11)
			m := gzb.Message{
				Stream:  os.Getenv("STREAM"),
				Topic:   os.Getenv("TOPIC"),
				Content: fmt.Sprintf("```bash\n$ %s\n%s\n$ %s\n", cloudEventData.IODocument.Input, cloudEventData.IODocument.Output, cloudEventData.IODocument.PS1),
			}

			_, err = bot.Message(m)
			if err != nil {
				fmt.Println(err)
				return nil
			}
		}

		return nil
	}
	return nil
}

func IsZulipStateCmd(input string) bool {
	return strings.Contains(input, "zulip on") || strings.Contains(input, "zulip off")
}

func ConvertState(input string) string {
	if strings.Contains(input, "zulip on") {
		return entity.StateZulipOn
	} else {
		return entity.StateZulipOff
	}
}
