BEGIN;

DROP INDEX IF EXISTS idx_unique_user_event_booking;
DROP INDEX IF EXISTS idx_bookings_event_id;
DROP INDEX IF EXISTS idx_bookings_expires_at;

DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS users;

COMMIT;