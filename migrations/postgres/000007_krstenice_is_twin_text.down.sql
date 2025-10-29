BEGIN;

ALTER TABLE krstenice
    ALTER COLUMN is_twin DROP DEFAULT,
    ALTER COLUMN is_twin TYPE BOOLEAN
        USING CASE
            WHEN lower(trim(is_twin)) IN ('true','1','yes','y','da','да') THEN TRUE
            WHEN lower(trim(is_twin)) IN ('false','0','no','n','ne','не') THEN FALSE
            ELSE FALSE
        END,
    ALTER COLUMN is_twin SET DEFAULT FALSE;

COMMIT;
