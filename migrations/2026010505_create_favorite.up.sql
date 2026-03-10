CREATE TABLE Favorites(
    userId UUID,
    underwriterId CHAR(2) NOT NULL,
    CONSTRAINT favorite_pkey PRIMARY KEY (userId, underwriterId),
    CONSTRAINT fk_favorite_users
        FOREIGN KEY (userId)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);