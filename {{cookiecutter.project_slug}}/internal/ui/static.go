// Package ui embeds the built React frontend (web/dist copied here by
// `mise run build` / the Dockerfile) so the final Go binary serves the whole app.
package ui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

// Assets returns the embedded frontend build rooted at dist/.
// A placeholder dist/index.html is committed so `go build` and `go test`
// work before the frontend has ever been built.
func Assets() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic("ui: embedded dist missing: " + err.Error())
	}
	return sub
}
