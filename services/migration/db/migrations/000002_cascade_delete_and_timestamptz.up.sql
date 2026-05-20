ALTER TABLE metadata DROP CONSTRAINT metadata_code_fkey;

ALTER TABLE metadata
  ADD CONSTRAINT metadata_code_fkey
  FOREIGN KEY (code) REFERENCES redirection(code)
  ON DELETE CASCADE;

ALTER TABLE metadata ALTER COLUMN created_at TYPE TIMESTAMPTZ;

ALTER TABLE redirection ALTER COLUMN expire_at TYPE TIMESTAMPTZ;