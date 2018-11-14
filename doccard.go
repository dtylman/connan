package main

import (
	"fmt"
	"path/filepath"

	"github.com/blevesearch/bleve/search"
	"github.com/dtylman/connan/db"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
)

//DocumentCard ...
type DocumentCard struct {
	Element *gowd.Element
	doc     *db.Document
	hit     *search.DocumentMatch
}

//NewDocumentCard creates a new card
func NewDocumentCard(doc *db.Document, hit *search.DocumentMatch) *DocumentCard {
	card := new(DocumentCard)
	card.doc = doc
	card.hit = hit
	card.Element = bootstrap.NewElement("div", "card m-3 col-md-3")
	header := card.Element.AddElement(bootstrap.NewElement("div", "card-header"))
	card.populateHeader(header)
	subtitle := card.Element.AddElement(bootstrap.NewElement("small", "text-center"))
	subtitle.AddElement(gowd.NewText(doc.GetField("type")))
	cardbody := card.Element.AddElement(bootstrap.NewElement("div", "card-body"))
	card.populateBody(cardbody)
	return card
}

func (dc *DocumentCard) populateHeader(header *gowd.Element) {
	linkTitle := header.AddElement(bootstrap.NewLinkButton(filepath.Base(dc.doc.Path)))
	//link.OnEvent(gowd.OnClick)
	linkTitle.SetAttribute("Title", "Open")

	linkOpen := header.AddElement(bootstrap.NewLinkButton(""))
	linkOpen.SetClass("mx-2")
	linkOpen.AddElement(bootstrap.NewElement("i", "fas fa-link"))
	linkOpen.SetAttribute("Title", "View")

	span := header.AddElement(bootstrap.NewElement("span", "badge badge-pill badge-info float-right"))
	span.AddElement(gowd.NewText(fmt.Sprintf("%.2f", dc.hit.Score)))
}

func (dc *DocumentCard) populateBody(cardbody *gowd.Element) {
	if dc.doc.IsImage() {
		linkImage := cardbody.AddElement(bootstrap.NewLinkButton(""))
		linkImage.SetClass("float-right")
		linkImage.SetAttribute("Title", "View Image")

		image := linkImage.AddElement(bootstrap.NewElement("img", "img-thumbnail thumb-search"))
		image.SetAttribute("src", fmt.Sprintf("file:///%s", dc.doc.Path))
	}

	hit := cardbody.AddElement(bootstrap.NewLinkButton(""))
	hit.SetAttribute("Title", "View Content")
	p := hit.AddElement(bootstrap.NewElement("p", "card-text"))

	for fragmentField, fragments := range dc.hit.Fragments {
		p.AddHTML(fragmentField+": ", nil)
		for _, fragment := range fragments {
			p.AddHTML(fragment, nil)
		}
	}
}

/*

        <div class="card m-3 col-md-3">
            <div class="card-header">
                <a href="#" title="poop">baba.doc</a>
                <a href="#" class="mx-2"><i class="fas fa-link" title="lala"></i></a>
                <span class="badge badge-pill badge-info float-right">2.34</span>
            </div>
            <small class="text-center">open document version 1.0</small>
            <div class="card-body">
                <a href="#" class="float-right" title="googoo!">
                    <img class="img-thumbnail thumb-search" src="https://thenypost.files.wordpress.com/2017/03/shutterstock.jpg?quality=90&strip=all&w=618&h=410&crop=1"
                        alt="text" /></a>
                <a href="#" title="boo goo!">
                    <p class="card-text">Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum
                        Lorem ipsum Lorem ipsum</p>
                </a>
            </div>
		</div>*/
