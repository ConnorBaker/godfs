package encode

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/connorbaker/godfs/database"
	"github.com/klauspost/reedsolomon"
)

// TODO: Hard-coded output directory.
var outputDir, _ = filepath.Abs("./encoded_shards")

func SplitAndEncode(
	data []byte,
	dataShards int,
	parityShards int,
) ([][]byte, error) {
	r, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		return nil, fmt.Errorf("failed to create reed-solomon encoder: %w", err)
	}

	// Split into equal-length shards.
	shards, err := r.Split(data)
	if err != nil {
		return nil, fmt.Errorf("failed to split data into shards: %w", err)
	}

	// Encode parity
	err = r.Encode(shards)
	if err != nil {
		return nil, fmt.Errorf("failed to encode parity: %w", err)
	}

	return shards, nil
}

func WriteShard(shard []byte, shardDBEntry database.Shard) error {
	shardPath := filepath.Join(outputDir, string(shardDBEntry.Hash))
	shardFile, err := os.OpenFile(shardPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open shard file: %w", err)
	} else {
		defer shardFile.Close()
	}

	_, err = shardFile.Write(shard)
	if err != nil {
		return fmt.Errorf("failed to write shard to file: %w", err)
	}

	return nil
}

func WriteShards(shards [][]byte, shardDBEntries []database.Shard) error {
	// Create output directory if it doesn't exist.
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		log.Println("Creating output directory...")
		err := os.Mkdir(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
		log.Println("Created output directory.")
	}

	for i, shard := range shards {
		log.Println("Writing shard", i, "with hash", shardDBEntries[i].Hash, "...")
		if err := WriteShard(shard, shardDBEntries[i]); err != nil {
			return fmt.Errorf("failed to write shard: %w", err)
		}
		log.Println("Wrote shard", i, "with hash", shardDBEntries[i].Hash, ".")
	}

	return nil
}
