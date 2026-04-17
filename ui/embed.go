package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist
var files embed.FS

func Handler() http.Handler {
	stripped, err := fs.Sub(files, "dist")
	if err != nil {
		panic("failed to sub ui dist: " + err.Error())
	}

	fileServer := http.FileServer(http.FS(stripped))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// use the original embed.FS directly to check existence
		// since it supports Stat reliably
		path := "dist" + r.URL.Path
		_, err := files.Open(path)
		if err != nil {
			// file not found — serve index.html for React router
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}