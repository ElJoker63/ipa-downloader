package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/inconshreveable/go-update"
	"github.com/majd/ipa-downloader/v2/backend/events"
	"github.com/majd/ipa-downloader/v2/backend/models"
	"github.com/majd/ipa-downloader/v2/pkg/version"
)

const (
	GithubRepo = "ElJoker63/ipatool-1"
	ReleaseAPI = "https://api.github.com/repos/%s/releases/latest"
)

type UpdateService interface {
	CheckForUpdate() (*models.UpdateInfo, error)
	ApplyUpdate(downloadURL string) error
}

type updateService struct {
	emitter events.Emitter
}

func NewUpdateService(emitter events.Emitter) UpdateService {
	return &updateService{
		emitter: emitter,
	}
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"download_url"`
	} `json:"assets"`
}

func (s *updateService) CheckForUpdate() (*models.UpdateInfo, error) {
	resp, err := http.Get(fmt.Sprintf(ReleaseAPI, GithubRepo))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode github release: %w", err)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := version.Version

	if latestVersion == currentVersion {
		return &models.UpdateInfo{Available: false, CurrentVersion: currentVersion}, nil
	}

	// Find the correct asset for the current platform
	downloadURL := ""
	platform := runtime.GOOS
	arch := runtime.GOARCH

	expectedName := fmt.Sprintf("ipa-downloader-%s-%s", platform, arch)
	if platform == "windows" {
		expectedName += ".exe"
	}

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, expectedName) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return nil, fmt.Errorf("could not find compatible binary for %s/%s in release %s", platform, arch, release.TagName)
	}

	return &models.UpdateInfo{
		Available:      true,
		LatestVersion:  latestVersion,
		CurrentVersion: currentVersion,
		ReleaseNotes:   release.Body,
		DownloadURL:    downloadURL,
		Mandatory:      true, // User requested mandatory update
	}, nil
}

func (s *updateService) ApplyUpdate(downloadURL string) error {
	s.emitter.EmitLog("INFO", "Downloading update from "+downloadURL, "UpdateService")

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download update: status %d", resp.StatusCode)
	}

	// Track download progress if possible, but go-update takes a reader.
	// We can wrap the reader to emit events.
	progressReader := &progressReader{
		Reader:  resp.Body,
		Total:   resp.ContentLength,
		Emitter: s.emitter,
	}

	err = update.Apply(progressReader, update.Options{})
	if err != nil {
		s.emitter.EmitLog("ERROR", "Failed to apply update: "+err.Error(), "UpdateService")
		return err
	}

	s.emitter.EmitLog("SUCCESS", "Update applied successfully. Restarting...", "UpdateService")

	// Restart logic depends on platform, but usually we just exit and let the user (o standard wrapper) restart.
	// For Wails, we might want to tell the frontend to show a "Restart Now" button or just exit.
	// The user requested "reemplazando el .exe que estaba en ejecucion", which go-update does.

	os.Exit(0)
	return nil
}

type progressReader struct {
	io.Reader
	Total   int64
	Current int64
	Emitter events.Emitter
}

func (r *progressReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.Current += int64(n)
	if r.Total > 0 {
		percent := int(float64(r.Current) / float64(r.Total) * 100)
		r.Emitter.Emit("update:progress", percent)
	}
	return n, err
}
