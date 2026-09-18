// Package webdist 內嵌前端建置產物（由 make web 複製到 dist/）。
package webdist

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var files embed.FS

// FS 為 dist 目錄內容；前端尚未建置時只有 .gitkeep。
var FS, _ = fs.Sub(files, "dist")
