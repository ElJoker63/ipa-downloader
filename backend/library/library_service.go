package library

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/ElJoker63/ipa-downloader/v2/backend/events"
	"github.com/ElJoker63/ipa-downloader/v2/backend/models"
	"github.com/ElJoker63/ipa-downloader/v2/backend/storage"
	"howett.net/plist"
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

	GetDownloadedIPAs(downloadDir string) ([]models.DownloadedIPA, error)

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

// ----------------- Downloaded IPAs -----------------

func (s *libraryService) GetDownloadedIPAs(downloadDir string) ([]models.DownloadedIPA, error) {
	if downloadDir == "" {
		settings, err := s.storage.GetSettings()
		if err == nil && settings != nil && settings.DefaultDownloadFolder != "" {
			downloadDir = settings.DefaultDownloadFolder
		} else {
			home, _ := os.UserHomeDir()
			downloadDir = filepath.Join(home, "Downloads")
		}
	}

	entries, err := os.ReadDir(downloadDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read download folder: %w", err)
	}

	// Fetch history & favorites for artwork lookup cache
	downloads, _ := s.storage.GetAllDownloads()
	artworkMap := make(map[string]string)
	appIDMap := make(map[string]int64)

	for _, d := range downloads {
		if d.BundleID != "" {
			if d.ArtworkURL != "" {
				artworkMap[d.BundleID] = d.ArtworkURL
			}
			if d.AppID != 0 {
				appIDMap[d.BundleID] = d.AppID
			}
		}
	}

	favorites, _ := s.storage.GetFavorites()
	for _, f := range favorites {
		if f.BundleID != "" {
			if f.ArtworkURL != "" {
				artworkMap[f.BundleID] = f.ArtworkURL
			}
			if f.AppID != 0 {
				appIDMap[f.BundleID] = f.AppID
			}
		}
	}

	var results []models.DownloadedIPA

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".ipa") {
			continue
		}

		fullPath := filepath.Join(downloadDir, entry.Name())
		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}

		ipaItem := parseIPAFile(fullPath, fileInfo, artworkMap, appIDMap)
		if ipaItem != nil {
			results = append(results, *ipaItem)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ModTime.After(results[j].ModTime)
	})

	return results, nil
}

func parseIPAFile(filePath string, fileInfo os.FileInfo, artworkMap map[string]string, appIDMap map[string]int64) *models.DownloadedIPA {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil
	}
	defer r.Close()

	var infoPlistFile *zip.File
	var metadataPlistFile *zip.File

	for _, f := range r.File {
		parts := strings.Split(filepath.ToSlash(f.Name), "/")
		if len(parts) == 3 && parts[0] == "Payload" && strings.HasSuffix(parts[1], ".app") && parts[2] == "Info.plist" {
			infoPlistFile = f
		} else if f.Name == "iTunesMetadata.plist" {
			metadataPlistFile = f
		}
	}

	if infoPlistFile == nil {
		return nil
	}

	rc, err := infoPlistFile.Open()
	if err != nil {
		return nil
	}
	buf, err := io.ReadAll(rc)
	rc.Close()
	if err != nil {
		return nil
	}

	var plistData map[string]interface{}
	decoder := plist.NewDecoder(bytes.NewReader(buf))
	if err := decoder.Decode(&plistData); err != nil {
		return nil
	}

	bundleID := getStringVal(plistData, "CFBundleIdentifier")
	appName := getStringVal(plistData, "CFBundleDisplayName")
	if appName == "" {
		appName = getStringVal(plistData, "CFBundleName")
	}
	version := getStringVal(plistData, "CFBundleShortVersionString")
	shortVersion := version
	if version == "" {
		version = getStringVal(plistData, "CFBundleVersion")
		shortVersion = version
	}
	minimumOS := getStringVal(plistData, "MinimumOSVersion")

	var artworkURL string
	var appID int64

	if metadataPlistFile != nil {
		if mRc, mErr := metadataPlistFile.Open(); mErr == nil {
			if mBuf, mErr2 := io.ReadAll(mRc); mErr2 == nil {
				var metaData map[string]interface{}
				mDecoder := plist.NewDecoder(bytes.NewReader(mBuf))
				if mErr3 := mDecoder.Decode(&metaData); mErr3 == nil {
					if name := getStringVal(metaData, "itemName"); name != "" {
						appName = name
					}
					if idVal, ok := metaData["itemId"].(uint64); ok {
						appID = int64(idVal)
					} else if idVal, ok := metaData["itemId"].(int64); ok {
						appID = idVal
					}
					if art := getStringVal(metaData, "softwareIcon57x57URL"); art != "" {
						artworkURL = art
					}
				}
			}
			mRc.Close()
		}
	}

	if artworkURL == "" && bundleID != "" {
		artworkURL = artworkMap[bundleID]
	}
	if appID == 0 && bundleID != "" {
		appID = appIDMap[bundleID]
	}

	size := fileInfo.Size()

	return &models.DownloadedIPA{
		FilePath:      filePath,
		FileName:      fileInfo.Name(),
		BundleID:      bundleID,
		AppName:       appName,
		Version:       version,
		ShortVersion:  shortVersion,
		MinimumOS:     minimumOS,
		FileSizeBytes: size,
		FormattedSize: formatBytes(size),
		ModTime:       fileInfo.ModTime(),
		ArtworkURL:    artworkURL,
		AppID:         appID,
	}
}

func getStringVal(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
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
		return nil
	}

	// If file doesn't exist, we consider it "deleted" already
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
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

