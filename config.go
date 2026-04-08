package main

import (
	"encoding/json"
	"os"
)

var configfile = "config.json"

// EXPLINATION: parsing json file into struct
// TODO: consider typesafety
type Config struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	TimeOut int    `json:"timeout"`
	ACL     string `json:"ACL"`
}

// WARNING:  type cfgKey is a private type
// to avoid key collision, preserves typesaftey
var config Config

func loadConfig() error {

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
