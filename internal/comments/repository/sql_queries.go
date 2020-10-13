package repository

const (
	createComment = `INSERT INTO comments (author_id, news_id, message) VALUES ($1, $2, $3) RETURNING *`

	updateComment = `UPDATE comments SET message = $1, updated_at = CURRENT_TIMESTAMP WHERE comment_id = $2 RETURNING *`

	deleteComment = `DELETE FROM comments WHERE comment_id = $1`

	getCommentByID = `SELECT comment_id, author_id, news_id, message, likes, updated_at 
						FROM comments 
						WHERE comment_id = $1`

	getTotalCountByNewsId = `SELECT COUNT(comment_id) FROM comments WHERE news_id = $1`

	getCommentsByNewsId = `SELECT * FROM comments WHERE news_id = $1 ORDER BY updated_at OFFSET $2 LIMIT $3`
)
