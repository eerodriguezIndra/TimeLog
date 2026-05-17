package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	IntervalMinutes int      `json:"interval_minutes"`
	CSVPath         string   `json:"csv_path"`
	Autostart       bool     `json:"autostart"`
	Clients         []string `json:"clients"`
	Activities      []string `json:"activities"`

	path string
	mu   sync.RWMutex
}

// Snapshot es una copia inmutable y sin lock de Config, segura para pasar
// entre goroutines y leer sin sincronización.
type Snapshot struct {
	IntervalMinutes int
	CSVPath         string
	Autostart       bool
	Clients         []string
	Activities      []string
}

const (
	defaultIntervalMinutes = 60
	configFileName         = "config.json"
	defaultCSVFileName     = "timelog.csv"
)

var defaultActivities = []string{
	"Reunión",
	"Plan de trabajo",
	"Documentación",
	"Ejecución de operación",
	"Comité de cambio",
	"Respuesta de correo",
}

func configDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "TimeLog")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func defaultCSVPath() (string, error) {
	docs, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(docs, "Documents", defaultCSVFileName), nil
}

func Load() (*Config, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, configFileName)

	c := &Config{path: path}

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		csvPath, cerr := defaultCSVPath()
		if cerr != nil {
			return nil, cerr
		}
		c.IntervalMinutes = defaultIntervalMinutes
		c.CSVPath = csvPath
		c.Autostart = false
		c.Activities = append([]string(nil), defaultActivities...)
		c.Clients = []string{}
		if err := c.Save(); err != nil {
			return nil, err
		}
		return c, nil
	}

	if err := json.Unmarshal(data, c); err != nil {
		return nil, err
	}
	c.path = path

	if c.IntervalMinutes <= 0 {
		c.IntervalMinutes = defaultIntervalMinutes
	}
	if len(c.Activities) == 0 {
		c.Activities = append([]string(nil), defaultActivities...)
	}
	if c.CSVPath == "" {
		if csvPath, err := defaultCSVPath(); err == nil {
			c.CSVPath = csvPath
		}
	}
	return c, nil
}

func (c *Config) Save() error {
	c.mu.RLock()
	data, err := json.MarshalIndent(c, "", "  ")
	c.mu.RUnlock()
	if err != nil {
		return err
	}
	tmp := c.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, c.path)
}

func (c *Config) Update(fn func(*Config)) error {
	c.mu.Lock()
	fn(c)
	c.mu.Unlock()
	return c.Save()
}

func (c *Config) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		IntervalMinutes: c.IntervalMinutes,
		CSVPath:         c.CSVPath,
		Autostart:       c.Autostart,
		Activities:      append([]string(nil), c.Activities...),
		Clients:         append([]string(nil), c.Clients...),
	}
}

func (c *Config) AddClient(name string) error {
	return c.Update(func(cc *Config) {
		for _, existing := range cc.Clients {
			if existing == name {
				return
			}
		}
		cc.Clients = append(cc.Clients, name)
	})
}
