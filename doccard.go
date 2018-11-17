package main

import (
	"fmt"

	"github.com/blevesearch/bleve/search"
	"github.com/dtylman/connan/db"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
)

//DocumentCard ...
type DocumentCard struct {
	Element      *gowd.Element
	doc          *db.Document
	hit          *search.DocumentMatch
	linkImage    *gowd.Element
	linkContent  *gowd.Element
	linkDocument *gowd.Element
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
	linkTitle := header.AddElement(bootstrap.NewLinkButton(dc.doc.Name()))
	//link.OnEvent(gowd.OnClick)
	linkTitle.SetAttribute("Title", "Open")
	linkTitle.SetAttribute("onclick", fmt.Sprintf(`nw.Shell.openItem('%v');`, dc.doc.Path))

	linkFolder := header.AddElement(bootstrap.NewLinkButton(""))
	linkFolder.SetClass("mx-2")
	linkFolder.AddElement(bootstrap.NewElement("i", "fas fa-folder"))
	linkFolder.SetAttribute("Title", "Show In Folder")
	linkFolder.SetAttribute("onclick", fmt.Sprintf(`nw.Shell.showItemInFolder('%v');`, dc.doc.Path))

	dc.linkContent = header.AddElement(bootstrap.NewLinkButton(""))
	dc.linkContent.SetClass("mx-2")
	dc.linkContent.AddElement(bootstrap.NewElement("i", "fas fa-eye"))
	dc.linkContent.SetAttribute("Title", "View")

	span := header.AddElement(bootstrap.NewElement("span", "badge badge-pill badge-info float-right"))
	span.AddElement(gowd.NewText(fmt.Sprintf("%.2f", dc.hit.Score)))
}

func (dc *DocumentCard) populateBody(cardbody *gowd.Element) {
	if dc.doc.IsImage() {
		dc.linkImage = cardbody.AddElement(bootstrap.NewLinkButton(""))
		dc.linkImage.SetClass("float-right")
		dc.linkImage.SetAttribute("Title", "View Image")

		image := dc.linkImage.AddElement(bootstrap.NewElement("img", "img-thumbnail thumb-search"))
		image.SetAttribute("src", fmt.Sprintf("file:///%s", dc.doc.Path))
	}

	dc.linkDocument = cardbody.AddElement(bootstrap.NewLinkButton(""))
	dc.linkDocument.SetAttribute("Title", "View Content")
	for fragmentField, fragments := range dc.hit.Fragments {
		if fragmentField == "path" || fragmentField == "fields.type" || fragmentField == "fields.mime" {
			continue
		}
		p := dc.linkDocument.AddElement(bootstrap.NewElement("p", "card-text small"))
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
