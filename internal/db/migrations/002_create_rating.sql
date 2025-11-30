--评分表

CREATE TABLE IF NOT EXISTS ratings(
    movie_title TEXT NOT NULL,
    rater_id TEXT NOT NULL,
    rating NUMERIC(2,1) NOT NULL CHECK (rating >= 0.5 AND rating <= 5.0 AND MOD((rating * 2)::INTEGER, 1) = 0),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (movie_title, rater_id),
    FOREIGN KEY (movie_title) REFERENCES movies(title) ON DELETE CASCADE
);