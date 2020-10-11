package repository

const (
	createComment = `INSERT INTO comments (author_id, news_id, message) VALUES ($1, $2, $3) RETURNING *`
)
