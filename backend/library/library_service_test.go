package library

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/ElJoker63/ipa-downloader/v2/backend/events"
	"github.com/ElJoker63/ipa-downloader/v2/backend/storage"
)

func TestGetDownloadedIPAs(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ipa_test_dir_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a dummy .ipa file inside tempDir
	ipaPath := filepath.Join(tempDir, "SampleApp_renamed.ipa")
	ipaFile, err := os.Create(ipaPath)
	if err != nil {
		t.Fatalf("failed to create ipa file: %v", err)
	}

	zipWriter := zip.NewWriter(ipaFile)

	// Write Info.plist inside zip Payload/SampleApp.app/Info.plist
	infoPlistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleIdentifier</key>
	<string>com.example.sampleapp</string>
	<key>CFBundleDisplayName</key>
	<string>Sample App</string>
	<key>CFBundleShortVersionString</key>
	<string>2.4.0</string>
	<key>MinimumOSVersion</key>
	<string>15.0</string>
</dict>
</plist>`

	w, err := zipWriter.Create("Payload/SampleApp.app/Info.plist")
	if err != nil {
		t.Fatalf("failed to create plist entry in zip: %v", err)
	}
	_, err = w.Write([]byte(infoPlistContent))
	if err != nil {
		t.Fatalf("failed to write plist content: %v", err)
	}

	_ = zipWriter.Close()
	_ = ipaFile.Close()

	// Initialize sqlite storage in memory/temp
	dataDir, err := os.MkdirTemp("", "ipa_test_db_*")
	if err != nil {
		t.Fatalf("failed to create db dir: %v", err)
	}
	defer os.RemoveAll(dataDir)

	store, err := storage.NewSQLiteStorage(dataDir)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}
	defer store.Close()

	emitter := events.NewEmitter()
	service := NewLibraryService(store, emitter)

	// Test GetDownloadedIPAs
	ipas, err := service.GetDownloadedIPAs(tempDir)
	if err != nil {
		t.Fatalf("GetDownloadedIPAs returned error: %v", err)
	}

	if len(ipas) != 1 {
		t.Fatalf("expected 1 downloaded IPA, got %d", len(ipas))
	}

	item := ipas[0]
	if item.BundleID != "com.example.sampleapp" {
		t.Errorf("expected BundleID 'com.example.sampleapp', got '%s'", item.BundleID)
	}
	if item.AppName != "Sample App" {
		t.Errorf("expected AppName 'Sample App', got '%s'", item.AppName)
	}
	if item.Version != "2.4.0" {
		t.Errorf("expected Version '2.4.0', got '%s'", item.Version)
	}
	if item.MinimumOS != "15.0" {
		t.Errorf("expected MinimumOS '15.0', got '%s'", item.MinimumOS)
	}
	if item.FileName != "SampleApp_renamed.ipa" {
		t.Errorf("expected FileName 'SampleApp_renamed.ipa', got '%s'", item.FileName)
	}
}
