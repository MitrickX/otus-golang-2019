package scheduler

import "github.com/mitrickx/otus-golang-2019/25/calendar/internal/domain/entities"

type Queue interface {
	AddEvent(event entities.Event) error
	ReadEvent() (entities.Event, error)
}
