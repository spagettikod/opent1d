package envctx

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/spagettikod/opent1d/datastore"
	"github.com/spagettikod/opent1d/scraper"
)

type Context struct {
	DB             datastore.Store
	Logger         zerolog.Logger
	Scraper        *scraper.LibreLinkupScraper
	ScrapeInterval time.Duration
}

func NewContext(db datastore.Store, log zerolog.Logger) *Context {
	return &Context{
		DB:             db,
		Logger:         log,
		Scraper:        nil,
		ScrapeInterval: 6 * time.Hour,
	}
}
