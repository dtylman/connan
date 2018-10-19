package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dtylman/connan/indexer"
	"github.com/dtylman/gowd"
)

type app struct {
	body    *gowd.Element
	em      gowd.ElementsMap
	content *gowd.Element
	pages   map[string]string
}

var options struct {
	libFolder string
	tesseract string
}

func (a *app) loadPage(name string) error {
	html, err := ioutil.ReadFile(name + ".html")
	if err != nil {
		return err
	}
	a.pages[name] = string(html)
	return nil
}

func (a *app) loadPages(names ...string) error {
	for _, page := range names {
		err := a.loadPage(page)
		if err != nil {
			return err
		}
	}
	return nil
}
func (a *app) init() error {
	a.pages = make(map[string]string)
	err := a.loadPages("body", "search", "indexer", "backup")
	if err != nil {
		return err
	}

	a.em = gowd.NewElementMap()
	a.body, err = gowd.ParseElement(a.pages["body"], a.em)
	if err != nil {
		return err
	}

	a.content = a.em["content"]
	a.em["pageSearch"].OnEvent(gowd.OnClick, a.pageSearchClicked)
	a.em["pageIndexer"].OnEvent(gowd.OnClick, a.pageIndexerClicked)
	a.em["pageBackup"].OnEvent(gowd.OnClick, a.pageBackupClicked)

	a.pageSearchClicked(nil, nil)
	return nil
}

func (a *app) run() error {
	err := a.init()
	if err != nil {
		return err
	}
	//start the ui loop
	return gowd.Run(a.body)
}

func (a *app) pageIndexerClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.content.RemoveElements()
	a.content.AddHTML(a.pages["indexer"], a.em)
	a.em["inputConnanFolder"].SetValue(options.libFolder)
	a.em["inputTesseract"].SetValue(options.tesseract)
	a.em["indexerStart"].OnEvent(gowd.OnClick, a.indexerStartClicked)
}

func (a *app) pageBackupClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.content.RemoveElements()
	a.content.AddHTML(a.pages["backup"], a.em)
}

func (a *app) handleSearchRequest(sender *gowd.Element, event *gowd.EventElement) {
	txtSearch := a.em["txtSearch"]
	term := txtSearch.GetValue()
	txtSearch.AutoFocus()
	txtSearch.SetValue("")
	a.content.AddElement(gowd.NewText(term))

	// 	a.content.AddHTML(`<h2>0 Found</h2>
	// <a href="#">lala</a>
	// <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et
	// 	dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex
	// 	ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
	// 	fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt
	// 	mollit anim id est laborum.</p></br>
	// <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et
	// 	dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex
	// 	ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
	// 	fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt
	// 	mollit anim id est laborum.</p>
	// <div class="line"></div>
	// <h2>Lorem Ipsum Dolor</h2>
	// <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et
	// 	dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex
	// 	ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
	// 	fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt
	// 	mollit anim id est laborum.</p>
	// <div class="line"></div>
	// <h2>Lorem Ipsum Dolor</h2>
	// <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et
	// 	dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex
	// 	ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
	// 	fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt
	// 	mollit anim id est laborum.</p>
	// <div class="line"></div>
	// <h3>Lorem Ipsum Dolor</h3>
	// <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et
	// 	dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex
	// 	ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
	// 	fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt
	// 	mollit anim id est laborum.</p>
	// </div>`, a.em)
}

func (a *app) indexerStartClicked(sender *gowd.Element, event *gowd.EventElement) {
	options.libFolder = a.em["inputConnanFolder"].GetValue()
	options.tesseract = a.em["inputTesseract"].GetValue()
	err := indexer.Start()
	if err != nil {
		gowd.Alert(fmt.Sprintf("Error: %v", err))
	}
	// a.em["sidebar"].Hide()
	// a.em["sidebarCollapse"].Hide()
	// a.content.RemoveElements()
	// stopBtn := bootstrap.NewButton(bootstrap.ButtonPrimary, "Stop")
	// a.body.Disable()
	// a.content.AddElement(stopBtn)
	// go func() {
	// 	for i := 0; i <= 100; i++ {
	// 		time.Sleep(time.Second / 2)
	// 		js := fmt.Sprintf(`document.getElementById('progress').style.width="%v%%";
	// 	document.getElementById('progress').innerHTML = "%v%%";`, i, i)
	// 		gowd.ExecJSNow(js)
	// 	}
	// }()
}

func (a *app) pageSearchClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.content.RemoveElements()

	a.content.AddHTML(a.pages["search"], a.em)
	a.em["doSearch"].OnEvent(gowd.OnClick, a.handleSearchRequest)
}

func main() {
	logfile, err := os.Create("connan.log")
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println(err)
	} else {
		defer logfile.Close()
		log.SetOutput(logfile)
	}
	a := new(app)
	err = a.run()
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
}
