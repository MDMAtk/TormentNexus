package main

import (
	"embed"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var startupURL string
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "tormentnexus://") {
			startupURL = arg
			break
		}
	}

	// Create an instance of the app structure
	app := NewApp(startupURL)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "TormentNexus Desktop",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 59, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
