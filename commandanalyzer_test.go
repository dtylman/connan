package main

import (
	"testing"

	"github.com/dtylman/connan/db"
	"github.com/stretchr/testify/assert"
)

func TestCondition_Match(t *testing.T) {
	db, err := db.Open("connan.testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	c := Condition{Expression: "*.github.com", Field: "Location"}
	doc, err := db.NewDocument("commandanalyzer_test.go")
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
