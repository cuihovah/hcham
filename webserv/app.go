package webserv

import (
	"io/ioutil"
	"net/http"
)

import (
	"html/template"
)

type Pages struct {
	P map[string]*template.Template
}

func (p *Pages) RenderPage(w http.ResponseWriter, name string, data interface{}) {
	t := p.P[name]
	t.Execute(w, data)
}

func (p *Pages) SetPage(name string, filepath string) error {
	t, err := ioutil.ReadFile(filepath)
	tmpl := template.New(name)
	tmpl, err = tmpl.Parse(string(t))
	if err != nil {
		return err
	}
	p.P[name] = tmpl
	return nil
}

func NewPages() *Pages {
	pages := &Pages{}
	pages.P = make(map[string]*template.Template)
	return pages
}
