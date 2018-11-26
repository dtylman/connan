package main

import (
	"io/ioutil"

	"github.com/dtylman/connan/db"
)

//PlainTextAnalyzer Analyzes 'text/plain'
type PlainTextAnalyzer struct {
}

//Process...
func (ta PlainTextAnalyzer) Process(path string, doc *db.Document) error {
	if doc.Mime() == "text/plain" {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		doc.SetContent(string(data))
	}
	return nil
}

//Name ...
func (ta PlainTextAnalyzer) Name() string {
	return "PlainTextAnalyzer"
}
