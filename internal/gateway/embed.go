package gateway

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var frontendFS embed.FS

// GetFrontendFS returns the embedded frontend filesystem
func GetFrontendFS() http.FileSystem {
	subFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		panic(err)
	}
	return http.FS(subFS)
}
