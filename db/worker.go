package db

import (
	"log"
	"os"
	"sync"
	"time"
)

//Worker db worker indexes items in the db queue.
type Worker struct {
	db      *DB
	running bool
	mutex   sync.Mutex
}

//NewWorker creates a new worker
func NewWorker(db *DB) *Worker {
	w := new(Worker)
	w.db = db
	return w
}

//Start starts the worker
func (w *Worker) Start() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if !w.running {
		w.running = true
		go w.work()
	}
}

//Stop stops the worker
func (w *Worker) Stop() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.running = false
}

//IsRunning is true if worker is running
func (w *Worker) IsRunning() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.running
}

func (w *Worker) work() {
	log.Println("Worker started")
	defer func() {
		log.Println("Worker stopped")
		w.Stop()
	}()
	lastItem := time.Now()
	for true {
		if !w.IsRunning() {
			break
		}
		item := w.db.Queue.Pop()
		if item == nil { //queue empty
			if time.Since(lastItem) > time.Second*5 {
				break
			}
			time.Sleep(time.Second / 2)
		} else {
			lastItem = time.Now()
			err := w.process(*item)
			if err != nil {
				log.Printf("Failed to process '%v': %v", *item, err)
			}
		}
	}
}

func (w *Worker) process(path string) error {
	log.Printf("Processing '%v'", path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	doc, err := w.db.Document(path)
	if err != nil {
		doc, err = NewDocument(path, fileInfo)
	}
	if err != nil {
		return err
	}
	doc.UpdateFileInfo(fileInfo)
	updated := doc.Analyze(w.db.analyzers)
	if updated {
		return w.db.Save(doc)
	}
	return nil
}
