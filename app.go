package main

import (
	"context"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
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
	for _, analyzer := range a.config.Analyzers {
		a.db.AddDocumentAnalyzer(analyzer)
	}
	a.db.AddDocumentAnalyzer(PlainTextAnalyzer{})
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

func (a *App) close(ctx context.Context, s os.Signal) {
	err := a.config.Save()
	if err != nil {
		log.Println(err)
	}
	a.indexer.Stop()
	a.indexer.Close()
}

func (a *App) run() error {
	defer a.close(nil, nil)
	//start the ui loop
	return gowd.Run(a.ui.body)
}

func (a *App) pageIndexerClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.ui.em["input-connan-folder"].SetValue(a.config.LibFolder)
	a.ui.content.SetElement(a.ui.pageIndexer)
}

func (a *App) pageBackupClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.ui.content.SetElement(a.ui.pageBackup)
}

func (a *App) scrollToLastPosition() {
	input := a.ui.em["scroll-value"]
	if input != nil {
		val := input.GetValue()
		if val != "" {
			gowd.ExecJS(fmt.Sprintf(`$(window.scrollTo(0,%v))`, val))
		}
	}
}

func (a *App) buttonSearchMoreClicked(sender *gowd.Element, event *gowd.EventElement) {
	if !a.hasMoreResults() {
		return
	}
	req := a.results.Request
	req.From += req.Size
	var err error
	a.results, err = a.indexer.db.Bleve.Search(req)
	if err != nil {
		gowd.Alert(fmt.Sprintf("%v", err))
		return
	}
	a.renderSearchResults(a.ui.em["div-search-results"])
}

//hasMoreResults is true if we can query the same search again with new results
func (a *App) hasMoreResults() bool {
	if a.results == nil {
		return false
	}
	return uint64((a.results.Request.From + a.results.Hits.Len())) <= a.results.Total
}

func (a *App) renderSearchResults(div *gowd.Element) {
	if a.hasMoreResults() {
		btnsearchmore := a.ui.em["button-search-more"]
		btnsearchmore.OnEvent(gowd.OnClick, a.buttonSearchMoreClicked)
		gowd.ExecJS("attach_scroll_event('button-search-more')")
	} else {
		gowd.ExecJS(`$(window).unbind("scroll");`)
	}
	for _, hit := range a.results.Hits {
		doc := a.db.Document(hit.ID)
		card := NewDocumentCard(doc, hit)
		div.AddElement(card.Element)
		card.linkContent.Object = doc
		card.linkContent.OnEvent(gowd.OnClick, a.documentContentClicked)
		if card.linkImage != nil {
			card.linkImage.Object = doc
			card.linkImage.OnEvent(gowd.OnClick, a.documentImageClicked)
		}
	}
	a.scrollToLastPosition()
}

func (a App) documentImageClicked(sender *gowd.Element, event *gowd.EventElement) {
	doc := sender.Object.(*db.Document)
	body := bootstrap.NewElement("div", "")

	img := body.AddElement(bootstrap.NewElement("img", "img-fluid"))
	img.SetAttribute("src", fmt.Sprintf("file://%s", doc.Path))
	a.showModal(doc.Name(), body)
}

func (a App) documentContentClicked(sender *gowd.Element, event *gowd.EventElement) {
	doc := sender.Object.(*db.Document)
	body := bootstrap.NewElement("div", "")

	body.AddElement(gowd.NewText(doc.GetField("content")))
	a.showModal(doc.Name(), body)
}

func (a *App) showModal(title string, body *gowd.Element) {
	a.scrollToLastPosition()
	a.ui.em["modal-body"].SetElement(body)
	a.ui.em["modal-title"].SetText(title)
	gowd.ExecJS(`$('#modal-default').modal({keyboard:true,focus:true});`)
}

func (a *App) buttonSearchGoClicked(sender *gowd.Element, event *gowd.EventElement) {
	input := a.ui.em["input-search"]
	term := input.GetValue()
	input.AutoFocus()
	input.SetValue("")
	a.ui.em["scroll-value"].SetValue("")
	var q query.Query
	if term == "" || term == "*" {
		q = bleve.NewMatchAllQuery()
	} else {
		q = bleve.NewQueryStringQuery(term)
	}
	req := bleve.NewSearchRequest(q)
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
	a.ui.em["div-search-summary"].SetElement(gowd.NewStyledText(summary, gowd.Heading4))

	a.renderSearchResults(divresults)
}

func (a *App) pageSearchClicked(sender *gowd.Element, event *gowd.EventElement) {
	a.ui.content.SetElement(a.ui.pageSearch)
}
