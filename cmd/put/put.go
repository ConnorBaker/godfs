package put

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/connorbaker/godfs/database"
	"github.com/connorbaker/godfs/encode"
	"github.com/connorbaker/godfs/utils"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func PutFile(
	path string,
	fileInfo fs.FileInfo,
	dataShards int,
	parityShards int,
	db *gorm.DB,
) error {
	// Read the file
	data, err := utils.ReadBytes(path, fileInfo.Size())
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Encode the file
	shards, err := encode.SplitAndEncode(data, dataShards, parityShards)
	if err != nil {
		return fmt.Errorf("failed to encode file: %w", err)
	}

	// Create arguments for a transaction
	transaction := database.CreateTransaction(path, fileInfo, data, shards, dataShards)

	// Write the shards to disk
	if err := encode.WriteShards(shards, transaction.Shards); err != nil {
		return fmt.Errorf("failed to write shards to disk: %w", err)
	}

	// Update the database
	if err := database.DoTransaction(transaction, db); err != nil {
		return fmt.Errorf("failed to put file in database: %w", err)
	}

	return nil
}

func action(ctx *cli.Context) error {
	path, fileInfo, err := parsePositionalArgs(ctx)
	if err != nil {
		return err
	}

	dataShards := ctx.Int("data-shards")
	parityShards := ctx.Int("parity-shards")

	db, err := database.OpenDB(database.DB_PATH)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Make sure the database has the files table
	if !database.TablesExist(db) {
		return fmt.Errorf("database does not have all tables -- make sure to run init first")
	}

	// Try to find the file in the database
	_, err = database.GetFileByPath(path, db)
	if err == nil {
		return fmt.Errorf("file already exists in database: %s", path)
	} else if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get file from database: %w", err)
	}

	return PutFile(path, fileInfo, dataShards, parityShards, db)
}

func parsePositionalArgs(ctx *cli.Context) (string, fs.FileInfo, error) {
	if ctx.NArg() != 1 {
		return "", nil, fmt.Errorf("expected 1 argument, got %d", ctx.NArg())
	}

	absPath, err := filepath.Abs(ctx.Args().First())
	if err != nil {
		return "", nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	fileInfo, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return "", nil, fmt.Errorf("input file does not exist %w", err)
	}

	return absPath, fileInfo, nil
}

func dataShardsFlagAction(_ *cli.Context, data int) error {
	if data < 1 {
		return fmt.Errorf("data shards must be at least 1")
	} else if data > 256 {
		return fmt.Errorf("data shards must be below 257")
	}
	return nil
}

func dataShardsFlag() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    "data-shards",
		Aliases: []string{"k"},
		Usage:   "Number of shards to split the data into, must be in [1, 256].",
		Value:   4,
		Action:  dataShardsFlagAction,
	}
}

func parityShardsFlagAction(_ *cli.Context, parity int) error {
	if parity < 1 {
		return fmt.Errorf("parity shards must be at least 1")
	}
	return nil
}

func parityShardsFlag() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    "parity-shards",
		Aliases: []string{"m"},
		Usage:   "Number of parity shards",
		Value:   2,
		Action:  parityShardsFlagAction,
	}
}

func PutCommand() *cli.Command {
	return &cli.Command{
		Name:      "put",
		Usage:     "Encode file and put it in the database",
		ArgsUsage: "<path>",
		Action:    action,
		Flags: []cli.Flag{
			dataShardsFlag(),
			parityShardsFlag(),
		},
	}
}
