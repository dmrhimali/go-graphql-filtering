CREATE TABLE IF NOT EXISTS Author_Friend_Author(
    author_id INT REFERENCES Author(ID),
    friend_author_id INT REFERENCES Author(ID),
    PRIMARY KEY (author_id, friend_author_id)
)
