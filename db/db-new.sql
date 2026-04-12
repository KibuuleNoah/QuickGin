-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- ─── Tables ───────────────────────────────────────────────────

CREATE TABLE user (
  id TEXT PRIMARY KEY,
  identifier TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  verified INTEGER NOT NULL DEFAULT 0,
  name TEXT NOT NULL,
  updated_at INTEGER NOT NULL DEFAULT (strftime('%s','now')),
  created_at INTEGER NOT NULL DEFAULT (strftime('%s','now')),
  last_login_at INTEGER
);

CREATE TABLE article (
  id TEXT PRIMARY KEY,
  user_id TEXT,
  title TEXT,
  content TEXT,
  updated_at INTEGER DEFAULT (strftime('%s','now')),
  created_at INTEGER DEFAULT (strftime('%s','now')),
  FOREIGN KEY (user_id) REFERENCES user(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- ─── Triggers for updated_at ──────────────────────────────────

CREATE TRIGGER update_user_updated_at
AFTER UPDATE ON user
FOR EACH ROW
BEGIN
  UPDATE user
  SET updated_at = strftime('%s','now')
  WHERE id = OLD.id;
END;

CREATE TRIGGER update_article_updated_at
AFTER UPDATE ON article
FOR EACH ROW
BEGIN
  UPDATE article
  SET updated_at = strftime('%s','now')
  WHERE id = OLD.id;
END;
