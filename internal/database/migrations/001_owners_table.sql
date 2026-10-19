-- Write your migrate up statements here

CREATE TABLE OWNERS (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


---- create above / drop below ----


DROP TABLE IF EXISTS OWNERS;


-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
