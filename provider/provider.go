package provider

import "errors"

var (
	// ErrNoSuchProvider is thrown when there is no provider
	ErrNoSuchProvider = errors.New("no such provider")
)

// Payload is all the info we need about a prisoners location
type Payload struct {
	// CountryCode represents the code, for example "NL"
	CountryCode string
	// GeoHash is latitude and longitude combined in a hash
	GeoHash string
}

// Provider is able to return a prisoners location data
type Provider interface {
	// Fetch the prisoners location data
	Lookup(ip string) (Payload, error)
}

// New creates a new provider using the given provider name
func New(provider string) (Provider, error) {
	switch provider {
	case "freeGeoIP":
		pr := new(freeGeoIP)
		return pr, nil
	}
	return nil, ErrNoSuchProvider
}
