-- Seed Data (Admin User, Genres)
-- UP

-- Insert default genres
INSERT INTO genres (name) VALUES
    ('Action'),
    ('Comedy'),
    ('Drama'),
    ('Horror'),
    ('Sci-Fi'),
    ('Romance'),
    ('Thriller'),
    ('Animation');

-- Insert initial admin user (password: admin123 - hashed with bcrypt cost 12)
-- Use this credentials to login: admin@movieapp.com / admin123
INSERT INTO users (email, password_hash, name, role)
VALUES (
    'admin@movieapp.com',
    '$2a$12$D0HUZlipxBfdNtDI7Ywk..UeGoQP9CCZgpTYoZT3xvN/kYPhs28GC',
    'System Admin',
    'admin'
)
ON CONFLICT (email) DO NOTHING;

-- DOWN
DELETE FROM users WHERE role = 'admin';
DELETE FROM genres;
