BEGIN;

ALTER TABLE krstenice
    ALTER COLUMN is_twin DROP DEFAULT,
    ALTER COLUMN is_twin TYPE VARCHAR(20)
        USING CASE
            WHEN is_twin IS TRUE THEN 'true'
            WHEN is_twin IS FALSE THEN 'false'
            ELSE ''
        END,
    ALTER COLUMN is_twin SET DEFAULT '';

COMMIT;
