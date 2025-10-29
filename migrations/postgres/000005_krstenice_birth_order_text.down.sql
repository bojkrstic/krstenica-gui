BEGIN;

ALTER TABLE krstenice
    ALTER COLUMN birth_order TYPE INTEGER USING CASE WHEN trim(birth_order) ~ '^[0-9]+$' THEN trim(birth_order)::INTEGER ELSE 0 END;

COMMIT;
