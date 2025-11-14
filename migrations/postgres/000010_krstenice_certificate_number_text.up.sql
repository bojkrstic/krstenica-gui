ALTER TABLE krstenice
    ALTER COLUMN number_of_certificate TYPE TEXT USING number_of_certificate::TEXT;
