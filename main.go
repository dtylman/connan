package main

import (
	"io/ioutil"
	"strings"

	"github.com/dtylman/gowd"
)

var body *gowd.Element

func mainloop() error {
	//creates a new bootstrap fluid container
	em := gowd.NewElementMap()
	innerHTML, err := ioutil.ReadFile("body.html")
	if err != nil {
		return err
	}
	body, err := gowd.ParseElement(strings.Replace(string(innerHTML), "\n", "", -1), em)
	if err != nil {
		return err
	}
	//start the ui loop
	return gowd.Run(body)
}

func main() {
	err := mainloop()
	if err != nil {
		panic(err)
	}
}
