package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/drzippie/pontecesures_net_mobile/internal/meteo"
)

// appVersion is reported in the User-Agent sent to pontecesures.net.
const appVersion = "0.1.0"

// The built frontend (frontend/dist) is embedded into the binary and served
// in-process on every platform (no localhost server on mobile).
//
//go:embed all:frontend/dist
var assets embed.FS

// chromeColour is the masthead/tab-bar colour (--color-chrome in style.css).
// On iOS the WebView is laid out inside the safe area, so the native
// background shows behind the status bar and the home indicator: painting it
// with the chrome colour makes the bars blend with the app header and tab bar.
var chromeColour = application.NewRGB(27, 58, 75)

// groundColour is the page background (--ground), used on desktop.
var groundColour = application.NewRGB(247, 246, 243)

func main() {
	userAgent := "PontecesuresNetMobile/" + appVersion + " (+https://pontecesures.net)"
	// MeteoGalicia sky icons are cached on disk and served by the app itself
	// under /meteo/icons/, so the weather view also works offline.
	icons := meteo.NewIconStore(func() string { return cacheDir("meteo-icons") }, userAgent)

	opts := application.Options{
		Name:        "Pontecesures.net",
		Description: "Las noticias de Pontecesures en tu móvil",
		Services: []application.Service{
			application.NewService(&ContentService{}),
			application.NewService(&WeatherService{icons: icons}),
			application.NewService(&SettingsService{}),
		},
		Assets: application.AssetOptions{
			Handler:    application.AssetFileServerFS(assets),
			Middleware: icons.Middleware,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		IOS: application.IOSOptions{
			BackgroundColour:   chromeColour,
			DisableLinkPreview: true,
			// Swipe from the left edge goes back (closes the article, returns
			// to the previous section): every view is a history entry.
			EnableBackForwardNavigationGestures: true,
			EnableInlineMediaPlayback:           true,
			ApplicationNameForUserAgent:         "PontecesuresNetMobile/" + appVersion,
		},
		Android: application.AndroidOptions{
			DisableOverscroll: true,
			BackgroundColour:  chromeColour,
		},
	}

	app := application.New(opts)
	bridgeLifecycle(app)

	// On mobile the window colour is what shows behind the system bars.
	windowColour := groundColour
	if application.System.IsMobile() {
		windowColour = chromeColour
	}

	// On desktop the window mimics a phone so the mobile layout can be
	// developed with `wails3 dev`. On mobile, size and title are ignored.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Pontecesures.net",
		Width:            420,
		Height:           860,
		MinWidth:         360,
		BackgroundColour: windowColour,
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
