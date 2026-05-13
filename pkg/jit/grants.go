package jit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const grantsFileName = "jit-grants.json"

// Grant represents an active JIT access grant.
type Grant struct {
	Group     string    `json:"group"`
	ExpiresAt time.Time `json:"access_expires_at"`
}

func grantsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "ok", grantsFileName), nil
}

// LoadGrants reads active (non-expired) grants from disk.
func LoadGrants() ([]Grant, error) {
	path, err := grantsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var grants []Grant
	if err := json.Unmarshal(data, &grants); err != nil {
		return nil, err
	}

	// Filter out expired grants
	now := time.Now()
	active := grants[:0]
	for _, g := range grants {
		if g.ExpiresAt.After(now) {
			active = append(active, g)
		}
	}

	return active, nil
}

// SaveGrant adds a grant and persists to disk.
func SaveGrant(grant Grant) error {
	grants, err := LoadGrants()
	if err != nil {
		grants = nil
	}

	// Replace if same group already exists
	found := false
	for i, g := range grants {
		if g.Group == grant.Group {
			grants[i] = grant
			found = true
			break
		}
	}
	if !found {
		grants = append(grants, grant)
	}

	return saveGrants(grants)
}

// RemoveGrant removes a grant by group name and persists to disk.
func RemoveGrant(group string) error {
	grants, err := LoadGrants()
	if err != nil {
		return err
	}

	filtered := grants[:0]
	for _, g := range grants {
		if g.Group != group {
			filtered = append(filtered, g)
		}
	}

	return saveGrants(filtered)
}

func saveGrants(grants []Grant) error {
	path, err := grantsPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.Marshal(grants)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
