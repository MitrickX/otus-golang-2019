package sql

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/domain/entities"
	"go.uber.org/zap"
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
	Host           string
	Port           string
	DbName         string
	User           string
	Password       string
	ConnectRetries int
}

func NewConfig(m map[string]string) (*Config, error) {
	keys := []string{"host", "port", "dbname", "user", "password"}
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return nil, fmt.Errorf("`%s` key is missing", key)
		}
	}

	connectRetries := 3
	if val, ok := m["connect_retries"]; ok && val != "" {
		var err error
		connectRetries, err = strconv.Atoi(m["connect_retries"])
		if err != nil {
			return nil, fmt.Errorf("connect_retries key error %w", err)
		}
	}

	return &Config{
		Host:           m["host"],
		Port:           m["port"],
		DbName:         m["dbname"],
		User:           m["user"],
		Password:       m["password"],
		ConnectRetries: connectRetries,
	}, nil

}

type EventRow struct {
	Id            int64
	Name          string
	StartTime     string  `db:"start_time"`
	EndTime       string  `db:"end_time"`
	BeforeMinutes *int64  `db:"before_minutes"`
	NotifiedTime  *string `db:"notified_time"`
}

type Storage struct {
	db      *sqlx.DB
	timeout time.Duration
	logger  *zap.SugaredLogger // for logging rare errors that must not be happened (like on rows.Close)
}

func NewStorage(cfg Config) (*Storage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	db, err := sqlx.Open("pgx", dsn) // *sql.DB
	if err != nil {
		return nil, fmt.Errorf("failed to load driver %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var connectErr error
	for i := 0; i < cfg.ConnectRetries; i++ {
		connectErr = db.PingContext(ctx)
		if connectErr == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if connectErr != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", connectErr)
	}

	return &Storage{
		db:      db,
		timeout: time.Duration(5) * time.Second,
	}, nil
}

func (s *Storage) AddEvent(event entities.Event) (int, error) {
	query := `INSERT INTO events(name, start_time, end_time, before_minutes, notified_time) 
				VALUES(:name, :start_time, :end_time, :before_minutes, :notified_time)
				RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	eventRow := convertEventToEventRow(event)

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
	query := `UPDATE events SET 
					name = :name, 
					start_time = :start_time,
					end_time = :end_time,
					before_minutes = :before_minutes,
					notified_time = :notified_time
				WHERE id = :id`

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	newEvent := entities.WithId(event, id)
	eventRow := convertEventToEventRow(newEvent)

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

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

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

func (s *Storage) GetEventsByPeriod(start *entities.DateTime, end *entities.DateTime) ([]entities.Event, error) {

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

//
func (s *Storage) GetEventsToNotify(start *entities.DateTime, end *entities.DateTime) ([]entities.Event, error) {
	// where statement params that will be glued by AND operator
	var where []string

	where = append(where, "before_minutes IS NOT NULL")
	where = append(where, "notified_time IS NULL")

	// bind params
	params := make(map[string]interface{})

	if start != nil {
		params["start_time"] = convertEventTimeToSqlDateTime(*start)
		where = append(where, "(start_time - make_interval(mins => before_minutes) >= :start_time)")
	}

	if end != nil {
		params["end_time"] = convertEventTimeToSqlDateTime(*end)
		where = append(where, "(start_time - make_interval(mins => before_minutes) <= :end_time)")
	}

	// build query
	whereStr := strings.Join(where, " AND ")
	query := buildSelectEventQuery(whereStr)

	// get events
	return s.getEvents(query, params)
}

func (s *Storage) MarkEventAsNotified(id int, time time.Time) error {
	event, err := s.GetEvent(id)
	if err != nil {
		return err
	}
	newEvent := event.Notified(time)
	return s.UpdateEvent(id, newEvent)
}

func (s *Storage) Count() (int, error) {
	query := `SELECT COUNT(*) FROM events`
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

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

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

// Not part of entities.Storage interface, convenient for integration tests, when need to fill data into db
func (s *Storage) InsertEvent(event entities.Event) (int, error) {
	query := `INSERT INTO events(id, name, start_time, end_time, before_minutes, notified_time) 
				VALUES(:id, :name, :start_time, :end_time, :before_minutes, :notified_time)
				RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	eventRow := convertEventToEventRow(event)

	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to insert event: %w", err)
	}

	var id int
	err = stmt.GetContext(ctx, &id, eventRow)
	if err != nil {
		return 0, fmt.Errorf("failed to insert event: %w", err)
	}

	return id, nil

}

// Get values from `pg_stat_user_tables` table for table 'events'
// It is not part of storage interface
func (s *Storage) GetStatValues(fields []string) (map[string]interface{}, error) {

	query := fmt.Sprintf(`SELECT %s FROM pg_stat_user_tables WHERE relname = :relname`, strings.Join(fields, ","))

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"relname": "events",
	})

	if err != nil {
		return nil, fmt.Errorf("failed get stats from table `pg_stat_user_tables`: %w", err)
	}

	result := make(map[string]interface{})

	if rows.Next() {
		err = rows.MapScan(result)
		if err != nil {
			return nil, fmt.Errorf("failed get stats from table `pg_stat_user_tables`: %w", err)
		}
	}

	return result, nil
}

// Get `n_live_tup` value from `pg_stat_user_tables` table for table 'events'
// It is not part of storage interface
func (s *Storage) GetStatValueNLiveTup() (int64, error) {

	field := "n_live_tup"
	fields := []string{field}

	var err error

	stat, err := s.GetStatValues(fields)
	if err != nil {
		return 0, err
	}

	val, ok := stat[field]
	if !ok {
		return 0, fmt.Errorf("unknown field `%s` in stat result", field)
	}

	var resVal int64
	switch v := val.(type) {
	case int64:
		resVal = v
	case int:
		resVal = int64(v)
	case int8:
		resVal = int64(v)
	case int16:
		resVal = int64(v)
	case int32:
		resVal = int64(v)
	case float32:
		resVal = int64(v)
	case float64:
		resVal = int64(v)
	default:
		err = fmt.Errorf("can't convert to int64 value %+v of type %T", val, val)
	}

	return resVal, err
}

func (s *Storage) getEvents(query string, arg interface{}) ([]entities.Event, error) {

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	defer cancel()

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

func convertSqlDateTimeToEventTime(dateTime string) (*entities.DateTime, error) {
	t, err := time.Parse(datetimeLayout, dateTime)
	if err != nil {
		return nil, err
	}
	eventTime := entities.ConvertFromTime(t)
	return &eventTime, nil
}

func convertEventTimeToSqlDateTime(eventTime entities.DateTime) string {
	return eventTime.Format(datetimeLayout)
}

func buildSelectEventQuery(where string) string {
	query := `SELECT 
					id, 
					name, 
					to_char(start_time, 'YYYY-MM-DD HH24::MI::SS') AS start_time, 
					to_char(end_time, 'YYYY-MM-DD HH24::MI::SS') AS end_time,
					before_minutes,
					to_char(notified_time, 'YYYY-MM-DD HH24::MI::SS') AS notified_time
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

	isNotifyingEnabled := false
	beforeMinutes := 0
	if eventRow.BeforeMinutes != nil {
		isNotifyingEnabled = true
		beforeMinutes = int(*eventRow.BeforeMinutes)
	}

	isNotified := false
	notifiedTime := time.Time{}
	if eventRow.NotifiedTime != nil {
		isNotified = true
		notifiedTime, err = time.Parse(datetimeLayout, *eventRow.NotifiedTime)
		if err != nil {
			return nil, fmt.Errorf("notified datetime preparing error: %w", err)
		}
	}

	event := entities.NewDetailedEventWithId(
		int(eventRow.Id),
		eventRow.Name,
		*start,
		*end,
		isNotifyingEnabled,
		beforeMinutes,
		isNotified,
		notifiedTime,
	)

	return &event, nil
}

func convertEventToEventRow(event entities.Event) EventRow {
	eventRow := EventRow{
		Id:        int64(event.Id()),
		Name:      event.Name(),
		StartTime: event.Start().Format(datetimeLayout),
		EndTime:   event.End().Format(datetimeLayout),
	}

	if event.IsNotifyingEnabled() {
		beforeMinutes := int64(event.BeforeMinutes())
		eventRow.BeforeMinutes = &beforeMinutes
	}

	if event.IsNotified() {
		notifiedTime := event.NotifiedTime().Format(datetimeLayout)
		eventRow.NotifiedTime = &notifiedTime
	}

	return eventRow
}
