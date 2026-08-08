package services

import (
	"context"
	"sync"

	"github.com/ElJoker63/ipatool-1/v2/backend/apple"
	"github.com/ElJoker63/ipatool-1/v2/backend/auth"
	"github.com/ElJoker63/ipatool-1/v2/backend/config"
	"github.com/ElJoker63/ipatool-1/v2/backend/device"
	"github.com/ElJoker63/ipatool-1/v2/backend/download"
	"github.com/ElJoker63/ipatool-1/v2/backend/events"
	"github.com/ElJoker63/ipatool-1/v2/backend/library"
	"github.com/ElJoker63/ipatool-1/v2/backend/models"
	"github.com/ElJoker63/ipatool-1/v2/backend/search"
	"github.com/ElJoker63/ipatool-1/v2/backend/storage"
	"github.com/ElJoker63/ipatool-1/v2/backend/update"
)

// AppService aggregates all domain services and is exposed to the Wails frontend.
type AppService struct {
	ctx             context.Context
	storage         storage.Storage
	emitter         events.Emitter
	appleClient     apple.Client
	authService     auth.AuthService
	searchService   search.SearchService
	downloadManager download.DownloadManager
	libraryService  library.LibraryService
	configService   config.ConfigService
	deviceService   device.DeviceService
	updateService   update.UpdateService
}

// NewAppService instantiates all sub-services and builds the service layer.
func NewAppService(dataDir string) (*AppService, error) {
	emitter := events.NewEmitter()

	store, err := storage.NewSQLiteStorage(dataDir)
	if err != nil {
		return nil, err
	}
	emitter.SetStorage(store)

	settings, _ := store.GetSettings()
	passphrase := ""
	if settings != nil {
		passphrase = settings.KeychainPassphrase
	}

	appleClient, err := apple.NewClient(passphrase)
	if err != nil {
		return nil, err
	}

	authService := auth.NewAuthService(appleClient, store, emitter)
	searchService := search.NewSearchService(appleClient, store, emitter)
	downloadManager := download.NewDownloadManager(appleClient, store, emitter)
	libraryService := library.NewLibraryService(store, emitter)
	configService := config.NewConfigService(store, emitter)
	deviceService := device.NewDeviceService(emitter)
	updateService := update.NewUpdateService(emitter)

	return &AppService{
		storage:         store,
		emitter:         emitter,
		appleClient:     appleClient,
		authService:     authService,
		searchService:   searchService,
		downloadManager: downloadManager,
		libraryService:  libraryService,
		configService:   configService,
		deviceService:   deviceService,
		updateService:   updateService,
	}, nil
}

// SetContext is called by the Wails lifecycle upon startup.
func (s *AppService) SetContext(ctx context.Context) {
	s.ctx = ctx
	s.emitter.SetContext(ctx)
	s.configService.SetContext(ctx)
	s.deviceService.SetContext(ctx)
	s.downloadManager.Start()
	s.deviceService.StartWatcher()
}

// Shutdown cleans up resources.
func (s *AppService) Shutdown() {
	s.deviceService.StopWatcher()
	s.downloadManager.Stop()
	_ = s.storage.Close()
}

// ----------------- Device Management Bindings -----------------

func (s *AppService) GetConnectedDevice() (*models.DeviceInfo, error) {
	return s.deviceService.GetConnectedDevice()
}

func (s *AppService) IsDeviceConnected() bool {
	return s.deviceService.IsDeviceConnected()
}

func (s *AppService) PairDevice() error {
	return s.deviceService.PairDevice()
}

func (s *AppService) ListInstalledApps(appType string) ([]models.InstalledApp, error) {
	apps, err := s.deviceService.ListInstalledApps(appType)
	if err != nil {
		return nil, err
	}

	// Enrich apps with artwork concurrently
	var wg sync.WaitGroup
	enrichedApps := make([]models.InstalledApp, len(apps))
	copy(enrichedApps, apps)

	// Limit concurrency to avoid hitting iTunes API rate limits too hard
	sem := make(chan struct{}, 5)

	for i := range enrichedApps {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			bundleID := enrichedApps[idx].BundleID
			// Skip system apps or apps that don't look like standard bundles if needed
			// But usually, we want artwork for all user apps at least.
			meta, err := s.searchService.Lookup(bundleID, "ios")
			if err == nil && meta != nil && meta.ArtworkURL != "" {
				enrichedApps[idx].ArtworkURL = meta.ArtworkURL
			}
		}(i)
	}

	wg.Wait()
	return enrichedApps, nil
}

func (s *AppService) InstallIPA(ipaPath string) error {
	return s.deviceService.InstallIPA(ipaPath)
}

func (s *AppService) UninstallApp(bundleID string) error {
	return s.deviceService.UninstallApp(bundleID)
}

func (s *AppService) ValidateIPA(ipaPath string) (*models.IPAInfo, error) {
	return s.deviceService.ValidateIPA(ipaPath)
}

func (s *AppService) SelectIPAFile() (string, error) {
	return s.deviceService.SelectIPAFile()
}

// ----------------- Auth Bindings -----------------

func (s *AppService) GetAccount() (*models.AccountProfile, error) {
	return s.authService.GetAccount()
}

func (s *AppService) Login(email, password, authCode string, remember bool) (*models.AccountProfile, error) {
	return s.authService.Login(email, password, authCode, remember)
}

func (s *AppService) Logout() error {
	return s.authService.Logout()
}

func (s *AppService) GetAuthStatus() string {
	return s.authService.GetStatus()
}

// ----------------- Search & Details Bindings -----------------

func (s *AppService) SearchApps(term string, platform string, limit int) ([]models.AppMetadata, error) {
	return s.searchService.Search(term, platform, limit)
}

func (s *AppService) LookupApp(bundleID string, platform string) (*models.AppMetadata, error) {
	return s.searchService.Lookup(bundleID, platform)
}

func (s *AppService) GetAppDetails(appID int64, bundleID string, platform string) (*models.AppDetailsOutput, error) {
	return s.searchService.GetAppDetails(appID, bundleID, platform)
}

func (s *AppService) GetSearchHistory(limit int) ([]models.SearchHistoryItem, error) {
	return s.searchService.GetSearchHistory(limit)
}

func (s *AppService) ClearSearchHistory() error {
	return s.searchService.ClearSearchHistory()
}

// ----------------- Download Bindings -----------------

func (s *AppService) QueueDownload(app models.AppMetadata, platform string, externalVersionID string, customOutputPath string) (*models.DownloadTask, error) {
	return s.downloadManager.QueueDownload(app, platform, externalVersionID, customOutputPath)
}

func (s *AppService) PauseDownload(id string) error {
	return s.downloadManager.PauseDownload(id)
}

func (s *AppService) ResumeDownload(id string) error {
	return s.downloadManager.ResumeDownload(id)
}

func (s *AppService) CancelDownload(id string) error {
	return s.downloadManager.CancelDownload(id)
}

func (s *AppService) RetryDownload(id string) error {
	return s.downloadManager.RetryDownload(id)
}

func (s *AppService) GetActiveDownloads() ([]models.DownloadTask, error) {
	return s.downloadManager.GetActiveDownloads()
}

func (s *AppService) GetAllDownloads() ([]models.DownloadTask, error) {
	return s.downloadManager.GetAllDownloads()
}

func (s *AppService) ClearCompletedDownloads() error {
	return s.downloadManager.ClearCompleted()
}

// ----------------- Library & Favorites Bindings -----------------

func (s *AppService) GetFavorites() ([]models.FavoriteApp, error) {
	return s.libraryService.GetFavorites()
}

func (s *AppService) SearchFavorites(query string) ([]models.FavoriteApp, error) {
	return s.libraryService.SearchFavorites(query)
}

func (s *AppService) AddFavorite(app models.FavoriteApp) error {
	return s.libraryService.AddFavorite(app)
}

func (s *AppService) RemoveFavorite(appID int64) error {
	return s.libraryService.RemoveFavorite(appID)
}

func (s *AppService) ToggleFavorite(app models.FavoriteApp) (bool, error) {
	return s.libraryService.ToggleFavorite(app)
}

func (s *AppService) GetDownloadHistory() ([]models.DownloadTask, error) {
	return s.libraryService.GetHistory()
}

func (s *AppService) DeleteHistoryItem(id string) error {
	return s.libraryService.DeleteHistoryItem(id)
}

func (s *AppService) ClearDownloadHistory() error {
	return s.libraryService.ClearHistory()
}

func (s *AppService) OpenFolder(path string) error {
	return s.libraryService.OpenFolder(path)
}

func (s *AppService) OpenFile(path string) error {
	return s.libraryService.OpenFile(path)
}

func (s *AppService) RevealInExplorer(path string) error {
	return s.libraryService.RevealInExplorer(path)
}

// ----------------- Settings & Logs Bindings -----------------

func (s *AppService) GetSettings() (*models.AppSettings, error) {
	return s.configService.GetSettings()
}

func (s *AppService) SaveSettings(settings models.AppSettings) error {
	return s.configService.SaveSettings(settings)
}

func (s *AppService) SelectDownloadDirectory(defaultPath string) (string, error) {
	return s.configService.SelectDownloadDirectory(defaultPath)
}

func (s *AppService) ClearAppCache() error {
	return s.configService.ClearAppCache()
}

func (s *AppService) GetCacheSize() (string, error) {
	return s.configService.GetCacheSize()
}

func (s *AppService) GetLogs(limit int) ([]models.LogEntry, error) {
	return s.storage.GetRecentLogs(limit)
}

func (s *AppService) ClearLogs() error {
	return s.storage.ClearLogs()
}

func (s *AppService) ExportLogs(destinationPath string) (string, error) {
	return s.configService.ExportLogs(destinationPath)
}

func (s *AppService) AddLog(level, message, context string) (*models.LogEntry, error) {
	entry, err := s.storage.AddLog(level, message, context)
	if err == nil && entry != nil {
		s.emitter.EmitLog(level, message, context)
	}
	return entry, err
}

// ----------------- Update Bindings -----------------

func (s *AppService) CheckForUpdate() (*models.UpdateInfo, error) {
	return s.updateService.CheckForUpdate()
}

func (s *AppService) ApplyUpdate(downloadURL string) error {
	return s.updateService.ApplyUpdate(downloadURL)
}
