CREATE EXTENSION citext;

DROP TABLE forum_users;
DROP TABLE votes;
DROP TABLE posts;
DROP TABLE threads;
DROP TABLE forums;
DROP TABLE users;

CREATE TABLE users (
    nickname CITEXT PRIMARY KEY,
    fullname CITEXT NOT NULL,
    about TEXT NOT NULL,
    email CITEXT UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS users_index_nickname ON users (nickname);

CREATE TABLE forums (
    title TEXT NOT NULL,
    nickname CITEXT NOT NULL,
    posts INT DEFAULT 0,
    threads INT DEFAULT 0,
    slug CITEXT PRIMARY KEY,
    FOREIGN KEY (nickname) REFERENCES users (nickname) 
);

/*1*/
CREATE INDEX IF NOT EXISTS forums_index_slug ON forums (slug);
/*CREATE INDEX IF NOT EXISTS forums_index_nickname ON forums (nickname);*/

CREATE TABLE threads (
    author CITEXT NOT NULL,
    create_date timestamptz DEFAULT now(),
    forum CITEXT NOT NULL,
    id SERIAL PRIMARY KEY,
    msg TEXT NOT NULL,
    slug CITEXT UNIQUE,
    title CITEXT NOT NULL,
    votes INT DEFAULT 0,
    FOREIGN KEY (forum) REFERENCES forums (slug),
    FOREIGN KEY (author) REFERENCES users (nickname) 
);

CREATE INDEX IF NOT EXISTS threads_index_id ON threads (id);
CREATE INDEX IF NOT EXISTS threads_index_slug ON threads (slug);
CREATE INDEX IF NOT EXISTS threads_index_forum ON threads (forum);

/*CREATE INDEX threads_index_author ON threads (author);*/
CREATE INDEX IF NOT EXISTS threads_index_forum_create_date ON threads (forum, create_date);
CREATE INDEX IF NOT EXISTS threads_index_create_date ON threads (create_date);

CREATE TABLE posts (
    author CITEXT NOT NULL,
    create_date timestamptz DEFAULT now(),
    forum CITEXT NOT NULL,
    id SERIAL PRIMARY KEY,
    is_edited BOOLEAN DEFAULT FALSE,
    msg TEXT NOT NULL,
    parent INT NOT NULL,
    thread INT NOT NULL,
    path BIGINT ARRAY,
    FOREIGN KEY (forum) REFERENCES forums (slug),
    FOREIGN KEY (author) REFERENCES users (nickname),
    FOREIGN KEY (thread) REFERENCES threads (id)
);

CREATE INDEX IF NOT EXISTS posts_index_thread_id_create_date_thread ON posts (id, create_date, thread);
CREATE INDEX IF NOT EXISTS posts_index_thread_id ON posts (thread, id);
CREATE INDEX IF NOT EXISTS posts_index_forum ON posts (forum);
CREATE INDEX IF NOT EXISTS posts_index_id ON posts (id);
CREATE INDEX IF NOT EXISTS posts_index_thread_parent_id ON posts (thread, (path[1]), id);
CREATE INDEX IF NOT EXISTS posts_index_thread_path ON posts (thread, path);
CREATE INDEX IF NOT EXISTS posts_index_thread_id_main ON posts (thread, id) WHERE parent = 0;

CREATE INDEX posts_index_thread_array_length ON posts (thread, (array_length(path, 1)));
CREATE INDEX posts_index_thread_array_length_path ON posts (thread, (array_length(path, 1)), (path[1]));
/*CREATE INDEX IF NOT EXISTS posts_index_parent ON posts ((path[1]));
CREATE INDEX IF NOT EXISTS posts_index_thread_created ON posts (thread, create_date);*/

CREATE TABLE votes (
  thread INT NOT NULL,
  voice INT NOT NULL,
  nickname CITEXT NOT NULL,
  FOREIGN KEY (thread) REFERENCES threads (id),
  FOREIGN KEY (nickname) REFERENCES users (nickname),
  UNIQUE (thread, nickname)
);

CREATE UNIQUE INDEX IF NOT EXISTS votex_index_thread_nickname ON votes (thread, nickname);

/*CREATE INDEX IF NOT EXISTS vote_index_thread ON votes (thread);
CREATE INDEX IF NOT EXISTS vote_index_nickname ON votes (nickname);*/

CREATE TABLE forum_users (
    forum CITEXT NOT NULL,
    nickname CITEXT NOT NULL,
    FOREIGN KEY (forum) REFERENCES forums (slug),
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    UNIQUE (forum, nickname)
);


CREATE INDEX IF NOT EXISTS forum_users_index_forum ON forum_users (forum);
CREATE INDEX IF NOT EXISTS forums_users_index_nickname ON forum_users (nickname);
CREATE INDEX IF NOT EXISTS forums_users_index_forum_nickname ON forum_users (forum, nickname);


