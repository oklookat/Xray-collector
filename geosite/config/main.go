package config

import (
	"encoding/json"
	"os"
)

func Load(path string) (*Config, error) {
	cfgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()

	dec := json.NewDecoder(cfgFile)
	cfg := &Config{}
	if err := dec.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, err
}

type Config struct {
	CachePath    string    `json:"cachePath"`
	MaxCacheSize int64     `json:"maxCacheSize"`
	Input        []*Input  `json:"input"`
	Output       []*Output `json:"output"`
}
