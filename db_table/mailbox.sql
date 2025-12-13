-- mailbox table
CREATE TABLE IF NOT EXISTS mailbox (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    tracking_id   TEXT NOT NULL UNIQUE,
    message_id    TEXT NOT NULL,
    maildir_path  TEXT NOT NULL,
    from_addr     TEXT NOT NULL,
    to_addr       TEXT,
    cc_addr       TEXT,
    bcc_addr      TEXT,
    subject       TEXT,
    body_text     TEXT NOT NULL,
    date_ts       INTEGER NOT NULL,
    status        TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_mailbox_tracking_id ON mailbox(tracking_id);
CREATE INDEX IF NOT EXISTS idx_mailbox_message_id  ON mailbox(message_id);
CREATE INDEX IF NOT EXISTS idx_mailbox_date_ts     ON mailbox(date_ts);
CREATE INDEX IF NOT EXISTS idx_mailbox_maildir_path ON mailbox(maildir_path);

-- mail_logs table
CREATE TABLE IF NOT EXISTS mail_logs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    mailbox_id  INTEGER NOT NULL,
    event       TEXT NOT NULL,
    details     TEXT,
    ts          INTEGER NOT NULL DEFAULT (strftime('%s','now')),
    FOREIGN KEY(mailbox_id) REFERENCES mailbox(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_mail_logs_mailbox_id ON mail_logs(mailbox_id);
CREATE INDEX IF NOT EXISTS idx_mail_logs_event      ON mail_logs(event);
CREATE INDEX IF NOT EXISTS idx_mail_logs_ts         ON mail_logs(ts);

-- Drop FTS table if exists (to handle added columns)
DROP TABLE IF EXISTS mailbox_fts;

-- FTS5 table linked to mailbox (only subject and body_text)
CREATE VIRTUAL TABLE mailbox_fts
USING fts5(
    subject,
    body_text,
    content='mailbox',
    content_rowid='id'
);

-- Rebuild initial FTS index
INSERT INTO mailbox_fts(mailbox_fts) VALUES('rebuild');

-- Trigger: AFTER INSERT
CREATE TRIGGER IF NOT EXISTS mailbox_ai AFTER INSERT ON mailbox
BEGIN
    INSERT INTO mailbox_fts(rowid, subject, body_text)
    VALUES (new.id, new.subject, new.body_text);
END;

-- Trigger: AFTER UPDATE
CREATE TRIGGER IF NOT EXISTS mailbox_au AFTER UPDATE ON mailbox
BEGIN
    UPDATE mailbox_fts
    SET subject = new.subject,
        body_text = new.body_text
    WHERE rowid = old.id;
END;

-- Trigger: AFTER DELETE
CREATE TRIGGER IF NOT EXISTS mailbox_ad AFTER DELETE ON mailbox
BEGIN
    DELETE FROM mailbox_fts WHERE rowid = old.id;
END;
