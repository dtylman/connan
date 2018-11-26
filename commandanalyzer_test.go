package main

import (
	"os"
	"testing"

	"github.com/dtylman/connan/db"
	"github.com/stretchr/testify/assert"
)

func TestCondition_Match(t *testing.T) {
	connandb, err := db.Open("connan.testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer connandb.Close()
	c := Condition{Pattern: "github.com", Field: "Location"}
	path := "commandanalyzer_test.go"
	fileInfo, err := os.Stat(path)
	assert.NoError(t, err)
	doc, err := db.NewDocument(path, fileInfo)
	assert.NoError(t, err)
	b, err := c.Match(doc)
	assert.False(t, b)
	assert.NoError(t, err)
	doc.SetField("Location", "nana")
	b, err = c.Match(doc)
	assert.False(t, b)
	assert.NoError(t, err)
	doc.SetField("Location", "api.github.com")
	b, err = c.Match(doc)
	assert.True(t, b)
	assert.NoError(t, err)
}
