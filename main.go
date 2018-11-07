package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	logfile, err := os.Create("connan.log")
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println(err)
	} else {
		defer logfile.Close()
		log.SetOutput(logfile)
	}
	a, err := NewApp()
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}
	err = a.run()
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
}
