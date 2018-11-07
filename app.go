package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dtylman/gowd"
)

//UI holds the application UI
type UI struct {
	body        *gowd.Element
	em          gowd.ElementsMap
	content     *gowd.Element
	pageBackup  *gowd.Element
	pageSearch  *gowd.Element
	pageIndexer *gowd.Element
}

//App is the application
type App struct {
	ui     UI
	config Options
}

//NewApp creates a new application
func NewApp() (*App, error) {
	a := new(App)
	err := a.config.Load()
	if err != nil {
		return nil, err
	}
	// a.indexer, err = db.NewIndexer(filepath.Join(a.config.LibFolder, "connandb"))
	// if err != nil {
	// 	return nil, err
	// }
	a.ui.em = gowd.NewElementMap()

	a.ui.body, err = a.loadPage("frontend/body.html")
	if err != nil {
		return nil, err
	}
	a.ui.em["button-menu-search"].OnEvent(gowd.OnClick, a.pageSearchClicked)
	a.ui.em["button-menu-indexer"].OnEvent(gowd.OnClick, a.pageIndexerClicked)
	a.ui.em["button-menu-backup"].OnEvent(gowd.OnClick, a.pageBackupClicked)

	a.ui.pageBackup, err = a.loadPage("frontend/backup.html")
	if err != nil {
		return nil, err
	}
	a.ui.pageIndexer, err = a.loadPage("frontend/indexer.html")
	if err != nil {
		return nil, err
	}
	a.ui.em["button-indexer-start"].OnEvent(gowd.OnClick, a.buttonIndexerStartClicked)
	a.ui.em["button-indexer-settings-save"].OnEvent(gowd.OnClick, a.buttonIndexerSaveClicked)
	a.ui.pageSearch, err = a.loadPage("frontend/search.html")
	if err != nil {
		return nil, err
	}
	a.ui.em["button-search-go"].OnEvent(gowd.OnClick, a.buttonSearchGoClicked)

	a.ui.content = a.ui.em["content"]
	a.pageSearchClicked(nil, nil)

	return a, nil
}

func (a *App) buttonIndexerStartClicked(sender *gowd.Element, event *gowd.EventElement) {

}

func (a *App) buttonIndexerSaveClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.config.LibFolder = a.ui.em["input-connan-folder"].GetValue()
	a.config.Tesseract = a.ui.em["input-tesseract"].GetValue()
	err := a.config.Save()
	if err != nil {
		gowd.Alert(fmt.Sprintf("Failed to save settings: %v", err))
	} else {
		gowd.Alert("Saved")
	}
}

func (a *App) loadPage(fileName string) (*gowd.Element, error) {
	html, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Loading %v: %v", fileName, err)
	}
	e, err := gowd.ParseElement(string(html), a.ui.em)
	if err != nil {
		return nil, fmt.Errorf("Loading %v: %v", fileName, err)
	}
	return e, nil
}

func (a *App) close() {
	err := a.config.Save()
	if err != nil {
		log.Println(err)
	}
	// err = a.indexer.Close()
	// if err != nil {
	// 	log.Println(err)
	// }
}

func (a *App) run() error {
	defer a.close()

	//start the ui loop
	return gowd.Run(a.ui.body)
}

func (a *App) pageIndexerClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.ui.em["input-connan-folder"].SetValue(a.config.LibFolder)
	a.ui.em["input-tesseract"].SetValue(a.config.Tesseract)
	a.ui.content.SetElement(a.ui.pageIndexer)
}

func (a *App) pageBackupClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.ui.content.SetElement(a.ui.pageBackup)
}

func (a *App) buttonSearchGoClicked(sender *gowd.Element, event *gowd.EventElement) {
	input := a.ui.em["input-search"]
	term := input.GetValue()
	input.AutoFocus()
	input.SetValue("")
	a.ui.content.AddElement(gowd.NewText(term))

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

func (a *App) indexerStartClicked(sender *gowd.Element, event *gowd.EventElement) {
	// options.LibFolder = a.em["inputConnanFolder"].GetValue()
	// options.Tesseract = a.em["inputTesseract"].GetValue()
	// err := indexer.Start()
	// if err != nil {
	// 	gowd.Alert(fmt.Sprintf("Error: %v", err))
	// }
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

func (a *App) pageSearchClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.ui.content.SetElement(a.ui.pageSearch)
}
