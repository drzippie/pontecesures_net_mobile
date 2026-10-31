package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
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
	opts := application.Options{
		Name:        "Pontecesures.net",
		Description: "Las noticias de Pontecesures en tu móvil",
		Services: []application.Service{
			application.NewService(&ContentService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		IOS: application.IOSOptions{
			BackgroundColour:            chromeColour,
			DisableLinkPreview:          true,
			EnableInlineMediaPlayback:   true,
			ApplicationNameForUserAgent: "PontecesuresNetMobile/" + appVersion,
		},
		Android: application.AndroidOptions{
			DisableOverscroll: true,
			BackgroundColour:  chromeColour,
		},
	}

	app := application.New(opts)

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
