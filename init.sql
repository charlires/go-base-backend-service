CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL
);

CREATE TABLE tracks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    album TEXT NOT NULL
);

CREATE TABLE playlists (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE playlist_tracks (
    playlist_id TEXT NOT NULL,
    track_id TEXT NOT NULL,
    PRIMARY KEY (playlist_id, track_id),
    FOREIGN KEY (playlist_id) REFERENCES playlists(id),
    FOREIGN KEY (track_id) REFERENCES tracks(id)
);

INSERT INTO users (id, username) VALUES
    ('user-001', 'carlos'),
    ('user-002', 'maria'),
    ('user-003', 'jordan'),
    ('user-004', 'sofia'),
    ('user-005', 'liam');

INSERT INTO tracks (id, title, artist, album) VALUES
    ('track-001', 'Midnight Run', 'The Signals', 'City Lights'),
    ('track-002', 'Golden Skies', 'Ava Harper', 'Sunrise Dreams'),
    ('track-003', 'Ocean Glass', 'Blue Harbor', 'Tides'),
    ('track-004', 'Static Hearts', 'Neon Avenue', 'After Hours'),
    ('track-005', 'Paper Planes', 'Luca Stone', 'Wanderlust'),
    ('track-006', 'Northern Glow', 'Elsa Reed', 'Polaris'),
    ('track-007', 'Velvet Echo', 'The Signals', 'City Lights'),
    ('track-008', 'Fireline', 'Atlas North', 'Wild Roads'),
    ('track-009', 'Slow Orbit', 'Mila Grey', 'Gravity'),
    ('track-010', 'Summer Static', 'Neon Avenue', 'After Hours');

INSERT INTO playlists (id, name, description, owner_id) VALUES
    ('playlist-001', 'Morning Boost', 'Upbeat tracks for a productive start.', 'user-001'),
    ('playlist-002', 'Late Night Drive', 'Moody songs for long evening rides.', 'user-002'),
    ('playlist-003', 'Focus Flow', 'Steady tracks for deep work sessions.', 'user-003'),
    ('playlist-004', 'Weekend Vibes', 'Easygoing favorites for relaxing weekends.', 'user-004'),
    ('playlist-005', 'Indie Mix', 'A handpicked collection of indie tracks.', 'user-005');

INSERT INTO playlist_tracks (playlist_id, track_id) VALUES
    ('playlist-001', 'track-001'),
    ('playlist-001', 'track-002'),
    ('playlist-001', 'track-006'),
    ('playlist-002', 'track-004'),
    ('playlist-002', 'track-008'),
    ('playlist-002', 'track-010'),
    ('playlist-003', 'track-003'),
    ('playlist-003', 'track-009'),
    ('playlist-004', 'track-005'),
    ('playlist-004', 'track-007'),
    ('playlist-005', 'track-002'),
    ('playlist-005', 'track-004'),
    ('playlist-005', 'track-010');

