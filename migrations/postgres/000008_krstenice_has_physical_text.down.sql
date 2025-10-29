BEGIN;

ALTER TABLE krstenice
    ALTER COLUMN has_physical_disability DROP DEFAULT,
    ALTER COLUMN has_physical_disability TYPE BOOLEAN
        USING CASE
            WHEN lower(trim(has_physical_disability)) IN ('true','1','yes','y','da','да') THEN TRUE
            WHEN lower(trim(has_physical_disability)) IN ('false','0','no','n','ne','не') THEN FALSE
            ELSE FALSE
        END,
    ALTER COLUMN has_physical_disability SET DEFAULT FALSE;

COMMIT;
