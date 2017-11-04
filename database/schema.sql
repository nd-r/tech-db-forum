SET SYNCHRONOUS_COMMIT = 'off';

DROP INDEX IF EXISTS users_email_index;
DROP INDEX IF EXISTS users_nickname_index;
DROP INDEX IF EXISTS forum_slug_index;
DROP INDEX IF EXISTS vote_nickname_threadid_index;
DROP INDEX IF EXISTS threadSlugIndex;
DROP INDEX IF EXISTS posts_thread_id_index;
DROP INDEX IF EXISTS posts_parents_index;
DROP INDEX IF EXISTS posts_thread_id_parents_index;
DROP INDEX IF EXISTS posts_parents;
DROP INDEX IF EXISTS forum_users_forum_id_user_id_index;
DROP INDEX IF EXISTS thread_forum_slug_index;
DROP INDEX IF EXISTS thread_forum_slug_created_index;

DROP TRIGGER IF EXISTS on_thread_insert
ON thread;
DROP TRIGGER IF EXISTS on_thread_insert_user
ON thread;
DROP TRIGGER IF EXISTS on_post_insert_user
ON post;

DROP FUNCTION IF EXISTS findForum(myslug TEXT );
DROP FUNCTION IF EXISTS thread_insert() CASCADE;
DROP FUNCTION IF EXISTS forum_users_update() CASCADE;

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS thread CASCADE;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS vote CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;

-- USERS
CREATE TABLE users (
  id       SERIAL PRIMARY KEY,
  about    TEXT DEFAULT NULL,
  email    TEXT NOT NULL UNIQUE,
  fullname TEXT NOT NULL,
  nickname TEXT NOT NULL UNIQUE
);

CREATE UNIQUE INDEX users_nickname_index
  ON users (lower(nickname));

CREATE UNIQUE INDEX users_email_index
  ON users (lower(email));

--
-- FORUM
--
CREATE TABLE forum (
  id        SERIAL PRIMARY KEY,
  slug      TEXT    NOT NULL UNIQUE,
  title     TEXT    NOT NULL,
  posts     BIGINT  NOT NULL DEFAULT 0,
  threads   INTEGER NOT NULL DEFAULT 0,
  moderator TEXT    NOT NULL
);

CREATE UNIQUE INDEX forum_slug_index
  ON forum (lower(slug));

--
-- THREAD
--
CREATE TABLE thread (
  id          SERIAL PRIMARY KEY,

  slug        TEXT UNIQUE DEFAULT NULL,
  title       TEXT NOT NULL,
  message     TEXT NOT NULL,
  forum_slug  TEXT NOT NULL,
  user_nick   TEXT NOT NULL,
  created     TIMESTAMP WITH TIME ZONE,
  votes_count INTEGER     DEFAULT 0
);

CREATE UNIQUE INDEX threadSlugIndex
  ON thread (lower(slug));

CREATE INDEX thread_forum_slug_index
  ON thread (lower(forum_slug));

CREATE INDEX thread_forum_slug_created_index
  ON thread (lower(forum_slug), created);
CREATE INDEX thread_forum_slug_created_asc_index
  ON thread (lower(forum_slug), created ASC);

CREATE FUNCTION thread_insert()
  RETURNS TRIGGER AS '
BEGIN
  UPDATE forum
  SET
    threads = forum.threads + 1
  WHERE slug = NEW.forum_slug;
  RETURN NULL;
END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_thread_insert
AFTER INSERT
  ON thread
FOR EACH ROW EXECUTE PROCEDURE thread_insert();

--
-- POST
--
CREATE TABLE post (
  id         SERIAL PRIMARY KEY,

  user_nick  TEXT       NOT NULL,
  message    TEXT       NOT NULL,
  created    TIMESTAMP WITH TIME ZONE,
  forum_slug TEXT       NOT NULL,
  thread_id  INTEGER    NOT NULL REFERENCES thread,
  is_edited  BOOLEAN    NOT NULL DEFAULT FALSE,
  parent     INTEGER             DEFAULT 0,
  parents    INTEGER [] NOT NULL
);

CREATE UNIQUE INDEX posts_thread_id_index
  ON post (thread_id, id);


CREATE INDEX posts_parents_index
  ON post
  USING GIN (parents);

CREATE UNIQUE INDEX posts_thread_id_parents_index
  ON post (thread_id, parents);
CREATE UNIQUE INDEX posts_parents
  ON post (parent, thread_id, parents);

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
CREATE UNIQUE INDEX vote_nickname_threadid_index
  ON vote (user_id, thread_id);


CREATE TABLE forum_users (
  forumId  INTEGER,
  nickname TEXT,
  about    TEXT,
  email    TEXT,
  fullname TEXT

);

CREATE UNIQUE INDEX fu_forumid_usernick_index
  ON forum_users (forumid, lower(nickname));

CREATE INDEX fu_forumid
  ON forum_users (forumId);

CREATE FUNCTION forum_users_update()
  RETURNS TRIGGER AS 'BEGIN WITH userinfo AS (SELECT
                                                about,
                                                email,
                                                fullname,
                                                nickname
                                              FROM users
                                              WHERE lower(nickname) = lower(new.user_nick)) INSERT INTO forum_users
VALUES ((SELECT id
         FROM forum
         WHERE lower(slug) = lower(new.forum_slug)), (SELECT nickname
                                          FROM userinfo), (SELECT about
                                                           FROM userinfo), (SELECT email
                                                                            FROM userinfo), (SELECT fullname
                                                                                             FROM userinfo))
ON CONFLICT DO NOTHING;

RETURN NULL;
END;' LANGUAGE plpgsql;

CREATE TRIGGER on_thread_insert_user
AFTER INSERT ON thread
FOR EACH ROW EXECUTE PROCEDURE forum_users_update();

-- CREATE TRIGGER on_post_insert_user
-- AFTER INSERT ON post
-- FOR EACH ROW EXECUTE PROCEDURE forum_users_update();


