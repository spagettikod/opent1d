package event

import (
	"fmt"

	"github.com/spagettikod/opent1d/envctx"
	"github.com/spagettikod/opent1d/scraper"
)

func OnSettingsSaved(ctx *envctx.Context) {
	elog := ctx.Logger.With().Str("event", "OnSettingsSaved").Logger()
	elog.Debug().Msg("event processing started")
	if ctx.Scraper != nil && ctx.Scraper.IsRunning() {
		elog.Debug().Msg("found a running scraper, initiating shutdown")
		ctx.Scraper.Stop()
	}
	elog.Debug().Msg("creating a new scraper to use the new settings")
	var err error
	ctx.Scraper, err = setupScraper(ctx)
	if err != nil {
		elog.Err(err).Msg("failed to setup scraper")
		return
	}
	elog.Debug().Msg("starting scraper")
	ctx.Scraper.Start()
	elog.Debug().Msg("event processing finished")
}

func OnStartup(ctx *envctx.Context) {
	elog := ctx.Logger.With().Str("event", "OnStartup").Logger()
	elog.Debug().Msg("event processing started")
	elog.Debug().Msg("creating scraper")
	var err error
	ctx.Scraper, err = setupScraper(ctx)
	if err != nil {
		elog.Err(err).Msg("failed to setup scraper")
		return
	}
	elog.Debug().Msg("starting scraper")
	ctx.Scraper.Start()
	elog.Debug().Msg("event processing finished")
}

func setupScraper(ctx *envctx.Context) (*scraper.LibreLinkupScraper, error) {
	s, err := ctx.DB.GetSettings()
	if err != nil {
		return nil, err
	}
	if !s.IsValid() {
		return nil, fmt.Errorf("could not setup scraper, please update you LibreLinkUp settings")
	}
	return scraper.NewLibreLinkUpScraper(ctx.DB, s.LibreLinkUpUsername, s.LibreLinkUpPassword, s.LibreLinkUpRegion, ctx.Logger, ctx.ScrapeInterval)
}
