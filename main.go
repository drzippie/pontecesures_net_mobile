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

// Brand background, shown before the WebView paints (avoids a white flash).
var background = application.NewRGB(247, 246, 243) // --ground

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
			BackgroundColour:            background,
			DisableLinkPreview:          true,
			EnableInlineMediaPlayback:   true,
			ApplicationNameForUserAgent: "PontecesuresNetMobile/" + appVersion,
		},
		Android: application.AndroidOptions{
			DisableOverscroll: true,
			BackgroundColour:  background,
		},
	}

	app := application.New(opts)

	// On desktop the window mimics a phone so the mobile layout can be
	// developed with `wails3 dev`. On mobile, size and title are ignored.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Pontecesures.net",
		Width:            420,
		Height:           860,
		MinWidth:         360,
		BackgroundColour: background,
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
