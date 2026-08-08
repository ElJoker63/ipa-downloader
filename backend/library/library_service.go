package library

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ElJoker63/ipa-downloader/v2/backend/events"
	"github.com/ElJoker63/ipa-downloader/v2/backend/models"
	"github.com/ElJoker63/ipa-downloader/v2/backend/storage"
)

// LibraryService handles favorites, download history, and native file manager operations.
type LibraryService interface {
	GetFavorites() ([]models.FavoriteApp, error)
	SearchFavorites(query string) ([]models.FavoriteApp, error)
	AddFavorite(app models.FavoriteApp) error
	RemoveFavorite(appID int64) error
	ToggleFavorite(app models.FavoriteApp) (bool, error)

	GetHistory() ([]models.DownloadTask, error)
	DeleteHistoryItem(id string) error
	ClearHistory() error

	OpenFolder(path string) error
	OpenFile(path string) error
	RevealInExplorer(path string) error
	DeleteFile(path string) error
}


type libraryService struct {
	storage storage.Storage
	emitter events.Emitter
}

// NewLibraryService initializes the library and history manager.
func NewLibraryService(store storage.Storage, emitter events.Emitter) LibraryService {
	return &libraryService{
		storage: store,
		emitter: emitter,
	}
}

// ----------------- Favorites -----------------

func (s *libraryService) GetFavorites() ([]models.FavoriteApp, error) {
	return s.storage.GetFavorites()
}

func (s *libraryService) SearchFavorites(query string) ([]models.FavoriteApp, error) {
	if query == "" {
		return s.storage.GetFavorites()
	}
	return s.storage.SearchFavorites(query)
}

func (s *libraryService) AddFavorite(app models.FavoriteApp) error {
	if app.CreatedAt.IsZero() {
		app.CreatedAt = time.Now()
	}
	err := s.storage.AddFavorite(app)
	if err == nil {
		s.emitter.EmitLog("INFO", fmt.Sprintf("Added '%s' to favorites", app.Name), "LibraryService")
		s.emitter.Emit(events.EventFavoritesUpdated, app)
	}
	return err
}

func (s *libraryService) RemoveFavorite(appID int64) error {
	err := s.storage.RemoveFavorite(appID)
	if err == nil {
		s.emitter.EmitLog("INFO", fmt.Sprintf("Removed favorite app ID %d", appID), "LibraryService")
		s.emitter.Emit(events.EventFavoritesUpdated, map[string]int64{"removedAppId": appID})
	}
	return err
}

func (s *libraryService) ToggleFavorite(app models.FavoriteApp) (bool, error) {
	isFav, err := s.storage.IsFavorite(app.AppID)
	if err != nil {
		return false, err
	}

	if isFav {
		err = s.RemoveFavorite(app.AppID)
		return false, err
	}

	err = s.AddFavorite(app)
	return true, err
}

// ----------------- History -----------------

func (s *libraryService) GetHistory() ([]models.DownloadTask, error) {
	return s.storage.GetAllDownloads()
}

func (s *libraryService) DeleteHistoryItem(id string) error {
	task, _ := s.storage.GetDownload(id)
	err := s.storage.DeleteDownload(id)
	if err == nil && task != nil {
		s.emitter.EmitLog("INFO", fmt.Sprintf("Deleted download history record: %s", task.AppName), "LibraryService")
	}
	return err
}

func (s *libraryService) ClearHistory() error {
	err := s.storage.ClearCompletedDownloads()
	if err == nil {
		s.emitter.EmitLog("INFO", "Cleared completed download history", "LibraryService")
	}
	return err
}

// ----------------- Native File System -----------------

func (s *libraryService) OpenFolder(path string) error {
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, "Downloads")
	}

	fi, err := os.Stat(path)
	if err == nil && !fi.IsDir() {
		path = filepath.Dir(path)
	}

	return openPathNative(path)
}

func (s *libraryService) OpenFile(path string) error {
	return openPathNative(path)
}

func (s *libraryService) RevealInExplorer(path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("explorer.exe", "/select,", filepath.Clean(path))
		return cmd.Start()
	case "darwin":
		cmd := exec.Command("open", "-R", path)
		return cmd.Start()
	default:
		// Linux
		dir := filepath.Dir(path)
		cmd := exec.Command("xdg-open", dir)
		return cmd.Start()
	}
}

func (s *libraryService) DeleteFile(path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}

	err := os.Remove(path)
	if err == nil {
		s.emitter.EmitLog("INFO", fmt.Sprintf("File deleted from storage: %s", filepath.Base(path)), "LibraryService")
	}
	return err
}


func openPathNative(path string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "start", "", path)
		return cmd.Start()
	case "darwin":
		cmd := exec.Command("open", path)
		return cmd.Start()
	default:
		cmd := exec.Command("xdg-open", path)
		return cmd.Start()
	}
}
