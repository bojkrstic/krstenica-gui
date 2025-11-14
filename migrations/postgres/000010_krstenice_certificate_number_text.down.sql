ALTER TABLE krstenice
    ALTER COLUMN number_of_certificate TYPE INTEGER USING NULLIF(number_of_certificate, '')::INTEGER;
