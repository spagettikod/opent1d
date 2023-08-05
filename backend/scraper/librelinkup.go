package scraper

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/spagettikod/opent1d/datastore"
	"github.com/spagettikod/opent1d/librelinkup"
)

type LibreLinkupScraper struct {
	db        datastore.Store
	ticket    *librelinkup.Ticket
	username  string
	password  string
	region    string
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
	s.log.Debug().Msg("starting scrape")
	if s.ticket == nil {
		s.log.Debug().Msg("ticket is empty, trying to login")
		if err := s.login(); err != nil {
			s.log.Err(err).Msgf("error occured trying to login to LibreLinkUp, trying again in %v", s.interval)
			return
		}
	}
	scrapeLog := s.log.With().Str("patientID", s.patientID).Logger()
	scrapeLog.Debug().Msg("fetching graph data")
	_, graph, err := s.ticket.Graph(s.patientID)
	if err != nil {
		scrapeLog.Err(err).Msgf("error while fetching graph data, aborting")
	} else {
		cgms := []datastore.CGMEntry{}
		for _, bg := range graph {
			ts, err := librelinkup.ToTime(bg.FactoryTimestamp)
			if err != nil {
				scrapeLog.Err(err).Msgf("error while converting '%v' into a timestamp", bg.FactoryTimestamp)
			} else {
				cgms = append(cgms, datastore.NewCGMEntry(ts, datastore.Mmoll(bg.Value)))
			}
		}
		if err := s.db.SaveCGM(cgms...); err != nil {
			scrapeLog.Err(err).Msg("could not save CGM data to datastore")
		}
	}
	scrapeLog.Debug().Msgf("finished fetching data, sleeping for %v", s.interval)
}

func (scraper *LibreLinkupScraper) login() error {
	scraper.log.Debug().Msg("signing in to LibreLinkUp")
	endpoint, found := librelinkup.EndpointByRegion(scraper.region)
	if !found {
		return fmt.Errorf("invalid endpoint region '%s', can not connectp", scraper.region)
	}
	var err error
	if scraper.ticket, err = librelinkup.Login(scraper.username, scraper.password, endpoint); err != nil {
		return err
	}
	scraper.log.Debug().Msg("successfully signed into LibreLinkUp")
	scraper.log.Debug().Msg("fetching patient identifier")
	conn, err := scraper.ticket.Connections()
	if err != nil {
		scraper.log.Err(err).Msg("error while fetching connections")
		return err
	}
	if len(conn) != 1 {
		err := fmt.Errorf("expected to find one connection but found %v, quiting scraper", len(conn))
		scraper.log.Err(err).Send()
		return err
	}
	scraper.patientID = conn[0].PatientID
	scraper.log.Debug().Msgf("successfully fetched patient identifier, '%s'", scraper.patientID)
	return nil
}

func NewLibreLinkUpScraper(db datastore.Store, username, password, region string, logger zerolog.Logger, interval time.Duration) (*LibreLinkupScraper, error) {
	scraper := &LibreLinkupScraper{
		db:       db,
		username: username,
		password: password,
		region:   region,
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
		log:      logger.With().Str("scraper", "LibreLinkUp").Logger(),
		interval: interval,
	}
	scraper.log = scraper.log.With().Str("username", scraper.username).Str("region", scraper.region).Logger()
	scraper.log.Info().Msg("initializing scraper")
	return scraper, nil
}
