package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	fileName string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles((filepath.Join("templates", t.fileName))))
	})
	t.templ.Execute(w, r)
}
func main() {

	var addr = flag.String("a", ":8080", "run address")

	flag.Parse()

	r := newRoom()

	http.Handle("/", &templateHandler{fileName: "chat.html"})
	http.Handle("/room", r)

	go r.run()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
