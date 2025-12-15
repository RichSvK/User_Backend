CREATE TABLE watchlist (
    userid UUID NOT NULL,
    stock CHAR(4) NOT NULL,
    CONSTRAINT watchlist_pkey PRIMARY KEY (userid, stock)
);

ALTER TABLE watchlist
ADD CONSTRAINT fk_watchlist_users
FOREIGN KEY (userid)
REFERENCES users(id)
ON DELETE CASCADE
ON UPDATE CASCADE;