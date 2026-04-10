package main

import (
	"encoding/json"
	"os"
	"time"
)

// EXPLINATION: parsing json file into struct
// TODO: consider typesafety
type Config struct {
	Host    string        `json:"host"`
	Port    int           `json:"port"`
	Timeout time.Duration `json:"timeout"`
	ACL     string        `json:"ACL"`
}

// WARNING:  type cfgKey is a private type
// to avoid key collision, preserves typesaftey
var config Config
var configfile = "config.json"

func loadConfig() error {

	loadDefaultConfig()

	jsonData, err := os.ReadFile(configfile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return err
	}

	return nil
}

// EXPLINATION: default config values, overwritten when json with specified key is inserted
func loadDefaultConfig() {
	config = Config{
		Host:    "0.0.0.0",
		Port:    3333,
		Timeout: 3 * time.Second,
		ACL:     "https://easylist.to/easylist/easylist.txt",
	}
}
