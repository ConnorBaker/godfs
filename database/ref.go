package database

// A reference from a file to shards.
// Each file is split into a number of shards.
// We must re-assemble the shards, in the original order, to get the original file.
type Ref struct {
	FileHash  SriHash `gorm:"primaryKey;type:text"`
	ShardHash SriHash `gorm:"primaryKey;type:text"`
	File      File    `gorm:"foreignKey:FileHash;references:Hash"`
	Shard     Shard   `gorm:"foreignKey:ShardHash;references:Hash"`
}
