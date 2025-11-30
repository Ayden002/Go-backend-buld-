--创建电影表

CREATE TABLE if NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    release_date DATE NOT NULL,
    genre TEXT NOT NULL,
    distributor TEXT,
    budget BIGINT,
    mpa_rating TEXT,

    -- 票房信息字段
    boxoffice_revenue_worldwide BIGINT,
    boxoffice_revenue_opening_weekend_usa BIGINT,
    boxoffice_currency TEXT,
    boxoffice_source TEXT,
    boxoffice_last_updated TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);