// Package ui provides embedded static files for the web UI.
package ui

import (
	"embed"
)

//go:embed "html" "static"
var Files embed.FS
