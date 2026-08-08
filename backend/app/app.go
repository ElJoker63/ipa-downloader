package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/majd/ipa-downloader/v2/backend/services"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App manages the desktop application lifecycle and window events.
type App struct {
	ctx     context.Context
	service *services.AppService
}

// NewApp creates a new App application struct.
func NewApp() (*App, *services.AppService, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dataDir := filepath.Join(home, ".ipa-downloader")

	appService, err := services.NewAppService(dataDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create app service: %w", err)
	}

	app := &App{
		service: appService,
	}

	return app, appService, nil
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.service.SetContext(ctx)
	_, _ = a.service.AddLog("INFO", "IPA Downloader Desktop started successfully", "AppLifecycle")
}

// DomReady is called after the front-end has finished loading.
func (a *App) DomReady(ctx context.Context) {
	_, _ = a.service.AddLog("DEBUG", "Frontend DOM is ready", "AppLifecycle")
}

// BeforeClose is called before the application terminates.
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// Shutdown is called at application termination.
func (a *App) Shutdown(ctx context.Context) {
	a.service.Shutdown()
}

// ShowNativeNotification displays a system desktop notification.
func (a *App) ShowNativeNotification(title, message string) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "notification:show", map[string]string{
			"title":   title,
			"message": message,
		})
	}
}
