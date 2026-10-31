package main

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/drzippie/pontecesures_net_mobile/internal/settings"
)

// SettingsService stores the user's display preferences (language, theme,
// text size). It is the place where profile management will live once
// accounts exist.
type SettingsService struct {
	store *settings.Store
}

// ServiceStartup resolves the settings file location.
func (s *SettingsService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	s.store = settings.NewStore(settingsPath())
	return nil
}

// settingsPath: persistent data, not a cache (preferences must survive
// cache clean-ups).
func settingsPath() string {
	return dataPath("settings.json")
}

// Get returns the saved settings (empty fields mean "default").
func (s *SettingsService) Get() settings.Settings {
	return s.store.Load()
}

// Save stores the settings.
func (s *SettingsService) Save(v settings.Settings) error {
	return s.store.Save(v)
}
