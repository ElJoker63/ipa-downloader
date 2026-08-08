package apple

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/byteness/keyring"
	cookiejar "github.com/juju/persistent-cookiejar"
	"github.com/majd/ipatool/v2/backend/models"
	"github.com/majd/ipatool/v2/pkg/appstore"
	"github.com/majd/ipatool/v2/pkg/keychain"
	"github.com/majd/ipatool/v2/pkg/util/machine"
	"github.com/majd/ipatool/v2/pkg/util/operatingsystem"
)

const (
	KeychainServiceName = "ipatool"
	ConfigDirectoryName = ".ipatool"
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
			return "ipatool-default-keychain-passphrase", nil
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

func (c *client) GetAccount() (*models.AccountProfile, error) {
	info, err := c.appstore.AccountInfo()
	if err != nil {
		return &models.AccountProfile{IsLoggedIn: false}, nil
	}

	country := extractCountryFromStorefront(info.Account.StoreFront)

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

	country := extractCountryFromStorefront(out.Account.StoreFront)

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
	for _, vid := range listOut.ExternalVersionIdentifiers {
		meta, err := c.appstore.GetVersionMetadata(appstore.GetVersionMetadataInput{
			Account:   accInfo.Account,
			App:       appstore.App{ID: app.ID, BundleID: app.BundleID},
			VersionID: vid,
		})
		if err != nil {
			versions = append(versions, models.VersionInfo{
				ExternalVersionID: vid,
				DisplayVersion:    vid,
			})
			continue
		}

		versions = append(versions, models.VersionInfo{
			ExternalVersionID: vid,
			DisplayVersion:    meta.DisplayVersion,
			ReleaseDate:       meta.ReleaseDate,
			FormattedDate:     meta.ReleaseDate.Format("2006-01-02"),
		})
	}

	return versions, nil
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

func extractCountryFromStorefront(sf string) string {
	parts := strings.Split(sf, "-")
	if len(parts) > 0 {
		code := strings.TrimSpace(parts[0])
		return code
	}
	return "US"
}

func convertAppToMetadata(a appstore.App) models.AppMetadata {
	formattedPrice := "Free"
	if a.Price > 0 {
		formattedPrice = fmt.Sprintf("$%.2f", a.Price)
	}

	return models.AppMetadata{
		ID:             a.ID,
		BundleID:       a.BundleID,
		Name:           a.Name,
		Version:        a.Version,
		Price:          a.Price,
		FormattedPrice: formattedPrice,
	}
}
