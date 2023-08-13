BEGIN;

ALTER TABLE user_authentications ADD CONSTRAINT uq_user_authentications_provider_social_id UNIQUE (provider, social_id);

ALTER TABLE user_authentications ADD CONSTRAINT uq_user_authentications_user_id_provider UNIQUE (user_id, provider);

COMMIT;