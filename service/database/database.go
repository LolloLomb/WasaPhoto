/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CreateUser(name string) error
	SetMyNickname(newName string, uid int) error
	GetId(username string) (int, error)
	GetUsername(uid int) (string, error)
	FollowUser(uid int, followedUid int) error
	GetFollowing(uid int) ([]UserTuple, error)
	GetFollowers(followedUid int) ([]UserTuple, error)
	UnfollowUser(uid int, followedUid int) error
	BanUser(uid int, bannedUid int) error
	BanExists(uid int, bannedUid int) (bool, error)
	IdExists(uid int) (bool, error)
	UnbanUser(uid int, bannedUid int) error
	CreatePhoto(uid int) (int, error)
	PhotoIdExists(photoId int) (bool, error)
	LikePhoto(uid int, photoId int) error
	UnlikePhoto(uid int, photoId int) error
	CommentPhoto(uid int, photoId int, text string) error
	UncommentPhoto(photoId int, commentId int) error
	PostsAmount(uid int) (int, error)
	DeletePhoto(photoId int) error
	GetStream(uid int) ([]int, error)
	FollowExists(uid int, followedUid int) (bool, error)
	IsPhotoOwner(uid int, photoId int) (bool, error)
	IsCommentOwner(uid int, commentId int) (bool, error)
	SearchUsers(uidCaller int, query string) ([]string, error)

	Ping() error
}

type UserTuple struct {
	Username string
	Token    int
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string

	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='user';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createTables(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec("PRAGMA foreign_key=ON;")
	if err != nil {
		return fmt.Errorf("error managind database : pragma foreign_key = on failed")
	}
	userQuery := `CREATE TABLE IF NOT EXISTS user (
		uid INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		UNIQUE(uid, username)
	);`
	_, err = db.Exec(userQuery)
	if err != nil {
		return fmt.Errorf("error creating users structure: %w", err)
	}
	photoQuery := `CREATE TABLE IF NOT EXISTS photo (photoId INTEGER PRIMARY KEY AUTOINCREMENT,
					upload_date DATETIME, 
					uid INTEGER, 
					FOREIGN KEY(uid) REFERENCES user(uid) ON DELETE CASCADE);`
	_, err = db.Exec(photoQuery)
	if err != nil {
		return fmt.Errorf("error creating photo structure: %w", err)
	}
	commentQuery := `CREATE TABLE IF NOT EXISTS comment (commentId INTEGER PRIMARY KEY AUTOINCREMENT,
					commentText TEXT,
					uploadDate DATETIME,
					uid INTEGER, 
					photoId INTEGER,
					FOREIGN KEY(uid) REFERENCES user(uid) ON DELETE CASCADE,
					FOREIGN KEY(photoId) REFERENCES photo(photoId) ON DELETE CASCADE);`
	_, err = db.Exec(commentQuery)
	if err != nil {
		return fmt.Errorf("error creating comment structure: %w", err)
	}
	followQuery := `CREATE TABLE IF NOT EXISTS follow (
		uid INTEGER,
		followedUid INTEGER, 
		PRIMARY KEY (uid, followedUid),
		FOREIGN KEY (uid) REFERENCES user(uid) ON DELETE CASCADE,
		FOREIGN KEY (followedUid) REFERENCES user(uid) ON DELETE CASCADE
	);`

	_, err = db.Exec(followQuery)
	if err != nil {
		return fmt.Errorf("error creating follow structure: %w", err)
	}
	likeQuery := `CREATE TABLE IF NOT EXISTS like (uid INTEGER,
					likedPhotoId INTEGER, 
					PRIMARY KEY (uid, likedPhotoId),
					FOREIGN KEY (uid) REFERENCES user(uid),
					FOREIGN KEY (likedPhotoId) REFERENCES photo(photoId));`
	_, err = db.Exec(likeQuery)
	if err != nil {
		return fmt.Errorf("error creating like structure: %w", err)
	}
	banQuery := `CREATE TABLE IF NOT EXISTS ban (
		uid INTEGER,
		bannedUid INTEGER,
		PRIMARY KEY (uid, bannedUid),
		FOREIGN KEY (uid) REFERENCES user(uid),
		FOREIGN KEY (bannedUid) REFERENCES user(uid));`
	_, err = db.Exec((banQuery))
	if err != nil {
		return fmt.Errorf("error creating ban structure: %w", err)
	}
	return nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
