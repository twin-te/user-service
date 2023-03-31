BEGIN;

CREATE TYPE user_authentications_provider_enum AS ENUM (
    'Google',
    'Twitter',
    'Apple'
);

CREATE TABLE users (
    id uuid NOT NULL,
    "createdAt" timestamp without time zone DEFAULT NOW() NOT NULL,
    CONSTRAINT "PK_a3ffb1c0c8416b9fc6f907b7433" PRIMARY KEY (id)
);

CREATE TABLE user_authentications (
    id serial NOT NULL,
    user_id uuid,
    provider user_authentications_provider_enum NOT NULL,
    social_id varchar NOT NULL,
    CONSTRAINT "PK_5357fb1162b50b926c77290c8bc" PRIMARY KEY (id),
    CONSTRAINT "FK_163ff5c9a502621798f57606e80" FOREIGN KEY (user_id) REFERENCES users(id)
);

COMMIT;