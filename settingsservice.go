package main

import (
	"context"
	"os"
	"path/filepath"

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

// settingsPath: app private storage on mobile, user config dir on desktop.
// Not a cache directory: preferences must survive cache clean-ups.
func settingsPath() string {
	if dir := application.Mobile.StoragePath(); dir != "" {
		return filepath.Join(dir, "settings.json")
	}
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "pontecesures_net_mobile", "settings.json")
	}
	return ""
}

// Get returns the saved settings (empty fields mean "default").
func (s *SettingsService) Get() settings.Settings {
	return s.store.Load()
}

// Save stores the settings.
func (s *SettingsService) Save(v settings.Settings) error {
	return s.store.Save(v)
}
