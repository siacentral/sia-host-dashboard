package web

import (
	"embed"
	"io/fs"
	"net/http"
)

var (
	//go:embed dist
	dist embed.FS
	// Assets http.FileSystem from the embedded web assets
	Assets = func() http.FileSystem {
		dist, err := fs.Sub(dist, "dist")
		if err != nil {
			panic(err)
		}
		return http.FS(dist)
	}()
)
