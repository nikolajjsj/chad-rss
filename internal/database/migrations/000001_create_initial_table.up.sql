CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username text NOT NULL UNIQUE,
  password text NOT NULL,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE feeds (
  id SERIAL PRIMARY KEY,
  nid text NOT NULL UNIQUE,
  url text NOT NULL UNIQUE,
  title text NOT NULL,
  summary text,
  authors text,
  image text,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE articles (
  id SERIAL PRIMARY KEY,
  nid text NOT NULL UNIQUE,
  rss_id text NOT NULL UNIQUE,
  url text NOT NULL,
  title text NOT NULL,
  summary text,
  content text,
  authors text,
  media text,
  published_at TIMESTAMP,

  feed_id INTEGER REFERENCES feeds(id),

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
