package services

import "time"

type ClientClock interface {
	Now() time.Time
}
