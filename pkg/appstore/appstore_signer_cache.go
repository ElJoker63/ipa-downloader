package appstore

import "sync"

// signerCache holds a single SAP ActionSigner shared across every AppStore
// call made during a session.
//
// Building a SAP signer starts an in-process CPU emulator (Unicorn) with its
// own ~100MB+ guest memory space and performs a certificate handshake with
// Apple over the network. Every AppStore method that talks to the private
// download/purchase API (Login, OwnedApps, Purchase, Download, ListVersions,
// GetVersionMetadata) used to build and immediately tear down its own signer,
// so browsing the app for a while could spin up and discard dozens of these
// heavy signers back to back — the main driver behind the app's high memory
// use. Caching one signer per authenticated identity and reusing it for the
// life of the session avoids repeating that cost on every single request.
type signerCache struct {
	mu     sync.Mutex
	signer ActionSigner
	guid   string
	config SAPConfig
}

// sharedSigner returns the cached ActionSigner for the given guid/config,
// creating one via actionSignerFactory the first time it's needed or
// whenever the account identity or SAP configuration changes. The returned
// signer's Sign method is safe for concurrent use by multiple callers (e.g.
// several downloads running at once), so callers should not Close it — it is
// closed automatically when it's replaced or when closeSharedSigner is
// called (typically on logout).
//
// A nil, nil result means no signer factory is configured; callers that
// require a signer should treat that as an error themselves, matching the
// previous per-call behavior.
func (t *appstore) sharedSigner(guid string, machineID []byte, config SAPConfig) (ActionSigner, error) {
	if t.actionSignerFactory == nil {
		return nil, nil
	}

	t.signerState.mu.Lock()
	defer t.signerState.mu.Unlock()

	if t.signerState.signer != nil && t.signerState.guid == guid && t.signerState.config == config {
		return t.signerState.signer, nil
	}

	if t.signerState.signer != nil {
		_ = t.signerState.signer.Close()
		t.signerState.signer = nil
	}

	signer, err := t.actionSignerFactory(config, machineID)
	if err != nil {
		return nil, err
	}

	if signer == nil {
		return nil, nil
	}

	t.signerState.signer = signer
	t.signerState.guid = guid
	t.signerState.config = config

	return signer, nil
}

// closeSharedSigner closes and discards the cached signer, if any. It is
// called when a session ends (Revoke) so a later Login starts a fresh
// handshake instead of reusing a signer tied to the previous account.
func (t *appstore) closeSharedSigner() error {
	t.signerState.mu.Lock()
	defer t.signerState.mu.Unlock()

	if t.signerState.signer == nil {
		return nil
	}

	err := t.signerState.signer.Close()
	t.signerState.signer = nil
	t.signerState.guid = ""
	t.signerState.config = SAPConfig{}

	return err
}
