package web

import (
	"embed"
)

//go:generate mkdir -p dist
//go:generate touch dist/dummy
//go:embed dist/*
// Assets embeds the frontend inside the bend binary
var Assets embed.FS
