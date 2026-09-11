package apple

import (
	"encoding/json"
	"errors"
	"fmt"
	gohttp "net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ElJoker63/ipa-downloader/v2/backend/models"
	"github.com/ElJoker63/ipa-downloader/v2/pkg/appstore"
	"github.com/ElJoker63/ipa-downloader/v2/pkg/keychain"
	"github.com/ElJoker63/ipa-downloader/v2/pkg/util/machine"
	"github.com/ElJoker63/ipa-downloader/v2/pkg/util/operatingsystem"
	"github.com/ElJoker63/ipa-downloader/v2/pkg/util/storefront"
	"github.com/byteness/keyring"
	cookiejar "github.com/juju/persistent-cookiejar"
)

const (
	KeychainServiceName = "ipa-downloader"
	ConfigDirectoryName = ".ipa-downloader"
	CookieJarFileName   = "cookiejar"
)

// Client encapsulates all interactions with Apple App Store APIs.
type Client interface {
	GetAppStore() appstore.AppStore
	GetAccount() (*models.AccountProfile, error)
	Login(email, password, authCode, endpoint string) (*models.AccountProfile, error)
	Revoke() error
	GetBag() (appstore.BagOutput, error)
	Search(term string, platform models.Platform, limit int64) ([]models.AppMetadata, error)
	Lookup(bundleID string, platform models.Platform) (*models.AppMetadata, error)
	Purchase(app models.AppMetadata) error
	ListVersions(app models.AppMetadata) ([]models.VersionInfo, error)
	GetPurchasedApps(page, limit int) (*models.PurchasedAppsOutput, error)
	GetKeychain() keychain.Keychain
}

type client struct {
	appstore appstore.AppStore
	keychain keychain.Keychain
	machine  machine.Machine
	os       operatingsystem.OperatingSystem
}

// NewClient initializes the Apple API client with persistent cookiejar and secure keychain.
func NewClient(passphrase string) (Client, error) {
	osUtil := operatingsystem.New()
	mach := machine.New(machine.Args{OS: osUtil})

	configDir := filepath.Join(mach.HomeDirectory(), ConfigDirectoryName)
	_ = osUtil.MkdirAll(configDir, 0700)

	cookieJar, err := cookiejar.New(&cookiejar.Options{
		Filename: filepath.Join(configDir, CookieJarFileName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cookiejar: %w", err)
	}

	ring, err := keyring.Open(keyring.Config{
		AllowedBackends: []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.SecretServiceBackend,
			keyring.FileBackend,
		},
		ServiceName: KeychainServiceName,
		FileDir:     configDir,
		FilePasswordFunc: func(s string) (string, error) {
			if passphrase != "" {
				return passphrase, nil
			}
			return "ipa-downloader-default-keychain-passphrase", nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	kc := keychain.New(keychain.Args{Keyring: ring})
	store := appstore.NewAppStore(appstore.Args{
		CookieJar:       cookieJar,
		OperatingSystem: osUtil,
		Keychain:        kc,
		Machine:         mach,
	})

	return &client{
		appstore: store,
		keychain: kc,
		machine:  mach,
		os:       osUtil,
	}, nil
}

func (c *client) GetAppStore() appstore.AppStore {
	return c.appstore
}

func (c *client) GetKeychain() keychain.Keychain {
	return c.keychain
}

func (c *client) GetAccount() (*models.AccountProfile, error) {
	info, err := c.appstore.AccountInfo()
	if err != nil {
		return &models.AccountProfile{IsLoggedIn: false}, nil
	}

	country := storefront.Format(info.Account.StoreFront)

	return &models.AccountProfile{
		Name:                info.Account.Name,
		Email:               info.Account.Email,
		StoreFront:          info.Account.StoreFront,
		StoreFrontCountry:   country,
		DirectoryServicesID: info.Account.DirectoryServicesID,
		Pod:                 info.Account.Pod,
		IsLoggedIn:          true,
	}, nil
}

func (c *client) Login(email, password, authCode, endpoint string) (*models.AccountProfile, error) {
	if endpoint == "" {
		bag, err := c.appstore.Bag(appstore.BagInput{})
		if err == nil {
			endpoint = bag.AuthEndpoint
		}
	}

	out, err := c.appstore.Login(appstore.LoginInput{
		Email:    email,
		Password: password,
		AuthCode: authCode,
		Endpoint: endpoint,
	})
	if err != nil {
		return nil, err
	}

	country := storefront.Format(out.Account.StoreFront)

	return &models.AccountProfile{
		Name:                out.Account.Name,
		Email:               out.Account.Email,
		StoreFront:          out.Account.StoreFront,
		StoreFrontCountry:   country,
		DirectoryServicesID: out.Account.DirectoryServicesID,
		Pod:                 out.Account.Pod,
		IsLoggedIn:          true,
	}, nil
}

func (c *client) Revoke() error {
	return c.appstore.Revoke()
}

func (c *client) GetBag() (appstore.BagOutput, error) {
	return c.appstore.Bag(appstore.BagInput{})
}

func (c *client) Search(term string, platform models.Platform, limit int64) ([]models.AppMetadata, error) {
	var account appstore.Account
	accInfo, err := c.appstore.AccountInfo()
	if err == nil && accInfo.Account.StoreFront != "" {
		account = accInfo.Account
	} else {
		account = appstore.Account{
			StoreFront: "143441-1,29", // Default US Storefront
		}
	}

	appPlatform, err := parseApplePlatform(platform)
	if err != nil {
		appPlatform = appstore.PlatformIPhone
	}

	if limit <= 0 {
		limit = 15
	}

	out, err := c.appstore.Search(appstore.SearchInput{
		Account:  account,
		Term:     term,
		Limit:    limit,
		Platform: appPlatform,
	})
	if err != nil {
		return nil, err
	}

	var results []models.AppMetadata
	for _, raw := range out.Results {
		meta := convertAppToMetadata(raw)
		results = append(results, meta)
	}

	return results, nil
}

func (c *client) Lookup(bundleID string, platform models.Platform) (*models.AppMetadata, error) {
	var account appstore.Account
	accInfo, err := c.appstore.AccountInfo()
	if err == nil && accInfo.Account.StoreFront != "" {
		account = accInfo.Account
	} else {
		account = appstore.Account{
			StoreFront: "143441-1,29", // Default US Storefront
		}
	}

	appPlatform, err := parseApplePlatform(platform)
	if err != nil {
		appPlatform = appstore.PlatformIPhone
	}

	out, err := c.appstore.Lookup(appstore.LookupInput{
		Account:  account,
		BundleID: bundleID,
		Platform: appPlatform,
	})
	if err != nil {
		return nil, err
	}

	meta := convertAppToMetadata(out.App)
	return &meta, nil
}

func (c *client) Purchase(app models.AppMetadata) error {
	accInfo, err := c.appstore.AccountInfo()
	if err != nil {
		return fmt.Errorf("authentication required: %w", err)
	}

	err = c.appstore.Purchase(appstore.PurchaseInput{
		Account: accInfo.Account,
		App: appstore.App{
			ID:       app.ID,
			BundleID: app.BundleID,
			Name:     app.Name,
			Price:    app.Price,
			Version:  app.Version,
		},
		StoreFront: accInfo.Account.StoreFront,
		Platform:   appstore.PlatformIPhone,
	})
	if err != nil && !errors.Is(err, appstore.ErrLicenseAlreadyExists) {
		return err
	}

	return nil
}

func (c *client) ListVersions(app models.AppMetadata) ([]models.VersionInfo, error) {
	accInfo, err := c.appstore.AccountInfo()
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	listOut, err := c.appstore.ListVersions(appstore.ListVersionsInput{
		Account: accInfo.Account,
		App: appstore.App{
			ID:       app.ID,
			BundleID: app.BundleID,
			Name:     app.Name,
			Version:  app.Version,
			Price:    app.Price,
		},
	})
	if err != nil {
		return nil, err
	}

	var versions []models.VersionInfo
	// Return in reverse chronological order (newest on top, previous versions below)
	for i := len(listOut.ExternalVersionIdentifiers) - 1; i >= 0; i-- {
		vid := listOut.ExternalVersionIdentifiers[i]
		display := fmt.Sprintf("Build %s", vid)
		if vid == listOut.LatestExternalVersionID && app.Version != "" {
			display = fmt.Sprintf("v%s (Build %s)", app.Version, vid)
		}

		versions = append(versions, models.VersionInfo{
			ExternalVersionID: vid,
			DisplayVersion:    display,
		})
	}

	return versions, nil
}

func (c *client) GetPurchasedApps(page, limit int) (*models.PurchasedAppsOutput, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > appstore.MaxOwnedAppsLimit {
		limit = appstore.MaxOwnedAppsLimit
	}

	accInfo, err := c.appstore.AccountInfo()
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	out, err := c.appstore.OwnedApps(appstore.OwnedAppsInput{
		Account: accInfo.Account,
		Page:    page,
		Limit:   limit,
	})
	if err != nil {
		return nil, err
	}

	countryCode := storefront.GetCountry(accInfo.Account.StoreFront).ID
	if countryCode == "" || countryCode == "Unknown" {
		countryCode = "US"
	}

	// Batch lookup artwork icons & metadata from iTunes for the current page
	var idStrs []string
	for _, app := range out.Results {
		if app.ID > 0 {
			idStrs = append(idStrs, strconv.FormatInt(app.ID, 10))
		}
	}

	appLookupMap := make(map[int64]models.AppMetadata)
	if len(idStrs) > 0 {
		lookupURL := fmt.Sprintf("https://%s%s?id=%s&country=%s", "itunes.apple.com", "/lookup", strings.Join(idStrs, ","), countryCode)
		req, rErr := gohttp.NewRequest("GET", lookupURL, nil)
		if rErr == nil {
			req.Header.Set("User-Agent", "Mozilla/5.0")
			resp, doErr := gohttp.DefaultClient.Do(req)
			if doErr == nil && resp.StatusCode == gohttp.StatusOK {
				defer resp.Body.Close()
				var searchRes struct {
					ResultCount int `json:"resultCount"`
					Results     []struct {
						TrackID        int64   `json:"trackId"`
						BundleID       string  `json:"bundleId"`
						TrackName      string  `json:"trackName"`
						ArtistName     string  `json:"artistName"`
						ArtworkURL512  string  `json:"artworkUrl512"`
						ArtworkURL100  string  `json:"artworkUrl100"`
						ArtworkURL60   string  `json:"artworkUrl60"`
						PrimaryGenre   string  `json:"primaryGenreName"`
						Price          float64 `json:"price"`
						FormattedPrice string  `json:"formattedPrice"`
						Version        string  `json:"version"`
					} `json:"results"`
				}
				if jsonErr := json.NewDecoder(resp.Body).Decode(&searchRes); jsonErr == nil {
					for _, item := range searchRes.Results {
						artwork := item.ArtworkURL512
						if artwork == "" {
							artwork = item.ArtworkURL100
						}
						if artwork == "" {
							artwork = item.ArtworkURL60
						}
						appLookupMap[item.TrackID] = models.AppMetadata{
							ID:             item.TrackID,
							BundleID:       item.BundleID,
							Name:           item.TrackName,
							Developer:      item.ArtistName,
							ArtworkURL:     artwork,
							ArtworkURL512:  item.ArtworkURL512,
							PrimaryGenre:   item.PrimaryGenre,
							Price:          item.Price,
							FormattedPrice: item.FormattedPrice,
						}
					}
				}
			}
		}

		// Fallback for any missing items without country filter
		var missingIDs []string
		for _, idStr := range idStrs {
			id, _ := strconv.ParseInt(idStr, 10, 64)
			if _, ok := appLookupMap[id]; !ok {
				missingIDs = append(missingIDs, idStr)
			}
		}
		if len(missingIDs) > 0 {
			globalURL := fmt.Sprintf("https://%s%s?id=%s", "itunes.apple.com", "/lookup", strings.Join(missingIDs, ","))
			if gReq, gErr := gohttp.NewRequest("GET", globalURL, nil); gErr == nil {
				gReq.Header.Set("User-Agent", "Mozilla/5.0")
				if gResp, gDoErr := gohttp.DefaultClient.Do(gReq); gDoErr == nil && gResp.StatusCode == gohttp.StatusOK {
					defer gResp.Body.Close()
					var gSearchRes struct {
						ResultCount int `json:"resultCount"`
						Results     []struct {
							TrackID        int64   `json:"trackId"`
							BundleID       string  `json:"bundleId"`
							TrackName      string  `json:"trackName"`
							ArtistName     string  `json:"artistName"`
							ArtworkURL512  string  `json:"artworkUrl512"`
							ArtworkURL100  string  `json:"artworkUrl100"`
							ArtworkURL60   string  `json:"artworkUrl60"`
							PrimaryGenre   string  `json:"primaryGenreName"`
							Price          float64 `json:"price"`
							FormattedPrice string  `json:"formattedPrice"`
							Version        string  `json:"version"`
						} `json:"results"`
					}
					if jsonErr := json.NewDecoder(gResp.Body).Decode(&gSearchRes); jsonErr == nil {
						for _, item := range gSearchRes.Results {
							artwork := item.ArtworkURL512
							if artwork == "" {
								artwork = item.ArtworkURL100
							}
							if artwork == "" {
								artwork = item.ArtworkURL60
							}
							appLookupMap[item.TrackID] = models.AppMetadata{
								ID:             item.TrackID,
								BundleID:       item.BundleID,
								Name:           item.TrackName,
								Developer:      item.ArtistName,
								ArtworkURL:     artwork,
								ArtworkURL512:  item.ArtworkURL512,
								PrimaryGenre:   item.PrimaryGenre,
								Price:          item.Price,
								FormattedPrice: item.FormattedPrice,
							}
						}
					}
				}
			}
		}
	}

	var results []models.AppMetadata
	for _, raw := range out.Results {
		meta := convertAppToMetadata(raw)
		if enriched, ok := appLookupMap[raw.ID]; ok {
			if enriched.ArtworkURL != "" {
				meta.ArtworkURL = enriched.ArtworkURL
				meta.ArtworkURL512 = enriched.ArtworkURL512
			}
			if enriched.Developer != "" {
				meta.Developer = enriched.Developer
			}
			if enriched.PrimaryGenre != "" {
				meta.PrimaryGenre = enriched.PrimaryGenre
			}
			if enriched.FormattedPrice != "" {
				meta.FormattedPrice = enriched.FormattedPrice
			}
		}
		results = append(results, meta)
	}

	return &models.PurchasedAppsOutput{
		Count:      out.Count,
		TotalCount: out.TotalCount,
		Page:       out.Page,
		Results:    results,
	}, nil
}

func parseApplePlatform(p models.Platform) (appstore.Platform, error) {
	switch p {
	case models.PlatformiPadOS:
		return appstore.PlatformIPad, nil
	case models.PlatformtvOS:
		return appstore.PlatformAppleTV, nil
	case models.PlatformIOS, "":
		return appstore.PlatformIPhone, nil
	default:
		return appstore.ParsePlatform(string(p))
	}
}

func convertAppToMetadata(a appstore.App) models.AppMetadata {
	formattedPrice := "Free"
	if a.Price > 0 {
		formattedPrice = fmt.Sprintf("$%.2f", a.Price)
	}

	artwork := a.ArtworkURL512
	if artwork == "" {
		artwork = a.ArtworkURL100
	}
	if artwork == "" {
		artwork = a.ArtworkURL60
	}

	purchaseDateStr := ""
	if !a.PurchaseDate.IsZero() {
		purchaseDateStr = a.PurchaseDate.Format(time.RFC3339)
	}

	return models.AppMetadata{
		ID:             a.ID,
		BundleID:       a.BundleID,
		Name:           a.Name,
		Version:        a.Version,
		Price:          a.Price,
		FormattedPrice: formattedPrice,
		ArtworkURL:     artwork,
		PurchaseDate:   purchaseDateStr,
	}
}
