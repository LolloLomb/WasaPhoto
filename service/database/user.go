package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var ErrUserAlreadyInTheDb error = errors.New("user already in the db")
var ErrNameAlreadyTaken error = errors.New("name already taken")
var ErrCannotFollowHimself error = errors.New("user id and following id are the same")
var ErrAlreadyFollowing error = errors.New("already following")
var ErrAlreadyBanned error = errors.New("already banned")
var ErrCannotBanHimself error = errors.New("user cannot ban himself")
var ErrAlreadyLiked error = errors.New("already liked")
var ErrPhotoDontExists error = errors.New("photo do not exists")
var ErrCommentDontExists error = errors.New("error do not exists")

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

func (db *appdbimpl) CreatePhoto(uid int) (int, error) {
	/* lastPhotoId, err := db.getLastPhotoId()
	if err != nil{
		return err
	} */
	formattedTime := time.Now().Format("2006-01-02T15:04:05Z")
	query := "INSERT INTO photo (upload_date, uid) VALUES (?,?);"
	_, err := db.c.Exec(query, formattedTime, uid)
	if err != nil {
		return -1, err
	}
	photoId, err := db.getLastPhotoId(uid)
	if err != nil {
		return photoId, err
	}
	return photoId, nil
}

func (db *appdbimpl) getLastPhotoId(uid int) (int, error) {
	query := "SELECT photoId FROM photo WHERE userId = ? ORDER BY photoId DESC LIMIT 1;"

	// Eseguire la query
	rows, err := db.c.Query(query, uid)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	// Scorrere i risultati della query
	var lastPhotoID int = -1
	for rows.Next() {
		if err := rows.Scan(&lastPhotoID); err != nil {
			return 0, err
		}
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	// Se non ci sono errori
	if lastPhotoID != -1 {
		return lastPhotoID, nil
	}

	// Se non ci sono risultati (tabella vuota), restituire un valore predefinito o un errore
	return lastPhotoID, errors.New("nessun risultato trovato per l'utente specificato")
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
	var uid int
	err := row.Scan(&uid)
	if err != nil {
		return 0, err
	}
	return uid, err
}

func (db *appdbimpl) GetUsername(uid int) (string, error) {
	query := "SELECT username FROM user WHERE uid = ?;"
	row := db.c.QueryRow(query, uid)
	var username string
	err := row.Scan(&username)
	if err != nil {
		return "", err
	}
	return username, err
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

func (db *appdbimpl) IdExists(uid int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM user WHERE uid = ?)"
	err := db.c.QueryRow(query, uid).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) PhotoIdExists(photoId int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM photo WHERE photoId = ?)"
	err := db.c.QueryRow(query, photoId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) CommentIdExists(commentId int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM comment WHERE commentId = ?)"
	err := db.c.QueryRow(query, commentId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) BanExists(uid int, bannedUid int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM ban WHERE uid = ? AND bannedUid = ?)"

	err := db.c.QueryRow(query, uid, bannedUid).Scan(&exists)
	if err != nil {
		return false, err
	}
	// fmt.Print(exists, uid, bannedUid)
	return exists, nil
}

func (db *appdbimpl) LikeExists(uid int, photoId int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM like WHERE uid = ? AND likedPhotoId = ?)"

	err := db.c.QueryRow(query, uid, photoId).Scan(&exists)
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

func (db *appdbimpl) GetFollowers(followedUid int) ([]int, error) {
	query := "SELECT uid FROM follow WHERE followedUid=?;"
	rows, err := db.c.Query(query, followedUid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var array []int

	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		value, _ := strconv.Atoi(uid)
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

	// Verifica se la tupla è già presente nella tabella ban
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

func (db *appdbimpl) LikePhoto(uid int, photoId int) error {
	// Verifica se la tupla è già presente nella tabella like
	exists, err := db.LikeExists(uid, photoId)
	if err != nil {
		return err
	}

	// Se la tupla è già presente, restituisci un errore specifico
	if exists {
		return ErrAlreadyLiked
	}

	// Inserisci la tupla nella tabella like
	query := "INSERT INTO like (uid, likedPhotoId) VALUES (?,?);"
	_, err = db.c.Exec(query, uid, photoId)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UnlikePhoto(uid int, photoId int) error {
	query := "DELETE FROM like WHERE uid=? AND likedPhotoId=?;"
	result, err := db.c.Exec(query, uid, photoId)

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

func (db *appdbimpl) UnbanUser(uid int, bannedUid int) error {
	query := "DELETE FROM ban WHERE uid=? AND bannedUid=?;"
	result, err := db.c.Exec(query, uid, bannedUid)

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

func (db *appdbimpl) CommentPhoto(uid int, photoId int, text string) error {
	// Inserisci la tupla nella tabella comment
	formattedTime := time.Now().Format("2006-01-02T15:04:05Z")
	query := "INSERT INTO comment (uid, commentText, photoId, uploadDate) VALUES (?,?,?);"
	_, err := db.c.Exec(query, uid, text, photoId, formattedTime)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UncommentPhoto(photoId int, commentId int) error {
	// gli errori possibili sono 1 non esiste la foto, 2 non esiste il commento
	valid, _ := db.PhotoIdExists(photoId)
	if !valid {
		return ErrPhotoDontExists
	}
	valid, _ = db.CommentIdExists(commentId)
	if !valid {
		return ErrCommentDontExists
	}

	query := "DELETE FROM comment WHERE commentId = ?;"
	result, err := db.c.Exec(query, commentId)

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

func (db *appdbimpl) DeletePhoto(photoId int) error {
	// gli errori possibili sono 1 non esiste la foto
	valid, _ := db.PhotoIdExists(photoId)
	if !valid {
		return ErrPhotoDontExists
	}

	query := "DELETE FROM photo WHERE photoId = ?;"
	result, err := db.c.Exec(query, photoId)

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

func (db *appdbimpl) PostsAmount(uid int) (int, error) {
	// restituisce il numero di post di un utente dato un uid
	valid, err := db.IdExists(uid)
	if !valid {
		return -1, err
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM photo WHERE uid = '%s'", strconv.Itoa(uid))

	// Esegui la query
	var conteggio int
	err = db.c.QueryRow(query).Scan(&conteggio)
	if err != nil {
		return -1, err
	}
	return conteggio, err
}

func (db *appdbimpl) GetStream(uid int) ([]int, error) {
	orderedPhotos, err := db.getOrderedPhotos(uid)
	if err != nil {
		return nil, err
	}
	return orderedPhotos, nil
}

func (db *appdbimpl) IsPhotoOwner(uid int, photoId int) (bool, error) {
	var realOwner string
	query := "SELECT uid FROM photo WHERE photoId = ?;"
	err := db.c.QueryRow(query, photoId).Scan(&realOwner)
	if err != nil {
		return false, err
	}
	realOwnerId, _ := db.GetId(realOwner)
	if uid == realOwnerId {
		return true, err
	}
	return false, nil
}

func (db *appdbimpl) IsCommentOwner(uid int, photoId int) (bool, error) {
	var realOwner string
	query := "SELECT uid FROM comment WHERE commentId = ?;"
	err := db.c.QueryRow(query, photoId).Scan(&realOwner)
	if err != nil {
		return false, err
	}
	realOwnerId, _ := db.GetId(realOwner)
	if uid == realOwnerId {
		return true, err
	}
	return false, nil
}

func (db *appdbimpl) getOrderedPhotos(uid int) ([]int, error) {
	query := `SELECT * FROM photo WHERE uid IN (SELECT followedUid FROM follow WHERE uid = ?) ORDER BY upload_date DESC LIMIT 10`

	rows, err := db.c.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderedPhotos []int

	for rows.Next() {
		var photo int
		err := rows.Scan(photo)
		if err != nil {
			return nil, err
		}
		orderedPhotos = append(orderedPhotos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderedPhotos, nil
}
