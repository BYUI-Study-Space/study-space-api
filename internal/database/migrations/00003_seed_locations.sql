-- +goose Up
-- Sample study spots. Hours and capacities are placeholders; adjust them to the real campus data.
INSERT INTO locations (name, building, room, description, has_wifi, group_friendly, noise_level, capacity, open_time, close_time) VALUES
    ('Library Quiet Floor',   'David O. McKay Library', '3rd Floor',     'Silent individual study with large desks and outlets.',       TRUE,  FALSE, 'quiet',         60, '07:00', '23:00'),
    ('Library Group Rooms',   'David O. McKay Library', '2nd Floor',     'Bookable rooms with whiteboards and a shared screen.',        TRUE,  TRUE,  'collaborative',  8, '07:00', '23:00'),
    ('Manwaring Commons',     'Manwaring Center',       'Main Floor',    'Open tables near food options; good for casual group work.', TRUE,  TRUE,  'moderate',       40, '07:00', '22:00'),
    ('Hinckley Study Lounge', 'Hinckley Building',      '1st Floor',     'Soft seating and small tables between classes.',             TRUE,  TRUE,  'moderate',       15, '08:00', '21:00'),
    ('Math Study Center',     'Romney Building',        'Ground Floor',  'Tutors on hand and tables for problem sets.',                TRUE,  TRUE,  'collaborative',  30, '08:00', '20:00'),
    ('Gardens Picnic Tables', 'Thomas E. Ricks Gardens','Outdoors',      'Fresh air and shade; no Wi-Fi coverage.',                    FALSE, TRUE,  'collaborative',  12, '08:00', '20:00');

-- +goose Down
DELETE FROM locations;
