BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    total_seats INTEGER NOT NULL CHECK (total_seats > 0),
    available_seats INTEGER NOT NULL CHECK (available_seats >= 0),
    booking_ttl INTEGER NOT NULL CHECK (booking_ttl > 0), 
    requires_payment BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CHECK (available_seats <= total_seats)
);


CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('pending', 'confirmed', 'canceled')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    CHECK (expires_at > created_at)
);


CREATE INDEX idx_bookings_expires_at ON bookings (expires_at);


CREATE INDEX idx_bookings_event_id ON bookings (event_id);


CREATE UNIQUE INDEX idx_unique_user_event_booking
ON bookings (event_id, user_id)
WHERE status = 'pending';

COMMIT;
