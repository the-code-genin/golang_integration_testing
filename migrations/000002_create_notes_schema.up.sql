CREATE TABLE IF NOT EXISTS core.notes (
    id UUID NOT NULL,
    title VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,

    CONSTRAINT notes_pkey PRIMARY KEY (id)
);