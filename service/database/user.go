package database

import (
	"database/sql"
	"errors"
	"strconv"
)

var ErrUserAlreadyInTheDb error = errors.New("user already in the db")
var ErrNameAlreadyTaken error = errors.New("name already taken")
var ErrCannotFollowHimself error = errors.New("user id and following id are the same")
var ErrAlreadyFollowing error = errors.New("already following")
var ErrAlreadyBanned error = errors.New("already banned")
var ErrCannotBanHimself error = errors.New("user cannot ban himself")

func (db *appdbimpl) CreateUser(username string) error {
	flag, _ := db.UsernameExists(username)
	if flag {
		return ErrUserAlreadyInTheDb
	}
	_, err := db.c.Exec("INSERT INTO user (username) VALUES (?);", username)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) SetMyNickname(newName string, uid int) error {
	flag, _ := db.UsernameExists(newName)
	if flag {
		return ErrNameAlreadyTaken
	}
	_, err := db.c.Exec("UPDATE user SET username=? WHERE uid=?;", newName, uid)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) UsernameExists(username string) (bool, error) {
	// controllo se ne esiste almeno uno
	query := "SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)"

	// eseguo
	row := db.c.QueryRow(query, username)

	var exists bool

	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *appdbimpl) GetId(username string) (int, error) {
	query := "SELECT uid FROM user WHERE username = ?;"
	row := db.c.QueryRow(query, username)
	var uid string
	err := row.Scan(&uid)
	if err != nil {
		return 0, err
	}
	result, _ := strconv.Atoi(uid)
	return result, err
}

func (db *appdbimpl) FollowUser(uid int, followedUid int) error {
	if uid == followedUid {
		return ErrCannotFollowHimself
	}

	// Verifica se la tupla è già presente nella tabella follow
	exists, err := db.FollowExists(uid, followedUid)
	if err != nil {
		return err
	}

	// Se la tupla è già presente, restituisci un errore specifico
	if exists {
		return ErrAlreadyFollowing
	}

	// Inserisci la tupla nella tabella follow
	query := "INSERT INTO follow (uid, followedUid) VALUES (?,?);"
	_, err = db.c.Exec(query, uid, followedUid)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) FollowExists(uid int, followedUid int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM follow WHERE uid = ? AND followedUid = ?)"

	err := db.c.QueryRow(query, uid, followedUid).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (db *appdbimpl) BanExists(uid int, bannedUid int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM ban WHERE uid = ? AND bannedUid = ?)"

	err := db.c.QueryRow(query, uid, bannedUid).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (db *appdbimpl) GetFollowing(uid int) ([]int, error) {
	query := "SELECT followedUid FROM follow WHERE uid=?;"
	rows, err := db.c.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var array []int

	for rows.Next() {
		var followedUid string
		if err := rows.Scan(&followedUid); err != nil {
			return nil, err
		}
		value, _ := strconv.Atoi(followedUid)
		array = append(array, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return array, nil
}

func (db *appdbimpl) UnfollowUser(uid int, followedUid int) error {
	query := "DELETE FROM follow WHERE uid=? AND followedUid=?;"
	result, err := db.c.Exec(query, uid, followedUid)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Nessuna riga è stata eliminata, quindi l'entry non è stata trovata
		return sql.ErrNoRows
	}

	return nil
}

func (db *appdbimpl) BanUser(uid int, bannedUid int) error {
	if uid == bannedUid {
		return ErrCannotBanHimself
	}

	// Verifica se la tupla è già presente nella tabella follow
	exists, err := db.BanExists(uid, bannedUid)
	if err != nil {
		return err
	}

	// Se la tupla è già presente, restituisci un errore specifico
	if exists {
		return ErrAlreadyBanned
	}

	// Inserisci la tupla nella tabella follow
	query := "INSERT INTO ban (uid, bannedUid) VALUES (?,?);"
	_, err = db.c.Exec(query, uid, bannedUid)
	if err != nil {
		return err
	}

	return nil
}
