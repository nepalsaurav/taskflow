import notmuch

maildir_path = "~/Maildir"


db = notmuch.Database(maildir_path, mode=notmuch.Database.MODE_READ_WRITE)


print(db)
