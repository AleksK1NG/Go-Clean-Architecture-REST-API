package repository

const (
	createNews = `INSERT INTO news (author_id, title, content, image_url, category, created_at) 
					VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($4, ''), now()) 
					RETURNING *`

	updateNews = `UPDATE news 
					SET title = COALESCE(NULLIF($1, ''), title),
						content = COALESCE(NULLIF($2, ''), content), 
					    image_url = COALESCE(NULLIF($3, ''), image_url), 
					    category = COALESCE(NULLIF($4, ''), category), 
					    updated_at = now() 
					RETURNING *`

	getNewsByID = `SELECT u.first_name,
       u.last_name,
       u.avatar,
       u.login_date,
       u.role,
       u.user_id,
       u.updated_at as user_updated_at,
       n.news_id,
       n.title,
       n.image_url,
       n.content,
       n.category,
       n.updated_at as news_updated_at
FROM news n
         LEFT JOIN users u on u.user_id = n.author_id
WHERE news_id = $1`

	deleteNews = `DELETE FROM news WHERE news_id = $1`

	getTotalCount = `SELECT COUNT(news_id) FROM news`

	getNews = `SELECT news_id, author_id, title, content, image_url, category, updated_at, created_at 
				FROM news 
				ORDER BY created_at, updated_at OFFSET $1 LIMIT $2`
)
