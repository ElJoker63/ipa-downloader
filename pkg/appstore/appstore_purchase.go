package appstore

import (
	"errors"
	"fmt"
	gohttp "net/http"

	"github.com/ElJoker63/ipa-downloader/v2/pkg/http"
)

var (
	ErrPasswordTokenExpired   = errors.New("password token is expired")
	ErrLicenseAlreadyExists   = errors.New("license already exists")
	ErrSubscriptionRequired   = errors.New("subscription required")
	ErrTemporarilyUnavailable = errors.New("item is temporarily unavailable")
)

type PurchaseInput struct {
	Account    Account
	App        App
	StoreFront string
	Country    string
	Platform   Platform
}

func (t *appstore) Purchase(input PurchaseInput) error {
	macAddr, err := t.machine.MacAddress()
	if err != nil {
		return fmt.Errorf("failed to get mac address: %w", err)
	}

	guid, machineID, err := machineIdentity(macAddr)
	if err != nil {
		return fmt.Errorf("failed to get machine identity: %w", err)
	}

	if input.App.Price > 0 {
		return errors.New("purchasing paid apps is not supported")
	}

	targetStoreFront := input.StoreFront
	if targetStoreFront == "" && input.Country != "" {
		targetStoreFront = StoreFrontForCountry(input.Account.StoreFront, input.Country)
	}
	if targetStoreFront == "" {
		targetStoreFront = input.Account.StoreFront
	}

	var signer ActionSigner
	if t.actionSignerFactory != nil {
		bag, bagErr := t.bag(guid)
		if bagErr == nil {
			s, sErr := t.sharedSigner(guid, machineID, bag.SAPConfig)
			if sErr == nil && s != nil {
				signer = s
			}
		}
	}

	err = t.purchaseWithParams(input.Account, input.App, targetStoreFront, guid, PricingParameterAppStore, signer)
	if err == nil {
		return nil
	}

	if targetStoreFront != input.Account.StoreFront && input.Account.StoreFront != "" {
		if fbErr := t.purchaseWithParams(input.Account, input.App, input.Account.StoreFront, guid, PricingParameterAppStore, signer); fbErr == nil {
			return nil
		}
	}

	if input.Platform != PlatformMacOS && errors.Is(err, ErrTemporarilyUnavailable) {
		err = t.purchaseWithParams(input.Account, input.App, targetStoreFront, guid, PricingParameterAppleArcade, signer)
		if err != nil {
			return fmt.Errorf("failed to purchase item with param '%s': %w", PricingParameterAppleArcade, err)
		}

		return nil
	}

	return fmt.Errorf("failed to purchase item with param '%s': %w", PricingParameterAppStore, err)
}

type purchaseResult struct {
	FailureType     string `plist:"failureType,omitempty"`
	CustomerMessage string `plist:"customerMessage,omitempty"`
	JingleDocType   string `plist:"jingleDocType,omitempty"`
	Status          int    `plist:"status,omitempty"`
	MAllowed        *bool  `plist:"m-allowed,omitempty"`
}

func (t *appstore) purchaseWithParams(acc Account, app App, storeFront string, guid string, pricingParameters string, signer ActionSigner) error {
	req := t.purchaseRequest(acc, app, storeFront, guid, pricingParameters, signer)
	res, err := t.purchaseClient.Send(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if res.Data.FailureType == FailureTypeTemporarilyUnavailable {
		return ErrTemporarilyUnavailable
	}

	if res.Data.CustomerMessage == CustomerMessageSubscriptionRequired {
		return ErrSubscriptionRequired
	}

	if res.Data.FailureType == FailureTypePasswordTokenExpired ||
		res.Data.FailureType == FailureTypeSignInRequired ||
		res.Data.FailureType == FailureTypeDeviceVerificationFailed ||
		res.Data.CustomerMessage == CustomerMessagePasswordChanged {
		return ErrPasswordTokenExpired
	}

	if res.Data.FailureType == FailureTypeLicenseAlreadyExists {
		return ErrLicenseAlreadyExists
	}

	if res.Data.FailureType != "" && res.Data.CustomerMessage != "" {
		return NewErrorWithMetadata(fmt.Errorf("apple purchase error (%s): %s", res.Data.FailureType, res.Data.CustomerMessage), res)
	}

	if res.Data.FailureType != "" {
		return NewErrorWithMetadata(fmt.Errorf("apple purchase error (%s)", res.Data.FailureType), res)
	}

	if res.StatusCode == gohttp.StatusInternalServerError {
		return ErrLicenseAlreadyExists
	}

	if res.Data.JingleDocType != "purchaseSuccess" || res.Data.Status != 0 {
		if res.Data.CustomerMessage != "" {
			return NewErrorWithMetadata(errors.New(res.Data.CustomerMessage), res)
		}
		return NewErrorWithMetadata(fmt.Errorf("failed to purchase app (status: %d, docType: %s)", res.Data.Status, res.Data.JingleDocType), res)
	}

	return nil
}

func (t *appstore) purchaseRequest(acc Account, app App, storeFront, guid string, pricingParameters string, signer ActionSigner) http.Request {
	podPrefix := ""
	if acc.Pod != "" {
		podPrefix = "p" + acc.Pod + "-"
	}

	return http.Request{
		URL:            fmt.Sprintf("https://%s%s%s", podPrefix, PrivateAppStoreAPIDomain, PrivateAppStoreAPIPathPurchase),
		Method:         http.MethodPOST,
		ResponseFormat: http.ResponseFormatXML,
		ActionSigner:   signer,
		Headers: map[string]string{
			"Content-Type":        "application/x-apple-plist",
			"iCloud-DSID":         acc.DirectoryServicesID,
			"X-Dsid":              acc.DirectoryServicesID,
			"X-Apple-Store-Front": storeFront,
			"X-Token":             acc.PasswordToken,
		},
		Payload: &http.XMLPayload{
			Content: map[string]interface{}{
				"appExtVrsId":               "0",
				"hasAskedToFulfillPreorder": "true",
				"buyWithoutAuthorization":   "true",
				"hasDoneAgeCheck":           "true",
				"guid":                      guid,
				"needDiv":                   "0",
				"origPage":                  fmt.Sprintf("Software-%d", app.ID),
				"origPageLocation":          "Buy",
				"price":                     "0",
				"pricingParameters":         pricingParameters,
				"productType":               "C",
				"salableAdamId":             app.ID,
				"serialNumber":              "0",
			},
		},
	}
}
