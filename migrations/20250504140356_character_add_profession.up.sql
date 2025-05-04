ALTER TABLE characters
ADD COLUMN profession TEXT CHECK (name <> '');
