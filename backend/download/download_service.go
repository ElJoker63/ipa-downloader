package download

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ElJoker63/ipa-downloader/v2/backend/apple"
	"github.com/ElJoker63/ipa-downloader/v2/backend/events"
	"github.com/ElJoker63/ipa-downloader/v2/backend/models"
	"github.com/ElJoker63/ipa-downloader/v2/backend/storage"
	"github.com/ElJoker63/ipa-downloader/v2/pkg/appstore"
)

// DownloadManager manages concurrent downloads with live progress, pause, resume, and cancellation.
type DownloadManager interface {
	Start()
	Stop()
	QueueDownload(app models.AppMetadata, platform string, externalVersionID string, customOutputPath string) (*models.DownloadTask, error)
	QueueFirmwareDownload(deviceName string, fw models.Firmware) (*models.DownloadTask, error)
	PauseDownload(id string) error
	ResumeDownload(id string) error
	CancelDownload(id string) error
	RetryDownload(id string) error
	GetActiveDownloads() ([]models.DownloadTask, error)
	GetAllDownloads() ([]models.DownloadTask, error)
	ClearCompleted() error
}

type downloadTaskHandle struct {
	task   models.DownloadTask
	ctx    context.Context
	cancel context.CancelFunc
}

type downloadManager struct {
	appleClient apple.Client
	storage     storage.Storage
	emitter     events.Emitter
	httpClient  *http.Client
	queueChan   chan string
	handles     map[string]*downloadTaskHandle
	maxWorkers  int
	stopChan    chan struct{}
	mu          sync.RWMutex
}

// NewDownloadManager initializes the download queue and background workers.
func NewDownloadManager(client apple.Client, store storage.Storage, emitter events.Emitter) DownloadManager {
	settings, err := store.GetSettings()
	maxWorkers := 3
	if err == nil && settings.MaxConcurrentDownloads >= 1 && settings.MaxConcurrentDownloads <= 10 {
		maxWorkers = settings.MaxConcurrentDownloads
	}

	dm := &downloadManager{
		appleClient: client,
		storage:     store,
		emitter:     emitter,
		httpClient:  &http.Client{Timeout: 0}, // No timeout for large downloads
		queueChan:   make(chan string, 100),
		handles:     make(map[string]*downloadTaskHandle),
		maxWorkers:  maxWorkers,
		stopChan:    make(chan struct{}),
	}

	return dm
}

func (m *downloadManager) Start() {
	for i := 0; i < m.maxWorkers; i++ {
		go m.worker(i)
	}

	// Re-queue any unfinished downloads on startup
	active, err := m.storage.GetActiveDownloads()
	if err == nil {
		for _, task := range active {
			if task.Status == models.DownloadStatusDownloading || task.Status == models.DownloadStatusQueued {
				m.queueChan <- task.ID
			}
		}
	}
}

func (m *downloadManager) Stop() {
	close(m.stopChan)
}

func (m *downloadManager) QueueDownload(app models.AppMetadata, platform string, externalVersionID string, customOutputPath string) (*models.DownloadTask, error) {
	settings, err := m.storage.GetSettings()
	destFolder := ""
	if err == nil && settings != nil {
		destFolder = settings.DefaultDownloadFolder
	}
	if destFolder == "" {
		home, _ := os.UserHomeDir()
		destFolder = filepath.Join(home, "Downloads")
	}

	if customOutputPath != "" {
		destFolder = customOutputPath
	}

	fileName := fmt.Sprintf("%s_%d_%s.ipa", sanitizeFilename(app.BundleID), app.ID, sanitizeFilename(app.Version))
	destPath := filepath.Join(destFolder, fileName)

	taskID := fmt.Sprintf("%d-%d", app.ID, time.Now().UnixNano())
	now := time.Now()

	task := models.DownloadTask{
		ID:                taskID,
		Type:              "app",
		AppID:             app.ID,
		BundleID:          app.BundleID,
		AppName:           app.Name,
		Developer:         app.Developer,
		Version:           app.Version,
		ArtworkURL:        app.ArtworkURL,
		DestinationPath:   destPath,
		Status:            models.DownloadStatusQueued,
		TotalBytes:        app.FileSizeBytes,
		DownloadedBytes:   0,
		Progress:          0.0,
		ExternalVersionID: externalVersionID,
		Platform:          platform,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := m.storage.SaveDownload(task); err != nil {
		return nil, fmt.Errorf("failed to save download to database: %w", err)
	}

	m.emitter.EmitLog("INFO", fmt.Sprintf("Queued download for %s (%s)", app.Name, app.BundleID), "DownloadManager")
	m.emitter.Emit(events.EventDownloadStatus, task)

	m.queueChan <- taskID
	return &task, nil
}

func (m *downloadManager) QueueFirmwareDownload(deviceName string, fw models.Firmware) (*models.DownloadTask, error) {
	settings, err := m.storage.GetSettings()
	destFolder := ""
	if err == nil && settings != nil {
		destFolder = settings.DefaultDownloadFolder
	}
	if destFolder == "" {
		home, _ := os.UserHomeDir()
		destFolder = filepath.Join(home, "Downloads")
	}

	destPath := filepath.Join(destFolder, fw.Filename)
	taskID := fmt.Sprintf("fw-%s-%d", fw.BuildID, time.Now().UnixNano())
	now := time.Now()

	task := models.DownloadTask{
		ID:              taskID,
		Type:            "firmware",
		URL:             fw.URL,
		AppName:         fmt.Sprintf("%s Firmware", deviceName),
		Version:         fw.Version,
		ArtworkURL:      "https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png",
		DestinationPath: destPath,
		Status:          models.DownloadStatusQueued,
		TotalBytes:      fw.Size,
		DownloadedBytes: 0,
		Progress:        0.0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := m.storage.SaveDownload(task); err != nil {
		return nil, err
	}

	m.emitter.EmitLog("INFO", fmt.Sprintf("Queued firmware download: %s", fw.Filename), "DownloadManager")
	m.emitter.Emit(events.EventDownloadStatus, task)

	m.queueChan <- taskID
	return &task, nil
}

func (m *downloadManager) PauseDownload(id string) error {
	m.mu.Lock()
	handle, ok := m.handles[id]
	if ok && handle.cancel != nil {
		handle.cancel()
		delete(m.handles, id)
	}
	m.mu.Unlock()

	task, err := m.storage.GetDownload(id)
	if err != nil {
		return err
	}

	task.Status = models.DownloadStatusPaused
	task.SpeedBytesPerSec = 0
	task.FormattedSpeed = "0 KB/s"
	task.ETASeconds = 0
	task.FormattedETA = "--"

	_ = m.storage.UpdateDownloadProgress(id, models.DownloadStatusPaused, task.DownloadedBytes, task.TotalBytes, task.Progress, 0, 0, "")
	m.emitter.EmitLog("INFO", fmt.Sprintf("Paused download: %s", task.AppName), "DownloadManager")
	m.emitter.Emit(events.EventDownloadPaused, task)
	m.emitter.EmitDownloadProgress(*task)

	return nil
}

func (m *downloadManager) ResumeDownload(id string) error {
	task, err := m.storage.GetDownload(id)
	if err != nil {
		return err
	}

	task.Status = models.DownloadStatusQueued
	_ = m.storage.SaveDownload(*task)
	m.emitter.EmitLog("INFO", fmt.Sprintf("Resumed download: %s", task.AppName), "DownloadManager")
	m.emitter.Emit(events.EventDownloadStatus, task)

	m.queueChan <- id
	return nil
}

func (m *downloadManager) CancelDownload(id string) error {
	m.mu.Lock()
	handle, ok := m.handles[id]
	if ok && handle.cancel != nil {
		handle.cancel()
		delete(m.handles, id)
	}
	m.mu.Unlock()

	task, err := m.storage.GetDownload(id)
	if err != nil {
		return err
	}

	task.Status = models.DownloadStatusCancelled
	_ = m.storage.UpdateDownloadProgress(id, models.DownloadStatusCancelled, task.DownloadedBytes, task.TotalBytes, task.Progress, 0, 0, "Cancelled by user")

	// Clean up temporary file
	tmpPath := task.DestinationPath + ".tmp"
	_ = os.Remove(tmpPath)

	m.emitter.EmitLog("INFO", fmt.Sprintf("Cancelled download: %s", task.AppName), "DownloadManager")
	m.emitter.Emit(events.EventDownloadCancelled, task)
	m.emitter.EmitDownloadProgress(*task)

	return nil
}

func (m *downloadManager) RetryDownload(id string) error {
	task, err := m.storage.GetDownload(id)
	if err != nil {
		return err
	}

	task.Status = models.DownloadStatusQueued
	task.DownloadedBytes = 0
	task.Progress = 0.0
	task.Error = ""
	_ = m.storage.SaveDownload(*task)

	m.emitter.EmitLog("INFO", fmt.Sprintf("Retrying download: %s", task.AppName), "DownloadManager")
	m.emitter.Emit(events.EventDownloadStatus, task)

	m.queueChan <- id
	return nil
}

func (m *downloadManager) GetActiveDownloads() ([]models.DownloadTask, error) {
	return m.storage.GetActiveDownloads()
}

func (m *downloadManager) GetAllDownloads() ([]models.DownloadTask, error) {
	return m.storage.GetAllDownloads()
}

func (m *downloadManager) ClearCompleted() error {
	return m.storage.ClearCompletedDownloads()
}

func (m *downloadManager) worker(workerID int) {
	for {
		select {
		case <-m.stopChan:
			return
		case taskID := <-m.queueChan:
			m.processDownload(taskID)
		}
	}
}

func (m *downloadManager) processDownload(taskID string) {
	task, err := m.storage.GetDownload(taskID)
	if err != nil || task == nil || task.Status == models.DownloadStatusCancelled || task.Status == models.DownloadStatusPaused {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	handle := &downloadTaskHandle{
		task:   *task,
		ctx:    ctx,
		cancel: cancel,
	}

	m.mu.Lock()
	m.handles[taskID] = handle
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.handles, taskID)
		m.mu.Unlock()
	}()

	task.Status = models.DownloadStatusDownloading
	_ = m.storage.SaveDownload(*task)
	m.emitter.EmitLog("INFO", fmt.Sprintf("Starting download for %s...", task.AppName), "DownloadManager")
	m.emitter.Emit(events.EventDownloadStatus, task)

	if task.Type == "firmware" {
		err = m.executeFirmwareDownload(ctx, task)
	} else {
		err = m.executeAppDownload(ctx, task)
	}

	if err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			m.emitter.EmitLog("INFO", fmt.Sprintf("Download stopped: %s", task.AppName), "DownloadManager")
			return
		}

		task.Status = models.DownloadStatusFailed
		task.Error = err.Error()
		_ = m.storage.UpdateDownloadProgress(taskID, models.DownloadStatusFailed, task.DownloadedBytes, task.TotalBytes, task.Progress, 0, 0, err.Error())
		m.emitter.EmitLog("ERROR", fmt.Sprintf("Download failed for %s: %v", task.AppName, err), "DownloadManager")
		m.emitter.Emit(events.EventDownloadFailed, task)
		m.emitter.EmitNotification("Download Failed", fmt.Sprintf("%s: %v", task.AppName, err), "error")
		return
	}

	// Completed
	task.Status = models.DownloadStatusCompleted
	task.Progress = 100.0
	task.SpeedBytesPerSec = 0
	task.FormattedSpeed = "0 KB/s"
	task.ETASeconds = 0
	task.FormattedETA = "Done"
	task.DownloadedBytes = task.TotalBytes

	_ = m.storage.UpdateDownloadProgress(taskID, models.DownloadStatusCompleted, task.TotalBytes, task.TotalBytes, 100.0, 0, 0, "")
	m.emitter.EmitLog("SUCCESS", fmt.Sprintf("Successfully downloaded %s to %s", task.AppName, task.DestinationPath), "DownloadManager")
	m.emitter.Emit(events.EventDownloadCompleted, task)
	m.emitter.EmitNotification("Download Complete", fmt.Sprintf("%s is ready!", task.AppName), "success")
}

func (m *downloadManager) executeAppDownload(ctx context.Context, task *models.DownloadTask) error {
	appstoreCore := m.appleClient.GetAppStore()

	accInfo, err := appstoreCore.AccountInfo()
	if err != nil {
		return fmt.Errorf("account not authenticated: %w", err)
	}

	app := appstore.App{
		ID:       task.AppID,
		BundleID: task.BundleID,
		Name:     task.AppName,
		Version:  task.Version,
	}

	platform, err := appstore.ParsePlatform(task.Platform)
	if err != nil {
		platform = appstore.PlatformIPhone
	}

	settings, _ := m.storage.GetSettings()
	if settings == nil || settings.AutoAcquireLicense {
		_ = appstoreCore.Purchase(appstore.PurchaseInput{
			Account: accInfo.Account,
			App:     app,
		})
	}

	dir := filepath.Dir(task.DestinationPath)
	_ = os.MkdirAll(dir, 0755)

	out, err := m.streamAppDownload(ctx, appstoreCore, accInfo.Account, app, task, platform)
	if err != nil {
		if errors.Is(err, appstore.ErrLicenseRequired) {
			m.emitter.EmitLog("INFO", fmt.Sprintf("Acquiring free license for %s...", task.AppName), "DownloadManager")
			_ = appstoreCore.Purchase(appstore.PurchaseInput{
				Account: accInfo.Account,
				App:     app,
			})
			out, err = m.streamAppDownload(ctx, appstoreCore, accInfo.Account, app, task, platform)
		}
		if err != nil {
			return err
		}
	}

	task.Status = models.DownloadStatusSigning
	task.Progress = 100.0
	task.SpeedBytesPerSec = 0
	task.FormattedSpeed = "Signing FairPlay DRM..."
	task.ETASeconds = 0
	task.FormattedETA = "Finalizing .ipa..."
	_ = m.storage.UpdateDownloadProgress(task.ID, models.DownloadStatusSigning, task.TotalBytes, task.TotalBytes, 100.0, 0, 0, "")
	m.emitter.EmitLog("INFO", fmt.Sprintf("[%s] Download stream complete (100%%). Starting FairPlay SINF signing process...", task.AppName), "DownloadManager")
	m.emitter.EmitDownloadProgress(*task)

	if len(out.Sinfs) > 0 {
		m.emitter.EmitLog("INFO", fmt.Sprintf("[%s] Injecting %d FairPlay SINF DRM signature(s) into final .ipa package...", task.AppName, len(out.Sinfs)), "DownloadManager")
		err = appstoreCore.ReplicateSinf(appstore.ReplicateSinfInput{
			Sinfs:       out.Sinfs,
			PackagePath: task.DestinationPath,
		})
		if err != nil {
			m.emitter.EmitLog("WARN", fmt.Sprintf("[%s] FairPlay SINF replication warning: %v", task.AppName, err), "DownloadManager")
		} else {
			m.emitter.EmitLog("SUCCESS", fmt.Sprintf("[%s] FairPlay SINF signatures successfully injected into %s", task.AppName, filepath.Base(task.DestinationPath)), "DownloadManager")
		}
	}

	return nil
}

func (m *downloadManager) executeFirmwareDownload(ctx context.Context, task *models.DownloadTask) error {
	dir := filepath.Dir(task.DestinationPath)
	_ = os.MkdirAll(dir, 0755)

	req, err := http.NewRequestWithContext(ctx, "GET", task.URL, nil)
	if err != nil {
		return err
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	tmpPath := task.DestinationPath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer out.Close()

	lastUpdateTime := time.Now()
	var lastBytes int64 = 0
	var smoothedSpeed float64 = 0
	var downloaded int64 = 0
	total := resp.ContentLength
	if total <= 0 {
		total = task.TotalBytes
	}

	buffer := make([]byte, 64*1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, writeErr := out.Write(buffer[:n])
			if writeErr != nil {
				return writeErr
			}
			downloaded += int64(n)

			// Progress update
			now := time.Now()
			elapsed := now.Sub(lastUpdateTime).Seconds()
			if elapsed >= 0.2 || downloaded == total {
				deltaBytes := downloaded - lastBytes
				if elapsed > 0 {
					instantSpeed := float64(deltaBytes) / elapsed
					if smoothedSpeed == 0 {
						smoothedSpeed = instantSpeed
					} else {
						smoothedSpeed = 0.2*instantSpeed + 0.8*smoothedSpeed
					}
				}

				speedBps := int64(smoothedSpeed)
				var etaSec int64 = 0
				if speedBps > 0 && total > downloaded {
					etaSec = (total - downloaded) / speedBps
				}

				percent := (float64(downloaded) / float64(total)) * 100.0
				task.DownloadedBytes = downloaded
				task.TotalBytes = total
				task.Progress = percent
				task.SpeedBytesPerSec = speedBps
				task.FormattedSpeed = formatSpeed(speedBps)
				task.ETASeconds = etaSec
				task.FormattedETA = formatETA(etaSec)

				_ = m.storage.UpdateDownloadProgress(task.ID, models.DownloadStatusDownloading, downloaded, total, percent, speedBps, etaSec, "")
				m.emitter.EmitDownloadProgress(*task)

				lastUpdateTime = now
				lastBytes = downloaded
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	out.Close()
	_ = os.Remove(task.DestinationPath)
	return os.Rename(tmpPath, task.DestinationPath)
}

func (m *downloadManager) streamAppDownload(ctx context.Context, store appstore.AppStore, acc appstore.Account, app appstore.App, task *models.DownloadTask, platform appstore.Platform) (appstore.DownloadOutput, error) {
	lastUpdateTime := time.Now()
	var lastBytes int64 = 0
	var smoothedSpeed float64 = 0

	progressCallback := func(downloadedBytes int64, totalBytes int64) {
		now := time.Now()
		elapsed := now.Sub(lastUpdateTime).Seconds()

		var percent float64 = 0
		if totalBytes > 0 {
			percent = (float64(downloadedBytes) / float64(totalBytes)) * 100.0
		}

		if elapsed >= 0.2 || downloadedBytes == totalBytes {
			deltaBytes := downloadedBytes - lastBytes
			if elapsed > 0 {
				instantSpeed := float64(deltaBytes) / elapsed
				if smoothedSpeed == 0 {
					smoothedSpeed = instantSpeed
				} else {
					smoothedSpeed = 0.3*instantSpeed + 0.7*smoothedSpeed
				}
			}

			speedBps := int64(smoothedSpeed)
			var etaSec int64 = 0
			if speedBps > 0 && totalBytes > downloadedBytes {
				etaSec = (totalBytes - downloadedBytes) / speedBps
			}

			task.DownloadedBytes = downloadedBytes
			task.TotalBytes = totalBytes
			task.Progress = percent
			task.SpeedBytesPerSec = speedBps
			task.FormattedSpeed = formatSpeed(speedBps)
			task.ETASeconds = etaSec
			task.FormattedETA = formatETA(etaSec)

			_ = m.storage.UpdateDownloadProgress(task.ID, models.DownloadStatusDownloading, downloadedBytes, totalBytes, percent, speedBps, etaSec, "")
			m.emitter.EmitDownloadProgress(*task)

			lastUpdateTime = now
			lastBytes = downloadedBytes
		}
	}

	return store.Download(appstore.DownloadInput{
		Account:           acc,
		App:               app,
		OutputPath:        task.DestinationPath,
		ExternalVersionID: task.ExternalVersionID,
		Platform:          platform,
		ProgressCallback:  progressCallback,
	})
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer("/", "_", "\\", "_", ":", "_", "*", "_", "?", "_", "\"", "_", "<", "_", ">", "_", "|", "_", " ", "_")
	return replacer.Replace(name)
}

func formatSpeed(bytesPerSec int64) string {
	const unit = 1024
	if bytesPerSec < unit {
		return fmt.Sprintf("%d B/s", bytesPerSec)
	}
	div, exp := int64(unit), 0
	for n := bytesPerSec / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB/s", float64(bytesPerSec)/float64(div), "KMGTPE"[exp])
}

func formatETA(seconds int64) string {
	if seconds <= 0 {
		return "--"
	}
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%dm %ds", seconds/60, seconds%60)
	}
	return fmt.Sprintf("%dh %dm", seconds/3600, (seconds%3600)/60)
}
