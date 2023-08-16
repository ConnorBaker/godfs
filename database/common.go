package database

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB_PATH, _ = filepath.Abs("godfs.sqlite")

func OpenDB(path string) (*gorm.DB, error) {
	log.Println("Opening DB at", path, "...")
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	log.Println("DB opened successfully.")

	return db, nil
}

func TablesExist(db *gorm.DB) bool {
	log.Println("Checking if tables exist...")
	tablesExist := true
	for _, model := range []interface{}{
		&Ref{},
		&Shard{},
		&File{},
	} {
		tableExists := db.Migrator().HasTable(model)
		if tableExists {
			log.Printf("Table for model %T exists.\n", model)
		} else {
			log.Printf("Table for model %T does not exist.\n", model)
		}
		tablesExist = tablesExist && tableExists
	}
	return tablesExist
}

func MigrateTables(db *gorm.DB) error {
	log.Println("Migrating tables...")
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

	// Create / migrate tables
	for _, model := range []interface{}{
		&Ref{},
		&Shard{},
		&File{},
	} {
		log.Printf("Migrating table for model %T...", model)
		if err := rollbackOnError(tx.AutoMigrate(model)); err != nil {
			return fmt.Errorf("failed to migrate table: %w", err)
		}
		log.Printf("Table for model %T migrated.", model)
	}

	log.Println("Committing transaction...")
	if err := rollbackOnError(tx.Commit().Error); err != nil {
		return err
	}
	log.Println("Transaction committed.")

	log.Println("Tables migrated successfully.")

	return nil
}
