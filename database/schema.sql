-- SET SYNCHRONOUS_COMMIT = 'off';
CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP INDEX IF EXISTS usersEmailIndex;
DROP INDEX IF EXISTS usersNicknameIndex;
DROP INDEX IF EXISTS forumSlugIndex;
DROP INDEX IF EXISTS voteNicknameAndThreadIdIndex;
DROP INDEX IF EXISTS threadSlugIndex;

DROP TRIGGER IF EXISTS on_vote_insert
ON vote;
DROP TRIGGER IF EXISTS on_vote_update
ON vote;
DROP TRIGGER IF EXISTS on_thread_insert
ON thread;
DROP TRIGGER IF EXISTS on_thread_insert_user
ON thread;
DROP TRIGGER IF EXISTS on_post_insert_user
ON post;

DROP FUNCTION IF EXISTS vote_insert() CASCADE;
DROP FUNCTION IF EXISTS vote_update() CASCADE;
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
  email    VARCHAR(50) NOT NULL UNIQUE,
  fullname VARCHAR(50) NOT NULL,
  nickname VARCHAR(50) NOT NULL UNIQUE
);

CREATE UNIQUE INDEX usersNicknameIndex
  ON users (lower(nickname));

CREATE UNIQUE INDEX usersEmailIndex
  ON users (lower(email));

--
-- FORUM
--
CREATE TABLE forum (
  id        SERIAL PRIMARY KEY,
  slug      VARCHAR(50)  NOT NULL UNIQUE,
  title     VARCHAR(100) NOT NULL,
  posts     BIGINT       NOT NULL DEFAULT 0,
  threads   INTEGER      NOT NULL DEFAULT 0,
  moderator VARCHAR(50)  NOT NULL
);

CREATE UNIQUE INDEX forumSlugIndex
  ON forum (lower(slug));

--
-- THREAD
--
CREATE TABLE thread (
  id          SERIAL PRIMARY KEY,

  slug        VARCHAR(50) UNIQUE DEFAULT NULL,
  title       VARCHAR(100) NOT NULL,
  message     TEXT         NOT NULL,
  forum_slug  VARCHAR(100) NOT NULL,
  user_nick   VARCHAR(50)  NOT NULL,
  created     TIMESTAMP WITH TIME ZONE,
  votes_count INTEGER            DEFAULT 0
);

CREATE UNIQUE INDEX threadSlugIndex
  ON thread (lower(slug));

--
-- POST
--
CREATE TABLE post (
  id         SERIAL PRIMARY KEY,

  user_nick  VARCHAR(50)               NOT NULL,
  message    TEXT                      NOT NULL,
  created    TIMESTAMP WITH TIME ZONE,
  forum_slug VARCHAR(100)              NOT NULL,
  thread_id  INTEGER REFERENCES thread NOT NULL,
  is_edited  BOOLEAN                   NOT NULL DEFAULT FALSE,
  parent     INTEGER                            DEFAULT 0,
  parents    INTEGER []                NOT NULL
);

--
-- VOTE
--
CREATE TABLE vote (
  id        SERIAL PRIMARY KEY,

  nickname  VARCHAR(50) NOT NULL,
  thread_id INTEGER REFERENCES thread,
  voice     INTEGER
);
CREATE UNIQUE INDEX voteNicknameAndThreadIdIndex
  ON vote (nickname, thread_id);

CREATE FUNCTION vote_insert()
  RETURNS TRIGGER AS '
BEGIN
  UPDATE thread
  SET
    votes_count = thread.votes_count + NEW.voice
  WHERE id = NEW.thread_id;
  RETURN NULL;
END;
' LANGUAGE plpgsql;

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

CREATE TRIGGER on_vote_insert
AFTER INSERT
  ON vote
FOR EACH ROW EXECUTE PROCEDURE vote_insert();

CREATE FUNCTION vote_update()
  RETURNS TRIGGER AS '
BEGIN

  IF OLD.voice = NEW.voice
  THEN
    RETURN NULL;
  END IF;

  UPDATE thread
  SET
    votes_count = thread.votes_count +
                  CASE WHEN NEW.voice = -1
                    THEN -2
                  ELSE 2 END
  WHERE id = NEW.thread_id;
  RETURN NULL;
END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_vote_update
AFTER UPDATE
  ON vote
FOR EACH ROW EXECUTE PROCEDURE vote_update();

-- CREATE FUNCTION on_vote()
--   RETURNS TRIGGER AS 'BEGIN IF(TG_OP = 'UPDATE') THEN IF  END;' LANGUAGE plpgsql;

CREATE TABLE forum_users (
  forumId INTEGER REFERENCES forum,
  userId  INTEGER REFERENCES users
);

CREATE FUNCTION forum_users_update()
  RETURNS TRIGGER AS 'BEGIN INSERT INTO forum_users (forumId, userId) VALUES ((SELECT id
                                                                               FROM forum
                                                                               WHERE
                                                                                 lower(slug) = lower(NEW.forum_slug)),
                                                                              (SELECT id
                                                                               FROM users
                                                                               WHERE lower(nickname) =
                                                                                     lower(NEW.user_nick)));
  RETURN NULL;
END;' LANGUAGE plpgsql;

CREATE TRIGGER on_thread_insert_user
AFTER INSERT ON thread
FOR EACH ROW EXECUTE PROCEDURE forum_users_update();

CREATE TRIGGER on_post_insert_user
AFTER INSERT ON post
FOR EACH ROW EXECUTE PROCEDURE forum_users_update();
