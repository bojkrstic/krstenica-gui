BEGIN;

-- Dodaj privremenu kolonu
ALTER TABLE krstenice ADD COLUMN book_temp VARCHAR(100);

-- Kopiraj podatke iz originalne kolone u privremenu
UPDATE krstenice SET book_temp = book::VARCHAR;

-- Obriši originalnu kolonu
ALTER TABLE krstenice DROP COLUMN book;

-- Preimenuj privremenu kolonu u originalno ime
ALTER TABLE krstenice RENAME COLUMN book_temp TO book;

-- Dodaj NOT NULL ograničenje (ako je postojalo)
ALTER TABLE krstenice ALTER COLUMN book SET NOT NULL;

COMMIT;