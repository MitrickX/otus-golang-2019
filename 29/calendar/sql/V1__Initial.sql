CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
	before_minutes INT NULL DEFAULT NULL,
	notified_time  TIMESTAMP NULL DEFAULT NULL
);
CREATE INDEX start_idx ON events USING btree (start_time, end_time);