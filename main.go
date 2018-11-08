package main

import (
	"fmt"
	"log"
)

func main() {
	err := AppLog.Open()
	if err != nil {
		panic(err)
	}
	defer AppLog.Close()
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
