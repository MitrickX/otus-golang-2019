CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL,
    contact_id INT NOT NULL,
    before_days SMALLINT NOT NULL DEFAULT 0,
    before_minutes INT NOT NULL DEFAULT 0
);
CREATE INDEX event_id_idx ON notifications (event_id);