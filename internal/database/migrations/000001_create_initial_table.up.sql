CREATE TABLE users (
  id   INTEGER PRIMARY KEY,
  username text NOT NULL,
  password text NOT NULL

  -- created_at timestamp default now(),
  -- updated_at timestamp default now()
);

CREATE TABLE feed (
  id   INTEGER PRIMARY KEY,
  feed_id text NOT NULL,
  url text NOT NULL,
  title text NOT NULL,
  summary text,
  authors text,
  image text
  
  -- created_at timestamp default now(),
  -- updated_at timestamp default now()
);

CREATE TABLE article (
  id   INTEGER PRIMARY KEY,
  article_id text NOT NULL,
  url text NOT NULL,
  title text NOT NULL,
  summary text,
  content text,
  authors text,
  media text,

  feed_id text NOT NULL

  -- created_at timestamp default now(),
  -- updated_at timestamp default now()
);

