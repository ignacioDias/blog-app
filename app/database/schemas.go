package database

const createSchema = `
	CREATE TABLE IF NOT EXISTS users
	(
		username TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		password TEXT
	);
	CREATE TABLE IF NOT EXISTS posts
	(
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		author TEXT REFERENCES users(username) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS user_follows 
	(
		follower_username TEXT NOT NULL,
		followed_username TEXT NOT NULL,
		PRIMARY KEY (follower_username, followed_username),
		FOREIGN KEY (follower_username) REFERENCES users(username) ON DELETE CASCADE,
		FOREIGN KEY (followed_username) REFERENCES users(username) ON DELETE CASCADE,
		CHECK (follower_username <> followed_username)
	);
	`

var insertPostSchema = `INSERT INTO posts(title, content, author) VALUES($1, $2, $3) RETURNING id`

var insertUserSchema = `INSERT INTO users(username, email, password) VALUES($1, $2, $3)`

var insertFollowSchema = `INSERT INTO user_follows (follower_username, followed_username) VALUES ($1, $2)`

var removeFollowSchema = `DELETE FROM user_follows WHERE follower_username = $1 AND followed_username = $2`

var getFollowersSchema = `SELECT follower_username FROM user_follows WHERE followed_username = $1`

var getFollowingSchema = `SELECT followed_username FROM user_follows WHERE follower_username = $1`

var checkFollowSchema = `SELECT EXISTS(SELECT 1 FROM user_follows WHERE follower_username = $1 AND followed_username = $2)`
