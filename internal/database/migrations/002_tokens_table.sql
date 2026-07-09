-- Write your migrate up statements here


CREATE TABLE TOKENS (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token TEXT NOT NULL,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (owner_id) REFERENCES OWNERS(id)
);

---- create above / drop below ----

DROP TABLE IF EXISTS TOKENS;

