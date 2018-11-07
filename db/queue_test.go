package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue_Clear(t *testing.T) {
	q := Queue{}
	item := q.Pop()
	assert.Nil(t, item)
	q.Clear()
	q.Add("lala")
	assert.NotEmpty(t, q.items)
	item = q.Pop()
	assert.NotNil(t, item)
	assert.Empty(t, q.items)
	assert.EqualValues(t, "lala", *item)
	q.Add("lala")
	q.Clear()
	assert.Empty(t, q.items)
}
