package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocument_String(t *testing.T) {
	d := new(Document)
	d.Fields = make(map[string]string)
	d.Path = "p"
	assert.EqualValues(t, d.Path, d.GetField("Path"))
	assert.Empty(t, d.GetField("p"))
	d.SetField("Path", "p1")
	assert.EqualValues(t, "p1", d.Path)
	assert.EqualValues(t, d.Path, d.GetField("Path"))
}

func TestDocument_IsImage(t *testing.T) {
	d := new(Document)
	d.Fields = make(map[string]string)
	assert.False(t, d.IsImage())
	d.SetField("mime", "jpeg")
	assert.False(t, d.IsImage())
	d.SetField("mime", "image/jpeg")
	assert.True(t, d.IsImage())
}
