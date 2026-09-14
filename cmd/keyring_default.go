//go:build !ios

package cmd

import (
	"github.com/ElJoker63/ipa-downloader/v2/pkg/keychain"
	"github.com/byteness/keyring"
)

func openKeyring(config keyring.Config) (keychain.Keyring, error) {
	return keyring.Open(config) //nolint:wrapcheck
}
