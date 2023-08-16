package database

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
)

type Algorithm string

const (
	Sha256 Algorithm = "sha256"
	Sha512 Algorithm = "sha512"
)

// TODO: Don't use a fixed size for the hash.
type hash struct {
	Algorithm Algorithm
	Value     [64]byte
}

func mkHash(algorithm Algorithm, data []byte) hash {
	switch algorithm {
	case Sha256:
		sum := sha256.Sum256(data)
		var extended [64]byte
		copy(extended[:32], sum[:])
		return hash{
			Algorithm: algorithm,
			Value:     extended,
		}
	case Sha512:
		return hash{
			Algorithm: algorithm,
			Value:     sha512.Sum512(data),
		}
	default:
		log.Panicln("Unknown algorithm:", algorithm)
		// This is unreachable.
		return hash{}
	}
}

type SriHash string

// Converts a s of the form "<type>-<base64>" to a SriHash.
func MkSriHash(algorithm Algorithm, data []byte) SriHash {
	hash := mkHash(algorithm, data)
	var hashLength int

	switch algorithm {
	case Sha256:
		hashLength = 32
	case Sha512:
		hashLength = 64
	default:
		log.Panicln("Unknown algorithm:", algorithm)
	}

	base64Hash := base64.URLEncoding.EncodeToString(hash.Value[:hashLength])
	return SriHash(fmt.Sprintf("%s-%s", algorithm, base64Hash))
}
