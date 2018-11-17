package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/thecodeteam/goodbye"
)

func main() {
	ctx := context.Background()
	defer goodbye.Exit(ctx, -1)

	goodbye.Notify(ctx)

	goodbye.Register(func(ctx context.Context, sig os.Signal) {
		log.Printf("1: %[1]d: %[1]s\n", sig)
	})

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
	goodbye.Register(a.close)
	err = a.run()
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
}
