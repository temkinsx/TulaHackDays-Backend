-- Schema: achievements
CREATE TABLE IF NOT EXISTS achievements
(
    id
    SERIAL
    PRIMARY
    KEY,
    name
    VARCHAR
(
    255
) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    icon VARCHAR
(
    50
),
    points INTEGER NOT NULL DEFAULT 0,
    type VARCHAR
(
    100
) NOT NULL,
    condition TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS user_achievements
(
    id
    SERIAL
    PRIMARY
    KEY,
    user_id
    UUID
    NOT
    NULL
    REFERENCES
    users
(
    id
) ON DELETE CASCADE,
    achievement_id INTEGER NOT NULL REFERENCES achievements
(
    id
)
  ON DELETE CASCADE,
    unlocked_at TIMESTAMP
  WITH TIME ZONE DEFAULT NOW(),
    UNIQUE
(
    user_id,
    achievement_id
)
    );
