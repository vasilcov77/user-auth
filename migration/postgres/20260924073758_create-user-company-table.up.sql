BEGIN;

CREATE TABLE IF NOT EXISTS "user_company"
(
    user_id     UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    company_id  UUID NOT NULL REFERENCES "company"(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    PRIMARY KEY (user_id, company_id)
    );

-- Индексы для ускорения запросов по обеим сторонам связи
CREATE INDEX IF NOT EXISTS idx_user_company_user_id ON "user_company"(user_id);
CREATE INDEX IF NOT EXISTS idx_user_company_company_id ON "user_company"(company_id);

COMMIT;