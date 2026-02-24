CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'USER' CHECK (role IN ('USER', 'ADMIN')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE movies(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    duration INTEGER NOT NULL CHECK (duration > 0),
    genre VARCHAR(100) NOT NULL,
    poster_url TEXT,
    created_at TIMESTAMPZ NOT NULL DEFAULT now()
)
-- Index on title for fast search, Without this, "WHERE title LIKE '%avengers%'" does a full table scan ,
-- With this index, it goes directly to matching rows
CREATE INDEX idx_movies_title
ON movies USING gin(to_tsvector('english', title));

CREATE TABLE showtimes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    starts_at  TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    room VARCHAR(50) NOT NULL
);
-- Composite index: "show me all showtimes for movie X sorted by time" Composite index works best when you filter by BOTH columns
CREATE INDEX idx_showtimes_movie_starts ON showtimes (movie_id, starts_at);

CREATE TABLE seats(
    id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    showtime_id UUID NOT NULL REFERENCES showtimes(id) ON DELETE CASCADE,
    row TEXT NOT NULL,
    number INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE' CHECK (status IN ('AVAILABLE', 'LOCKED', 'RESERVED')),
    UNIQUE (showtime_id, row, number)
)
-- Fast lookup for "give me all AVAILABLE seats for showtime X"
CREATE INDEX idx_seats_showtime_status ON seats (showtime_id, status)

CREATE TABLE reservation(
    id UUID PRIMARY KEY gen_random_uuid(),
    user_id NOT NULL REFERENCES users(id),
    total_price DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'CONFIRMED', 'CANCELLED')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ NOT NULL  -- background job cancels if still PENDING after this

);

CREATE INDEX idx_reservations_user ON reservations (user_id);
-- Composite: background job queries "WHERE status = PENDING AND expires_at < NOW()"
CREATE INDEX idx_reservations_cleanup ON reservations (status, expires_at);
-- RESERVATION_SEATS (join table)
-- LEARN: Why this table?
--   One reservation can have MANY seats.
--   One seat appears in ONE reservation at a time (enforced by seat.status).
--   This is a many-to-many relationship join table.
CREATE TABLE reservation_seats(
    id UUID PRIMARY KEY gen_random_uuid(),
    reservation_id UUID NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
    seat_id        UUID NOT NULL REFERENCES seats(id),
    price          DECIMAL(10, 2) NOT NULL,
    UNIQUE (reservation_id, seat_id) 
);

CREATE INDEX idx_res_seats_reservation ON reservation_seats (reservation_id);
-- admin seed remains