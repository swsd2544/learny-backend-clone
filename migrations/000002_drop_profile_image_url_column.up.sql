BEGIN;

ALTER TABLE users DROP COLUMN profile_image_url;

COMMIT;