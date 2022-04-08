CREATE DATABASE IF NOT EXISTS log;

DROP TABLE IF EXISTS log.user_creation;

CREATE TABLE IF NOT EXISTS log.user_creation (
    date Date DEFAULT toDate(0),
    dt DateTime,
    msg String
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (date, dt)
SETTINGS index_granularity = 1024;
