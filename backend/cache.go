package main

import (
	"time"

	"github.com/spagettikod/opent1d/librelinkup"
)

type Cache struct {
	Timestamp time.Time
}

type LibreTicketCache struct {
	Cache
	Ticket *librelinkup.Ticket
}

type ActiveSensorCache struct {
	Cache
	ActiveSensor librelinkup.ActiveSensor
}
