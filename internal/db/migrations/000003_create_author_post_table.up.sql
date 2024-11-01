CREATE TABLE IF NOT EXISTS Author_Post(
    author_id INT NOT NULL,
    post_id INT NOT NULL,
    PRIMARY KEY (author_id, post_id),
    FOREIGN KEY (author_id) REFERENCES Author(ID),
    FOREIGN KEY (post_id) REFERENCES Post(ID)
)
