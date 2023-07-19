package librelinkup

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spagettikod/opent1d/ctx"
	"github.com/spagettikod/opent1d/datastore"
)

func StartScraper(db datastore.Store, log zerolog.Logger) error {
	log.Info().Msg("setting up LibreLinkUp scraper")
	if ctx.IsDebug() {
		log.Debug().Msg("loading settings from database")
	}
	s, err := db.GetSettings()
	if err != nil {
		return err
	}
	if err := validateLibreLinkUpSettings(s); err != nil {
		return err
	}
	slog := log.With().Str("username", s.LibreLinkUpUsername).Str("region", s.LibreLinkUpRegion).Logger()
	var ticket *Ticket
	if ctx.IsDebug() {
		slog.Debug().Msg("signing in to LibreLinkUp")
	}
	endpoint, found := EndpointByRegion(s.LibreLinkUpRegion)
	if !found {
		return fmt.Errorf("invalid endpoint region '%s', can not connect to LibreLinkUp", s.LibreLinkUpRegion)
	}
	if ticket, err = Login(s.LibreLinkUpUsername, s.LibreLinkUpPassword, endpoint); err != nil {
		if errors.Is(err, ErrWrongRegionEndpoint) {
			slog.Info().Msg("endpoint was incorrect, trying to resolve correct user endpoint")
			endpoint, err := FindEndpoint(s.LibreLinkUpUsername, s.LibreLinkUpPassword)
			if err != nil {
				return err
			}
			if ctx.IsDebug() {
				slog.Debug().Msgf("got correct endpoint region '%s' from LibreLinkUp, saving to settings for future use", endpoint.Region)
			}
			s.LibreLinkUpRegion = endpoint.Region
			slog = log.With().Str("username", s.LibreLinkUpUsername).Str("endpoint", s.LibreLinkUpRegion).Logger()
			if err := db.SaveSettings(s); err != nil {
				return err
			}
			if ctx.IsDebug() {
				slog.Debug().Msg("trying to sign in to LibreLinkUp again with the new endpoint")
			}
			if ticket, err = Login(s.LibreLinkUpUsername, s.LibreLinkUpPassword, endpoint); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	slog.Info().Msg("successfully signed into LibreLinkUp")
	go scrape(db, ticket, slog)
	return nil
}

func validateLibreLinkUpSettings(s datastore.Settings) error {
	if strings.TrimSpace(s.LibreLinkUpUsername) == "" {
		return fmt.Errorf("LibreLinkUp username is empty, can not sign in to LibreLinkUp")
	}
	if strings.TrimSpace(s.LibreLinkUpPassword) == "" {
		return fmt.Errorf("LibreLinkUp password for username '%s' is empty, can not sign in to LibreLinkUp", s.LibreLinkUpUsername)
	}
	if strings.TrimSpace(s.LibreLinkUpRegion) == "" {
		return fmt.Errorf("LibreLinkUp region for username '%s' is empty, can not sign in to LibreLinkUp", s.LibreLinkUpRegion)
	}
	return nil
}

func scrape(db datastore.Store, t *Ticket, slog zerolog.Logger) {
	slog.Info().Msg("starting LibreLinkUp scraper")
	conn, err := t.Connections()
	if err != nil {
		log.Err(err).Msg("error while fetching connections from LibreLinkUp, quiting scraper")
		return
	}
	if len(conn) != 1 {
		log.Error().Msgf("expected to find one connection but found %v, quiting scraper", len(conn))
		return
	}
	patientId := conn[0].PatientID
	scrapeLog := slog.With().Str("patientId", patientId).Logger()
	for {
		_, graph, err := t.Graph(patientId)
		if err != nil {
			scrapeLog.Err(err).Msgf("error while fetching LibreLinkUp graph data")
		} else {
			cgms := []datastore.CGMEntry{}
			for _, bg := range graph {
				ts, err := ToTime(bg.FactoryTimestamp)
				if err != nil {
					scrapeLog.Err(err).Msgf("error while converting %v into a timestamp", bg.FactoryTimestamp)
				} else {
					cgms = append(cgms, datastore.NewCGMEntry(ts, datastore.Mmoll(bg.Value)))
				}
			}
			if err := db.SaveCGM(cgms...); err != nil {
				scrapeLog.Err(err).Msg("could not save CGM data to datastore")
			}
		}
		if ctx.IsDebug() {
			scrapeLog.Debug().Msg("finished scraping, sleeping for one hour")
		}
		time.Sleep(time.Hour)
	}
}
