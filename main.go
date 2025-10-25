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
var storageApi *StorageAPI
var updateApi *UpdateAPI

func main() {
	println("Starting LOADVAL application...")

	if err := InitLogger(); err != nil {
		println("Warning: Failed to initialize logger:", err.Error())
	}
	defer CloseLogger()

	LogInfo("=== LOADVAL Application Starting ===")
	LogInfo("Initializing application components...")

	appobj = app.NewApp()
	valorantApi = NewValorantAPI()
	storageApi = NewStorageAPI()
	updateApi = NewUpdateAPI()

	LogInfo("App objects created successfully")
	LogInfo("Log file location: " + GetLogFilePath())

	err := wails.Run(&options.App{
		Title:     "LOADVAL",
		Width:     1280,
		Height:    800,
		MinWidth:  800,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		OnStartup:        startup,
		OnShutdown:       shutdown,
		Frameless:        true,
		Bind: []interface{}{
			appobj,
			valorantApi,
			storageApi,
			updateApi,
		},
	})

	if err != nil {
		LogError("Application error: " + err.Error())
		println("Error:", err.Error())
	} else {
		LogInfo("Application closed successfully")
		println("Application closed successfully")
	}
}

func startup(ctx context.Context) {
	app.Startup(appobj, ctx)
	LogInfo("Application startup complete")
}

func shutdown(ctx context.Context) {
	LogInfo("Application shutting down...")
}
