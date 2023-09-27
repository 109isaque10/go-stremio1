package stremio

import "time"

// Manifest describes the capabilities of the addon.
// See https://github.com/Stremio/stremio-addon-sdk/blob/f6f1f2a8b627b9d4f2c62b003b251d98adadbebe/docs/api/responses/manifest.md
type Manifest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	// One of the following is required
	// Note: Can only have one in code because of how Go (de-)serialization works
	//Resources     []string       `json:"resources,omitempty"`
	ResourceItems []ResourceItem `json:"resources,omitempty"`

	Types    []string      `json:"types"` // Stremio supports "movie", "series", "channel" and "tv"
	Catalogs []CatalogItem `json:"catalogs"`

	// Optional
	IDprefixes    []string      `json:"idPrefixes,omitempty"`
	Background    string        `json:"background,omitempty"` // URL
	Logo          string        `json:"logo,omitempty"`       // URL
	ContactEmail  string        `json:"contactEmail,omitempty"`
	BehaviorHints BehaviorHints `json:"behaviorHints,omitempty"`
}

// clone returns a deep copy of m.
// We're not using one of the deep copy libraries because only few are maintained and even they have issues.
func (m Manifest) clone() Manifest {
	var resourceItems []ResourceItem
	if m.ResourceItems != nil {
		resourceItems = make([]ResourceItem, len(m.ResourceItems))
		for i, resourceItem := range m.ResourceItems {
			resourceItems[i] = resourceItem.clone()
		}
	}

	var types []string
	if m.Types != nil {
		types = make([]string, len(m.Types))
		for i, t := range m.Types {
			types[i] = t
		}
	}

	var catalogs []CatalogItem
	if m.Catalogs != nil {
		catalogs = make([]CatalogItem, len(m.Catalogs))
		for i, catalog := range m.Catalogs {
			catalogs[i] = catalog.clone()
		}
	}

	var idPrefixes []string
	if m.IDprefixes != nil {
		idPrefixes = make([]string, len(m.IDprefixes))
		for i, idPrefix := range m.IDprefixes {
			idPrefixes[i] = idPrefix
		}
	}

	return Manifest{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Version:     m.Version,

		ResourceItems: resourceItems,

		Types:    types,
		Catalogs: catalogs,

		IDprefixes:    idPrefixes,
		Background:    m.Background,
		Logo:          m.Logo,
		ContactEmail:  m.ContactEmail,
		BehaviorHints: m.BehaviorHints,
	}
}

type ResourceItem struct {
	Name  string   `json:"name"`
	Types []string `json:"types"` // Stremio supports "movie", "series", "channel" and "tv"

	// Optional
	IDprefixes []string `json:"idPrefixes,omitempty"`
}

func (ri ResourceItem) clone() ResourceItem {
	var types []string
	if ri.Types != nil {
		types = make([]string, len(ri.Types))
		for i, t := range ri.Types {
			types[i] = t
		}
	}

	var idPrefixes []string
	if ri.IDprefixes != nil {
		idPrefixes = make([]string, len(ri.IDprefixes))
		for i, idPrefix := range ri.IDprefixes {
			idPrefixes[i] = idPrefix
		}
	}

	return ResourceItem{
		Name:  ri.Name,
		Types: types,

		IDprefixes: idPrefixes,
	}
}

type BehaviorHints struct {
	// Note: Must include `omitempty`, otherwise it will be included if this struct is used in another one, even if the field of the containing struct is marked as `omitempty`
	Adult        bool `json:"adult,omitempty"`
	P2P          bool `json:"p2p,omitempty"`
	Configurable bool `json:"configurable,omitempty"`
	// If you set this to true, it will be true for the "/manifest.json" endpoint, but false for the "/:userData/manifest.json" endpoint, because otherwise Stremio won't show the "Install" button in its UI.
	ConfigurationRequired bool `json:"configurationRequired,omitempty"`
}

// CatalogItem represents a catalog.
type CatalogItem struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`

	// Optional
	Extra []ExtraItem `json:"extra,omitempty"`
}

func (ci CatalogItem) clone() CatalogItem {
	var extras []ExtraItem
	if ci.Extra != nil {
		extras = make([]ExtraItem, len(ci.Extra))
		for i, extra := range ci.Extra {
			extras[i] = extra.clone()
		}
	}

	return CatalogItem{
		Type: ci.Type,
		ID:   ci.ID,
		Name: ci.Name,

		Extra: extras,
	}
}

type ExtraItem struct {
	Name string `json:"name"`

	// Optional
	IsRequired   bool     `json:"isRequired,omitempty"`
	Options      []string `json:"options,omitempty"`
	OptionsLimit int      `json:"optionsLimit,omitempty"`
}

func (ei ExtraItem) clone() ExtraItem {
	var options []string
	if ei.Options != nil {
		options = make([]string, len(ei.Options))
		for i, option := range ei.Options {
			options[i] = option
		}
	}

	return ExtraItem{
		Name: ei.Name,

		IsRequired:   ei.IsRequired,
		Options:      options,
		OptionsLimit: ei.OptionsLimit,
	}
}

// MetaPreviewItem represents a meta preview item and is meant to be used within catalog responses.
// See https://github.com/Stremio/stremio-addon-sdk/blob/f6f1f2a8b627b9d4f2c62b003b251d98adadbebe/docs/api/responses/meta.md#meta-preview-object
type MetaPreviewItem struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Poster string `json:"poster"` // URL

	// Optional
	PosterShape string `json:"posterShape,omitempty"`

	// Optional, used for the "Discover" page sidebar
	Genres      []string       `json:"genres,omitempty"`   // Will be replaced by Links at some point
	Director    []string       `json:"director,omitempty"` // Will be replaced by Links at some point
	Cast        []string       `json:"cast,omitempty"`     // Will be replaced by Links at some point
	Links       []MetaLinkItem `json:"links,omitempty"`    // For genres, director, cast and potentially more. Not fully supported by Stremio yet!
	IMDbRating  string         `json:"imdbRating,omitempty"`
	ReleaseInfo string         `json:"releaseInfo,omitempty"` // E.g. "2000" for movies and "2000-2014" or "2000-" for TV shows
	Description string         `json:"description,omitempty"`

	Background     string          `json:"background,omitempty"` // URL
	Logo           string          `json:"logo,omitempty"`       // URL
	Videos         []VideoItem     `json:"videos,omitempty"`
	Year           string          `json:"year,omitempty"`
	Writer         []string        `json:"writer,omitempty"`
	Country        string          `json:"country,omitempty"`
	Runtime        string          `json:"runtime,omitempty"`
	TrailerStreams []TrailerStream `json:"trailerStreams,omitempty"`
	Trailers       []Trailers      `json:"trailers,omitempty"`
	Slug           string          `json:"slug,omitempty"`
	Status         string          `json:"status,omitempty"`
	IMDBId         string          `json:"imdb_id,omitempty"`
}

type Trailers struct {
	Source string `json:"source,omitempty"`
	Type   string `json:"type,omitempty"`
}

type TrailerStream struct {
	Title       string `json:"title,omitempty"`
	YouTubeID   string `json:"ytId,omitempty"`
	Url         string `json:"url,omitempty"`
	InfoHash    string `json:"infoHash,omitempty"`
	FileIndex   int    `json:"fileIdx,omitempty"`
	ExternalUrl string `json:"externalUrl,omitempty"`
}

type MetaItemBehaviorHints struct {
	DefaultVideoID     *string `json:"defaultVideoId"`
	HasScheduledVideos bool    `json:"hasScheduledVideos,omitempty"`
}

// MetaItem represents a meta item and is meant to be used when info for a specific item was requested.
// See https://github.com/Stremio/stremio-addon-sdk/blob/f6f1f2a8b627b9d4f2c62b003b251d98adadbebe/docs/api/responses/meta.md
type MetaItem struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`

	// Optional
	Genres         []string        `json:"genres,omitempty"`   // Will be replaced by Links at some point
	Director       []string        `json:"director,omitempty"` // Will be replaced by Links at some point
	Cast           []string        `json:"cast,omitempty"`     // Will be replaced by Links at some point
	Links          []MetaLinkItem  `json:"links,omitempty"`    // For genres, director, cast and potentially more. Not fully supported by Stremio yet!
	Poster         string          `json:"poster,omitempty"`   // URL
	PosterShape    string          `json:"posterShape,omitempty"`
	Background     string          `json:"background,omitempty"` // URL
	Logo           string          `json:"logo,omitempty"`       // URL
	Description    string          `json:"description,omitempty"`
	ReleaseInfo    string          `json:"releaseInfo,omitempty"` // E.g. "2000" for movies and "2000-2014" or "2000-" for TV shows
	IMDbRating     string          `json:"imdbRating,omitempty"`
	Released       string          `json:"released,omitempty"` // Must be ISO 8601, e.g. "2010-12-06T05:00:00.000Z"
	Videos         []VideoItem     `json:"videos,omitempty"`
	Runtime        string          `json:"runtime,omitempty"`
	Slug           string          `json:"slug,omitempty"`
	Status         string          `json:"status,omitempty"`
	TrailerStreams []TrailerStream `json:"trailerStreams,omitempty"`
	Language       string          `json:"language,omitempty"`
	Country        string          `json:"country,omitempty"`
	Awards         string          `json:"awards,omitempty"`
	Website        string          `json:"website,omitempty"` // URL

	// TODO: behaviorHints
	//BehaviorHints MetaItemBehaviorHints `json:"behaviorHints,omitempty"`
}

// MetaLinkItem links to a page within Stremio.
// It will at some point replace the usage of `genres`, `director` and `cast`.
// Note: It's not fully supported by Stremio yet (not fully on PC and not at all on Android)!
type MetaLinkItem struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	URL      string `json:"url"` //  // URL. Can be "Meta Links" (see https://github.com/Stremio/stremio-addon-sdk/blob/f6f1f2a8b627b9d4f2c62b003b251d98adadbebe/docs/api/responses/meta.links.md)
}

type VideoItem struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Released time.Time `json:"released"` // Must be ISO 8601, e.g. "2010-12-06T05:00:00.000Z"

	Number      int       `json:"number,omitempty"`
	Description string    `json:"description,omitempty"`
	FirstAired  time.Time `json:"firstAired,omitempty"`

	// Optional
	Thumbnail string       `json:"thumbnail,omitempty"` // URL
	Streams   []StreamItem `json:"streams,omitempty"`
	Available bool         `json:"available,omitempty"`
	Episode   int          `json:"episode,omitempty"`
	Season    int          `json:"season,omitempty"`
	Trailer   string       `json:"trailer,omitempty"` // Youtube ID
	Overview  string       `json:"overview,omitempty"`

	// TO DELETE:
	Rating string `json:"rating,omitempty"`
	TvdbId int    `json:"tvdb_id,omitempty"`
}

type Subtitles struct {
	Id       string `json:"id,omitempty"`
	URL      string `json:"url,omitempty"`
	Language string `json:"language,omitempty"`
}

type ProxyHeaders struct {
	Request map[string]string `json:"request,omitempty"`
}

type StreamItemBehaviorHints struct {
	CountryWhitelist string       `json:"countryWhitelist,omitempty"`
	NotWebReady      bool         `json:"notWebReady,omitempty"`
	BingeGroup       string       `json:"bingeGroup,omitempty"`
	ProxyHeaders     ProxyHeaders `json:"proxyHeaders,omitempty"`
}

// StreamItem represents a stream for a MetaItem.
// See https://github.com/Stremio/stremio-addon-sdk/blob/d1915074439bf152c0c0f1a7603ccf93c05a1f89/docs/api/responses/stream.md
type StreamItem struct {
	// One of the following is required
	URL         string `json:"url,omitempty"` // URL
	YoutubeID   string `json:"ytId,omitempty"`
	InfoHash    string `json:"infoHash,omitempty"`
	ExternalURL string `json:"externalUrl,omitempty"` // URL

	// Optional
	Name        string      `json:"name,omitempty"`  // Usually used for stream quality
	Title       string      `json:"title,omitempty"` // NOTE: Soon to be deprecated in favor of description
	Description string      `json:"description,omitempty"`
	FileIndex   uint8       `json:"fileIdx,omitempty"` // Only when using InfoHash
	Subtitles   []Subtitles `json:"subtitles,omitempty"`

	// TODO: subtitles
	BehaviorHints *StreamItemBehaviorHints `json:"behaviorHints,omitempty"`
}
