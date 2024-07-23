package main

import (
	"context"
	"fmt"
	"log"

	"sandslash/service"

	"github.com/jackc/pgx/v5"
)

func main() {
	sandslashConfig := service.LoadConfig()

	connString := sandslashConfig.CreateDatabaseConnectionUri()

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	sql := `DROP TABLE IF EXISTS goose_db_version;`

	_, err = conn.Exec(context.Background(), sql)
	if err != nil {
		log.Fatalf("Failed to execute SQL command: %v", err)
	}

	fmt.Println("Table 'goose_db_version' dropped successfully.")
}
