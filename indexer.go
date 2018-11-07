package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/dtylman/connan/db"
)

//Indexer is the connan indexer
type Indexer struct {
	queuing bool
	db      *db.DB
	worker  *db.Worker
}

//NewIndexer returns a new indexer
func NewIndexer(d *db.DB) (*Indexer, error) {
	i := new(Indexer)
	i.queuing = false
	i.db = d
	i.worker = db.NewWorker(i.db)
	return i, nil
}

//Start ...
func (i *Indexer) Start(root string) error {
	if !i.queuing {
		i.queuing = true
		log.Println("Indexer started")
		defer func() {
			i.queuing = false
			log.Println("Indexer stopped")
		}()
		i.worker.Start()
		return filepath.Walk(root, i.walk)
	}
	return nil
}

func (i *Indexer) walk(path string, info os.FileInfo, err error) error {
	if !i.queuing {
		return errors.New("indexed stopped")
	}
	if !info.IsDir() {
		log.Printf("Queuing '%v'", path)
		i.db.Queue.Add(path)
	}
	return nil
}

//Stop stop indexing
func (i *Indexer) Stop() {
	i.queuing = false
	i.worker.Stop()
}
