BEGIN;

ALTER TABLE krstenice
    ALTER COLUMN is_church_married DROP DEFAULT,
    ALTER COLUMN is_church_married TYPE VARCHAR(20)
        USING CASE
            WHEN is_church_married IS TRUE THEN 'true'
            WHEN is_church_married IS FALSE THEN 'false'
            ELSE ''
        END,
    ALTER COLUMN is_church_married SET DEFAULT '';

COMMIT;
