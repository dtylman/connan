package main

import (
	"encoding/json"
	"io/ioutil"
)

// Options application options
type Options struct {
	LibFolder string `json:"libfolder`
	Tesseract string `json:"tesseract`
}

var options Options

func (o *Options) load() error {
	data, err := ioutil.ReadFile("connan.config")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, o)
}

func (o *Options) save() error {
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("connan.config", data, 0644)
}
