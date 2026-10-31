package main

import (
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Lifecycle events for the frontend.
//
// The WebView's own visibilitychange does not fire on Android when the app
// goes to the background, and events sent to a paused WebView are only
// delivered when it resumes. So Go, which receives the native events right
// away, measures the time spent in the background and sends it with
// app:foreground.
const eventForeground = "app:foreground" // data: milliseconds spent in the background

func init() {
	application.RegisterEvent[int64](eventForeground)
}

type lifecycle struct {
	mu           sync.Mutex
	backgroundAt time.Time
}

func (l *lifecycle) background(*application.ApplicationEvent) {
	l.mu.Lock()
	if l.backgroundAt.IsZero() {
		l.backgroundAt = time.Now()
	}
	l.mu.Unlock()
}

func (l *lifecycle) foreground(app *application.App) func(*application.ApplicationEvent) {
	return func(*application.ApplicationEvent) {
		l.mu.Lock()
		var away time.Duration
		if !l.backgroundAt.IsZero() {
			away = time.Since(l.backgroundAt)
			l.backgroundAt = time.Time{}
		}
		l.mu.Unlock()
		app.Event.Emit(eventForeground, away.Milliseconds())
	}
}

// bridgeLifecycle wires the iOS and Android native events. On desktop none of
// these fire, which is fine.
func bridgeLifecycle(app *application.App) {
	l := &lifecycle{}
	app.Event.OnApplicationEvent(events.Android.ActivityStopped, l.background)
	app.Event.OnApplicationEvent(events.IOS.ApplicationDidEnterBackground, l.background)
	app.Event.OnApplicationEvent(events.Android.ActivityResumed, l.foreground(app))
	app.Event.OnApplicationEvent(events.IOS.ApplicationDidBecomeActive, l.foreground(app))
}
