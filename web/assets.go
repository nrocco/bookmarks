package web

import (
	"embed"
)

//go:embed dist/*
// Assets embeds the frontend inside the bend binary
var Assets embed.FS
