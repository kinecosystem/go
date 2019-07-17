package hal

import (
	"strings"
)

type Link struct {
	Href      string `json:"href,omitempty"`
	Templated bool   `json:"templated,omitempty"`
}

func (l *Link) PopulateTemplated() {
	l.Templated = strings.Contains(l.Href, "{")
}

func NewLink(href string) Link {
	l := Link{Href: href}
	l.PopulateTemplated()
	return l
}

func NewLinkPtr(href string) *Link {
	l := NewLink(href)
	return &l
}
