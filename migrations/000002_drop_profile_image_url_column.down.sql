BEGIN;

ALTER TABLE users ADD profile_image_url text NOT NULL;

COMMIT;