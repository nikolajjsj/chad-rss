-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password) VALUES (
  ?, ?
) RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET username = ?, password = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: GetFeeds :many
SELECT 
    f.nid,
    f.url,
    f.title,
    f.summary,
    f.authors,
    f.image
FROM 
    users u
JOIN 
    user_feed uf ON u.id = uf.user_id
JOIN 
    feeds f ON uf.feed_id = f.id
WHERE 
    u.id = ?
ORDER BY
    f.title
LIMIT ? OFFSET ?;

-- name: GetFeedByID :one
SELECT 
    f.id,
    f.nid,
    f.url,
    f.title,
    f.summary,
    f.authors,
    f.image
FROM
    users u 
JOIN
    user_feed uf ON u.id = uf.user_id
JOIN
    feeds f ON uf.feed_id = f.id
WHERE
    u.id = ? AND f.nid = ?
LIMIT 1;

-- name: CreateFeed :one
INSERT OR IGNORE INTO feeds (nid, url, title, summary, authors, image) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: AddFeedToUser :exec
INSERT OR IGNORE INTO user_feed (user_id, feed_id) VALUES (
    ?, ?
);

-- name: RemoveFeedFromUser :exec
DELETE FROM user_feed
WHERE user_id = ? AND feed_id = ?;

-- name: GetUserFeedArticles :many
SELECT 
    a.nid,
    a.url,
    a.title,
    a.summary,
    a.content,
    a.authors,
    a.media,
    a.published_at
FROM
    users u
JOIN
    user_feed uf ON u.id = uf.user_id
JOIN
    feeds f ON uf.feed_id = f.id
JOIN    
    articles a ON f.id = a.feed_id
WHERE
    u.id = ? AND f.nid = ?
ORDER BY
    a.published_at DESC
LIMIT ? OFFSET ?;

-- name: CreateFeedArticles :many
INSERT OR IGNORE INTO articles (rss_id, nid, url, title, summary, content, authors, media, published_at, feed_id) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetArticle :one
SELECT 
    a.nid,
    a.url,
    a.title,
    a.summary,
    a.content,
    a.authors,
    a.media,
    a.published_at
FROM
    users u
JOIN
    user_feed uf ON u.id = uf.user_id
JOIN
    feeds f ON uf.feed_id = f.id
JOIN
    articles a ON f.id = a.feed_id
WHERE
    u.id = ? AND a.nid = ?
LIMIT 1;
