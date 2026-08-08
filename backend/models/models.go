package models

import (
	"time"
)

// Platform represents the target Apple platform.
type Platform string

const (
	PlatformIOS    Platform = "ios"
	PlatformiPadOS Platform = "ipados"
	PlatformtvOS   Platform = "tvos"
)

// AccountProfile holds authenticated Apple ID account details.
type AccountProfile struct {
	Name                string `json:"name"`
	Email               string `json:"email"`
	StoreFront          string `json:"storeFront"`
	StoreFrontCountry   string `json:"storeFrontCountry"`
	DirectoryServicesID string `json:"directoryServicesId"`
	Pod                 string `json:"pod"`
	IsLoggedIn          bool   `json:"isLoggedIn"`
}

// AppMetadata represents rich metadata for an App Store application.
type AppMetadata struct {
	ID                    int64    `json:"id"`
	BundleID              string   `json:"bundleId"`
	Name                  string   `json:"name"`
	Developer             string   `json:"developer"`
	DeveloperID           int64    `json:"developerId,omitempty"`
	Version               string   `json:"version"`
	BuildNumber           string   `json:"buildNumber,omitempty"`
	Price                 float64  `json:"price"`
	FormattedPrice        string   `json:"formattedPrice"`
	Currency              string   `json:"currency,omitempty"`
	ArtworkURL            string   `json:"artworkUrl"`
	ArtworkURL512         string   `json:"artworkUrl512,omitempty"`
	Screenshots           []string `json:"screenshots"`
	IpadScreenshots       []string `json:"ipadScreenshots,omitempty"`
	Description           string   `json:"description"`
	ReleaseNotes          string   `json:"releaseNotes,omitempty"`
	ReleaseDate           string   `json:"releaseDate"`
	CurrentVersionDate    string   `json:"currentVersionDate,omitempty"`
	MinimumOSVersion      string   `json:"minimumOsVersion"`
	FileSizeBytes         int64    `json:"fileSizeBytes"`
	FormattedSize         string   `json:"formattedSize"`
	AverageUserRating     float64  `json:"averageUserRating"`
	UserRatingCount       int      `json:"userRatingCount"`
	ContentAdvisoryRating string   `json:"contentAdvisoryRating"`
	Genres                []string `json:"genres"`
	PrimaryGenre          string   `json:"primaryGenre"`
	SupportedPlatforms    []string `json:"supportedPlatforms"`
	IsFavorite            bool     `json:"isFavorite"`
}

// DownloadStatus represents the state of a download task.
type DownloadStatus string

const (
	DownloadStatusQueued      DownloadStatus = "queued"
	DownloadStatusDownloading DownloadStatus = "downloading"
	DownloadStatusPaused      DownloadStatus = "paused"
	DownloadStatusCompleted   DownloadStatus = "completed"
	DownloadStatusCancelled   DownloadStatus = "cancelled"
	DownloadStatusFailed      DownloadStatus = "failed"
)

// DownloadTask represents an active or historical download item.
type DownloadTask struct {
	ID                string         `json:"id"`
	AppID             int64          `json:"appId"`
	BundleID          string         `json:"bundleId"`
	AppName           string         `json:"appName"`
	Developer         string         `json:"developer"`
	Version           string         `json:"version"`
	ArtworkURL        string         `json:"artworkUrl"`
	DestinationPath   string         `json:"destinationPath"`
	Status            DownloadStatus `json:"status"`
	TotalBytes        int64          `json:"totalBytes"`
	DownloadedBytes   int64          `json:"downloadedBytes"`
	Progress          float64        `json:"progress"` // 0.0 - 100.0
	SpeedBytesPerSec  int64          `json:"speedBytesPerSec"`
	FormattedSpeed    string         `json:"formattedSpeed"`
	ETASeconds        int64          `json:"etaSeconds"`
	FormattedETA      string         `json:"formattedETA"`
	ExternalVersionID string         `json:"externalVersionId,omitempty"`
	Platform          string         `json:"platform"`
	Error             string         `json:"error,omitempty"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	CompletedAt       *time.Time     `json:"completedAt,omitempty"`
}

// FavoriteApp represents an app bookmarked by the user.
type FavoriteApp struct {
	AppID          int64     `json:"appId"`
	BundleID       string    `json:"bundleId"`
	Name           string    `json:"name"`
	Developer      string    `json:"developer"`
	Version        string    `json:"version"`
	Price          float64   `json:"price"`
	FormattedPrice string    `json:"formattedPrice"`
	ArtworkURL     string    `json:"artworkUrl"`
	PrimaryGenre   string    `json:"primaryGenre"`
	CreatedAt      time.Time `json:"createdAt"`
}

// SearchHistoryItem records a user search query.
type SearchHistoryItem struct {
	ID        int64     `json:"id"`
	Term      string    `json:"term"`
	Platform  string    `json:"platform"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"createdAt"`
}

// AppSettings defines configurable user preferences.
type AppSettings struct {
	Theme                  string `json:"theme"`                  // "dark", "light", "system"
	Language               string `json:"language"`               // "en", "es", "zh", etc.
	DefaultDownloadFolder  string `json:"defaultDownloadFolder"`  // path to save IPAs
	MaxConcurrentDownloads int    `json:"maxConcurrentDownloads"` // 1 - 10
	AutoCheckUpdates       bool   `json:"autoCheckUpdates"`
	AutoAcquireLicense     bool   `json:"autoAcquireLicense"` // whether to automatically purchase free license
	RememberCredentials    bool   `json:"rememberCredentials"`
	KeychainPassphrase     string `json:"keychainPassphrase,omitempty"`
	SearchLimit            int    `json:"searchLimit"`
}

// LogEntry represents an event emitted by the backend.
type LogEntry struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"` // "DEBUG", "INFO", "WARN", "ERROR", "SUCCESS"
	Message   string    `json:"message"`
	Context   string    `json:"context,omitempty"`
}

// VersionInfo represents an App Store version history item.
type VersionInfo struct {
	ExternalVersionID string    `json:"externalVersionId"`
	DisplayVersion    string    `json:"displayVersion"`
	ReleaseDate       time.Time `json:"releaseDate"`
	FormattedDate     string    `json:"formattedDate"`
}

// AppDetailsOutput returns comprehensive details for an app.
type AppDetailsOutput struct {
	Metadata       AppMetadata   `json:"metadata"`
	VersionHistory []VersionInfo `json:"versionHistory"`
	IsFavorite     bool          `json:"isFavorite"`
}
