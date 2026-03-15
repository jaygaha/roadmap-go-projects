-- Create Seats Table (for tracking all seats per showtime)
-- UP
CREATE TABLE IF NOT EXISTS seats (
    id SERIAL PRIMARY KEY,
    showtime_id INTEGER NOT NULL REFERENCES showtimes(id) ON DELETE CASCADE,
    seat_number VARCHAR(10) NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(showtime_id, seat_number)
);

CREATE INDEX idx_seats_showtime_id ON seats(showtime_id);
CREATE INDEX idx_seats_is_available ON seats(is_available);

-- DOWN
DROP TABLE IF EXISTS seats;
