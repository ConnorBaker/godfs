package decode

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/connorbaker/godfs/database"
	"github.com/connorbaker/godfs/utils"
	"github.com/klauspost/reedsolomon"
)

// TODO: Hard-coded output directory.
var inputDir, _ = filepath.Abs("./encoded_shards")

func verifyAndReconstruct(
	shards [][]byte,
	r reedsolomon.Encoder,
	shardsAreReconstructed bool,
) ([][]byte, error) {
	if shardsAreReconstructed {
		log.Println("Verifying reconstructed shards...")
	} else {
		log.Println("Verifying shards...")
	}
	ok, err := r.Verify(shards)

	if shardsAreReconstructed && err != nil {
		log.Println("Reconstructed shards failed verification:", err)
		return nil, err
	} else if shardsAreReconstructed && !ok {
		log.Println("Reconstructed shards failed verification.")
		return nil, fmt.Errorf("reconstructed shards failed verification")
	} else if shardsAreReconstructed && ok {
		log.Println("Reconstructed shards verified.")
		return shards, nil
	} else if !shardsAreReconstructed && err != nil {
		log.Println("Shards failed verification:", err)
	} else if !shardsAreReconstructed && !ok {
		log.Println("Shards failed verification.")
	} else if ok {
		log.Println("Shards verified.")
		return shards, nil
	}

	// If we get here, we need to reconstruct the data.
	log.Println("Reconstructing data...")
	err = r.Reconstruct(shards)
	if err != nil {
		log.Println("Reconstruction failed:", err)
		return nil, err
	}
	log.Println("Reconstruction succeeded.")
	return verifyAndReconstruct(shards, r, true)
}

func DecodeAndJoin(
	shards [][]byte,
	dataShards int,
	parityShards int,
	fileDBEntry database.File,
	path string,
) error {
	r, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		return fmt.Errorf("failed to create reed-solomon encoder: %w", err)
	}

	// Verify shards, reconstruct if necessary.
	shards, err = verifyAndReconstruct(shards, r, false)
	if err != nil {
		return fmt.Errorf("failed to verify and reconstruct shards: %w", err)
	}

	// Join shards.
	log.Println("Joining shards...")
	data := bytes.NewBuffer(make([]byte, 0, fileDBEntry.Size))
	err = r.Join(data, shards, int(fileDBEntry.Size))
	if err != nil {
		log.Println("Failed to join shards:", err)
		return err
	}
	log.Println("Joined shards.")

	// Verify the hash matches the file.
	log.Println("Verifying hash...")
	if hash := database.MkSriHash(database.Sha512, data.Bytes()); hash != fileDBEntry.Hash {
		return fmt.Errorf("hashes do not match: expected %s but got %s", fileDBEntry.Hash, hash)
	}
	log.Println("Hash verified.")

	// Write the file.
	log.Println("Writing file...")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	} else {
		defer file.Close()
	}
	_, err = file.Write(data.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	log.Println("Wrote file.")

	return nil
}

// Can return nil!
func ReadShard(shardDBEntry database.Shard) []byte {
	shardPath := filepath.Join(inputDir, string(shardDBEntry.Hash))
	log.Println("Reading shard", shardDBEntry.Hash, "...")
	shard, err := utils.ReadBytes(shardPath, shardDBEntry.Size)
	if err != nil {
		log.Println("Failed to read shard", shardDBEntry.Hash, ":", err)
		return nil
	}
	log.Println("Read shard", shardDBEntry.Hash, ".")

	return shard
}

func ReadShards(shardDBEntries []database.Shard) [][]byte {
	shards := make([][]byte, len(shardDBEntries))

	for i := 0; i < len(shards); i++ {
		shards[i] = ReadShard(shardDBEntries[i])
	}

	return shards
}
