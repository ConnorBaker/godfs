package init

import (
	"fmt"
	"log"

	"github.com/connorbaker/godfs/database"
	"github.com/urfave/cli/v2"
)

func action(ctx *cli.Context) error {
	if ctx.NArg() > 0 {
		return fmt.Errorf("no positional arguments allowed")
	}
	reinit := ctx.Bool("reinit")

	db, err := database.OpenDB(database.DB_PATH)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	tablesExist := database.TablesExist(db)
	if tablesExist && !reinit {
		return fmt.Errorf("database already initialized; to migrate existing tables, specify the 'reinit' flag")
	} else if tablesExist && reinit {
		log.Println("Will migrate existing tables to the latest version because the 'reinit' flag was specified")
	}

	return database.MigrateTables(db)
}

func InitCommand() *cli.Command {
	return &cli.Command{
		Name:   "init",
		Usage:  "Creates the database and initializes the file system",
		Action: action,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "reinit",
				Usage: "Reinitializes the database by migrating existing tables to the latest version",
			},
		},
	}
}
