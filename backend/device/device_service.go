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

	"github.com/ElJoker63/ipa-downloader/v2/backend/events"
	"github.com/ElJoker63/ipa-downloader/v2/backend/models"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"howett.net/plist"
)

// DeviceService manages interactions with USB-connected Apple devices.
type DeviceService interface {
	SetContext(ctx context.Context)
	StartWatcher()
	StopWatcher()
	GetConnectedDevices() ([]models.DeviceInfo, error)
	IsDeviceConnected(udid string) bool
	PairDevice(udid string) error
	ListInstalledApps(udid string, appType string) ([]models.InstalledApp, error)
	QueueInstall(udid string, ipaPath string) error
	UninstallApp(udid string, bundleID string) error
	ValidateIPA(ipaPath string) (*models.IPAInfo, error)
	SelectIPAFile() (string, error)
	SelectMultipleIPAFiles() ([]string, error)
}

type installTask struct {
	udid    string
	ipaPath string
}

type deviceService struct {
	emitter events.Emitter
	ctx     context.Context
	mu      sync.RWMutex

	usbmux         giDevice.Usbmux
	devices        map[string]giDevice.Device
	deviceInfos    map[string]*models.DeviceInfo
	watcherRunning bool
	stopChan       chan struct{}

	installQueue chan installTask
}

// NewDeviceService creates a new instance of DeviceService.
func NewDeviceService(emitter events.Emitter) DeviceService {
	s := &deviceService{
		emitter:      emitter,
		stopChan:     make(chan struct{}),
		devices:      make(map[string]giDevice.Device),
		deviceInfos:  make(map[string]*models.DeviceInfo),
		installQueue: make(chan installTask, 100),
	}
	go s.processInstallQueue()
	return s
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
	defer func() {
		if r := recover(); r != nil {
			s.emitter.EmitLog("ERROR", fmt.Sprintf("Panic recovered in device service: %v", r), "DeviceService")
		}
	}()

	usbmux, err := giDevice.NewUsbmux()
	if err != nil {
		return
	}

	rawDevices, err := usbmux.Devices()
	if err != nil {
		return
	}

	s.mu.Lock()
	s.usbmux = usbmux
	currentUdids := make(map[string]bool)

	for _, dev := range rawDevices {
		udid := dev.Properties().SerialNumber
		currentUdids[udid] = true

		if _, exists := s.devices[udid]; !exists {
			// New device connected
			s.devices[udid] = dev
			info, err := s.fetchDeviceInfo(dev)
			if err == nil {
				s.deviceInfos[udid] = info
				s.emitter.EmitLog("INFO", fmt.Sprintf("iOS Device connected: %s (%s)", info.Name, info.Model), "DeviceService")
				s.emitter.Emit("device:connected", info)
			}
		} else {
			// Update existing device info (battery/storage might change)
			info, err := s.fetchDeviceInfo(dev)
			if err == nil {
				s.deviceInfos[udid] = info
				s.emitter.Emit("device:updated", info)
			}
		}
	}

	// Detect disconnections
	for udid := range s.devices {
		if !currentUdids[udid] {
			s.emitter.EmitLog("INFO", fmt.Sprintf("iOS Device disconnected: %s", udid), "DeviceService")
			delete(s.devices, udid)
			delete(s.deviceInfos, udid)
			s.emitter.Emit("device:disconnected", udid)
		}
	}
	s.mu.Unlock()
}

func (s *deviceService) processInstallQueue() {
	for task := range s.installQueue {
		s.mu.RLock()
		dev, exists := s.devices[task.udid]
		info := s.deviceInfos[task.udid]
		s.mu.RUnlock()

		if !exists {
			s.emitter.EmitLog("ERROR", fmt.Sprintf("Queue: Device %s no longer connected", task.udid), "DeviceService")
			continue
		}

		s.runActualInstall(dev, info, task.ipaPath)
	}
}

func (s *deviceService) runActualInstall(dev giDevice.Device, info *models.DeviceInfo, ipaPath string) {
	s.emitter.EmitLog("INFO", fmt.Sprintf("Starting queued installation on %s: %s", info.Name, filepath.Base(ipaPath)), "DeviceService")

	s.emitInstallProgress("Preparing", 10, fmt.Sprintf("[%s] Validating IPA...", info.Name))
	ipaInfo, err := s.ValidateIPA(ipaPath)
	if err != nil {
		s.emitInstallProgress("Failed", 0, err.Error())
		return
	}

	s.emitInstallProgress("Copying", 30, fmt.Sprintf("[%s] Transferring %s...", info.Name, ipaInfo.BundleName))

	err = dev.AppInstall(ipaPath)
	if err != nil {
		s.emitInstallProgress("Failed", 0, err.Error())
		s.emitter.Emit(events.EventDeviceInstallFailed, map[string]string{"error": err.Error(), "ipa": ipaPath, "udid": info.UDID})
		s.emitter.EmitLog("ERROR", fmt.Sprintf("Failed to install %s on %s: %v", ipaInfo.BundleName, info.Name, err), "DeviceService")
		return
	}

	s.emitInstallProgress("Complete", 100, fmt.Sprintf("[%s] Successfully installed %s", info.Name, ipaInfo.BundleName))
	s.emitter.EmitLog("SUCCESS", fmt.Sprintf("App %s installed on %s", ipaInfo.BundleName, info.Name), "DeviceService")
	s.emitter.Emit(events.EventDeviceInstallComplete, map[string]interface{}{"info": ipaInfo, "udid": info.UDID})
}

func (s *deviceService) GetConnectedDevices() ([]models.DeviceInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []models.DeviceInfo
	for _, info := range s.deviceInfos {
		list = append(list, *info)
	}
	return list, nil
}

func (s *deviceService) IsDeviceConnected(udid string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.devices[udid]
	return exists
}

func (s *deviceService) PairDevice(udid string) error {
	s.mu.RLock()
	dev, exists := s.devices[udid]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("device %s not connected", udid)
	}

	_, err := dev.Pair()
	if err != nil {
		return fmt.Errorf("device pairing failed: %w", err)
	}

	s.emitter.EmitLog("SUCCESS", fmt.Sprintf("Device %s paired", udid), "DeviceService")
	s.checkDeviceConnection() // Refresh info
	return nil
}

func (s *deviceService) ListInstalledApps(udid string, appType string) ([]models.InstalledApp, error) {
	s.mu.RLock()
	dev, exists := s.devices[udid]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("device %s not connected", udid)
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
		return nil, fmt.Errorf("failed to browse apps: %w", err)
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

func (s *deviceService) QueueInstall(udid string, ipaPath string) error {
	s.mu.RLock()
	_, exists := s.devices[udid]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("device %s not connected", udid)
	}

	s.installQueue <- installTask{udid: udid, ipaPath: ipaPath}
	s.emitter.EmitLog("INFO", fmt.Sprintf("Added to installation queue: %s", filepath.Base(ipaPath)), "DeviceService")
	return nil
}

func (s *deviceService) UninstallApp(udid string, bundleID string) error {
	s.mu.RLock()
	dev, exists := s.devices[udid]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("device %s not connected", udid)
	}

	s.emitter.EmitLog("INFO", fmt.Sprintf("Uninstalling %s from %s", bundleID, udid), "DeviceService")

	err := dev.AppUninstall(bundleID)
	if err != nil {
		return fmt.Errorf("failed to uninstall: %w", err)
	}

	s.emitter.EmitLog("SUCCESS", fmt.Sprintf("Uninstalled %s", bundleID), "DeviceService")
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

func (s *deviceService) SelectMultipleIPAFiles() ([]string, error) {
	s.mu.RLock()
	ctx := s.ctx
	s.mu.RUnlock()

	if ctx == nil {
		return nil, fmt.Errorf("context not available")
	}

	paths, err := runtime.OpenMultipleFilesDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select IPA Files to Install",
		Filters: []runtime.FileFilter{
			{DisplayName: "iOS Application Archive (*.ipa)", Pattern: "*.ipa"},
		},
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func (s *deviceService) fetchDeviceInfo(dev giDevice.Device) (*models.DeviceInfo, error) {
	props := dev.Properties()

	info := &models.DeviceInfo{
		UDID:         props.SerialNumber,
		SerialNumber: props.SerialNumber, // Default to props if specific key fails
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

	// Fetch more detailed IDs
	if val, err := dev.GetValue("", "SerialNumber"); err == nil {
		if str, ok := val.(string); ok {
			info.SerialNumber = str
		}
	}
	if val, err := dev.GetValue("", "InternationalMobileEquipmentIdentity"); err == nil {
		if str, ok := val.(string); ok {
			info.IMEI = str
		}
	}
	if val, err := dev.GetValue("", "InternationalMobileEquipmentIdentity2"); err == nil {
		if str, ok := val.(string); ok {
			info.IMEI2 = str
		}
	}
	if val, err := dev.GetValue("", "ModelNumber"); err == nil {
		if str, ok := val.(string); ok {
			info.ModelNumber = str
		}
	}
	if val, err := dev.GetValue("", "RegionInfo"); err == nil {
		if str, ok := val.(string); ok {
			info.RegionInfo = str
		}
	}
	if val, err := dev.GetValue("", "ActivationState"); err == nil {
		if str, ok := val.(string); ok {
			info.ActivationState = str
		}
	}
	// Jailbreak detection - primitive check
	if val, err := dev.GetValue("", "BrickState"); err == nil {
		if i, ok := val.(uint64); ok && i > 0 {
			info.IsJailbroken = true
		}
	}


	// Fetch Storage Info (Global domain)
	if val, err := dev.GetValue("", "TotalDiskCapacity"); err == nil {
		info.StorageTotal = getInt64FromVal(val)
	}
	if val, err := dev.GetValue("", "TotalDataAvailable"); err == nil {
		info.StorageFree = getInt64FromVal(val)
	}
	info.StorageUsed = info.StorageTotal - info.StorageFree

	// Fetch Battery Info (Requires trusted connection/session)
	if info.IsPaired {
		batteryInfo, err := dev.GetValue("com.apple.mobile.battery", "")
		if err == nil && batteryInfo != nil {
			if m, ok := batteryInfo.(map[string]interface{}); ok {
				info.BatteryLevel = int(getInt64FromVal(m["BatteryCurrentCapacity"]))
				if b, ok := m["IsCharging"].(bool); ok {
					info.BatteryCharging = b
				}

				// Health calculation: (FullChargeCapacity / DesignCapacity) * 100
				fullCap := getInt64FromVal(m["FullChargeCapacity"])
				designCap := getInt64FromVal(m["DesignCapacity"])
				if designCap > 0 {
					info.BatteryHealth = int((float64(fullCap) / float64(designCap)) * 100)
				}
				info.ChargeCycles = int(getInt64FromVal(m["CycleCount"]))
			}
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

func getInt64FromVal(val interface{}) int64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return v
	case uint64:
		return int64(v)
	case int:
		return int64(v)
	case float64:
		return int64(v)
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
