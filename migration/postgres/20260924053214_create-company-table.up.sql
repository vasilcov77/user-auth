BEGIN;

CREATE TABLE IF NOT EXISTS "company"
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    name            TEXT NOT NULL,
    active          BOOLEAN
);

COMMIT;