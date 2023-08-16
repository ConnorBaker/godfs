package database

import (
	"gorm.io/gorm"
)

// A file in the store.
// TODO: What if we have multiple copies of a file? Do we need both Hash and Path as primary keys?
type File struct {
	Size int64   // Size of the file in bytes
	Hash SriHash `gorm:"primaryKey;type:text;unique"` // Hash of the file
	Path string  `gorm:"type:text;unique"`            // Absolute path to the file
}

func PutFile(file File, db *gorm.DB) error {
	return db.Create(&file).Error
}

func GetFileByHash(hash SriHash, db *gorm.DB) (File, error) {
	var file File
	err := db.
		Model(&File{}).
		Where("hash = ?", hash).
		First(&file).
		Error
	return file, err
}

func GetFileByPath(path string, db *gorm.DB) (File, error) {
	var file File
	err := db.
		Model(&File{}).
		Where("path = ?", path).
		First(&file).
		Error
	return file, err
}

func GetPathsMatching(pattern string, db *gorm.DB) ([]string, error) {
	var paths []string
	err := db.
		Model(&File{}).
		Where("path LIKE ?", pattern).
		Pluck("path", &paths).
		Error
	return paths, err
}

func GetPaths(db *gorm.DB) ([]string, error) {
	return GetPathsMatching("%", db)
}

func GetNumberOfPaths(db *gorm.DB) (int64, error) {
	return GetNumberOfPathsMatching("%", db)
}

func GetNumberOfPathsMatching(pattern string, db *gorm.DB) (int64, error) {
	var count int64
	err := db.
		Model(&File{}).
		Where("path LIKE ?", pattern).
		Count(&count).
		Error
	return count, err
}
