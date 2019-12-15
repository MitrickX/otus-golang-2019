package notificaiton

import (
	"encoding/json"
	"github.com/mitrickx/otus-golang-2019/28/calendar/internal/domain/entities"
)

const dateTimeLayout = "2006-01-02 15:04"

// Event main info that will pushed into queue
type EventInfo struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
}

// Extract main event info from biz event entity
func extractEventInfo(event entities.Event) EventInfo {
	return EventInfo{
		Id:    event.Id(),
		Name:  event.Name(),
		Start: event.Start().Time().Format(dateTimeLayout),
		End:   event.End().Time().Format(dateTimeLayout),
	}
}

// serialize event info for queue
func serializeEvent(event EventInfo) ([]byte, error) {
	result, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// un-serialize event info after read from queue
func unSerializeEvent(msg []byte, event *EventInfo) error {
	err := json.Unmarshal(msg, event)
	if err != nil {
		return err
	}
	return nil
}
