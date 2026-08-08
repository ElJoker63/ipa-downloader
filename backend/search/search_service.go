package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ElJoker63/ipatool-1/v2/backend/apple"
	"github.com/ElJoker63/ipatool-1/v2/backend/events"
	"github.com/ElJoker63/ipatool-1/v2/backend/models"
	"github.com/ElJoker63/ipatool-1/v2/backend/storage"
)

// SearchService handles searching for apps, retrieving rich metadata, and managing search history.
type SearchService interface {
	Search(term string, platform string, limit int) ([]models.AppMetadata, error)
	Lookup(bundleID string, platform string) (*models.AppMetadata, error)
	GetAppDetails(appID int64, bundleID string, platform string) (*models.AppDetailsOutput, error)
	GetSearchHistory(limit int) ([]models.SearchHistoryItem, error)
	ClearSearchHistory() error
}

type searchService struct {
	appleClient apple.Client
	storage     storage.Storage
	emitter     events.Emitter
	httpClient  *http.Client
	mu          sync.RWMutex
}

// NewSearchService creates a new search service.
func NewSearchService(client apple.Client, store storage.Storage, emitter events.Emitter) SearchService {
	return &searchService{
		appleClient: client,
		storage:     store,
		emitter:     emitter,
		httpClient:  &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *searchService) Search(term string, platform string, limit int) ([]models.AppMetadata, error) {
	if term == "" {
		return []models.AppMetadata{}, nil
	}

	s.emitter.EmitLog("INFO", fmt.Sprintf("Searching for '%s' on %s...", term, platform), "SearchService")

	if limit <= 0 {
		limit = 15
	}

	// 1. Search via client with fallback to direct iTunes Search API
	results, err := s.appleClient.Search(term, models.Platform(platform), int64(limit))
	if err != nil || len(results) == 0 {
		s.emitter.EmitLog("INFO", fmt.Sprintf("Querying direct iTunes Search API for '%s'...", term), "SearchService")
		directResults, directErr := s.directITunesSearch(term, platform, limit)
		if directErr == nil && len(directResults) > 0 {
			results = directResults
		} else if err != nil {
			s.emitter.EmitLog("ERROR", fmt.Sprintf("Search failed: %v", err), "SearchService")
			return nil, err
		}
	}

	// 2. Enhance results with high-resolution artwork and developer info
	enhancedResults := make([]models.AppMetadata, 0, len(results))
	for _, app := range results {
		// Fetch rich metadata if available
		richMeta := s.enrichMetadata(app, platform)

		// Check if it is a favorite
		isFav, _ := s.storage.IsFavorite(richMeta.ID)
		richMeta.IsFavorite = isFav

		enhancedResults = append(enhancedResults, richMeta)

		// Cache app metadata for instant detail lookup
		_ = s.storage.CacheAppMetadata(richMeta)
	}

	// 3. Record search history in background
	_ = s.storage.AddSearchHistory(term, platform, len(enhancedResults))

	s.emitter.EmitLog("SUCCESS", fmt.Sprintf("Found %d results for '%s'", len(enhancedResults), term), "SearchService")
	return enhancedResults, nil
}

func (s *searchService) directITunesSearch(term string, platform string, limit int) ([]models.AppMetadata, error) {
	entity := "software"
	switch platform {
	case "ipados":
		entity = "iPadSoftware"
	case "tvos":
		entity = "tvSoftware"
	}

	searchURL := fmt.Sprintf("https://itunes.apple.com/search?media=software&entity=%s&term=%s&country=US&limit=%d", entity, url.QueryEscape(term), limit)
	resp, err := s.httpClient.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iTunes API responded with HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data iTunesLookupResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	results := make([]models.AppMetadata, 0, len(data.Results))
	for _, item := range data.Results {
		artwork := item.ArtworkURL512
		if artwork == "" {
			artwork = item.ArtworkURL100
		}
		if artwork == "" {
			artwork = item.ArtworkURL60
		}

		priceStr := item.FormattedPrice
		if priceStr == "" {
			if item.Price == 0 {
				priceStr = "Free"
			} else {
				priceStr = fmt.Sprintf("$%.2f", item.Price)
			}
		}

		var sizeBytes int64
		fmt.Sscanf(item.FileSizeBytes, "%d", &sizeBytes)

		results = append(results, models.AppMetadata{
			ID:                    item.TrackID,
			BundleID:              item.BundleID,
			Name:                  item.TrackName,
			Developer:             item.ArtistName,
			DeveloperID:           item.ArtistID,
			Version:               item.Version,
			Price:                 item.Price,
			FormattedPrice:        priceStr,
			Currency:              item.Currency,
			ArtworkURL:            artwork,
			ArtworkURL512:         item.ArtworkURL512,
			Screenshots:           item.ScreenshotUrls,
			IpadScreenshots:       item.IpadScreenshotUrls,
			Description:           item.Description,
			ReleaseNotes:          item.ReleaseNotes,
			ReleaseDate:           item.ReleaseDate,
			CurrentVersionDate:    item.CurrentVersionDate,
			MinimumOSVersion:      item.MinimumOSVersion,
			FileSizeBytes:         sizeBytes,
			FormattedSize:         formatByteSize(sizeBytes),
			AverageUserRating:     item.AverageUserRating,
			UserRatingCount:       item.UserRatingCount,
			ContentAdvisoryRating: item.ContentAdvisoryRating,
			Genres:                item.Genres,
			PrimaryGenre:          item.PrimaryGenreName,
			SupportedPlatforms:    item.SupportedDevices,
		})
	}

	return results, nil
}

func (s *searchService) Lookup(bundleID string, platform string) (*models.AppMetadata, error) {
	s.emitter.EmitLog("INFO", fmt.Sprintf("Looking up bundle ID '%s'...", bundleID), "SearchService")

	// Check cache
	if cached, err := s.storage.GetCachedAppMetadata(bundleID); err == nil && cached != nil {
		isFav, _ := s.storage.IsFavorite(cached.ID)
		cached.IsFavorite = isFav
		return cached, nil
	}

	app, err := s.appleClient.Lookup(bundleID, models.Platform(platform))
	if err != nil {
		s.emitter.EmitLog("ERROR", fmt.Sprintf("Lookup failed: %v", err), "SearchService")
		return nil, err
	}

	rich := s.enrichMetadata(*app, platform)
	isFav, _ := s.storage.IsFavorite(rich.ID)
	rich.IsFavorite = isFav

	_ = s.storage.CacheAppMetadata(rich)
	return &rich, nil
}

func (s *searchService) GetAppDetails(appID int64, bundleID string, platform string) (*models.AppDetailsOutput, error) {
	s.emitter.EmitLog("INFO", fmt.Sprintf("Fetching full app details for %s (ID: %d)...", bundleID, appID), "SearchService")

	var meta *models.AppMetadata
	var err error

	if bundleID != "" {
		meta, err = s.Lookup(bundleID, platform)
	}
	if meta == nil || err != nil {
		// Fallback query by ID
		meta = &models.AppMetadata{
			ID:       appID,
			BundleID: bundleID,
		}
	}

	// Fetch full iTunes metadata (screenshots, full description, minimum OS, advisory)
	detailedMeta := s.fetchFullITunesMetadata(meta.ID, meta.BundleID)
	if detailedMeta != nil {
		meta = detailedMeta
	}

	isFav, _ := s.storage.IsFavorite(meta.ID)
	meta.IsFavorite = isFav

	// Fetch version history if available
	var versions []models.VersionInfo
	if vList, err := s.appleClient.ListVersions(*meta); err == nil && len(vList) > 0 {
		versions = vList
	} else if meta.Version != "" {
		versions = append(versions, models.VersionInfo{
			ExternalVersionID: meta.Version,
			DisplayVersion:    fmt.Sprintf("v%s (Current)", meta.Version),
		})
	}

	return &models.AppDetailsOutput{
		Metadata:       *meta,
		VersionHistory: versions,
		IsFavorite:     isFav,
	}, nil
}

func (s *searchService) GetSearchHistory(limit int) ([]models.SearchHistoryItem, error) {
	return s.storage.GetSearchHistory(limit)
}

func (s *searchService) ClearSearchHistory() error {
	return s.storage.ClearSearchHistory()
}

// enrichMetadata queries the iTunes search API to enrich app metadata with artwork, description, etc.
func (s *searchService) enrichMetadata(app models.AppMetadata, platform string) models.AppMetadata {
	full := s.fetchFullITunesMetadata(app.ID, app.BundleID)
	if full != nil {
		if full.Name == "" {
			full.Name = app.Name
		}
		if full.Version == "" {
			full.Version = app.Version
		}
		return *full
	}

	if app.ArtworkURL == "" {
		app.ArtworkURL = "https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png"
	}
	return app
}

type iTunesLookupResponse struct {
	ResultCount int `json:"resultCount"`
	Results     []struct {
		TrackID               int64    `json:"trackId"`
		BundleID              string   `json:"bundleId"`
		TrackName             string   `json:"trackName"`
		ArtistName            string   `json:"artistName"`
		ArtistID              int64    `json:"artistId"`
		Version               string   `json:"version"`
		Price                 float64  `json:"price"`
		FormattedPrice        string   `json:"formattedPrice"`
		Currency              string   `json:"currency"`
		ArtworkURL60          string   `json:"artworkUrl60"`
		ArtworkURL100         string   `json:"artworkUrl100"`
		ArtworkURL512         string   `json:"artworkUrl512"`
		ScreenshotUrls        []string `json:"screenshotUrls"`
		IpadScreenshotUrls    []string `json:"ipadScreenshotUrls"`
		Description           string   `json:"description"`
		ReleaseNotes          string   `json:"releaseNotes"`
		ReleaseDate           string   `json:"releaseDate"`
		CurrentVersionDate    string   `json:"currentVersionReleaseDate"`
		MinimumOSVersion      string   `json:"minimumOsVersion"`
		FileSizeBytes         string   `json:"fileSizeBytes"`
		AverageUserRating     float64  `json:"averageUserRating"`
		UserRatingCount       int      `json:"userRatingCount"`
		ContentAdvisoryRating string   `json:"contentAdvisoryRating"`
		Genres                []string `json:"genres"`
		PrimaryGenreName      string   `json:"primaryGenreName"`
		SupportedDevices      []string `json:"supportedDevices"`
	} `json:"results"`
}

func (s *searchService) fetchFullITunesMetadata(appID int64, bundleID string) *models.AppMetadata {
	endpoint := ""
	if appID > 0 {
		endpoint = fmt.Sprintf("https://itunes.apple.com/lookup?id=%d", appID)
	} else if bundleID != "" {
		endpoint = fmt.Sprintf("https://itunes.apple.com/lookup?bundleId=%s", url.QueryEscape(bundleID))
	} else {
		return nil
	}

	resp, err := s.httpClient.Get(endpoint)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var data iTunesLookupResponse
	if err := json.Unmarshal(body, &data); err != nil || data.ResultCount == 0 {
		return nil
	}

	item := data.Results[0]
	artwork := item.ArtworkURL512
	if artwork == "" {
		artwork = item.ArtworkURL100
	}
	if artwork == "" {
		artwork = item.ArtworkURL60
	}

	priceStr := item.FormattedPrice
	if priceStr == "" {
		if item.Price == 0 {
			priceStr = "Free"
		} else {
			priceStr = fmt.Sprintf("$%.2f", item.Price)
		}
	}

	var sizeBytes int64
	fmt.Sscanf(item.FileSizeBytes, "%d", &sizeBytes)
	formattedSize := formatByteSize(sizeBytes)

	return &models.AppMetadata{
		ID:                    item.TrackID,
		BundleID:              item.BundleID,
		Name:                  item.TrackName,
		Developer:             item.ArtistName,
		DeveloperID:           item.ArtistID,
		Version:               item.Version,
		Price:                 item.Price,
		FormattedPrice:        priceStr,
		Currency:              item.Currency,
		ArtworkURL:            artwork,
		ArtworkURL512:         item.ArtworkURL512,
		Screenshots:           item.ScreenshotUrls,
		IpadScreenshots:       item.IpadScreenshotUrls,
		Description:           item.Description,
		ReleaseNotes:          item.ReleaseNotes,
		ReleaseDate:           item.ReleaseDate,
		CurrentVersionDate:    item.CurrentVersionDate,
		MinimumOSVersion:      item.MinimumOSVersion,
		FileSizeBytes:         sizeBytes,
		FormattedSize:         formattedSize,
		AverageUserRating:     item.AverageUserRating,
		UserRatingCount:       item.UserRatingCount,
		ContentAdvisoryRating: item.ContentAdvisoryRating,
		Genres:                item.Genres,
		PrimaryGenre:          item.PrimaryGenreName,
		SupportedPlatforms:    item.SupportedDevices,
	}
}

func formatByteSize(bytes int64) string {
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
