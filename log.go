package main

import (
	"fmt"
	"log"
	"os"
)

//Log is the application log
type Log struct {
	file *os.File
	// the last log message
	Last string
}

//AppLog is the application logger
var AppLog Log

//Open opens the application log
func (ap *Log) Open() error {
	var err error
	AppLog.file, err = os.Create("connan.log")
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(ap)
	return err
}

//Close closes the log
func (ap *Log) Close() {
	err := ap.file.Close()
	if err != nil {
		fmt.Println(err)
	}
}

//Write writer imp
func (ap *Log) Write(b []byte) (n int, err error) {
	ap.Last = string(b)
	return ap.file.Write(b)
}
