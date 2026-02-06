CREATE TABLE Favorite(
    userId UUID,
    underwriterId CHAR(2) NOT NULL,
    CONSTRAINT watchlist_pkey PRIMARY KEY (userid, underwriterId)
);

ALTER TABLE Favorite
ADD CONSTRAINT fk_favorite_users
FOREIGN KEY (userId)
REFERENCES users(id)
ON DELETE CASCADE
ON UPDATE CASCADE;