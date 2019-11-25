package sql

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/mitrickx/otus-golang-2019/23/calendar/internal/domain/entities"
	"go.uber.org/zap"
	"strings"
	"time"
)

const (
	datetimeLayout = "2006-01-02 15:04:05"
)

var ErrorNotFound = errors.New("event not found")

type ErrorEventListErrors struct {
	errs []error
}

func (e *ErrorEventListErrors) Get(i int) error {
	return e.errs[i]
}

func (e *ErrorEventListErrors) Error() string {
	buffer := strings.Builder{}
	buffer.WriteString("some errors happened when get event list: ")
	first := true
	for _, er := range e.errs {
		if first {
			buffer.WriteString(fmt.Sprintf("%s", er))
			first = false
		} else {
			buffer.WriteString(fmt.Sprintf(", %s", er))
		}

	}
	return buffer.String()
}

type Config struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	Timeout  time.Duration
}

func NewConfig(m map[string]string) (*Config, error) {
	keys := []string{"host", "port", "dbname", "user", "password"}
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return nil, fmt.Errorf("`%s` key is missing", key)
		}
	}
	return &Config{
		Host:     m["host"],
		Port:     m["port"],
		DbName:   m["dbname"],
		User:     m["user"],
		Password: m["password"],
	}, nil

}

type EventRow struct {
	Id        int64
	Name      string
	StartTime string `db:"start_time"`
	EndTime   string `db:"end_time"`
}

type Storage struct {
	db      *sqlx.DB
	timeout time.Duration
	logger  *zap.SugaredLogger // for logging rare errors that must not be happened (like on rows.Close)
}

func NewStorage(cfg Config) (*Storage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	db, err := sqlx.Open("pgx", dsn) // *sql.DB
	if err != nil {
		return nil, fmt.Errorf("failed to load driver %w", err)
	}

	var timeout time.Duration
	if cfg.Timeout == 0 {
		timeout = 5 * time.Second
	}

	ctx, _ := context.WithTimeout(context.Background(), timeout)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return &Storage{
		db:      db,
		timeout: timeout,
	}, nil
}

func (s *Storage) AddEvent(event entities.Event) (int, error) {
	query := `INSERT INTO events(name, start_time, end_time) 
				VALUES(:name, :start_time, :end_time)
				RETURNING id`

	ctx, _ := context.WithTimeout(context.Background(), s.timeout)

	/*map[string]string{
		"name":       event.Name(),
		"start_date": event.Start().Format(dateLayout),
		"start_time": event.Start().Format(timeLayout),
		"end_date":   event.End().Format(dateLayout),
		"end_time":   event.End().Format(timeLayout),
	}*/

	eventRow := &EventRow{
		Name:      event.Name(),
		StartTime: event.Start().Format(datetimeLayout),
		EndTime:   event.End().Format(datetimeLayout),
	}

	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to add event: %w", err)
	}

	var id int
	err = stmt.GetContext(ctx, &id, eventRow)
	if err != nil {
		return 0, fmt.Errorf("failed to add event: %w", err)
	}

	return id, nil

}

func (s *Storage) UpdateEvent(id int, event entities.Event) error {
	query := `UPDATE events SET name = :name, 
					start_time = :start_time,
					end_time = :end_time
				WHERE id = :id`

	ctx, _ := context.WithTimeout(context.Background(), s.timeout)

	/*
		map[string]interface{}{
				"id":         id,
				"name":       event.Name(),
				"start_date": event.Start().Format(dateLayout),
				"start_time": event.Start().Format(timeLayout),
				"end_date":   event.End().Format(dateLayout),
				"end_time":   event.End().Format(timeLayout),
			}
	*/

	eventRow := &EventRow{
		Id:        int64(id),
		Name:      event.Name(),
		StartTime: event.Start().Format(datetimeLayout),
		EndTime:   event.End().Format(datetimeLayout),
	}

	result, err := s.db.NamedExecContext(ctx, query, eventRow)
	if err != nil {
		return err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return ErrorNotFound
	}

	return nil
}

func (s *Storage) DeleteEvent(id int) error {
	query := `DELETE FROM events WHERE id = $1`

	ctx, _ := context.WithTimeout(context.Background(), s.timeout)

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return ErrorNotFound
	}

	return nil
}

func (s *Storage) GetEvent(id int) (entities.Event, error) {

	query := buildSelectEventQuery("id = :id")

	events, err := s.getEvents(query, map[string]interface{}{
		"id": id,
	})

	if len(events) > 0 {
		return events[0], nil
	}

	emptyEvent := entities.Event{}
	if err == nil {
		return emptyEvent, entities.StorageErrorEventNotFound
	} else if innerErr, ok := err.(*ErrorEventListErrors); ok {
		return emptyEvent, innerErr.Get(0)
	} else {
		return emptyEvent, err
	}

}

func (s *Storage) GetAllEvents() ([]entities.Event, error) {
	query := buildSelectEventQuery("")
	return s.getEvents(query, map[string]interface{}{})
}

func (s *Storage) GetEventsByPeriod(start *entities.EventTime, end *entities.EventTime) ([]entities.Event, error) {

	// where statement params that will be glued by AND operator
	var where []string

	// bind params
	params := make(map[string]interface{})

	if start != nil {
		params["start_time"] = convertEventTimeToSqlDateTime(*start)
		where = append(where, "start_time >= :start_time")
	}

	if end != nil {
		params["end_time"] = convertEventTimeToSqlDateTime(*end)
		where = append(where, "start_time <= :end_time")
	}

	// build query
	whereStr := strings.Join(where, " AND ")
	query := buildSelectEventQuery(whereStr)

	// get events
	return s.getEvents(query, params)
}

func (s *Storage) Count() (int, error) {
	query := `SELECT COUNT(*) FROM events`
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)
	row := s.db.QueryRowContext(ctx, query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, nil
	}

	return count, nil
}

func (s *Storage) ClearAll() error {
	query := `DELETE FROM events`
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) getEvents(query string, arg interface{}) ([]entities.Event, error) {

	ctx, _ := context.WithTimeout(context.Background(), s.timeout)

	var rows *sqlx.Rows
	var err error

	rows, err = s.db.NamedQueryContext(ctx, query, arg)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil && s.logger != nil {
			s.logger.Errorf("error on rows.Close: %s\n", err)
		}
	}()

	var events []entities.Event
	var errList []error

	for rows.Next() {
		eventRow := &EventRow{}
		err := rows.StructScan(eventRow)
		if err != nil {
			errList = append(errList, err)
			continue
		}

		event, err := convertEventRowToEvent(eventRow)
		if err != nil {
			errList = append(errList, err)
			continue
		}

		events = append(events, *event)
	}

	if errList != nil {
		resErr := &ErrorEventListErrors{
			errs: errList,
		}
		return events, resErr
	}

	return events, nil
}

func convertSqlDateTimeToEventTime(dateTime string) (*entities.EventTime, error) {
	t, err := time.Parse(datetimeLayout, dateTime)
	if err != nil {
		return nil, err
	}
	eventTime := entities.ConvertFromTime(t)
	return &eventTime, nil
}

func convertEventTimeToSqlDateTime(eventTime entities.EventTime) string {
	return eventTime.Format(datetimeLayout)
}

func buildSelectEventQuery(where string) string {
	query := `SELECT 
					id, 
					name, 
					to_char(start_time, 'YYYY-MM-DD HH24::MI::SS') AS start_time, 
					to_char(end_time, 'YYYY-MM-DD HH24::MI::SS') AS end_time
				FROM events `
	if where == "" {
		return query
	} else {
		return query + " WHERE " + where
	}
}

func convertEventRowToEvent(eventRow *EventRow) (*entities.Event, error) {

	start, err := convertSqlDateTimeToEventTime(eventRow.StartTime)
	if err != nil {
		return nil, fmt.Errorf("start datetime preparing error: %w", err)
	}
	end, err := convertSqlDateTimeToEventTime(eventRow.EndTime)
	if err != nil {
		return nil, fmt.Errorf("end datetime preparing error: %w", err)
	}

	event := entities.NewEventWithId(int(eventRow.Id), eventRow.Name, *start, *end)
	return &event, nil
}
