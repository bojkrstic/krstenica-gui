BEGIN;

ALTER TABLE krstenice
    ALTER COLUMN birth_order TYPE TEXT USING birth_order::text;

COMMIT;
