-- +goose Up
-- Mailbox table
CREATE TABLE mailbox (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  tracking_id TEXT NOT NULL UNIQUE,
  message_id TEXT NOT NULL UNIQUE,
  maildir_path TEXT NOT NULL UNIQUE,
  date_ts INTEGER NOT NULL,
  from_addr TEXT,
  to_addr TEXT,
  cc_addr TEXT,
  bcc_addr TEXT,
  subject TEXT
);

-- Indexes for fast lookup
CREATE INDEX idx_tracking_id ON mailbox (tracking_id);

CREATE INDEX idx_message_id ON mailbox (message_id);

CREATE INDEX idx_maildir_path ON mailbox (maildir_path);

CREATE INDEX idx_from_addr ON mailbox (from_addr);

CREATE INDEX idx_to_addr ON mailbox (to_addr);

-- FTS table for subjects
CREATE VIRTUAL TABLE email_subjects USING fts5 (tracking_id, subject);

-- +goose StatementBegin
CREATE TRIGGER trg_mailbox_insert AFTER INSERT ON mailbox BEGIN
INSERT INTO
  email_subjects (tracking_id, subject)
VALUES
  (NEW.tracking_id, NEW.subject);

END;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER trg_mailbox_delete AFTER DELETE ON mailbox BEGIN
DELETE FROM email_subjects
WHERE
  tracking_id = OLD.tracking_id;

END;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER trg_mailbox_update AFTER
UPDATE ON mailbox WHEN OLD.subject != NEW.subject BEGIN
UPDATE email_subjects
SET
  subject = NEW.subject
WHERE
  tracking_id = NEW.tracking_id;

END;

-- +goose StatementEnd
