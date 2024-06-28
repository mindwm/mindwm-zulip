package neo4j

import (
	"context"
	"fmt"
	"function/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repository struct {
	driver neo4j.DriverWithContext
}

func NewRepository(driver neo4j.DriverWithContext) *Repository {
	return &Repository{driver: driver}
}

func (r Repository) GetZulipBotState(ctx context.Context) (entity.ZulipBotState, error) {
	const op = "Neo4j.Repository.GetZulipBotState"

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := "MATCH (n:ZulipBotState) RETURN n LIMIT 1"
	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return entity.ZulipBotState{}, fmt.Errorf("%s: %v", op, err)
	}

	if result.Next(context.Background()) {
		record := result.Record()
		fmt.Println(record)
		node, found := record.Get("n")
		if !found {
			return entity.ZulipBotState{}, nil
		}

		props := node.(neo4j.Node).Props
		state := entity.ZulipBotState{
			ID:    node.(neo4j.Node).ElementId,
			State: entity.State(props["state"].(string)),
		}

		return state, nil
	}

	return entity.ZulipBotState{}, result.Err()
}

func (r Repository) CreateZulipBotState(ctx context.Context, state string) error {
	const op = "Neo4j.Repository.CreateZulipBotState"

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := "CREATE (n:ZulipBotState{state: $state})"
	params := map[string]interface{}{
		"state": state,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return result.Err()
}

func (r Repository) UpdateZulipBotState(ctx context.Context, state string) error {
	const op = "Neo4j.Repository.UpdateZulipBotState"

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := "MATCH (n:ZulipBotState) SET n.state = $state"
	params := map[string]interface{}{
		"state": state,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return result.Err()
}
