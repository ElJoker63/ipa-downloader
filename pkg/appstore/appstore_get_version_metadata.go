package appstore

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ElJoker63/ipa-downloader/v2/pkg/http"
)

type GetVersionMetadataInput struct {
	Account   Account
	App       App
	VersionID string
}

type GetVersionMetadataOutput struct {
	DisplayVersion string
	ReleaseDate    time.Time
}

func (t *appstore) GetVersionMetadata(input GetVersionMetadataInput) (GetVersionMetadataOutput, error) {
	macAddr, err := t.machine.MacAddress()
	if err != nil {
		return GetVersionMetadataOutput{}, fmt.Errorf("failed to get mac address: %w", err)
	}

	guid, machineID, err := machineIdentity(macAddr)
	if err != nil {
		guid = strings.ReplaceAll(strings.ToUpper(macAddr), ":", "")
	}

	var signer ActionSigner
	if t.actionSignerFactory != nil {
		bag, bagErr := t.bag(guid)
		if bagErr == nil {
			s, sErr := t.actionSignerFactory(bag.SAPConfig, machineID)
			if sErr == nil && s != nil {
				signer = s
				defer signer.Close()
			}
		}
	}

	req := t.getVersionMetadataRequest(input.Account, input.App, guid, input.VersionID, signer)
	res, err := t.downloadClient.Send(req)

	if err != nil {
		return GetVersionMetadataOutput{}, fmt.Errorf("failed to send http request: %w", err)
	}

	if res.Data.FailureType == FailureTypePasswordTokenExpired || res.Data.FailureType == FailureTypeSignInRequired {
		return GetVersionMetadataOutput{}, ErrPasswordTokenExpired
	}

	if res.Data.FailureType == FailureTypeLicenseNotFound {
		return GetVersionMetadataOutput{}, ErrLicenseRequired
	}

	if res.Data.FailureType != "" && res.Data.CustomerMessage != "" {
		return GetVersionMetadataOutput{}, NewErrorWithMetadata(fmt.Errorf("received error: %s", res.Data.CustomerMessage), res)
	}

	if res.Data.FailureType != "" {
		return GetVersionMetadataOutput{}, NewErrorWithMetadata(fmt.Errorf("received error: %s", res.Data.FailureType), res)
	}

	if len(res.Data.Items) == 0 {
		return GetVersionMetadataOutput{}, NewErrorWithMetadata(errors.New("invalid response"), res)
	}

	item := res.Data.Items[0]

	// Do not fall back to item.Metadata here. The App Store download API can
	// return stale version and release date values, so the IPA Info.plist is the
	// source of truth and failures should be visible to callers.
	metadata, err := t.readVersionMetadataFromIPA(item.URL)
	if err != nil {
		return GetVersionMetadataOutput{}, fmt.Errorf("failed to read version metadata: %w", err)
	}

	return GetVersionMetadataOutput(metadata), nil
}

func (t *appstore) getVersionMetadataRequest(acc Account, app App, guid string, version string, signer ActionSigner) http.Request {
	payload := map[string]interface{}{
		"creditDisplay":     "",
		"guid":              guid,
		"salableAdamId":     app.ID,
		"externalVersionId": version,
		"appExtVrsId":       version,
		"pricingParameters": "STDQ",
		"serialNumber":      "0",
	}

	podPrefix := ""
	if acc.Pod != "" {
		podPrefix = "p" + acc.Pod + "-"
	}

	headers := map[string]string{
		"Content-Type": "application/x-apple-plist",
		"iCloud-DSID":  acc.DirectoryServicesID,
		"X-Dsid":       acc.DirectoryServicesID,
	}
	if acc.StoreFront != "" {
		headers["X-Apple-Store-Front"] = acc.StoreFront
	}
	if acc.PasswordToken != "" {
		headers["X-Token"] = acc.PasswordToken
	}

	return http.Request{
		URL:            fmt.Sprintf("https://%s%s%s?guid=%s", podPrefix, PrivateAppStoreAPIDomain, PrivateAppStoreAPIPathDownload, guid),
		Method:         http.MethodPOST,
		ResponseFormat: http.ResponseFormatXML,
		ActionSigner:   signer,
		Headers:        headers,
		Payload: &http.XMLPayload{
			Content: payload,
		},
	}
}
