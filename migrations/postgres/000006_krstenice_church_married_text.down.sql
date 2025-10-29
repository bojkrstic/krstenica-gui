BEGIN;

ALTER TABLE krstenice
    ALTER COLUMN is_church_married DROP DEFAULT,
    ALTER COLUMN is_church_married TYPE BOOLEAN
        USING CASE
            WHEN lower(trim(is_church_married)) IN ('true','1','yes','y','da','да') THEN TRUE
            WHEN lower(trim(is_church_married)) IN ('false','0','no','n','ne','не') THEN FALSE
            ELSE FALSE
        END,
    ALTER COLUMN is_church_married SET DEFAULT TRUE;

COMMIT;
