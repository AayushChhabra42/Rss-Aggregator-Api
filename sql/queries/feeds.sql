-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url,user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
Select * from feeds;

-- name: GetNextFeedstoFetch :many
Select * from feeds
order by last_fetched_at asc nulls first
limit $1;

-- name: MarkFeedsFetched :one
Update feeds
set last_fetched_at=Now(),
updated_at=Now()
WHERE id=$1
RETURNING *;
