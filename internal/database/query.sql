-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password) VALUES (
  $1, $2
) RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET username = $1, password = $2
WHERE id = $3
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: FeedsCount :one
SELECT COUNT(*) FROM feeds;

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
    u.id = $1
ORDER BY
    f.title
LIMIT $2 OFFSET $3;

-- name: GetAllFeeds :many
SELECT 
    id,
    nid,
    url
FROM
    feeds;

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
    u.id = $1 AND f.nid = $2
LIMIT 1;

-- name: CreateFeed :one
INSERT INTO feeds (nid, url, title, summary, authors, image) 
VALUES ($1, $2, $3, $4, $5, $6) 
  ON CONFLICT DO NOTHING 
RETURNING *;

-- name: AddFeedToUser :exec
INSERT INTO user_feed (user_id, feed_id) 
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveFeedFromUser :exec
DELETE FROM user_feed
WHERE user_id = $1 AND feed_id = $2;

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
    u.id = $1 AND f.nid = $2
ORDER BY
    a.published_at DESC
LIMIT $3 OFFSET $4;

-- name: CreateFeedArticles :many
INSERT INTO articles (rss_id, nid, url, title, summary, content, authors, media, published_at, feed_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
  ON CONFLICT DO NOTHING
RETURNING *;

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
    u.id = $1 AND a.nid = $2
LIMIT 1;
