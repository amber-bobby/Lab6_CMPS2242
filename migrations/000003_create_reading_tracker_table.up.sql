CREATE TABLE reading_tracker (
    tracker_id SERIAL PRIMARY KEY,
    book_id INT,
    start_date DATE,
    finish_date DATE,
    status VARCHAR(50),
    rating INT,
    FOREIGN KEY (book_id) REFERENCES books(book_id)
);