package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/majd/ipatool/v2/backend/models"
	_ "modernc.org/sqlite"
)

// Storage defines the interface for local persistence.
type Storage interface {
	Close() error

	// Downloads
	SaveDownload(task models.DownloadTask) error
	UpdateDownloadProgress(id string, status models.DownloadStatus, downloadedBytes, totalBytes int64, progress float64, speed int64, eta int64, errStr string) error
	GetDownload(id string) (*models.DownloadTask, error)
	GetAllDownloads() ([]models.DownloadTask, error)
	GetActiveDownloads() ([]models.DownloadTask, error)
	DeleteDownload(id string) error
	ClearCompletedDownloads() error

	// Favorites
	AddFavorite(app models.FavoriteApp) error
	RemoveFavorite(appID int64) error
	IsFavorite(appID int64) (bool, error)
	GetFavorites() ([]models.FavoriteApp, error)
	SearchFavorites(query string) ([]models.FavoriteApp, error)

	// Search History
	AddSearchHistory(term, platform string, resultCount int) error
	GetSearchHistory(limit int) ([]models.SearchHistoryItem, error)
	ClearSearchHistory() error

	// Settings
	GetSettings() (*models.AppSettings, error)
	SaveSettings(settings models.AppSettings) error

	// App Cache
	CacheAppMetadata(app models.AppMetadata) error
	GetCachedAppMetadata(bundleID string) (*models.AppMetadata, error)
	ClearAppCache() error
	GetCacheSizeBytes() (int64, error)

	// Logs
	AddLog(level, message, context string) (*models.LogEntry, error)
	GetRecentLogs(limit int) ([]models.LogEntry, error)
	ClearLogs() error
}

type sqliteStorage struct {
	db     *sql.DB
	dbPath string
	mu     sync.RWMutex
}

// NewSQLiteStorage initializes the SQLite database at the specified directory.
func NewSQLiteStorage(dataDir string) (Storage, error) {
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "ipatool.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite single-writer safe mode

	s := &sqliteStorage{
		db:     db,
		dbPath: dbPath,
	}

	if err := s.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return s, nil
}

func (s *sqliteStorage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *sqliteStorage) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS downloads (
		id TEXT PRIMARY KEY,
		app_id INTEGER NOT NULL,
		bundle_id TEXT NOT NULL,
		app_name TEXT NOT NULL,
		developer TEXT NOT NULL,
		version TEXT NOT NULL,
		artwork_url TEXT NOT NULL,
		destination_path TEXT NOT NULL,
		status TEXT NOT NULL,
		total_bytes INTEGER DEFAULT 0,
		downloaded_bytes INTEGER DEFAULT 0,
		progress REAL DEFAULT 0.0,
		speed_bytes INTEGER DEFAULT 0,
		eta_seconds INTEGER DEFAULT 0,
		external_version_id TEXT,
		platform TEXT DEFAULT 'ios',
		error TEXT,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		completed_at TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_downloads_status ON downloads(status);
	CREATE INDEX IF NOT EXISTS idx_downloads_created ON downloads(created_at DESC);

	CREATE TABLE IF NOT EXISTS favorites (
		app_id INTEGER PRIMARY KEY,
		bundle_id TEXT NOT NULL,
		name TEXT NOT NULL,
		developer TEXT NOT NULL,
		version TEXT NOT NULL,
		price REAL DEFAULT 0,
		formatted_price TEXT,
		artwork_url TEXT NOT NULL,
		primary_genre TEXT,
		created_at TIMESTAMP NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_favorites_name ON favorites(name);

	CREATE TABLE IF NOT EXISTS search_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		term TEXT NOT NULL,
		platform TEXT NOT NULL,
		count INTEGER DEFAULT 0,
		created_at TIMESTAMP NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_search_created ON search_history(created_at DESC);

	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS app_cache (
		bundle_id TEXT PRIMARY KEY,
		data_json TEXT NOT NULL,
		cached_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp TIMESTAMP NOT NULL,
		level TEXT NOT NULL,
		message TEXT NOT NULL,
		context TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp DESC);
	`

	_, err := s.db.Exec(schema)
	return err
}

// ----------------- Downloads -----------------

func (s *sqliteStorage) SaveDownload(task models.DownloadTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `
	INSERT INTO downloads (
		id, app_id, bundle_id, app_name, developer, version, artwork_url,
		destination_path, status, total_bytes, downloaded_bytes, progress,
		speed_bytes, eta_seconds, external_version_id, platform, error,
		created_at, updated_at, completed_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		status = excluded.status,
		total_bytes = excluded.total_bytes,
		downloaded_bytes = excluded.downloaded_bytes,
		progress = excluded.progress,
		speed_bytes = excluded.speed_bytes,
		eta_seconds = excluded.eta_seconds,
		error = excluded.error,
		updated_at = excluded.updated_at,
		completed_at = excluded.completed_at;
	`

	_, err := s.db.Exec(
		query,
		task.ID, task.AppID, task.BundleID, task.AppName, task.Developer, task.Version, task.ArtworkURL,
		task.DestinationPath, string(task.Status), task.TotalBytes, task.DownloadedBytes, task.Progress,
		task.SpeedBytesPerSec, task.ETASeconds, task.ExternalVersionID, task.Platform, task.Error,
		task.CreatedAt, task.UpdatedAt, task.CompletedAt,
	)
	return err
}

func (s *sqliteStorage) UpdateDownloadProgress(id string, status models.DownloadStatus, downloadedBytes, totalBytes int64, progress float64, speed int64, eta int64, errStr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	var completedAt *time.Time
	if status == models.DownloadStatusCompleted {
		completedAt = &now
	}

	query := `
	UPDATE downloads SET
		status = ?,
		downloaded_bytes = ?,
		total_bytes = ?,
		progress = ?,
		speed_bytes = ?,
		eta_seconds = ?,
		error = ?,
		updated_at = ?,
		completed_at = CASE WHEN ? IS NOT NULL THEN ? ELSE completed_at END
	WHERE id = ?;
	`

	_, err := s.db.Exec(query, string(status), downloadedBytes, totalBytes, progress, speed, eta, errStr, now, completedAt, completedAt, id)
	return err
}

func (s *sqliteStorage) GetDownload(id string) (*models.DownloadTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
	SELECT id, app_id, bundle_id, app_name, developer, version, artwork_url,
	       destination_path, status, total_bytes, downloaded_bytes, progress,
	       speed_bytes, eta_seconds, external_version_id, platform, error,
	       created_at, updated_at, completed_at
	FROM downloads WHERE id = ?;
	`

	row := s.db.QueryRow(query, id)
	return s.scanDownload(row)
}

func (s *sqliteStorage) GetAllDownloads() ([]models.DownloadTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
	SELECT id, app_id, bundle_id, app_name, developer, version, artwork_url,
	       destination_path, status, total_bytes, downloaded_bytes, progress,
	       speed_bytes, eta_seconds, external_version_id, platform, error,
	       created_at, updated_at, completed_at
	FROM downloads ORDER BY created_at DESC;
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.DownloadTask
	for rows.Next() {
		task, err := s.scanDownloadRows(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (s *sqliteStorage) GetActiveDownloads() ([]models.DownloadTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
	SELECT id, app_id, bundle_id, app_name, developer, version, artwork_url,
	       destination_path, status, total_bytes, downloaded_bytes, progress,
	       speed_bytes, eta_seconds, external_version_id, platform, error,
	       created_at, updated_at, completed_at
	FROM downloads WHERE status IN ('queued', 'downloading', 'paused') ORDER BY created_at ASC;
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.DownloadTask
	for rows.Next() {
		task, err := s.scanDownloadRows(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (s *sqliteStorage) DeleteDownload(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM downloads WHERE id = ?", id)
	return err
}

func (s *sqliteStorage) ClearCompletedDownloads() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM downloads WHERE status IN ('completed', 'cancelled', 'failed')")
	return err
}

func (s *sqliteStorage) scanDownload(row *sql.Row) (*models.DownloadTask, error) {
	var task models.DownloadTask
	var status string
	var externalVersionID, errStr sql.NullString
	var completedAt sql.NullTime

	err := row.Scan(
		&task.ID, &task.AppID, &task.BundleID, &task.AppName, &task.Developer, &task.Version, &task.ArtworkURL,
		&task.DestinationPath, &status, &task.TotalBytes, &task.DownloadedBytes, &task.Progress,
		&task.SpeedBytesPerSec, &task.ETASeconds, &externalVersionID, &task.Platform, &errStr,
		&task.CreatedAt, &task.UpdatedAt, &completedAt,
	)
	if err != nil {
		return nil, err
	}

	task.Status = models.DownloadStatus(status)
	if externalVersionID.Valid {
		task.ExternalVersionID = externalVersionID.String
	}
	if errStr.Valid {
		task.Error = errStr.String
	}
	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	return &task, nil
}

func (s *sqliteStorage) scanDownloadRows(rows *sql.Rows) (*models.DownloadTask, error) {
	var task models.DownloadTask
	var status string
	var externalVersionID, errStr sql.NullString
	var completedAt sql.NullTime

	err := rows.Scan(
		&task.ID, &task.AppID, &task.BundleID, &task.AppName, &task.Developer, &task.Version, &task.ArtworkURL,
		&task.DestinationPath, &status, &task.TotalBytes, &task.DownloadedBytes, &task.Progress,
		&task.SpeedBytesPerSec, &task.ETASeconds, &externalVersionID, &task.Platform, &errStr,
		&task.CreatedAt, &task.UpdatedAt, &completedAt,
	)
	if err != nil {
		return nil, err
	}

	task.Status = models.DownloadStatus(status)
	if externalVersionID.Valid {
		task.ExternalVersionID = externalVersionID.String
	}
	if errStr.Valid {
		task.Error = errStr.String
	}
	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	return &task, nil
}

// ----------------- Favorites -----------------

func (s *sqliteStorage) AddFavorite(app models.FavoriteApp) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `
	INSERT INTO favorites (app_id, bundle_id, name, developer, version, price, formatted_price, artwork_url, primary_genre, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(app_id) DO UPDATE SET
		version = excluded.version,
		price = excluded.price,
		formatted_price = excluded.formatted_price,
		artwork_url = excluded.artwork_url;
	`

	_, err := s.db.Exec(query, app.AppID, app.BundleID, app.Name, app.Developer, app.Version, app.Price, app.FormattedPrice, app.ArtworkURL, app.PrimaryGenre, app.CreatedAt)
	return err
}

func (s *sqliteStorage) RemoveFavorite(appID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM favorites WHERE app_id = ?", appID)
	return err
}

func (s *sqliteStorage) IsFavorite(appID int64) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM favorites WHERE app_id = ?", appID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *sqliteStorage) GetFavorites() ([]models.FavoriteApp, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT app_id, bundle_id, name, developer, version, price, formatted_price, artwork_url, primary_genre, created_at FROM favorites ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.FavoriteApp
	for rows.Next() {
		var item models.FavoriteApp
		var primaryGenre sql.NullString
		if err := rows.Scan(&item.AppID, &item.BundleID, &item.Name, &item.Developer, &item.Version, &item.Price, &item.FormattedPrice, &item.ArtworkURL, &primaryGenre, &item.CreatedAt); err != nil {
			return nil, err
		}
		if primaryGenre.Valid {
			item.PrimaryGenre = primaryGenre.String
		}
		list = append(list, item)
	}

	return list, nil
}

func (s *sqliteStorage) SearchFavorites(query string) ([]models.FavoriteApp, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	likePattern := "%" + query + "%"
	rows, err := s.db.Query(`
		SELECT app_id, bundle_id, name, developer, version, price, formatted_price, artwork_url, primary_genre, created_at
		FROM favorites
		WHERE name LIKE ? OR bundle_id LIKE ? OR developer LIKE ? OR primary_genre LIKE ?
		ORDER BY name ASC
	`, likePattern, likePattern, likePattern, likePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.FavoriteApp
	for rows.Next() {
		var item models.FavoriteApp
		var primaryGenre sql.NullString
		if err := rows.Scan(&item.AppID, &item.BundleID, &item.Name, &item.Developer, &item.Version, &item.Price, &item.FormattedPrice, &item.ArtworkURL, &primaryGenre, &item.CreatedAt); err != nil {
			return nil, err
		}
		if primaryGenre.Valid {
			item.PrimaryGenre = primaryGenre.String
		}
		list = append(list, item)
	}

	return list, nil
}

// ----------------- Search History -----------------

func (s *sqliteStorage) AddSearchHistory(term, platform string, resultCount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("INSERT INTO search_history (term, platform, count, created_at) VALUES (?, ?, ?, ?)", term, platform, resultCount, time.Now())
	return err
}

func (s *sqliteStorage) GetSearchHistory(limit int) ([]models.SearchHistoryItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 {
		limit = 20
	}

	rows, err := s.db.Query("SELECT id, term, platform, count, created_at FROM search_history ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.SearchHistoryItem
	for rows.Next() {
		var item models.SearchHistoryItem
		if err := rows.Scan(&item.ID, &item.Term, &item.Platform, &item.Count, &item.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (s *sqliteStorage) ClearSearchHistory() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM search_history")
	return err
}

// ----------------- Settings -----------------

func (s *sqliteStorage) GetSettings() (*models.AppSettings, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	homeDir, _ := os.UserHomeDir()
	defaultDownloads := filepath.Join(homeDir, "Downloads")

	settings := &models.AppSettings{
		Theme:                  "system",
		Language:               "en",
		DefaultDownloadFolder:  defaultDownloads,
		MaxConcurrentDownloads: 3,
		AutoCheckUpdates:       true,
		AutoAcquireLicense:     true,
		RememberCredentials:    true,
		SearchLimit:            15,
	}

	rows, err := s.db.Query("SELECT key, value FROM settings")
	if err != nil {
		return settings, nil
	}
	defer rows.Close()

	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err == nil {
			switch k {
			case "theme":
				settings.Theme = v
			case "language":
				settings.Language = v
			case "defaultDownloadFolder":
				settings.DefaultDownloadFolder = v
			case "maxConcurrentDownloads":
				var limit int
				fmt.Sscanf(v, "%d", &limit)
				if limit >= 1 && limit <= 10 {
					settings.MaxConcurrentDownloads = limit
				}
			case "autoCheckUpdates":
				settings.AutoCheckUpdates = (v == "true")
			case "autoAcquireLicense":
				settings.AutoAcquireLicense = (v == "true")
			case "rememberCredentials":
				settings.RememberCredentials = (v == "true")
			case "keychainPassphrase":
				settings.KeychainPassphrase = v
			case "searchLimit":
				var limit int
				fmt.Sscanf(v, "%d", &limit)
				if limit > 0 {
					settings.SearchLimit = limit
				}
			}
		}
	}

	return settings, nil
}

func (s *sqliteStorage) SaveSettings(settings models.AppSettings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	pairs := map[string]string{
		"theme":                  settings.Theme,
		"language":               settings.Language,
		"defaultDownloadFolder":  settings.DefaultDownloadFolder,
		"maxConcurrentDownloads": fmt.Sprintf("%d", settings.MaxConcurrentDownloads),
		"autoCheckUpdates":       fmt.Sprintf("%t", settings.AutoCheckUpdates),
		"autoAcquireLicense":     fmt.Sprintf("%t", settings.AutoAcquireLicense),
		"rememberCredentials":    fmt.Sprintf("%t", settings.RememberCredentials),
		"keychainPassphrase":     settings.KeychainPassphrase,
		"searchLimit":            fmt.Sprintf("%d", settings.SearchLimit),
	}

	for k, v := range pairs {
		_, err := tx.Exec("INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value", k, v)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ----------------- App Cache -----------------

func (s *sqliteStorage) CacheAppMetadata(app models.AppMetadata) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(app)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO app_cache (bundle_id, data_json, cached_at)
		VALUES (?, ?, ?)
		ON CONFLICT(bundle_id) DO UPDATE SET
			data_json = excluded.data_json,
			cached_at = excluded.cached_at;
	`, app.BundleID, string(data), time.Now())
	return err
}

func (s *sqliteStorage) GetCachedAppMetadata(bundleID string) (*models.AppMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var dataJSON string
	var cachedAt time.Time
	err := s.db.QueryRow("SELECT data_json, cached_at FROM app_cache WHERE bundle_id = ?", bundleID).Scan(&dataJSON, &cachedAt)
	if err != nil {
		return nil, err
	}

	// Cache valid for 24 hours
	if time.Since(cachedAt) > 24*time.Hour {
		return nil, nil
	}

	var app models.AppMetadata
	if err := json.Unmarshal([]byte(dataJSON), &app); err != nil {
		return nil, err
	}

	return &app, nil
}

func (s *sqliteStorage) ClearAppCache() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM app_cache")
	return err
}

func (s *sqliteStorage) GetCacheSizeBytes() (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fi, err := os.Stat(s.dbPath)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

// ----------------- Logs -----------------

func (s *sqliteStorage) AddLog(level, message, context string) (*models.LogEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	res, err := s.db.Exec("INSERT INTO logs (timestamp, level, message, context) VALUES (?, ?, ?, ?)", now, level, message, context)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	return &models.LogEntry{
		ID:        id,
		Timestamp: now,
		Level:     level,
		Message:   message,
		Context:   context,
	}, nil
}

func (s *sqliteStorage) GetRecentLogs(limit int) ([]models.LogEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 {
		limit = 100
	}

	rows, err := s.db.Query("SELECT id, timestamp, level, message, context FROM logs ORDER BY id DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.LogEntry
	for rows.Next() {
		var item models.LogEntry
		var ctx sql.NullString
		if err := rows.Scan(&item.ID, &item.Timestamp, &item.Level, &item.Message, &ctx); err != nil {
			return nil, err
		}
		if ctx.Valid {
			item.Context = ctx.String
		}
		list = append(list, item)
	}

	// Reverse to chronological order
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	return list, nil
}

func (s *sqliteStorage) ClearLogs() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM logs")
	return err
}
