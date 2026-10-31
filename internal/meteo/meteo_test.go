package meteo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const sample = `{"predConcello":{"idConcello":36044,"nome":"Pontecesures","listaPredDiaConcello":[
{"ceo":{"manha":101,"noite":201,"tarde":111},"ceoDia":101,"dataPredicion":"2026-10-31T00:00:00","nivelAviso":null,
 "pchoiva":{"manha":5,"noite":5,"tarde":60},"tMax":21,"tMin":11,"tmaxFranxa":{"manha":18,"noite":17,"tarde":21},
 "tminFranxa":{"manha":11,"noite":12,"tarde":17},"uvMax":3,"vento":{"manha":301,"noite":299,"tarde":-9999}},
{"ceo":{"manha":-9999,"noite":-9999,"tarde":-9999},"ceoDia":-9999,"dataPredicion":"2026-11-01T00:00:00","nivelAviso":"2",
 "pchoiva":{"manha":-9999,"noite":-9999,"tarde":-9999},"tMax":-9999,"tMin":9,"uvMax":-9999,"vento":{"manha":-9999,"noite":-9999,"tarde":-9999}}]}}`

func TestParse(t *testing.T) {
	f, err := parse([]byte(sample), time.Unix(0, 0), false)
	if err != nil {
		t.Fatal(err)
	}
	if f.Place != "Pontecesures" || len(f.Days) != 2 {
		t.Fatalf("forecast = %+v", f)
	}
	d := f.Days[0]
	if d.Date != "2026-10-31" || *d.Sky != 101 || *d.Max != 21 || *d.Min != 11 || *d.UV != 3 || d.Warning != nil {
		t.Fatalf("day0 = %+v", d)
	}
	aft := d.Periods[1]
	if aft.Key != "afternoon" || *aft.Sky != 111 || *aft.Rain != 60 || aft.Wind != nil || *aft.Max != 21 {
		t.Fatalf("afternoon = %+v", aft)
	}
	d1 := f.Days[1]
	if d1.Sky != nil || d1.Max != nil || *d1.Min != 9 || d1.Warning == nil || *d1.Warning != 2 || d1.Periods[0].Max != nil {
		t.Fatalf("day1 = %+v", d1)
	}
}

func TestFreshCacheAndOfflineFallback(t *testing.T) {
	hits, up := 0, true
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if r.URL.Query().Get("idConc") != "36044" {
			t.Errorf("query = %v", r.URL.Query())
		}
		if !up {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(sample))
	}))
	defer srv.Close()

	now := time.Now()
	c := NewClient(Options{CacheDir: t.TempDir(), Endpoint: srv.URL})
	c.now = func() time.Time { return now }
	ctx := context.Background()

	if f, err := c.Forecast(ctx, PontecesuresID); err != nil || f.Stale {
		t.Fatalf("first: %+v %v", f, err)
	}
	// Within the fresh window: served from cache, no request.
	if _, err := c.Forecast(ctx, PontecesuresID); err != nil || hits != 1 {
		t.Fatalf("fresh cache: hits=%d err=%v", hits, err)
	}
	// Expired and the network is down: stale copy.
	now = now.Add(2 * time.Hour)
	up = false
	f, err := c.Forecast(ctx, PontecesuresID)
	if err != nil || !f.Stale || hits != 2 {
		t.Fatalf("offline: %+v %v hits=%d", f, err, hits)
	}
}
