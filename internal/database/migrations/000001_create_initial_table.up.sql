CREATE TABLE users (
  id   INTEGER PRIMARY KEY AUTOINCREMENT,
  username text NOT NULL UNIQUE,
  password text NOT NULL,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE feeds (
  id   INTEGER PRIMARY KEY AUTOINCREMENT,
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
  id   INTEGER PRIMARY KEY AUTOINCREMENT,
  nid text NOT NULL UNIQUE,
  rss_id text NOT NULL UNIQUE,
  url text NOT NULL,
  title text NOT NULL,
  summary text,
  content text,
  authors text,
  media text,
  published_at TIMESTAMP,

  feed_id INTEGER NOT NULL,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY(feed_id) REFERENCES feeds(id)
);

