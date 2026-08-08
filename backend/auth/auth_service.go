package auth

import (
	"errors"
	"fmt"
	"sync"

	"github.com/majd/ipa-downloader/v2/backend/apple"
	"github.com/majd/ipa-downloader/v2/backend/events"
	"github.com/majd/ipa-downloader/v2/backend/models"
	"github.com/majd/ipa-downloader/v2/backend/storage"
	"github.com/majd/ipa-downloader/v2/pkg/appstore"
)

// AuthService manages Apple ID login, session state, and credentials.
type AuthService interface {
	GetAccount() (*models.AccountProfile, error)
	Login(email, password, authCode string, remember bool) (*models.AccountProfile, error)
	Logout() error
	GetStatus() string // "Connected", "Connecting", "Not Connected"
}

type authService struct {
	appleClient apple.Client
	storage     storage.Storage
	emitter     events.Emitter
	status      string
	mu          sync.RWMutex
}

// NewAuthService creates a new authentication service.
func NewAuthService(client apple.Client, store storage.Storage, emitter events.Emitter) AuthService {
	s := &authService{
		appleClient: client,
		storage:     store,
		emitter:     emitter,
		status:      "Not Connected",
	}

	// Check if already authenticated on startup
	if acc, err := s.GetAccount(); err == nil && acc.IsLoggedIn {
		s.status = "Connected"
	}

	return s
}

func (s *authService) GetAccount() (*models.AccountProfile, error) {
	acc, err := s.appleClient.GetAccount()
	if err != nil {
		s.mu.Lock()
		s.status = "Not Connected"
		s.mu.Unlock()
		return &models.AccountProfile{IsLoggedIn: false}, nil
	}

	s.mu.Lock()
	if acc.IsLoggedIn {
		s.status = "Connected"
	} else {
		s.status = "Not Connected"
	}
	s.mu.Unlock()

	return acc, nil
}

func (s *authService) Login(email, password, authCode string, remember bool) (*models.AccountProfile, error) {
	s.mu.Lock()
	s.status = "Connecting"
	s.mu.Unlock()

	s.emitter.EmitLog("INFO", fmt.Sprintf("Authenticating Apple ID: %s...", email), "AuthService")
	s.emitter.Emit(events.EventAuthStatus, map[string]string{"status": "Connecting"})

	bag, err := s.appleClient.GetBag()
	endpoint := ""
	if err == nil {
		endpoint = bag.AuthEndpoint
	}

	acc, err := s.appleClient.Login(email, password, authCode, endpoint)
	if err != nil {
		s.mu.Lock()
		s.status = "Not Connected"
		s.mu.Unlock()

		s.emitter.EmitLog("ERROR", fmt.Sprintf("Authentication failed: %v", err), "AuthService")
		s.emitter.Emit(events.EventAuthStatus, map[string]string{"status": "Not Connected", "error": err.Error()})

		if errors.Is(err, appstore.ErrAuthCodeRequired) {
			return nil, fmt.Errorf("2FA verification code is required")
		}
		return nil, err
	}

	s.mu.Lock()
	s.status = "Connected"
	s.mu.Unlock()

	// Update settings if remember requested
	if remember {
		if settings, err := s.storage.GetSettings(); err == nil {
			settings.RememberCredentials = true
			_ = s.storage.SaveSettings(*settings)
		}
	}

	s.emitter.EmitLog("SUCCESS", fmt.Sprintf("Successfully connected as %s (%s)", acc.Name, acc.Email), "AuthService")
	s.emitter.Emit(events.EventAuthStatus, map[string]string{"status": "Connected"})
	s.emitter.Emit(events.EventAuthAccount, acc)
	s.emitter.EmitNotification("Connected", fmt.Sprintf("Logged in as %s", acc.Email), "success")

	return acc, nil
}

func (s *authService) Logout() error {
	s.emitter.EmitLog("INFO", "Revoking App Store session...", "AuthService")

	err := s.appleClient.Revoke()
	if err != nil {
		s.emitter.EmitLog("WARN", fmt.Sprintf("Revoke warning: %v", err), "AuthService")
	}

	s.mu.Lock()
	s.status = "Not Connected"
	s.mu.Unlock()

	s.emitter.EmitLog("INFO", "Logged out successfully", "AuthService")
	s.emitter.Emit(events.EventAuthStatus, map[string]string{"status": "Not Connected"})
	s.emitter.Emit(events.EventAuthAccount, &models.AccountProfile{IsLoggedIn: false})
	s.emitter.EmitNotification("Logged Out", "App Store credentials revoked", "info")

	return nil
}

func (s *authService) GetStatus() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}
