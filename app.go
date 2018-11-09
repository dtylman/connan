package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/dtylman/connan/db"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
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
	ui      UI
	config  Options
	indexer *Indexer
	db      *db.DB
	results *bleve.SearchResult
}

//NewApp creates a new application
func NewApp() (*App, error) {
	a := new(App)
	err := a.config.Load()
	if err != nil {
		return nil, err
	}
	a.db, err = db.Open(a.config.DBFolder)
	if err != nil {
		return nil, err
	}
	a.indexer, err = NewIndexer(a.db)
	if err != nil {
		return nil, err
	}
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
	a.ui.em["button-indexer-stop"].OnEvent(gowd.OnClick, a.buttonIndexerStopClicked)
	a.ui.em["button-indexer-settings-save"].OnEvent(gowd.OnClick, a.buttonIndexerSaveClicked)
	a.ui.pageSearch, err = a.loadPage("frontend/search.html")
	if err != nil {
		return nil, err
	}
	a.ui.em["button-search-go"].OnEvent(gowd.OnClick, a.buttonSearchGoClicked)

	a.ui.content = a.ui.em["content"]
	a.pageSearchClicked(nil, nil)

	go a.progressUpdate()
	return a, nil
}

//progressUpdate updates progress bar when they are displayed
func (a *App) progressUpdate() {
	var value, total int
	var label string
	for true {
		time.Sleep(time.Second / 2)
		ip := a.ui.content.Find("progress-label-indexer")
		if ip != nil {
			if a.indexer.worker.IsRunning() {
				value = a.indexer.queued - a.indexer.db.Queue.Len()
				total = a.indexer.queued
				label = AppLog.Last
				gowd.ExecJSNow(fmt.Sprintf("set_progress('progress-bar-indexer',%v,%v,'progress-label-indexer','%v');",
					value, total, html.EscapeString(label)))
			}
		}
	}
}

func (a *App) buttonIndexerStopClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.indexer.Stop()
}

func (a *App) buttonIndexerStartClicked(sender *gowd.Element, event *gowd.EventElement) {
	err := a.indexer.Start(a.config.LibFolder)
	if err != nil {
		gowd.Alert(fmt.Sprintf("%v", err))
		return
	}
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
	a.indexer.Stop()
	a.indexer.Close()
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

	req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(term))
	req.Highlight = bleve.NewHighlightWithStyle("html")
	var err error
	a.results, err = a.indexer.db.Bleve.Search(req)
	if err != nil {
		gowd.Alert(fmt.Sprintf("%v", err))
		return
	}

	divresults := a.ui.em["div-search-results"]
	divresults.RemoveElements()
	summary := fmt.Sprintf("%d matches, showing %d through %d, took %s", a.results.Total, a.results.Request.From+1, a.results.Request.From+len(a.results.Hits), a.results.Took)
	divresults.AddElement(gowd.NewStyledText(summary, gowd.Heading4))
	for _, hit := range a.results.Hits {
		doc := a.db.Document(hit.ID)
		link := bootstrap.NewLinkButton(filepath.Base(hit.ID))
		link.SetAttribute("onclick", fmt.Sprintf("nw.Shell.openItem('%v');", doc.Path))
		header := bootstrap.NewElement("h4", "heading-small mb-4")
		header.AddElement(link)
		header.AddHTML("&nbsp;", nil)
		header.AddElement(gowd.NewStyledText(fmt.Sprintf("(%f)", hit.Score), gowd.ItalicText))
		dr := bootstrap.NewElement("div", "card-body", header)
		img := gowd.NewElement("img")
		img.SetClass("rounded pull-right img-thumbnail")
		img.SetAttribute("style", "height: 100px; width: 100px")
		img.SetAttribute("src", "file:///"+doc.Path)
		dr.AddElement(img)

		//link.OnEvent(gowd.OnClick, a.buttonSearchClicked)
		for fragmentField, fragments := range hit.Fragments {
			dr.AddElement(gowd.NewStyledText(fragmentField, gowd.StrongText))
			for _, fragment := range fragments {
				dr.AddHTML(fragment, nil)
			}
		}
		for otherFieldName, otherFieldValue := range hit.Fields {
			if _, ok := hit.Fragments[otherFieldName]; !ok {
				dr.AddElement(gowd.NewStyledText(otherFieldName, gowd.StrongText))
				dr.AddElement(gowd.NewText(fmt.Sprintf("%v", otherFieldValue)))
			}
		}
		divresults.AddElement(dr)
		divresults.AddElement(gowd.NewElement("hr"))

	}
}

func (a *App) pageSearchClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.ui.content.SetElement(a.ui.pageSearch)
}
