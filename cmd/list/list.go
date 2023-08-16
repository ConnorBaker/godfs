package list

import (
	"fmt"
	"log"

	"github.com/connorbaker/godfs/database"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func action(ctx *cli.Context) error {
	if ctx.NArg() > 0 {
		return fmt.Errorf("no positional arguments allowed")
	}

	pattern := ctx.String("pattern")
	db, err := database.OpenDB(database.DB_PATH)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Make sure the database has the files table
	if !database.TablesExist(db) {
		return fmt.Errorf("database does not have all tables -- make sure to run init first")
	}

	var paths []string
	if pattern == "" {
		paths, err = database.GetPaths(db)
	} else {
		paths, err = database.GetPathsMatching(pattern, db)
	}

	if len(paths) == 0 || err == gorm.ErrRecordNotFound {
		log.Println("no paths found")
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to get paths: %w", err)
	}

	for _, path := range paths {
		fmt.Println(path)
	}

	return nil
}

func ListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List files in the database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "pattern",
				Usage: "SQLite string pattern to match against",
			},
		},
		Action: action,
	}
}
