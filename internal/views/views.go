// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package views

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/http"
)

type View struct {
	assetsFS  FS
	viewsFS   FS
	templates *template.Template
}

type FS struct {
	FS   embed.FS
	Path string
}

func New(assetsFS, viewsFS FS) *View {
	return &View{
		assetsFS:  assetsFS,
		viewsFS:   viewsFS,
		templates: template.Must(template.ParseFS(viewsFS.FS, "views/*.gohtml")),
	}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	log.Printf("%s %s: rendering template %q\n", r.Method, r.URL.Path, name)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// in development, we want to reload the templates on each request
	t, err := template.ParseFS(v.viewsFS.FS, "views/*.gohtml")
	if err != nil {
		log.Printf("%s %s: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("%s %s: parsed components\n", r.Method, r.URL.Path)

	// parse into a buffer so that we can handle errors without writing to the response
	buf := &bytes.Buffer{}
	if err := t.ExecuteTemplate(buf, name, data); err != nil {
		log.Printf("%s %s: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//if err := v.templates.ExecuteTemplate(w, name, data); err != nil {
	//	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//	log.Printf("template error: %v", err)
	//}
}
