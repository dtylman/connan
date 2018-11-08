package db

import (
	"sync"
)

//Queue is a work queue
type Queue struct {
	mutex sync.Mutex
	items []string
}

//Add adds item to the queue
func (q *Queue) Add(path string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.items == nil {
		q.items = make([]string, 0)
	}
	q.items = append(q.items, path)
}

//Pop pops item from queue or nil if empty
func (q *Queue) Pop() *string {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	item := &q.items[0]
	q.items = q.items[1:]
	return item
}

//Clear ...
func (q *Queue) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.items = make([]string, 0)
}

//Len returns the number of items in the queue
func (q *Queue) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.items)
}
