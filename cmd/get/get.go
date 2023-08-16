package get

import (
	"fmt"
	"path/filepath"

	"github.com/connorbaker/godfs/database"
	"github.com/connorbaker/godfs/decode"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func GetFile(
	toPath string,
	file database.File,
	db *gorm.DB,
) error {
	// Query the database for the shards to use
	shardDBEntries, err := database.GetShards(file.Hash, db)
	if err != nil {
		return fmt.Errorf("failed to get shards from database: %w", err)
	}

	// Take note of the number of data and parity shards
	var dataShards, parityShards int
	for _, shard := range shardDBEntries {
		if shard.IsData {
			dataShards++
		} else {
			parityShards++
		}
	}

	// Read the shards from disk
	shards := decode.ReadShards(shardDBEntries)

	// Decode the file and write it to disk
	err = decode.DecodeAndJoin(shards, dataShards, parityShards, file, toPath)
	if err != nil {
		return fmt.Errorf("failed to decode file: %w", err)
	}

	return nil
}

func action(ctx *cli.Context) error {
	fromPath, toPath, err := parsePositionalArgs(ctx)
	if err != nil {
		return err
	}

	db, err := database.OpenDB(database.DB_PATH)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Make sure the database has the files table
	if !database.TablesExist(db) {
		return fmt.Errorf("database does not have all tables -- make sure to run init first")
	}

	// Try to get the file from the database
	file, err := database.GetFileByPath(fromPath, db)
	if err == gorm.ErrRecordNotFound {
		return fmt.Errorf("file does not exist in database")
	} else if err != nil {
		return fmt.Errorf("failed to get file from database: %w", err)
	}

	return GetFile(toPath, file, db)
}

func parsePositionalArgs(ctx *cli.Context) (string, string, error) {
	if ctx.NArg() != 2 {
		return "", "", fmt.Errorf("expected 2 arguments, got %d", ctx.NArg())
	}
	fromPath, toPath := ctx.Args().Get(0), ctx.Args().Get(1)

	absFromPath, err := filepath.Abs(fromPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	absToPath, err := filepath.Abs(toPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	return absFromPath, absToPath, nil
}

func GetCommand() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Retrieve a file from the database and write it to disk.",
		ArgsUsage: "<from path> <to path>",
		Action:    action,
	}
}
