package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ElJoker63/ipa-downloader/v2/backend/events"
	"github.com/ElJoker63/ipa-downloader/v2/backend/models"
	"github.com/ElJoker63/ipa-downloader/v2/pkg/version"
	"github.com/inconshreveable/go-update"
)

const (
	GithubRepo = "ElJoker63/ipa-downloader"
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
		BrowserDownloadURL string `json:"browser_download_url"`
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

	if !isNewer(latestVersion, currentVersion) {
		return &models.UpdateInfo{Available: false, CurrentVersion: currentVersion}, nil
	}

	// Find the correct asset for the current platform
	downloadURL := ""
	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "macos"
	}

	// The workflow generates names like:
	// ipa-downloader-1.0.0-windows.exe
	// ipa-downloader-1.0.0-macos-universal.zip
	// ipa-downloader-1.0.0-linux

	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, "ipa-downloader") && strings.Contains(name, latestVersion) && strings.Contains(name, osName) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return nil, fmt.Errorf("could not find compatible binary for %s in release %s", osName, release.TagName)
	}

	return &models.UpdateInfo{
		Available:      true,
		LatestVersion:  latestVersion,
		CurrentVersion: currentVersion,
		ReleaseNotes:   release.Body,
		DownloadURL:    downloadURL,
		Mandatory:      true,
	}, nil
}

func (s *updateService) ApplyUpdate(downloadURL string) error {
	exePath, err := os.Executable()
	if err == nil {
		if evalPath, err := filepath.EvalSymlinks(exePath); err == nil {
			exePath = evalPath
		}
	}

	s.emitter.EmitLog("INFO", "Downloading update from "+downloadURL, "UpdateService")

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download update: status %d", resp.StatusCode)
	}

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

	// Allow a brief moment for the event to reach the frontend
	time.Sleep(500 * time.Millisecond)

	if exePath != "" {
		if err := restartApp(exePath); err != nil {
			s.emitter.EmitLog("ERROR", "Failed to restart application: "+err.Error(), "UpdateService")
			return err
		}
	}

	os.Exit(0)
	return nil
}

func restartApp(exePath string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "darwin" && strings.Contains(exePath, ".app/Contents/MacOS") {
		idx := strings.Index(exePath, ".app")
		appBundle := exePath[:idx+4]
		cmd = exec.Command("open", "-n", appBundle)
	} else {
		cmd = exec.Command(exePath, os.Args[1:]...)
	}

	cmd.Env = os.Environ()
	setDetachFlags(cmd)

	return cmd.Start()
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

func isNewer(latest, current string) bool {
	if latest == current {
		return false
	}

	lParts := strings.Split(latest, ".")
	cParts := strings.Split(current, ".")

	for i := 0; i < len(lParts) && i < len(cParts); i++ {
		var lV, cV int
		fmt.Sscanf(lParts[i], "%d", &lV)
		fmt.Sscanf(cParts[i], "%d", &cV)

		if lV > cV {
			return true
		}
		if lV < cV {
			return false
		}
	}

	return len(lParts) > len(cParts)
}
