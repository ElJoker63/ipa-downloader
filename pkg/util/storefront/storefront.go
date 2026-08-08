package storefront

import "strings"

// Country represents a country associated with an Apple Storefront.
type Country struct {
	ID   string
	Name string
	Flag string
}

var storefronts = map[string]Country{
	"143441": {ID: "US", Name: "United States", Flag: "🇺🇸"},
	"143444": {ID: "GB", Name: "United Kingdom", Flag: "🇬🇧"},
	"143455": {ID: "CA", Name: "Canada", Flag: "🇨🇦"},
	"143460": {ID: "AU", Name: "Australia", Flag: "🇦🇺"},
	"143442": {ID: "FR", Name: "France", Flag: "🇫🇷"},
	"143443": {ID: "DE", Name: "Germany", Flag: "🇩🇪"},
	"143450": {ID: "IT", Name: "Italy", Flag: "🇮🇹"},
	"143454": {ID: "ES", Name: "Spain", Flag: "🇪🇸"},
	"143447": {ID: "JP", Name: "Japan", Flag: "🇯🇵"},
	"143462": {ID: "MX", Name: "Mexico", Flag: "🇲🇽"},
	"143463": {ID: "AR", Name: "Argentina", Flag: "🇦🇷"},
	"143465": {ID: "BR", Name: "Brazil", Flag: "🇧🇷"},
	"143471": {ID: "CL", Name: "Chile", Flag: "🇨🇱"},
	"143472": {ID: "CO", Name: "Colombia", Flag: "🇨🇴"},
	"143476": {ID: "PE", Name: "Peru", Flag: "🇵🇪"},
	"143477": {ID: "VE", Name: "Venezuela", Flag: "🇻🇪"},
	"143480": {ID: "RU", Name: "Russia", Flag: "🇷🇺"},
	"143467": {ID: "CN", Name: "China", Flag: "🇨🇳"},
	"143469": {ID: "IN", Name: "India", Flag: "🇮🇳"},
}

// GetCountry returns the Country associated with the given Storefront ID.
// The input can be a full storefront string like "143441-1,29" or just "143441".
func GetCountry(sf string) Country {
	id := sf
	if strings.Contains(sf, "-") {
		id = strings.Split(sf, "-")[0]
	}

	if country, ok := storefronts[id]; ok {
		return country
	}

	return Country{ID: "Unknown", Name: "Unknown (" + id + ")", Flag: "🏳️"}
}

// Format returns a human-readable string for the storefront.
func Format(sf string) string {
	c := GetCountry(sf)
	return c.Flag + " " + c.Name
}
