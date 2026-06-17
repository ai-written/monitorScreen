package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:      "Monitor Screen",
		Width:      1280,
		Height:     820,
		MinWidth:   960,
		MinHeight:  600,
		Frameless:  true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:    app.startup,
		OnShutdown:   app.shutdown,
		OnBeforeClose: app.beforeClose,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			BackdropType:         windows.Auto,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "monitor-screen-4f3b2c1a",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				if app.ctx != nil {
					runtime.WindowShow(app.ctx)
					runtime.WindowUnminimise(app.ctx)
				}
			},
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
