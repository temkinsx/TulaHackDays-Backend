-- Schema: users
CREATE TABLE IF NOT EXISTS users
(
    id
    UUID
    PRIMARY
    KEY
    DEFAULT
    gen_random_uuid
(
),
    email VARCHAR
(
    255
) UNIQUE NOT NULL,
    username VARCHAR
(
    100
) UNIQUE NOT NULL,
    password VARCHAR
(
    255
) NOT NULL,
    first_name VARCHAR
(
    100
),
    last_name VARCHAR
(
    100
),
    avatar VARCHAR
(
    500
),
    points INTEGER DEFAULT 0,
    level INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP
                         WITH TIME ZONE DEFAULT NOW()
    );
