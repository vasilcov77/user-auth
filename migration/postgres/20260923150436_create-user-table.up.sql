BEGIN;

CREATE TABLE IF NOT EXISTS "user"
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    email            VARCHAR(255) NOT NULL,
    password         VARCHAR(255) NOT NULL,
    phone            VARCHAR(20),
    active           BOOLEAN,
    verified_email   BOOLEAN,
    verified_phone   BOOLEAN
);

COMMIT;