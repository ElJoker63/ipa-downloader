package update

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"
)

func TestIsNewer(t *testing.T) {
	tests := []struct {
		latest   string
		current  string
		expected bool
	}{
		{"1.4.3", "1.4.2", true},
		{"1.4.2", "1.4.2", false},
		{"1.4.1", "1.4.2", false},
		{"1.5.0", "1.4.9", true},
		{"2.0.0", "1.9.9", true},
		{"1.4.2.1", "1.4.2", true},
		{"1.4.2", "1.4.2.1", false},
	}

	for _, tt := range tests {
		got := isNewer(tt.latest, tt.current)
		if got != tt.expected {
			t.Errorf("isNewer(%q, %q) = %v, expected %v", tt.latest, tt.current, got, tt.expected)
		}
	}
}

func TestFindExecutableInZip(t *testing.T) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	// Add Info.plist
	plistWriter, err := zw.Create("ipa-downloader-desktop.app/Contents/Info.plist")
	if err != nil {
		t.Fatalf("failed to create plist entry: %v", err)
	}
	expectedPlist := []byte("<plist>test</plist>")
	if _, err := plistWriter.Write(expectedPlist); err != nil {
		t.Fatalf("failed to write plist data: %v", err)
	}

	// Add Binary
	binWriter, err := zw.Create("ipa-downloader-desktop.app/Contents/MacOS/ipa-downloader-desktop")
	if err != nil {
		t.Fatalf("failed to create bin entry: %v", err)
	}
	expectedBin := []byte("binary content")
	if _, err := binWriter.Write(expectedBin); err != nil {
		t.Fatalf("failed to write binary data: %v", err)
	}

	if err := zw.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("failed to create zip reader: %v", err)
	}

	// Test exact match
	rc, plistData, err := findExecutableInZip(zr, "ipa-downloader-desktop")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rc.Close()

	binData, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("failed to read binary data: %v", err)
	}
	if string(binData) != string(expectedBin) {
		t.Errorf("binary content mismatch: got %q, want %q", binData, expectedBin)
	}
	if string(plistData) != string(expectedPlist) {
		t.Errorf("plist content mismatch: got %q, want %q", plistData, expectedPlist)
	}

	// Test fallback when targetName differs
	rcFallback, _, err := findExecutableInZip(zr, "other-name")
	if err != nil {
		t.Fatalf("unexpected error on fallback: %v", err)
	}
	defer rcFallback.Close()
	binDataFallback, _ := io.ReadAll(rcFallback)
	if string(binDataFallback) != string(expectedBin) {
		t.Errorf("fallback binary mismatch: got %q, want %q", binDataFallback, expectedBin)
	}

	// Test not found when zip has no binary
	emptyBuf := new(bytes.Buffer)
	emptyZw := zip.NewWriter(emptyBuf)
	_, _ = emptyZw.Create("some/file.txt")
	_ = emptyZw.Close()
	emptyZr, _ := zip.NewReader(bytes.NewReader(emptyBuf.Bytes()), int64(emptyBuf.Len()))

	_, _, err = findExecutableInZip(emptyZr, "nonexistent")
	if err == nil {
		t.Errorf("expected error for nonexistent binary, got nil")
	}
}
