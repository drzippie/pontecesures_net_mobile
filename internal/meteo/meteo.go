// Package meteo reads the short-term forecast for a municipality from
// MeteoGalicia (Xunta de Galicia), the same public source used by the
// website's weather page.
//
// Endpoint and code tables: "Servicio JSON coa predición a curto prazo por
// concello" (https://www.meteogalicia.gal/datosred/infoweb/meteo/docs/rss/JSON_Pred_Concello_gl.pdf).
// Sky codes are 101-125 (day) and 201-225 (night); wind codes are 299-332.
// Any value of -9999 means "not available" and is returned as null.
package meteo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/drzippie/pontecesures_net_mobile/internal/cache"
)

// PontecesuresID is the INE code MeteoGalicia uses for Pontecesures.
const PontecesuresID = 36044

const (
	endpoint = "https://servizos.meteogalicia.gal/mgrss/predicion/jsonPredConcellos.action"
	// freshFor: MeteoGalicia updates the forecast a few times a day; within
	// this window the cached copy is served without a request.
	freshFor   = 30 * time.Minute
	notAvail   = -9999
	maxBody    = 1 << 20
	reqTimeout = 15 * time.Second
)

// Forecast is the multi-day forecast for one municipality.
type Forecast struct {
	Place     string `json:"place"`
	Days      []Day  `json:"days"`
	FetchedAt string `json:"fetchedAt"` // RFC 3339
	// Stale is true when the network failed and an older cached copy is
	// returned.
	Stale bool `json:"stale"`
}

// Day is the forecast for one date. Pointer fields are null when
// MeteoGalicia has no value.
type Day struct {
	Date    string   `json:"date"` // YYYY-MM-DD
	Sky     *int     `json:"sky"`  // most representative sky code of the day
	Min     *int     `json:"min"`
	Max     *int     `json:"max"`
	UV      *int     `json:"uv"`
	Warning *int     `json:"warning"` // warning level, null when none
	Periods []Period `json:"periods"` // morning, afternoon, night
}

// Period is one of the three daily time slots.
type Period struct {
	Key  string `json:"key"` // "morning" | "afternoon" | "night"
	Sky  *int   `json:"sky"`
	Rain *int   `json:"rain"` // probability of rain, %
	Wind *int   `json:"wind"`
	Max  *int   `json:"max"`
	Min  *int   `json:"min"`
}

// Client fetches forecasts with an offline fallback cache.
type Client struct {
	http      *http.Client
	cache     *cache.Store
	endpoint  string
	userAgent string
	now       func() time.Time
}

// Options configures a Client.
type Options struct {
	CacheDir  string
	UserAgent string
	Endpoint  string // for tests
}

// NewClient builds a Client.
func NewClient(o Options) *Client {
	if o.Endpoint == "" {
		o.Endpoint = endpoint
	}
	return &Client{
		http:      &http.Client{Timeout: reqTimeout},
		cache:     cache.New(o.CacheDir),
		endpoint:  o.Endpoint,
		userAgent: o.UserAgent,
		now:       time.Now,
	}
}

// Forecast returns the forecast for a municipality (MeteoGalicia id).
func (c *Client) Forecast(ctx context.Context, municipality int) (Forecast, error) {
	u := c.endpoint + "?idConc=" + strconv.Itoa(municipality) + "&request_locale=gl"

	if e, ok := c.cache.Get(u); ok && c.now().Sub(e.StoredAt) < freshFor {
		return parse(e.Body, e.StoredAt, false)
	}

	body, err := c.get(ctx, u)
	if err == nil {
		f, perr := parse(body, c.now(), false)
		if perr == nil {
			c.cache.Put(u, cache.Entry{StoredAt: c.now(), Body: body})
			return f, nil
		}
		err = perr
	}
	if e, ok := c.cache.Get(u); ok {
		return parse(e.Body, e.StoredAt, true)
	}
	return Forecast{}, err
}

func (c *Client) get(ctx context.Context, u string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meteogalicia: HTTP %d", res.StatusCode)
	}
	return io.ReadAll(io.LimitReader(res.Body, maxBody))
}

// raw mirrors the MeteoGalicia payload.
type raw struct {
	PredConcello struct {
		Nome  string `json:"nome"`
		Lista []struct {
			Data      string          `json:"dataPredicion"`
			CeoDia    *int            `json:"ceoDia"`
			Ceo       franxas         `json:"ceo"`
			PChoiva   franxas         `json:"pchoiva"`
			Vento     franxas         `json:"vento"`
			TMax      *int            `json:"tMax"`
			TMin      *int            `json:"tMin"`
			TMaxF     franxas         `json:"tmaxFranxa"`
			TMinF     franxas         `json:"tminFranxa"`
			UVMax     *int            `json:"uvMax"`
			NivelAvis json.RawMessage `json:"nivelAviso"`
		} `json:"listaPredDiaConcello"`
	} `json:"predConcello"`
}

type franxas struct {
	Manha *int `json:"manha"`
	Tarde *int `json:"tarde"`
	Noite *int `json:"noite"`
}

var errEmpty = errors.New("meteogalicia: empty forecast")

func parse(body []byte, fetched time.Time, stale bool) (Forecast, error) {
	var r raw
	if err := json.Unmarshal(body, &r); err != nil {
		return Forecast{}, err
	}
	if len(r.PredConcello.Lista) == 0 {
		return Forecast{}, errEmpty
	}
	f := Forecast{
		Place:     r.PredConcello.Nome,
		FetchedAt: fetched.Format(time.RFC3339),
		Stale:     stale,
	}
	for _, d := range r.PredConcello.Lista {
		day := Day{
			Date:    strings.SplitN(d.Data, "T", 2)[0],
			Sky:     clean(d.CeoDia),
			Min:     clean(d.TMin),
			Max:     clean(d.TMax),
			UV:      clean(d.UVMax),
			Warning: warning(d.NivelAvis),
		}
		day.Periods = []Period{
			{Key: "morning", Sky: clean(d.Ceo.Manha), Rain: clean(d.PChoiva.Manha), Wind: clean(d.Vento.Manha), Max: clean(d.TMaxF.Manha), Min: clean(d.TMinF.Manha)},
			{Key: "afternoon", Sky: clean(d.Ceo.Tarde), Rain: clean(d.PChoiva.Tarde), Wind: clean(d.Vento.Tarde), Max: clean(d.TMaxF.Tarde), Min: clean(d.TMinF.Tarde)},
			{Key: "night", Sky: clean(d.Ceo.Noite), Rain: clean(d.PChoiva.Noite), Wind: clean(d.Vento.Noite), Max: clean(d.TMaxF.Noite), Min: clean(d.TMinF.Noite)},
		}
		f.Days = append(f.Days, day)
	}
	return f, nil
}

func clean(v *int) *int {
	if v == nil || *v == notAvail {
		return nil
	}
	return v
}

// warning accepts null, a number or a numeric string.
func warning(rawVal json.RawMessage) *int {
	s := strings.Trim(strings.TrimSpace(string(rawVal)), `"`)
	if s == "" || s == "null" {
		return nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 || n == notAvail {
		return nil
	}
	return &n
}
