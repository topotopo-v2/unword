CREATE TABLE IF NOT EXISTS words
(
    id
    UUID
    PRIMARY
    KEY,
    word
    TEXT
    NOT
    NULL,
    native_script
    TEXT,
    pronunciation
    TEXT
    NOT
    NULL,
    language
    TEXT
    NOT
    NULL,
    country
    TEXT
    NOT
    NULL,
    country_code
    CHAR
(
    2
) NOT NULL,
    definition TEXT NOT NULL,
    word_date DATE NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW
(
)
    );