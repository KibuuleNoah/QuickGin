-- Drop and recreate database
DROP DATABASE IF EXISTS quickt;
CREATE DATABASE quickt
  WITH ENCODING = 'UTF8'
       LC_COLLATE = 'en_US.UTF-8'
       LC_CTYPE   = 'en_US.UTF-8'
       TEMPLATE   = template0;

\connect quickt

-- ─── Extensions ───────────────────────────────────────────────
CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;

-- ─── Trigger Functions ────────────────────────────────────────
CREATE OR REPLACE FUNCTION created_at_column() RETURNS trigger
  LANGUAGE plpgsql AS $$
BEGIN
  NEW.created_at := EXTRACT(EPOCH FROM NOW())::bigint;
  NEW.updated_at := EXTRACT(EPOCH FROM NOW())::bigint;
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION update_at_column() RETURNS trigger
  LANGUAGE plpgsql AS $$
BEGIN
  NEW.updated_at := EXTRACT(EPOCH FROM NOW())::bigint;
  RETURN NEW;
END;
$$;


-- ─── One reusable function ────────────────────────────────────
CREATE OR REPLACE FUNCTION generate_short_id(seq_name TEXT)
RETURNS TEXT
LANGUAGE plpgsql
AS $$
DECLARE
  chars  TEXT    := '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
  base   INTEGER := 62;
  len    INTEGER := 8;
  result TEXT    := '';
  num    BIGINT;
  rem    INTEGER;
BEGIN
  num := nextval(seq_name);
  num := num + (random() * 1000000)::BIGINT * 1000000000;

  WHILE length(result) < len LOOP
    rem    := num % base;
    result := substr(chars, rem + 1, 1) || result;
    num    := num / base;
  END LOOP;

  RETURN result;
END;
$$;


-- ─── One sequence per table ───────────────────────────────────
CREATE SEQUENCE user_id_seq    START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE article_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;


-- ─── Tables ───────────────────────────────────────────────────
CREATE TABLE "user" (
  id  TEXT  PRIMARY KEY DEFAULT generate_short_id('user_id_seq'),
  identifier       VARCHAR      NOT NULL UNIQUE,
  password    VARCHAR      NOT NULL,
  verified    BOOLEAN      NOT NULL DEFAULT false,
  name        VARCHAR      NOT NULL,
  updated_at  BIGINT       NOT NULL DEFAULT 0,
  created_at  BIGINT       NOT NULL DEFAULT 0,
  last_login_at BIGINT
);

CREATE TABLE article (
  id  TEXT  PRIMARY KEY DEFAULT generate_short_id('article_id_seq'),
  user_id     TEXT      REFERENCES "user"(id) ON UPDATE CASCADE ON DELETE CASCADE,
  title       VARCHAR,
  content     TEXT,
  updated_at  BIGINT,
  created_at  BIGINT
);

-- ─── Triggers ─────────────────────────────────────────────────
CREATE TRIGGER create_user_created_at
  BEFORE INSERT ON "user"
  FOR EACH ROW EXECUTE PROCEDURE created_at_column();

CREATE TRIGGER update_user_updated_at
  BEFORE UPDATE ON "user"
  FOR EACH ROW EXECUTE PROCEDURE update_at_column();

CREATE TRIGGER create_article_created_at
  BEFORE INSERT ON article
  FOR EACH ROW EXECUTE PROCEDURE created_at_column();

CREATE TRIGGER update_article_updated_at
  BEFORE UPDATE ON article
  FOR EACH ROW EXECUTE PROCEDURE update_at_column();
