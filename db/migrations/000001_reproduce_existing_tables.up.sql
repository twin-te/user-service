BEGIN;

CREATE TYPE user_authentications_provider_enum AS ENUM ('Google', 'Twitter', 'Apple');

CREATE TABLE users (
    id uuid NOT NULL,
    "createdAt" timestamp without time zone DEFAULT NOW() NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE user_authentications (
    id serial NOT NULL,
    user_id uuid,
    provider user_authentications_provider_enum NOT NULL,
    social_id varchar NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

COMMIT;