CREATE TABLE users (
    id serial PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email text UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version INTEGER NOT NULL DEFAULT 1
);
