CREATE TABLE books (
    book_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    genre_id INT,
    published_year INT,
    FOREIGN KEY (genre_id) REFERENCES genres(genre_id)
);