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
	queued int
	db     *db.DB
	worker *db.Worker
}

//NewIndexer returns a new indexer
func NewIndexer(database *db.DB) (*Indexer, error) {
	i := new(Indexer)
	i.queued = 0
	i.db = database
	i.worker = db.NewWorker(i.db)
	return i, nil
}

//Start ...
func (i *Indexer) Start(root string) error {
	if i.worker.IsRunning() {
		return errors.New("Indexer running")
	}
	i.queued = 0
	i.db.Queue.Clear()
	log.Println("Indexer Started")
	defer func() {
		log.Printf("%v items queued", i.queued)
	}()
	i.worker.Start()
	return filepath.Walk(root, i.walk)
}

func (i *Indexer) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println(err)
		return nil
	}
	if !i.worker.IsRunning() {
		return errors.New("Worker stopped")
	}
	if !info.IsDir() {
		log.Printf("Queuing '%v'", path)
		i.db.Queue.Add(path)
		i.queued++
	}
	return nil
}

//Stop stop indexing
func (i *Indexer) Stop() {
	i.worker.Stop()
}

//Close closes the indexer
func (i *Indexer) Close() error {
	return i.db.Close()
}
