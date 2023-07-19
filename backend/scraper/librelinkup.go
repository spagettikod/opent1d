package scraper

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spagettikod/opent1d/datastore"
	"github.com/spagettikod/opent1d/librelinkup"
)

type LibreLinkupScraper struct {
	db        datastore.Store
	ticket    *librelinkup.Ticket
	patientID string
	running   bool
	log       zerolog.Logger
	interval  time.Duration
	stopCh    chan struct{}
	doneCh    chan struct{}
}

func (s *LibreLinkupScraper) IsRunning() bool {
	return s.running
}

func (s *LibreLinkupScraper) Stop() {
	s.log.Debug().Msg("stopping scraper")
	close(s.stopCh)
	<-s.doneCh
	s.running = false
}

func (s *LibreLinkupScraper) Start() {
	s.log.Debug().Msgf("starting scraper")
	go s.run()
	s.running = true
}

func (s *LibreLinkupScraper) run() {
	defer close(s.doneCh)

	s.scrape()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.scrape()
		case <-s.stopCh:
			s.log.Debug().Msgf("received stop signal, stopping scrape")
			return
		}
	}
}

func (s *LibreLinkupScraper) scrape() {
	s.log.Debug().Msg("fetching graph data")
	_, graph, err := s.ticket.Graph(s.patientID)
	if err != nil {
		s.log.Err(err).Msgf("error while fetching graph data, aborting")
	} else {
		cgms := []datastore.CGMEntry{}
		for _, bg := range graph {
			ts, err := librelinkup.ToTime(bg.FactoryTimestamp)
			if err != nil {
				s.log.Err(err).Msgf("error while converting '%v' into a timestamp", bg.FactoryTimestamp)
			} else {
				cgms = append(cgms, datastore.NewCGMEntry(ts, datastore.Mmoll(bg.Value)))
			}
		}
		if err := s.db.SaveCGM(cgms...); err != nil {
			s.log.Err(err).Msg("could not save CGM data to datastore")
		}
	}
	s.log.Debug().Msgf("finished fetching data, sleeping for %v", s.interval)
}

func NewLibreLinkUpScraper(db datastore.Store, logger zerolog.Logger, interval time.Duration) (*LibreLinkupScraper, error) {
	scraper := &LibreLinkupScraper{
		db:       db,
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
		log:      logger.With().Str("scraper", "LibreLinkUp").Logger(),
		interval: interval,
	}
	scraper.log.Info().Msg("initializing scraper")
	log.Debug().Msg("loading settings from database")
	s, err := scraper.db.GetSettings()
	if err != nil {
		return nil, err
	}
	if !s.IsValid() {
		return nil, fmt.Errorf("could not setup scraper, please update you LibreLinkUp settings")
	}
	scraper.log = scraper.log.With().Str("username", s.LibreLinkUpUsername).Str("region", s.LibreLinkUpRegion).Logger()
	scraper.log.Debug().Msg("signing in to LibreLinkUp")
	endpoint, found := librelinkup.EndpointByRegion(s.LibreLinkUpRegion)
	if !found {
		return nil, fmt.Errorf("invalid endpoint region '%s', can not connectp", s.LibreLinkUpRegion)
	}
	scraper.log = scraper.log.With().Str("endpoint", endpoint.Hostname).Logger()
	if scraper.ticket, err = librelinkup.Login(s.LibreLinkUpUsername, s.LibreLinkUpPassword, endpoint); err != nil {
		return nil, err
	}
	scraper.log.Debug().Msg("successfully signed into LibreLinkUp")
	scraper.log.Debug().Msg("fetching patient identifier")
	conn, err := scraper.ticket.Connections()
	if err != nil {
		scraper.log.Err(err).Msg("error while fetching connections")
		return nil, err
	}
	if len(conn) != 1 {
		err := fmt.Errorf("expected to find one connection but found %v, quiting scraper", len(conn))
		scraper.log.Err(err).Send()
		return nil, err
	}
	scraper.patientID = conn[0].PatientID
	scraper.log.Debug().Msgf("successfully fetched patient identifier, '%s'", scraper.patientID)
	scraper.log = scraper.log.With().Str("patientID", scraper.patientID).Logger()
	return scraper, nil
}
