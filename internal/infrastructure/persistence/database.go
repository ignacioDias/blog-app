package persistence

import (
	"fmt"
	"log"
	"postapi/internal/domain"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	db                   *sqlx.DB
	UserRepository       domain.UserRepository
	PostRepository       domain.PostRepository
	ProfileRepository    domain.ProfileRepository
	UserFollowRepository domain.UserFollowRepository
}

func (d *DB) Open() error {
	pg, err := sqlx.Open("postgres", pgConnStr)
	if err != nil {
		return err
	}
	log.Println("Connected to Database!")
	pg.MustExec(createSchema)
	d.db = pg

	d.UserRepository = &UserRepositoryImpl{db: d.db}
	d.PostRepository = &PostRepositoryImpl{db: d.db}
	d.ProfileRepository = &ProfileRepositoryImpl{db: d.db}
	d.UserFollowRepository = &UserFollowRepositoryImpl{db: d.db}

	return nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

var (
	dbUsername = "postgres"
	dbPassword = "postgres"
	dbHost     = "localhost"
	dbTable    = "postgres"
	dbPort     = "5432"
	pgConnStr  = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUsername, dbTable, dbPassword)
)

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
	CREATE TABLE IF NOT EXISTS profiles
	(
		username TEXT PRIMARY KEY,
		description TEXT,
		profile_picture TEXT,
		FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE
	);
	`

var insertPostSchema = `INSERT INTO posts(title, content, author) VALUES($1, $2, $3) RETURNING id`

var insertUserSchema = `INSERT INTO users(username, email, password) VALUES($1, $2, $3)`

var insertFollowSchema = `INSERT INTO user_follows (follower_username, followed_username) VALUES ($1, $2)`

var removeFollowSchema = `DELETE FROM user_follows WHERE follower_username = $1 AND followed_username = $2`

var getFollowersSchema = `SELECT follower_username FROM user_follows WHERE followed_username = $1`

var getFollowingSchema = `SELECT followed_username FROM user_follows WHERE follower_username = $1`

var insertProfileSchema = `INSERT INTO profiles(username, description, profile_picture) VALUES($1, $2, $3)`

var getProfileSchema = `SELECT * FROM profiles WHERE username = $1`

var updateProfileSchema = `UPDATE profiles SET description = $2, profile_picture = $3 WHERE username = $1`
