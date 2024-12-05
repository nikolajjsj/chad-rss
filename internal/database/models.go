package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
}

type Feed struct {
	gorm.Model
	NID      string `gorm:"unique"`
	URL      string `gorm:"unique"`
	Title    string
	Summary  string
	Authors  string
	Image    string
	Articles []Article
	UserID   uint
	User     User
}

type Article struct {
	gorm.Model
	NID         string `gorm:"unique"`
	RSSID       string `gorm:"unique"`
	URL         string `gorm:"unique"`
	Title       string
	Summary     string
	Content     string
	Authors     string
	Media       string
	PublishedAt time.Time
	FeedID      uint
	Feed        Feed
}

// CREATE TABLE users (
//   id SERIAL PRIMARY KEY,
//   username text NOT NULL UNIQUE,
//   password text NOT NULL,
//
//   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );
//
// CREATE TABLE feeds (
//   id SERIAL PRIMARY KEY,
//   nid text NOT NULL UNIQUE,
//   url text NOT NULL UNIQUE,
//   title text NOT NULL,
//   summary text,
//   authors text,
//   image text,
//
//   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );
//
// CREATE TABLE articles (
//   id SERIAL PRIMARY KEY,
//   nid text NOT NULL UNIQUE,
//   rss_id text NOT NULL UNIQUE,
//   url text NOT NULL,
//   title text NOT NULL,
//   summary text,
//   content text,
//   authors text,
//   media text,
//   published_at TIMESTAMP,
//
//   feed_id INTEGER REFERENCES feeds(id),
//
//   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );
//
// CREATE TABLE user_feed (
//   user_id INTEGER REFERENCES users(id),
//   feed_id INTEGER REFERENCES feeds(id),
//   PRIMARY KEY (user_id, feed_id)
// );
