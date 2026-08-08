package config

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/majd/ipatool/v2/backend/events"
	"github.com/majd/ipatool/v2/backend/models"
	"github.com/majd/ipatool/v2/backend/storage"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ConfigService manages settings, preferences, directories, and cache cleanup.
type ConfigService interface {
	SetContext(ctx context.Context)
	GetSettings() (*models.AppSettings, error)
	SaveSettings(settings models.AppSettings) error
	SelectDownloadDirectory(defaultPath string) (string, error)
	ClearAppCache() error
	GetCacheSize() (string, error)
	ExportLogs(destinationPath string) (string, error)
}

type configService struct {
	storage storage.Storage
	emitter events.Emitter
	ctx     context.Context
	mu      sync.RWMutex
}

// NewConfigService creates a new configuration service.
func NewConfigService(store storage.Storage, emitter events.Emitter) ConfigService {
	return &configService{
		storage: store,
		emitter: emitter,
	}
}

func (s *configService) SetContext(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ctx = ctx
}

func (s *configService) GetSettings() (*models.AppSettings, error) {
	return s.storage.GetSettings()
}

func (s *configService) SaveSettings(settings models.AppSettings) error {
	err := s.storage.SaveSettings(settings)
	if err == nil {
		s.emitter.EmitLog("INFO", "Updated application settings", "ConfigService")
	}
	return err
}

func (s *configService) SelectDownloadDirectory(defaultPath string) (string, error) {
	s.mu.RLock()
	ctx := s.ctx
	s.mu.RUnlock()

	if ctx == nil {
		return "", fmt.Errorf("runtime context not ready")
	}

	selected, err := runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{
		DefaultDirectory: defaultPath,
		Title:            "Select Default IPA Download Folder",
	})
	if err != nil {
		return "", err
	}

	if selected != "" {
		settings, err := s.storage.GetSettings()
		if err == nil && settings != nil {
			settings.DefaultDownloadFolder = selected
			_ = s.storage.SaveSettings(*settings)
		}
	}

	return selected, nil
}

func (s *configService) ClearAppCache() error {
	err := s.storage.ClearAppCache()
	if err == nil {
		s.emitter.EmitLog("INFO", "Cleared offline app metadata cache", "ConfigService")
		s.emitter.EmitNotification("Cache Cleared", "Offline application metadata cache has been emptied.", "info")
	}
	return err
}

func (s *configService) GetCacheSize() (string, error) {
	bytes, err := s.storage.GetCacheSizeBytes()
	if err != nil {
		return "0 KB", err
	}
	return formatBytes(bytes), nil
}

func (s *configService) ExportLogs(destinationPath string) (string, error) {
	logs, err := s.storage.GetRecentLogs(5000)
	if err != nil {
		return "", err
	}

	s.mu.RLock()
	ctx := s.ctx
	s.mu.RUnlock()

	if destinationPath == "" && ctx != nil {
		selected, err := runtime.SaveFileDialog(ctx, runtime.SaveDialogOptions{
			DefaultFilename: "ipatool-diagnostics.log",
			Title:           "Export Application Logs",
			Filters: []runtime.FileFilter{
				{DisplayName: "Log Files (*.log)", Pattern: "*.log"},
				{DisplayName: "Text Files (*.txt)", Pattern: "*.txt"},
			},
		})
		if err != nil || selected == "" {
			return "", err
		}
		destinationPath = selected
	}

	file, err := os.Create(destinationPath)
	if err != nil {
		return "", fmt.Errorf("failed to create log file: %w", err)
	}
	defer file.Close()

	for _, entry := range logs {
		line := fmt.Sprintf("[%s] [%s] [%s] %s\n", entry.Timestamp.Format("2006-01-02 15:04:05.000"), entry.Level, entry.Context, entry.Message)
		_, _ = file.WriteString(line)
	}

	s.emitter.EmitLog("SUCCESS", fmt.Sprintf("Exported %d log entries to %s", len(logs), destinationPath), "ConfigService")
	s.emitter.EmitNotification("Logs Exported", fmt.Sprintf("Successfully exported logs to %s", destinationPath), "success")
	return destinationPath, nil
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
