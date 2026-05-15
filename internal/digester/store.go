package digester

import (
	"encoding/json"
	"errors"
	"os"
)

// Store persists and retrieves Digest values from a JSON file on disk.
type Store struct {
	path string
}

// NewStore returns a Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the digest to the store file, creating or overwriting it.
func (s *Store) Save(d Digest) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Load reads the digest from the store file. If the file does not exist,
// it returns a zero Digest and ErrNotFound.
func (s *Store) Load() (Digest, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Digest{}, ErrNotFound
		}
		return Digest{}, err
	}
	var d Digest
	if err := json.Unmarshal(data, &d); err != nil {
		return Digest{}, err
	}
	return d, nil
}

// ErrNotFound is returned by Load when no stored digest exists yet.
var ErrNotFound = errors.New("digester: no stored digest found")
