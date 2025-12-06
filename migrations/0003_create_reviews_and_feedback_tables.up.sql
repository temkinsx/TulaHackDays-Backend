-- Schema: reviews & feedback
CREATE TABLE IF NOT EXISTS reviews
(
    id
    SERIAL
    PRIMARY
    KEY,
    place_id
    UUID
    NOT
    NULL
    REFERENCES
    places
(
    id
) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users
(
    id
)
  ON DELETE CASCADE,
    title VARCHAR
(
    255
),
    content TEXT NOT NULL,
    images TEXT[],
    rating DOUBLE PRECISION NOT NULL CHECK
(
    rating
    >=
    1
    AND
    rating
    <=
    5
),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP
  WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP
  WITH TIME ZONE DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS place_ratings
(
    id
    SERIAL
    PRIMARY
    KEY,
    place_id
    UUID
    NOT
    NULL
    REFERENCES
    places
(
    id
) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users
(
    id
)
  ON DELETE CASCADE,
    cleanliness DOUBLE PRECISION NOT NULL CHECK
(
    cleanliness
    >=
    1
    AND
    cleanliness
    <=
    5
),
    service DOUBLE PRECISION NOT NULL CHECK
(
    service
    >=
    1
    AND
    service
    <=
    5
),
    quality DOUBLE PRECISION NOT NULL CHECK
(
    quality
    >=
    1
    AND
    quality
    <=
    5
),
    value DOUBLE PRECISION NOT NULL CHECK
(
    value
    >=
    1
    AND
    value
    <=
    5
),
    overall DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP
  WITH TIME ZONE DEFAULT NOW(),
    UNIQUE
(
    place_id,
    user_id
)
    );

CREATE TABLE IF NOT EXISTS comments
(
    id
    SERIAL
    PRIMARY
    KEY,
    review_id
    INTEGER
    NOT
    NULL
    REFERENCES
    reviews
(
    id
) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users
(
    id
)
  ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP
  WITH TIME ZONE DEFAULT NOW()
    );
