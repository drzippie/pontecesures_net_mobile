package main

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/drzippie/pontecesures_net_mobile/internal/meteo"
)

// WeatherService exposes the MeteoGalicia forecast for Pontecesures.
type WeatherService struct {
	client *meteo.Client
	// icons is shared with the asset server middleware (main.go).
	icons *meteo.IconStore
}

// ServiceStartup initialises the MeteoGalicia client.
func (s *WeatherService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	s.client = meteo.NewClient(meteo.Options{
		CacheDir:  cacheDir("meteo-cache"),
		UserAgent: "PontecesuresNetMobile/" + appVersion + " (+https://pontecesures.net)",
	})
	return nil
}

// Forecast returns the 3-4 day forecast for Pontecesures.
func (s *WeatherService) Forecast() (meteo.Forecast, error) {
	ctx, cancel := reqCtx()
	defer cancel()
	f, err := s.client.Forecast(ctx, meteo.PontecesuresID)
	if err == nil && s.icons != nil && !f.Stale {
		// Download this forecast's icons now so they also show offline.
		go s.icons.Prefetch(f)
	}
	return f, err
}
