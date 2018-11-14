package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//OptionsFileName holds the default options config file
var OptionsFileName = "connan.config.json"

// Options application options
type Options struct {
	LibFolder string            `json:"libFolder"`
	DBFolder  string            `json:dbFolder`
	Analyzers []CommandAnalyzer `json:"analyzers"`
}

//Load loads options from file
func (o *Options) Load() error {
	_, err := os.Stat(OptionsFileName)
	if err == os.ErrNotExist {
		return nil
	}
	data, err := ioutil.ReadFile(OptionsFileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, o)
}

//Save saves options to file
func (o *Options) Save() error {
	data, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(OptionsFileName, data, 0644)
}
