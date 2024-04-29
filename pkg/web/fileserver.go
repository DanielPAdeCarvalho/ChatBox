package web

import (
	"net/http"
)

func ServeStaticFiles() {
	staticFiles := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFiles))
	htmlFiles := http.FileServer(http.Dir("templates"))
	http.Handle("/", htmlFiles)
}
