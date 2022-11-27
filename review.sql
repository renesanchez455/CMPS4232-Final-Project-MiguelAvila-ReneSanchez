DROP TABLE IF EXISTS reviews;

CREATE TABLE reviews (
    review_id serial PRIMARY KEY,
    insertion_date date NOT NULL DEFAULT NOW(),
    author_name text NOT NULL,
    book_name text NOT NULL,
    review text NOT NULL
);

INSERT INTO reviews(author_name, book_name, review)
VALUES ('Joe Shmoe', 'Life', 'Amazing book');