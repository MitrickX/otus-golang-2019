CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);
CREATE INDEX start_idx ON events USING btree (start_time, end_time);