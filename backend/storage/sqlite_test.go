package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ElJoker63/ipatool-1/v2/backend/models"
)

func TestSQLiteStorage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ipa-downloader-sqlite-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := NewSQLiteStorage(tempDir)
	if err != nil {
		t.Fatalf("failed to init sqlite storage: %v", err)
	}
	defer store.Close()

	// 1. Test Settings
	settings, err := store.GetSettings()
	if err != nil {
		t.Fatalf("GetSettings failed: %v", err)
	}
	if settings.Theme != "system" {
		t.Errorf("expected default theme 'system', got %s", settings.Theme)
	}

	settings.Theme = "dark"
	settings.MaxConcurrentDownloads = 5
	if err := store.SaveSettings(*settings); err != nil {
		t.Fatalf("SaveSettings failed: %v", err)
	}

	updatedSettings, err := store.GetSettings()
	if err != nil || updatedSettings.Theme != "dark" || updatedSettings.MaxConcurrentDownloads != 5 {
		t.Fatalf("unexpected settings after save: %+v", updatedSettings)
	}

	// 2. Test Downloads
	now := time.Now()
	task := models.DownloadTask{
		ID:              "test-task-1",
		AppID:           324684580,
		BundleID:        "com.spotify.client",
		AppName:         "Spotify",
		Developer:       "Spotify AB",
		Version:         "8.9.0",
		ArtworkURL:      "https://example.com/icon.png",
		DestinationPath: filepath.Join(tempDir, "Spotify.ipa"),
		Status:          models.DownloadStatusDownloading,
		TotalBytes:      1000000,
		DownloadedBytes: 250000,
		Progress:        25.0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := store.SaveDownload(task); err != nil {
		t.Fatalf("SaveDownload failed: %v", err)
	}

	saved, err := store.GetDownload("test-task-1")
	if err != nil || saved.AppName != "Spotify" {
		t.Fatalf("GetDownload failed or mismatch: %+v, err: %v", saved, err)
	}

	// Update progress
	if err := store.UpdateDownloadProgress("test-task-1", models.DownloadStatusCompleted, 1000000, 1000000, 100.0, 50000, 0, ""); err != nil {
		t.Fatalf("UpdateDownloadProgress failed: %v", err)
	}

	completed, err := store.GetDownload("test-task-1")
	if err != nil || completed.Status != models.DownloadStatusCompleted {
		t.Fatalf("expected completed status, got: %+v", completed)
	}

	// 3. Test Favorites
	fav := models.FavoriteApp{
		AppID:          324684580,
		BundleID:       "com.spotify.client",
		Name:           "Spotify",
		Developer:      "Spotify AB",
		Version:        "8.9.0",
		Price:          0,
		FormattedPrice: "Free",
		ArtworkURL:     "https://example.com/icon.png",
		PrimaryGenre:   "Music",
		CreatedAt:      now,
	}

	if err := store.AddFavorite(fav); err != nil {
		t.Fatalf("AddFavorite failed: %v", err)
	}

	isFav, err := store.IsFavorite(324684580)
	if err != nil || !isFav {
		t.Fatalf("expected isFavorite=true, got %v (err: %v)", isFav, err)
	}

	favList, err := store.GetFavorites()
	if err != nil || len(favList) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(favList))
	}

	if err := store.RemoveFavorite(324684580); err != nil {
		t.Fatalf("RemoveFavorite failed: %v", err)
	}
	isFavAfter, _ := store.IsFavorite(324684580)
	if isFavAfter {
		t.Fatalf("expected isFavorite=false after removal")
	}

	// 4. Test Search History
	if err := store.AddSearchHistory("spotify", "ios", 5); err != nil {
		t.Fatalf("AddSearchHistory failed: %v", err)
	}
	hist, err := store.GetSearchHistory(10)
	if err != nil || len(hist) != 1 || hist[0].Term != "spotify" {
		t.Fatalf("unexpected search history: %+v", hist)
	}

	// 5. Test Logs
	entry, err := store.AddLog("INFO", "Test log message", "TestContext")
	if err != nil || entry == nil {
		t.Fatalf("AddLog failed: %v", err)
	}
	logs, err := store.GetRecentLogs(10)
	if err != nil || len(logs) != 1 || logs[0].Message != "Test log message" {
		t.Fatalf("unexpected logs: %+v", logs)
	}
}
