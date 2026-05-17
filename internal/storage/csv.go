package storage

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Entry struct {
	Timestamp   time.Time
	Client      string
	Activity    string
	Description string
}

type CSVStore struct {
	mu sync.Mutex
}

func New() *CSVStore {
	return &CSVStore{}
}

var header = []string{"fecha", "hora", "cliente", "actividad", "descripcion"}

func (s *CSVStore) Append(path string, entry Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	info, err := os.Stat(path)
	needsHeader := false
	switch {
	case os.IsNotExist(err):
		needsHeader = true
	case err != nil:
		return err
	case info.Size() == 0:
		needsHeader = true
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if needsHeader {
		if err := w.Write(header); err != nil {
			return err
		}
	}

	row := []string{
		entry.Timestamp.Format("2006-01-02"),
		entry.Timestamp.Format("15:04:05"),
		entry.Client,
		entry.Activity,
		entry.Description,
	}
	return w.Write(row)
}
