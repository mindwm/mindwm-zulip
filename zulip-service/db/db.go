package db

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
)

func InitNeo4j(ctx context.Context) neo4j.DriverWithContext {
	driver, err := neo4j.NewDriverWithContext(
		os.Getenv("DB_URI"),
		neo4j.BasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), ""))
	if err != nil {
		log.Fatal(err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return driver
}
