-- Schema: places
CREATE TABLE IF NOT EXISTS places
(
    id
    UUID
    PRIMARY
    KEY
    DEFAULT
    gen_random_uuid
(
),
    name VARCHAR
(
    255
) NOT NULL,
    description TEXT,
    type VARCHAR
(
    50
) NOT NULL,
    address VARCHAR
(
    500
) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    phone VARCHAR
(
    50
),
    website VARCHAR
(
    500
),
    images TEXT[],
    average_rating DOUBLE PRECISION DEFAULT 0,
    review_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_by_id UUID NOT NULL REFERENCES users
(
    id
),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP
                         WITH TIME ZONE DEFAULT NOW()
    );
