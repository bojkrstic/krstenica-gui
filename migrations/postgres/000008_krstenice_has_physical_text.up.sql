BEGIN;

ALTER TABLE krstenice
    ALTER COLUMN has_physical_disability DROP DEFAULT,
    ALTER COLUMN has_physical_disability TYPE VARCHAR(20)
        USING CASE
            WHEN has_physical_disability IS TRUE THEN 'true'
            WHEN has_physical_disability IS FALSE THEN 'false'
            ELSE ''
        END,
    ALTER COLUMN has_physical_disability SET DEFAULT '';

COMMIT;
