package cmd

import (
	"fmt"

	"github.com/ElJoker63/ipa-downloader/v2/pkg/keychain"
	gokeychain "github.com/byteness/go-keychain"
	"github.com/byteness/keyring"
)

// iosKeyring fixes the update query in keyring v1.9.0, which includes attributes
// that SecItemUpdate rejects on iOS (kSecMatchLimit and kSecReturnAttributes).
type iosKeyring struct {
	keyring.Keyring
	service string
}

func openKeyring(config keyring.Config) (keychain.Keyring, error) {
	ring, err := keyring.Open(config)
	if err != nil {
		return nil, fmt.Errorf("open iOS keyring: %w", err)
	}

	return &iosKeyring{Keyring: ring, service: config.ServiceName}, nil
}

func (k *iosKeyring) Set(item keyring.Item) error {
	query := gokeychain.NewItem()
	query.SetSecClass(gokeychain.SecClassGenericPassword)
	query.SetService(k.service)
	query.SetAccount(item.Key)

	attributes := gokeychain.NewItem()
	attributes.SetData(item.Data)
	attributes.SetLabel(item.Label)
	attributes.SetDescription(item.Description)

	err := gokeychain.UpdateItem(query, attributes)
	if err == gokeychain.ErrorItemNotFound {
		return k.Keyring.Set(item) //nolint:wrapcheck
	}

	if err != nil {
		return fmt.Errorf("update iOS keyring item: %w", err)
	}

	return nil
}
