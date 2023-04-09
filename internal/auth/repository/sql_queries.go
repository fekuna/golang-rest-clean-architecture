package repository

const (
	createUserQuery = `
		INSERT INTO users (first_name, last_name, email, password, role, about, avatar, phone_number, address, city, gender, postcode, birthday, created_at, updated_at, login_date) VALUES($1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'user'), $6, $7, $8, $9, $10, $11, $12, $13, now(), now(), now()) RETURNING *
	`

	findUserByEmail = `
		SELECT user_id, first_name, last_name, email, role, about, avatar, phone_number, address, city, gender, postcode, birthday, created_at, updated_at, login_date, password FROM users WHERE email = $1
	`
)
