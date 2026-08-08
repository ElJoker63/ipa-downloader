package main

import (
	"embed"
	"fmt"

	"github.com/majd/ipatool/v2/backend/app"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

// runDesktopApp initializes and starts the Wails desktop application.
func runDesktopApp() error {
	desktopApp, appService, err := app.NewApp()
	if err != nil {
		return fmt.Errorf("failed to initialize desktop application: %w", err)
	}

	err = wails.Run(&options.App{
		Title:             "IPATool Desktop",
		Width:             1180,
		Height:            780,
		MinWidth:          960,
		MinHeight:         640,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 11, G: 15, B: 25, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  desktopApp.Startup,
		OnDomReady: desktopApp.DomReady,
		OnShutdown: desktopApp.Shutdown,
		Bind: []interface{}{
			appService,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			BackdropType:         windows.Mica,
			DisableWindowIcon:    false,
		},
	})

	return err
}
