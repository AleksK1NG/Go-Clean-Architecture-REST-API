package repository

const (
	createUserQuery = `INSERT INTO users (first_name, last_name, email, password, role, about, avatar, phone_number, address,
	               		city, gender, postcode, birthday, created_at, updated_at, login_date)
						VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'user'), $6, $7, $8, $9, $10, $11, $12, $13, now(), now(), now()) 
						RETURNING *`

	updateUserQuery = `UPDATE users 
						SET first_name = COALESCE(NULLIF($1, ''), first_name),
						    last_name = COALESCE(NULLIF($2, ''), last_name),
						    email = COALESCE(NULLIF($3, ''), email),
						    role = COALESCE(NULLIF($4, ''), role),
						    about = COALESCE(NULLIF($5, ''), about),
						    avatar = COALESCE(NULLIF($6, ''), avatar),
						    phone_number = COALESCE(NULLIF($7, ''), phone_number),
						    address = COALESCE(NULLIF($8, ''), address),
						    city = COALESCE(NULLIF($9, ''), city),
						    gender = COALESCE(NULLIF($10, ''), gender),
						    postcode = COALESCE(NULLIF($11, 0), postcode),
						    birthday = COALESCE(NULLIF($12, '')::date, birthday),
						    updated_at = now()
						WHERE user_id = $13
						RETURNING *
						`

	deleteUserQuery = `DELETE FROM users WHERE user_id = $1`

	getUserQuery = `SELECT user_id, first_name, last_name, email, role, about, avatar, phone_number, 
       				 address, city, gender, postcode, birthday, created_at, updated_at, login_date  
					 FROM users 
					 WHERE user_id = $1`

	getTotalCount = `SELECT COUNT(user_id) FROM users 
						WHERE first_name ILIKE '%' || $1 || '%' or last_name ILIKE '%' || $1 || '%'`

	findUsers = `SELECT user_id, first_name, last_name, email, role, about, avatar, phone_number, address,
	              city, gender, postcode, birthday, created_at, updated_at, login_date 
				  FROM users 
				  WHERE first_name ILIKE '%' || $1 || '%' or last_name ILIKE '%' || $1 || '%'
				  ORDER BY first_name, last_name
				  OFFSET $2 LIMIT $3
				  `

	getTotal = `SELECT COUNT(user_id) FROM users`

	getUsers = `SELECT user_id, first_name, last_name, email, role, about, avatar, phone_number, 
       			 address, city, gender, postcode, birthday, created_at, updated_at, login_date
				 FROM users 
				 ORDER BY COALESCE(NULLIF($1, ''), first_name) OFFSET $2 LIMIT $3`

	findUserByEmail = `SELECT user_id, first_name, last_name, email, role, about, avatar, phone_number, 
       			 		address, city, gender, postcode, birthday, created_at, updated_at, login_date, password
				 		FROM users 
				 		WHERE email = $1`
)
