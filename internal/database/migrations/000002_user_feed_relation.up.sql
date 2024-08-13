CREATE TABLE user_feed (
  user_id INTEGER NOT NULL,
  feed_id INTEGER NOT NULL,
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(feed_id) REFERENCES feeds(id),
  PRIMARY KEY (user_id, feed_id)
);

