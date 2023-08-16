package database

import (
	"io/fs"
	"log"

	"gorm.io/gorm"
)

type Transaction struct {
	File   File
	Shards []Shard
	Refs   []Ref
}

func DoTransaction(transaction *Transaction, db *gorm.DB) error {
	// Start the transaction
	tx := db.Begin()

	// Function to handle rollback in case of an error
	rollbackOnError := func(err error) error {
		if err != nil {
			log.Println("Rolling back transaction...")
			tx.Rollback()
			log.Println("Transaction rolled back.")
			return err
		}
		return nil
	}

	// Defer rollback if there's a panic (unrecoverable error)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	log.Println("Writing file entry to database...")
	if err := rollbackOnError(tx.Create(transaction.File).Error); err != nil {
		return err
	}
	log.Println("File entry written to database.")

	log.Println("Writing shard entries to database...")
	if err := rollbackOnError(tx.CreateInBatches(transaction.Shards, len(transaction.Shards)).Error); err != nil {
		return err
	}
	log.Println("Shard entries written to database.")

	log.Println("Writing ref entries to database...")
	if err := rollbackOnError(tx.CreateInBatches(transaction.Refs, len(transaction.Refs)).Error); err != nil {
		return err
	}
	log.Println("Ref entries written to database.")

	log.Println("Committing transaction...")
	if err := rollbackOnError(tx.Commit().Error); err != nil {
		return err
	}
	log.Println("Transaction committed.")

	return nil
}

func CreateTransaction(
	path string,
	fileInfo fs.FileInfo,
	data []byte,
	shards [][]byte,
	dataShards int,
) *Transaction {
	// Create the file entry
	fileDBEntry := File{
		Path: path,
		Size: fileInfo.Size(),
		Hash: MkSriHash(Sha512, data),
	}

	// Create the shard and ref entries
	shardDBEntries := make([]Shard, len(shards))
	refDBEntries := make([]Ref, len(shards))
	for i, shard := range shards {
		shardDBEntries[i] = Shard{
			DeviceID: 0,
			Size:     int64(len(shard)),
			Hash:     MkSriHash(Sha512, shard),
			Number:   i,
			IsData:   i < dataShards,
		}
		refDBEntries[i] = Ref{
			FileHash:  fileDBEntry.Hash,
			ShardHash: shardDBEntries[i].Hash,
		}
	}

	return &Transaction{
		File:   fileDBEntry,
		Shards: shardDBEntries,
		Refs:   refDBEntries,
	}
}
