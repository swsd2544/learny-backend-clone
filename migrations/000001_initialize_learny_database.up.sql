BEGIN;

-- noinspection SpellCheckingInspection
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS characters (
    id bigserial PRIMARY KEY,
    image_url text NOT NULL,
    rarity text NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    version bigint NOT NULL DEFAULT 1
);

INSERT INTO characters (image_url, rarity) VALUES
    ('https://www.google.com/url?sa=i&url=https%3A%2F%2Fwww.amongusavatarcreator.com%2F&psig=AOvVaw3lmrBdcPO-WvmnkBdpGb1d&ust=1670508872605000&source=images&cd=vfe&ved=0CBAQjRxqFwoTCNiE1dHY5_sCFQAAAAAdAAAAABAE', 'common');

-- noinspection SqlResolve
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username text NOT NULL,
    firstname text NOT NULL,
    lastname text NOT NULL,
    email citext UNIQUE NOT NULL,
    hash_password bytea NOT NULL,
    profile_image_url text NOT NULL,
    coin bigint NOT NULL,
    role text NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    version bigint NOT NULL DEFAULT 1,
    character_id bigint REFERENCES characters(id)
);

CREATE TABLE IF NOT EXISTS classes (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    description text NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    version bigint NOT NULL DEFAULT 1,
    teacher_id bigint REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS enrollments (
    user_id bigint REFERENCES users(id),
    class_id bigint REFERENCES classes(id),
    PRIMARY KEY (user_id, class_id)
);

COMMIT;