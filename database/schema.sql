DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS thread CASCADE;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS vote CASCADE;

CREATE TABLE users (
    nickname varchar(30) PRIMARY KEY,

    email varchar(40) UNIQUE,
    fullname varchar(30) NOT NULL,
    about text
);

CREATE TABLE forum (
    slug varchar(30) PRIMARY KEY,

    title varchar(30) NOT NULL,
    posts BIGINT,
    threads INT,
    user_nickname varchar(40) NOT NULL REFERENCES users
);

CREATE TABLE thread (
    id SERIAL PRIMARY KEY,

    slug varchar(30) UNIQUE,
    title varchar(30) NOT NULL,
    message TEXT NOT NULL,
    forum_slug varchar(30) NOT NULL REFERENCES forum,
    author_nickname varchar(30) NOT NULL REFERENCES users,
    created TIMESTAMP WITH TIME ZONE,
    votes_count INTEGER
);

CREATE TABLE post (
    id SERIAL PRIMARY KEY,

    author_nickname varchar(30) NOT NULL REFERENCES users,
    message TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE,
    forum_slug varchar(30) REFERENCES forum,
    thread_id INTEGER REFERENCES thread,
    is_edited BOOLEAN NOT NULL,
    parent INTEGER
);

CREATE TABLE vote (
    id SERIAL PRIMARY KEY,

    nickname varchar(30) REFERENCES users,
    slug_id INTEGER REFERENCES thread
);
