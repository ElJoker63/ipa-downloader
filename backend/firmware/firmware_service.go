package firmware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/ElJoker63/ipa-downloader/v2/backend/events"
	"github.com/ElJoker63/ipa-downloader/v2/backend/models"
)

const (
	IPSWAPI = "https://api.ipsw.me/v2.1/firmwares.json"
)

type FirmwareService interface {
	GetDevices() ([]models.AppleHardware, error)
	GetDeviceDetails(identifier string) (*models.AppleHardware, error)
}

type firmwareService struct {
	emitter    events.Emitter
	cache      []models.AppleHardware
	lastUpdate time.Time
	mu         sync.RWMutex
}

func NewFirmwareService(emitter events.Emitter) FirmwareService {
	return &firmwareService{
		emitter: emitter,
	}
}

type ipswResponse struct {
	Devices map[string]struct {
		Name      string            `json:"name"`
		Platform  string            `json:"platform"`
		Firmwares []models.Firmware `json:"firmwares"`
	} `json:"devices"`
}

func (s *firmwareService) GetDevices() ([]models.AppleHardware, error) {
	s.mu.RLock()
	if len(s.cache) > 0 && time.Since(s.lastUpdate) < 6*time.Hour {
		defer s.mu.RUnlock()
		return s.cache, nil
	}
	s.mu.RUnlock()

	return s.refreshCache()
}

func (s *firmwareService) refreshCache() ([]models.AppleHardware, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.cache) > 0 && time.Since(s.lastUpdate) < 6*time.Hour {
		return s.cache, nil
	}

	s.emitter.EmitLog("INFO", "Fetching Apple firmware list from IPSW.me...", "FirmwareService")

	resp, err := http.Get(IPSWAPI)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch firmwares: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPSW.me API returned status %d", resp.StatusCode)
	}

	var data ipswResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode firmwares: %w", err)
	}

	var devices []models.AppleHardware
	for id, d := range data.Devices {
		devices = append(devices, models.AppleHardware{
			Identifier: id,
			Name:       d.Name,
			Platform:   d.Platform,
			Firmwares:  d.Firmwares,
		})
	}

	// Sort devices by name
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Name < devices[j].Name
	})

	s.cache = devices
	s.lastUpdate = time.Now()
	s.emitter.EmitLog("SUCCESS", fmt.Sprintf("Loaded %d Apple devices from IPSW.me", len(devices)), "FirmwareService")
	return s.cache, nil
}

func (s *firmwareService) GetDeviceDetails(identifier string) (*models.AppleHardware, error) {
	devices, err := s.GetDevices()
	if err != nil {
		return nil, err
	}

	for _, d := range devices {
		if d.Identifier == identifier {
			return &d, nil
		}
	}

	return nil, fmt.Errorf("device not found")
}
