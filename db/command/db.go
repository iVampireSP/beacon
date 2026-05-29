package command

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func printVersion(db *sql.DB) error {
	version, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}
	if version == 0 {
		fmt.Println("  Current version: none")
	} else {
		fmt.Printf("  Current version: %d\n", version)
	}
	return nil
}
