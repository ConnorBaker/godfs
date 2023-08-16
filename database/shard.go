package database

import (
	"gorm.io/gorm"
)

// Path to the shard is always its hash.
type Shard struct {
	DeviceID uint64  // ID of the device that the file is on
	Size     int64   // Size of the file in bytes
	Hash     SriHash `gorm:"primaryKey;type:text;unique"` // Hash of the file
	Number   int     // Number of the shard
	IsData   bool    // True if this is a data shard, false if it is a parity shard
}

// Given a FileHash, return all the Shards that make up that file, in order.
func GetShards(fileHash SriHash, db *gorm.DB) ([]Shard, error) {
	var shards []Shard
	err := db.
		Model(&Shard{}).
		Joins("JOIN refs ON refs.shard_hash = shards.hash").
		Where("refs.file_hash = ?", fileHash).
		Order("number ASC").
		Find(&shards).
		Error

	return shards, err
}
