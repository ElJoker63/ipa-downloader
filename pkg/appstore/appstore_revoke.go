package appstore

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/byteness/keyring"
)

func (t *appstore) Revoke() error {
	accountErr := t.keychain.Remove("account")
	if accountErr != nil {
		accountErr = fmt.Errorf("failed to remove account from keychain: %w", accountErr)
	}

	// Drop the cached SAP signer regardless of keychain removal succeeding,
	// so a subsequent login always starts from a clean handshake instead of
	// reusing a signer tied to the account that just logged out.
	closeErr := t.closeSharedSigner()
	if closeErr != nil {
		closeErr = fmt.Errorf("failed to close SAP action signer: %w", closeErr)
	}

	t.kbsyncCache.mu.Lock()
	defer t.kbsyncCache.mu.Unlock()

	t.kbsyncCache.entry = kbsyncCacheEntry{}

	cacheErr := t.keychain.Remove(kbsyncCacheKey)
	if errors.Is(cacheErr, keyring.ErrKeyNotFound) || errors.Is(cacheErr, fs.ErrNotExist) {
		cacheErr = nil
	} else if cacheErr != nil {
		cacheErr = fmt.Errorf("failed to remove kbsync cache from keychain: %w", cacheErr)
	}

	return errors.Join(accountErr, closeErr, cacheErr)
}
