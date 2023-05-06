CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE events (
  id SERIAL PRIMARY KEY,
  title VARCHAR,
  event_date DATE,
  start_time TIME WITH TIME ZONE,
  end_time TIME WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
  CONSTRAINT event_times_valid CHECK (end_time > start_time)
);

CREATE INDEX idx_events_event_date ON events (event_date);
CREATE INDEX idx_events_deleted_at ON events (deleted_at);
CREATE INDEX idx_events_title ON events (title);
CREATE OR REPLACE FUNCTION check_overlapping_events() RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM events e
        WHERE e.event_date = NEW.event_date AND e.deleted_at IS NULL
            AND (
                (e.start_time < NEW.start_time AND e.end_time > NEW.start_time)
                OR (e.start_time >= NEW.start_time AND e.start_time < NEW.end_time)
            )
            AND e.id <> NEW.id
    ) THEN
        RAISE EXCEPTION 'Event overlaps with another event on the same day';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_overlapping_events
BEFORE INSERT OR UPDATE ON events
FOR EACH ROW
EXECUTE FUNCTION check_overlapping_events();




