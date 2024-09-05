CREATE TABLE users
(
    id            varchar NOT NULL PRIMARY KEY,
    email         varchar NOT NULL UNIQUE,
    password_hash varchar NOT NULL,
    name          varchar NOT NULL
);
