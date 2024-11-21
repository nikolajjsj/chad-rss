CREATE TABLE user_feed (
  user_id INTEGER REFERENCES users(id),
  feed_id INTEGER REFERENCES feeds(id),
  PRIMARY KEY (user_id, feed_id)
);
