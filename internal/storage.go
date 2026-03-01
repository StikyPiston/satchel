package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Item struct {
	Name  string    `json:"name"`
	Size  int64     `json:"size"`
	Added time.Time `json:"added"`
}

type Index struct {
	NextID int64          `json:"next_id"`
	Items  map[int64]Item `json:"items"`
}

func baseDir() string {
	runtime := os.Getenv("XDG_RUNTIME_DIR")
	if runtime == "" {
		runtime = os.TempDir()
	}
	return filepath.Join(runtime, fmt.Sprintf("satchel-%d", os.Getuid()))
}

func ObjectsDir() string {
	return filepath.Join(baseDir(), "objects")
}

func indexPath() string {
	return filepath.Join(baseDir(), "index.json")
}

func Ensure() error {
	return os.MkdirAll(ObjectsDir(), 0755)
}

func Load() (*Index, error) {
	if err := Ensure(); err != nil {
		return nil, err
	}

	path := indexPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Index{
			NextID: 1,
			Items:  make(map[int64]Item),
		}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}

	if idx.Items == nil {
		idx.Items = make(map[int64]Item)
	}

	return &idx, nil
}

func Save(idx *Index) error {
	tmp := indexPath() + ".tmp"

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmp, indexPath())
}
