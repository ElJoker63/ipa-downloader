package device

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	giDevice "github.com/electricbubble/gidevice"
	"github.com/majd/ipa-downloader/v2/backend/events"
	"github.com/majd/ipa-downloader/v2/backend/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"howett.net/plist"
)

// DeviceService manages interactions with USB-connected Apple devices.
type DeviceService interface {
	SetContext(ctx context.Context)
	StartWatcher()
	StopWatcher()
	GetConnectedDevice() (*models.DeviceInfo, error)
	IsDeviceConnected() bool
	PairDevice() error
	ListInstalledApps(appType string) ([]models.InstalledApp, error)
	InstallIPA(ipaPath string) error
	UninstallApp(bundleID string) error
	ValidateIPA(ipaPath string) (*models.IPAInfo, error)
	SelectIPAFile() (string, error)
}

type deviceService struct {
	emitter events.Emitter
	ctx     context.Context
	mu      sync.RWMutex

	usbmux         giDevice.Usbmux
	cachedDevice   giDevice.Device
	cachedInfo     *models.DeviceInfo
	watcherRunning bool
	stopChan       chan struct{}
}

// NewDeviceService creates a new instance of DeviceService.
func NewDeviceService(emitter events.Emitter) DeviceService {
	return &deviceService{
		emitter:  emitter,
		stopChan: make(chan struct{}),
	}
}

func (s *deviceService) SetContext(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ctx = ctx
}

func (s *deviceService) StartWatcher() {
	s.mu.Lock()
	if s.watcherRunning {
		s.mu.Unlock()
		return
	}
	s.watcherRunning = true
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	go s.pollLoop()
}

func (s *deviceService) StopWatcher() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.watcherRunning {
		return
	}
	s.watcherRunning = false
	close(s.stopChan)
}

func (s *deviceService) pollLoop() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// Initial check
	s.checkDeviceConnection()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.checkDeviceConnection()
		}
	}
}

func (s *deviceService) checkDeviceConnection() {
	usbmux, err := giDevice.NewUsbmux()
	if err != nil {
		s.handleDisconnect()
		return
	}

	devices, err := usbmux.Devices()
	if err != nil || len(devices) == 0 {
		s.handleDisconnect()
		return
	}

	// Pick the first connected USB device
	dev := devices[0]

	s.mu.Lock()
	wasConnected := s.cachedDevice != nil
	s.cachedDevice = dev
	s.usbmux = usbmux
	s.mu.Unlock()

	info, err := s.fetchDeviceInfo(dev)
	if err != nil {
		s.handleDisconnect()
		return
	}

	s.mu.Lock()
	s.cachedInfo = info
	s.mu.Unlock()

	if !wasConnected {
		s.emitter.EmitLog("INFO", fmt.Sprintf("iOS Device connected: %s (%s)", info.Name, info.Model), "DeviceService")
		s.emitter.Emit(events.EventDeviceConnected, info)
	}
}

func (s *deviceService) handleDisconnect() {
	s.mu.Lock()
	wasConnected := s.cachedDevice != nil
	s.cachedDevice = nil
	s.cachedInfo = nil
	s.mu.Unlock()

	if wasConnected {
		s.emitter.EmitLog("INFO", "iOS Device disconnected", "DeviceService")
		s.emitter.Emit(events.EventDeviceDisconnected, nil)
	}
}

func (s *deviceService) fetchDeviceInfo(dev giDevice.Device) (*models.DeviceInfo, error) {
	props := dev.Properties()

	info := &models.DeviceInfo{
		UDID:         props.SerialNumber,
		SerialNumber: props.SerialNumber,
		IsConnected:  true,
		IsPaired:     true,
	}

	// Retrieve detailed Lockdown values
	if val, err := dev.GetValue("", "DeviceName"); err == nil {
		if str, ok := val.(string); ok {
			info.Name = str
		}
	}
	if val, err := dev.GetValue("", "ProductType"); err == nil {
		if str, ok := val.(string); ok {
			info.ProductType = str
			info.Model = formatProductType(str)
		}
	}
	if val, err := dev.GetValue("", "DeviceClass"); err == nil {
		if str, ok := val.(string); ok {
			info.DeviceClass = str
		}
	}
	if val, err := dev.GetValue("", "ProductVersion"); err == nil {
		if str, ok := val.(string); ok {
			info.IOSVersion = str
		}
	}
	if val, err := dev.GetValue("", "BuildVersion"); err == nil {
		if str, ok := val.(string); ok {
			info.BuildVersion = str
		}
	}
	if val, err := dev.GetValue("", "WiFiAddress"); err == nil {
		if str, ok := val.(string); ok {
			info.WiFiAddress = str
		}
	}

	if info.Name == "" {
		info.Name = props.SerialNumber
	}
	if info.Model == "" {
		info.Model = "iOS Device"
	}

	return info, nil
}

func formatProductType(productType string) string {
	modelsMap := map[string]string{
		"iPhone10,3": "iPhone X",
		"iPhone10,6": "iPhone X",
		"iPhone11,2": "iPhone XS",
		"iPhone11,4": "iPhone XS Max",
		"iPhone11,6": "iPhone XS Max",
		"iPhone11,8": "iPhone XR",
		"iPhone12,1": "iPhone 11",
		"iPhone12,3": "iPhone 11 Pro",
		"iPhone12,5": "iPhone 11 Pro Max",
		"iPhone12,8": "iPhone SE (2nd Gen)",
		"iPhone13,1": "iPhone 12 mini",
		"iPhone13,2": "iPhone 12",
		"iPhone13,3": "iPhone 12 Pro",
		"iPhone13,4": "iPhone 12 Pro Max",
		"iPhone14,2": "iPhone 13 Pro",
		"iPhone14,3": "iPhone 13 Pro Max",
		"iPhone14,4": "iPhone 13 mini",
		"iPhone14,5": "iPhone 13",
		"iPhone14,6": "iPhone SE (3rd Gen)",
		"iPhone14,7": "iPhone 14",
		"iPhone14,8": "iPhone 14 Plus",
		"iPhone15,2": "iPhone 14 Pro",
		"iPhone15,3": "iPhone 14 Pro Max",
		"iPhone15,4": "iPhone 15",
		"iPhone15,5": "iPhone 15 Plus",
		"iPhone16,1": "iPhone 15 Pro",
		"iPhone16,2": "iPhone 15 Pro Max",
	}

	if name, ok := modelsMap[productType]; ok {
		return name
	}
	return productType
}

func (s *deviceService) GetConnectedDevice() (*models.DeviceInfo, error) {
	s.mu.RLock()
	dev := s.cachedDevice
	info := s.cachedInfo
	s.mu.RUnlock()

	if dev == nil {
		// Attempt instant scan
		s.checkDeviceConnection()
		s.mu.RLock()
		info = s.cachedInfo
		s.mu.RUnlock()
	}

	if info == nil {
		return nil, fmt.Errorf("no iOS device connected via USB")
	}

	return info, nil
}

func (s *deviceService) IsDeviceConnected() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cachedDevice != nil
}

func (s *deviceService) PairDevice() error {
	s.mu.RLock()
	dev := s.cachedDevice
	s.mu.RUnlock()

	if dev == nil {
		return fmt.Errorf("no iOS device connected to pair")
	}

	_, err := dev.Pair()
	if err != nil {
		return fmt.Errorf("device pairing failed (please unlock device and trust this computer): %w", err)
	}

	s.emitter.EmitLog("SUCCESS", "Device successfully paired", "DeviceService")
	s.emitter.Emit(events.EventDevicePairStatus, true)
	return nil
}

func (s *deviceService) ListInstalledApps(appType string) ([]models.InstalledApp, error) {
	s.mu.RLock()
	dev := s.cachedDevice
	s.mu.RUnlock()

	if dev == nil {
		return nil, fmt.Errorf("no iOS device connected")
	}

	var targetType giDevice.ApplicationType
	switch strings.ToLower(appType) {
	case "system":
		targetType = giDevice.ApplicationTypeSystem
	case "user":
		targetType = giDevice.ApplicationTypeUser
	default:
		targetType = giDevice.ApplicationTypeAny
	}

	result, err := dev.InstallationProxyBrowse(
		giDevice.WithApplicationType(targetType),
		giDevice.WithReturnAttributes(
			"CFBundleIdentifier",
			"CFBundleDisplayName",
			"CFBundleName",
			"CFBundleShortVersionString",
			"CFBundleVersion",
			"ApplicationType",
			"MinimumOSVersion",
			"SignerIdentity",
			"StaticDiskUsage",
			"DynamicDiskUsage",
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to browse installed applications: %w", err)
	}

	var installedApps []models.InstalledApp

	for _, item := range result {
		if m, ok := item.(map[string]interface{}); ok {
				app := models.InstalledApp{
					BundleID:       getStringVal(m, "CFBundleIdentifier"),
					Version:        getStringVal(m, "CFBundleVersion"),
					ShortVersion:   getStringVal(m, "CFBundleShortVersionString"),
					AppType:        getStringVal(m, "ApplicationType"),
					MinimumOS:      getStringVal(m, "MinimumOSVersion"),
					SignerIdentity: getStringVal(m, "SignerIdentity"),
					Size:           getInt64Val(m, "StaticDiskUsage"),
					DynamicSize:    getInt64Val(m, "DynamicDiskUsage"),
				}

				appName := getStringVal(m, "CFBundleDisplayName")
				if appName == "" {
					appName = getStringVal(m, "CFBundleName")
				}
				if appName == "" {
					appName = app.BundleID
				}
				app.Name = appName

				if app.BundleID != "" {
					installedApps = append(installedApps, app)
				}
			}
	}

	return installedApps, nil
}

func (s *deviceService) InstallIPA(ipaPath string) error {
	s.mu.RLock()
	dev := s.cachedDevice
	s.mu.RUnlock()

	if dev == nil {
		return fmt.Errorf("no iOS device connected")
	}

	s.emitter.EmitLog("INFO", fmt.Sprintf("Starting IPA installation: %s", filepath.Base(ipaPath)), "DeviceService")

	s.emitInstallProgress("Preparing", 10, "Validating IPA package...")
	info, err := s.ValidateIPA(ipaPath)
	if err != nil {
		s.emitInstallProgress("Failed", 0, err.Error())
		return fmt.Errorf("IPA validation failed: %w", err)
	}

	s.emitInstallProgress("Copying", 30, fmt.Sprintf("Transferring %s (%s)...", info.BundleName, formatBytes(info.FileSizeBytes)))

	// Perform actual installation
	s.emitInstallProgress("Installing", 60, "Deploying package to device...")
	err = dev.AppInstall(ipaPath)
	if err != nil {
		s.emitInstallProgress("Failed", 0, err.Error())
		s.emitter.Emit(events.EventDeviceInstallFailed, map[string]string{"error": err.Error(), "ipa": ipaPath})
		return fmt.Errorf("IPA installation failed: %w", err)
	}

	s.emitInstallProgress("Complete", 100, fmt.Sprintf("Successfully installed %s (%s)", info.BundleName, info.Version))
	s.emitter.EmitLog("SUCCESS", fmt.Sprintf("App %s (%s) installed on device", info.BundleName, info.BundleID), "DeviceService")
	s.emitter.Emit(events.EventDeviceInstallComplete, info)

	return nil
}

func (s *deviceService) UninstallApp(bundleID string) error {
	s.mu.RLock()
	dev := s.cachedDevice
	s.mu.RUnlock()

	if dev == nil {
		return fmt.Errorf("no iOS device connected")
	}

	s.emitter.EmitLog("INFO", fmt.Sprintf("Uninstalling app: %s", bundleID), "DeviceService")

	err := dev.AppUninstall(bundleID)
	if err != nil {
		return fmt.Errorf("failed to uninstall app %s: %w", bundleID, err)
	}

	s.emitter.EmitLog("SUCCESS", fmt.Sprintf("Successfully uninstalled %s", bundleID), "DeviceService")
	return nil
}

func (s *deviceService) ValidateIPA(ipaPath string) (*models.IPAInfo, error) {
	r, err := zip.OpenReader(ipaPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open IPA file (invalid archive): %w", err)
	}
	defer r.Close()

	var infoPlistFile *zip.File
	var size int64

	for _, f := range r.File {
		size += f.FileInfo().Size()
		// Look for Payload/*.app/Info.plist
		parts := strings.Split(filepath.ToSlash(f.Name), "/")
		if len(parts) == 3 && parts[0] == "Payload" && strings.HasSuffix(parts[1], ".app") && parts[2] == "Info.plist" {
			infoPlistFile = f
			break
		}
	}

	if infoPlistFile == nil {
		return nil, fmt.Errorf("invalid IPA: missing Info.plist in Payload directory")
	}

	rc, err := infoPlistFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to read Info.plist from IPA: %w", err)
	}
	defer rc.Close()

	buf, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract Info.plist: %w", err)
	}

	var plistData map[string]interface{}
	decoder := plist.NewDecoder(bytes.NewReader(buf))
	if err := decoder.Decode(&plistData); err != nil {
		return nil, fmt.Errorf("failed to parse Info.plist: %w", err)
	}

	info := &models.IPAInfo{
		BundleID:      getStringVal(plistData, "CFBundleIdentifier"),
		BundleName:    getStringVal(plistData, "CFBundleDisplayName"),
		Version:       getStringVal(plistData, "CFBundleVersion"),
		ShortVersion:  getStringVal(plistData, "CFBundleShortVersionString"),
		MinimumOS:     getStringVal(plistData, "MinimumOSVersion"),
		FileSizeBytes: size,
		IsValid:       true,
	}

	if info.BundleName == "" {
		info.BundleName = getStringVal(plistData, "CFBundleName")
	}

	return info, nil
}

func (s *deviceService) SelectIPAFile() (string, error) {
	s.mu.RLock()
	ctx := s.ctx
	s.mu.RUnlock()

	if ctx == nil {
		return "", fmt.Errorf("context not available")
	}

	path, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select IPA File to Install",
		Filters: []runtime.FileFilter{
			{DisplayName: "iOS Application Archive (*.ipa)", Pattern: "*.ipa"},
		},
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

func (s *deviceService) emitInstallProgress(phase string, percent int, message string) {
	prog := models.DeviceInstallProgress{
		Phase:   phase,
		Percent: percent,
		Message: message,
	}
	s.emitter.Emit(events.EventDeviceInstallProgress, prog)
}

// Utility helper functions
func getStringVal(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getInt64Val(m map[string]interface{}, key string) int64 {
	if val, ok := m[key]; ok && val != nil {
		switch v := val.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case float64:
			return int64(v)
		}
	}
	return 0
}

func formatBytes(bytes int64) string {
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
