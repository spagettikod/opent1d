package main

import (
	"opent1d/librelinkup"
	"time"
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
