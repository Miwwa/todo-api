CREATE TABLE users
(
    id            varchar NOT NULL PRIMARY KEY,
    email         varchar NOT NULL UNIQUE,
    password_hash varchar NOT NULL,
    name          varchar NOT NULL
);

CREATE TABLE todos
(
    id          varchar NOT NULL PRIMARY KEY,
    user_id     varchar NOT NULL REFERENCES users (id),
    title       varchar NOT NULL,
    description varchar
);
