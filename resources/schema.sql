CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users, forum, thread, post, vote, forum_users CASCADE;

DROP FUNCTION IF EXISTS thread_insert();

--
-- USERS
--
CREATE TABLE users (
  id       SERIAL PRIMARY KEY,

  nickname CITEXT NOT NULL UNIQUE,
  email    CITEXT NOT NULL UNIQUE,

  about    TEXT DEFAULT NULL,
  fullname TEXT   NOT NULL
);

CREATE UNIQUE INDEX users_covering_index
  ON users (nickname, email, about, fullname);

CREATE UNIQUE INDEX users_nickname_index
  ON users (nickname);

CREATE UNIQUE INDEX users_email_index
  ON users (email);

--
-- FORUM
--
CREATE TABLE forum (
  id        SERIAL PRIMARY KEY,
  slug      CITEXT  NOT NULL UNIQUE,

  title     TEXT    NOT NULL,
  moderator CITEXT  NOT NULL,

  threads   INTEGER NOT NULL DEFAULT 0,
  posts     BIGINT  NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX forum_slug_index
  ON forum (slug);

CREATE UNIQUE INDEX forum_slug_id_index
  ON forum (slug, id);
--
-- THREAD
--
CREATE TABLE thread (
  id          SERIAL PRIMARY KEY,

  slug        CITEXT  DEFAULT NULL,
  title       TEXT    NOT NULL,
  message     TEXT    NOT NULL,

  forum_id    INTEGER NOT NULL,
  forum_slug  CITEXT  NOT NULL,

  user_id     INTEGER,
  user_nick   CITEXT  NOT NULL,

  created     TIMESTAMPTZ,
  votes_count INTEGER DEFAULT 0
);

CREATE FUNCTION thread_insert()
  RETURNS TRIGGER AS
$BODY$
BEGIN
  UPDATE forum
  SET
    threads = forum.threads + 1
  WHERE slug = NEW.forum_slug;
  RETURN NULL;
END;
$BODY$
LANGUAGE plpgsql;

CREATE TRIGGER on_thread_insert
AFTER INSERT
  ON thread
FOR EACH ROW EXECUTE PROCEDURE thread_insert();

CREATE UNIQUE INDEX thread_slug_index
  ON thread (slug);

CREATE UNIQUE INDEX thread_slug_id_index
  ON thread (slug, id);

CREATE UNIQUE INDEX thread_forum_id_created_index
  ON thread (forum_id, created);

CREATE UNIQUE INDEX thread_id_forum_slug_index
  ON thread (id, forum_slug);

CREATE UNIQUE INDEX thread_slug_forum_slug_index
  ON thread (slug, forum_slug);

CREATE UNIQUE INDEX thread_covering_index
  ON thread (forum_id, created, id, slug, title, message, forum_slug, user_nick, created, votes_count);

--
-- POST
--
CREATE TABLE post (
  id          SERIAL PRIMARY KEY,

  user_nick   TEXT      NOT NULL,

  message     TEXT      NOT NULL,
  created     TIMESTAMPTZ,

  forum_slug  TEXT      NOT NULL,
  thread_id   INTEGER   NOT NULL,

  parent      INTEGER            DEFAULT 0,
  parents     BIGINT [] NOT NULL,
  main_parent BIGINT    NOT NULL,

  is_edited   BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX posts_thread_id_index
  ON post (thread_id, id);
--
-- CREATE INDEX posts_parents_index
--   ON post
--   USING GIN (parents);

CREATE UNIQUE INDEX posts_thread_id_parents
  ON post (id, thread_id, parents);
--
CREATE UNIQUE INDEX posts_thread_id_parents_index
  ON post (thread_id, parents);

CREATE UNIQUE INDEX posts_parents
  ON post (parent, thread_id, parents)
  WHERE parent = 0;


-- CREATE UNIQUE INDEX parent_tree_1
--   ON post (parent, thread_id, main_parent)
--   WHERE parent = 0;
--
-- CREATE UNIQUE INDEX parent_tree_12
--   ON post (parent, thread_id, id, main_parent)
--   WHERE parent = 0;
-- CREATE UNIQUE INDEX parent_tree_13
--   ON post (parent, thread_id, id, main_parent)
--   WHERE parent = 0;
CREATE UNIQUE INDEX ON post (thread_id, id, parent, main_parent)
  WHERE parent = 0;

CREATE UNIQUE INDEX ON post(id, parents);

CREATE UNIQUE INDEX parent_tree_2
  ON post (main_parent, parents);

CREATE UNIQUE INDEX parent_tree_3
  ON post (main_parent, parents, id);
  
CREATE UNIQUE INDEX parent_tree_4
  ON post (id, main_parent);

--
-- VOTE
--
CREATE TABLE vote (
  id         SERIAL PRIMARY KEY,

  user_id    INTEGER NOT NULL,
  thread_id  INTEGER NOT NULL REFERENCES thread,

  voice      INTEGER,
  prev_voice INTEGER DEFAULT 0,
  CONSTRAINT unique_user_and_thread UNIQUE (user_id, thread_id)
);

CREATE UNIQUE INDEX vote_nickname_thread_id_index
  ON vote (user_id, thread_id);


CREATE TABLE forum_users (
  forumId  INTEGER,
  nickname TEXT,
  email    TEXT,

  about    TEXT,
  fullname TEXT
);

CREATE UNIQUE INDEX forum_users_forum_id_nickname_index2
  ON forum_users (forumId, lower(nickname));

CREATE UNIQUE INDEX forum_users_covering_index2
  ON forum_users (forumId, lower(nickname), nickname, email, about, fullname);
