BEGIN;

UPDATE users SET character_id = 1 WHERE character_id IS NULL;

ALTER TABLE users
    ALTER COLUMN character_id SET DEFAULT 1,
    ALTER COLUMN character_id SET NOT NULL;

COMMIT;