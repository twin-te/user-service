DROP TABLE IF EXISTS user_authentications;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_authentications_provider_enum;

CREATE TABLE users (
  id uuid NOT NULL, 
  "createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
  "deletedAt" TIMESTAMP,
  PRIMARY KEY(id)
);

CREATE TYPE user_authentications_provider_enum AS ENUM ('Google', 'Twitter', 'Apple');

CREATE TABLE user_authentications (
  id SERIAL NOT NULL, 
  user_id uuid, 
  provider user_authentications_provider_enum NOT NULL, 
  social_id VARCHAR NOT NULL, 
  PRIMARY KEY(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);
