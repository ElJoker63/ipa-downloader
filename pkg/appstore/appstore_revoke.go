package appstore

import (
	"fmt"
)

func (t *appstore) Revoke() error {
	err := t.keychain.Remove("account")

	// Drop the cached SAP signer regardless of keychain removal succeeding,
	// so a subsequent login always starts from a clean handshake instead of
	// reusing a signer tied to the account that just logged out.
	closeErr := t.closeSharedSigner()

	if err != nil {
		return fmt.Errorf("failed to remove account from keychain: %w", err)
	}

	if closeErr != nil {
		return fmt.Errorf("failed to close SAP action signer: %w", closeErr)
	}

	return nil
}
