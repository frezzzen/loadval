package main

import (
	"context"
	"embed"
	"go-wails-svelte5-template/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var appobj *app.App

var valorantApi *ValorantAPI

func main() {
	appobj = app.NewApp()
	valorantApi = NewValorantAPI()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "wails-template-svelte",
		Width:  1920,
		Height: 1080,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		OnStartup:        startup,
		Frameless:        true,
		Bind: []interface{}{
			appobj,
			valorantApi,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func startup(ctx context.Context) {
	app.Startup(appobj, ctx)
}
