package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	// "log"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

type User struct {
	Content string `json:"username"`
}

/*
type UserToken struct {
	Username string `json:"username"`
	Token    int    `json:"token"`
}
*/

type Profile struct {
	Username  string               `json:"username"`
	Id        int                  `json:"id"`
	Following []database.UserTuple `json:"following"`
	Followers []database.UserTuple `json:"followers"`
	Posts     []database.Photo     `json:"posts"`
}

type NewUsername struct {
	Content string `json:"newUsername"`
}

type Identifier struct {
	Content string `json:"identifier"`
}

type Response struct {
	SuccessMessage string `json:"success,omitempty"`
	ErrorMessage   string `json:"error,omitempty"`
}

var ErrForbidden error = errors.New("unauthorized")

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Decodifica il corpo JSON della richiesta in una struttura
	var requestBody NewUsername
	auth := r.Header.Get("Authorization")
	authConv, _ := strconv.Atoi(auth)
	uid, _ := strconv.Atoi(ps.ByName("uid"))

	if uid != authConv || auth == "" {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Gestisci errori di decodifica JSON
		response := Response{ErrorMessage: "Errore durante la decodifica del corpo JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// Estrai il nuovo username dalla struttura
	newUsername := requestBody.Content

	if newUsername == "" {
		response := Response{ErrorMessage: "Error occured while decoding JSON newUsername: value not found"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	err := rt.db.SetMyNickname(newUsername, uid)
	if errors.Is(err, database.ErrNameAlreadyTaken) {
		response := Response{ErrorMessage: "Nome utente già in uso"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	if err != nil {
		// Gestisci l'errore e invia una risposta JSON con il messaggio di errore
		response := Response{ErrorMessage: fmt.Sprintf("Errore durante l'impostazione del nuovo username: %v", err)}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	// Creare una risposta JSON di successo con il messaggio di log
	response := Response{SuccessMessage: fmt.Sprintf("Nuovo username impostato con successo. Risultato: %v", newUsername)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusOK)
}

func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Decodifica il corpo JSON della richiesta in una struttura
	var requestBody User
	bytes, _ := io.ReadAll(r.Body)
	// println(string(bytes))
	if err := json.Unmarshal(bytes, &requestBody); err != nil {
		// if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Gestisci errori di decodifica JSON
		response := Response{ErrorMessage: "Error occured while decoding JSON1"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	username := requestBody.Content
	if username == "" {
		response := Response{ErrorMessage: "Error occured while decoding JSON2"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	// println(username)
	err := rt.db.CreateUser(username)
	id, _ := rt.db.GetId(username)
	// println(username, id)
	// log.Printf("%v", err)
	if errors.Is(err, database.ErrUserAlreadyInTheDb) {
		// Gestisci l'errore e invia una risposta JSON con il messaggio di errore
		response := Response{SuccessMessage: fmt.Sprintf("ID: %v", id)}
		sendJSONResponse(w, response, 200)
		return
	}
	if err != nil {
		response := Response{ErrorMessage: "Errore"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	// Creare una risposta JSON di successo con il messaggio di log

	response := Response{SuccessMessage: fmt.Sprintf("ID: %v", id)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusCreated)
}

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Decodifica il corpo JSON della richiesta in una struttura
	var requestBody User

	auth := r.Header.Get("Authorization")
	uid := ps.ByName("uid")

	if uid != auth || auth == "" {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	bytes, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(bytes, &requestBody); err != nil {
		response := Response{ErrorMessage: "Error occured while decoding JSON1"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	username := requestBody.Content
	if username == "" {
		response := Response{ErrorMessage: "Error occured while decoding JSON2"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	int_uid, _ := strconv.Atoi(uid)

	valid, _ := rt.db.IdExists(int_uid)
	// Chiamo il database
	if !valid {
		response := Response{ErrorMessage: "Id utente loggato non valido"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	followedUid, err := rt.db.GetId(requestBody.Content)
	if err != nil {
		response := Response{ErrorMessage: "Username dell'utente da seguire non valido"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	valid, _ = rt.db.BanExists(int_uid, followedUid)
	if valid {
		// funziona al contrario perchè true vuol dire che è stato bannato
		response := Response{ErrorMessage: "Non puoi seguire un utente che hai bannato"}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}
	// se invece vuole seguire qualcuno che lo ha bannato
	valid, _ = rt.db.BanExists(followedUid, int_uid)
	if valid {
		response := Response{ErrorMessage: "Internal error"}
		sendJSONResponse(w, response, http.StatusForbidden)
	}

	err = rt.db.FollowUser(int_uid, followedUid)
	if errors.Is(err, database.ErrCannotFollowHimself) {
		response := Response{ErrorMessage: "User id e following user id coincidono"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrAlreadyFollowing) {
		response := Response{ErrorMessage: fmt.Sprintf("L'utente già segue il seguente following id: %d", followedUid)}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	if err != nil {
		response := Response{ErrorMessage: "Errore server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}
	// Creare una risposta JSON di successo con il messaggio di log
	response := Response{SuccessMessage: fmt.Sprintf("User followed successfully: %d", followedUid)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusCreated)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("uid")
	followedUid, _ := strconv.Atoi(ps.ByName("following_uid"))
	auth := r.Header.Get("Authorization")

	if uid != auth || auth == "" {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	int_uid, _ := strconv.Atoi(uid)

	err := rt.db.UnfollowUser(int_uid, followedUid)

	if errors.Is(err, sql.ErrNoRows) {
		response := Response{ErrorMessage: "You aren't following this user"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	if err != nil {
		response := Response{ErrorMessage: "Errore server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}
	// Creare una risposta JSON di successo con il messaggio di log
	response := Response{SuccessMessage: fmt.Sprintf("User successfully removed from following list: %d", followedUid)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusOK)
}

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Decodifica il corpo JSON della richiesta in una struttura
	var requestBody User

	auth := r.Header.Get("Authorization")
	uid := ps.ByName("uid")
	int_uid, _ := strconv.Atoi(uid)

	if uid != auth || auth == "" {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Gestisci errori di decodifica JSON
		response := Response{ErrorMessage: "Errore durante la decodifica del corpo JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	valid, _ := rt.db.IdExists(int_uid)
	if !valid {
		response := Response{ErrorMessage: "Id dell'utente loggato non valido"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	uidBanned, err := rt.db.GetId(requestBody.Content)

	if err != nil {
		response := Response{ErrorMessage: "Username dell'utente da bannare non valido"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	if flag, _ := rt.db.FollowExists(int_uid, uidBanned); flag {
		err = rt.db.UnfollowUser(int_uid, uidBanned)
		if err != nil {
			response := Response{ErrorMessage: "Errore interno del server, impossibile rimuovere l'utente dalla lista follow automaticamente"}
			sendJSONResponse(w, response, http.StatusInternalServerError)
			return
		}
	}

	if flag, _ := rt.db.FollowExists(uidBanned, int_uid); flag {
		err = rt.db.UnfollowUser(uidBanned, int_uid)
		if err != nil {
			response := Response{ErrorMessage: "Errore interno del server, impossibile rimuovere l'utente dalla lista follow automaticamente"}
			sendJSONResponse(w, response, http.StatusInternalServerError)
			return
		}
	}

	// fine controlli input, inizio controlli output
	err = rt.db.BanUser(int_uid, uidBanned)

	if errors.Is(err, database.ErrAlreadyBanned) {
		response := Response{ErrorMessage: "Utente già bannato"}
		sendJSONResponse(w, response, http.StatusNotModified)
		return
	}
	if errors.Is(err, database.ErrCannotBanHimself) {
		// Gestisci l'errore e invia una risposta JSON con il messaggio di errore
		response := Response{ErrorMessage: "L'utente non può bannarsi da solo"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	valid, _ = rt.db.BanExists(uidBanned, int_uid)
	if valid {
		response := Response{ErrorMessage: "Internal error"}
		sendJSONResponse(w, response, http.StatusForbidden)
	}

	if err != nil {
		response := Response{ErrorMessage: "Errore interno del server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	// Creare una risposta JSON di successo con il messaggio di log
	response := Response{SuccessMessage: fmt.Sprintf("Utente bannato con successo: %d", uidBanned)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusCreated)
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	uid := ps.ByName("uid")
	int_uid, _ := strconv.Atoi(uid)
	bannedUid, _ := strconv.Atoi(ps.ByName("banned_uid"))
	auth := r.Header.Get("Authorization")

	if uid != auth || auth == "" {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	err := rt.db.UnbanUser(int_uid, bannedUid)

	if errors.Is(err, sql.ErrNoRows) {
		response := Response{ErrorMessage: "You haven't banned this user"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	if err != nil {
		response := Response{ErrorMessage: "Errore server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}
	// Creare una risposta JSON di successo con il messaggio di log
	response := Response{SuccessMessage: fmt.Sprintf("User successfully removed from banned list: %d", bannedUid)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusOK)
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	uid, _ := strconv.Atoi(ps.ByName("uid"))
	username, err := rt.db.GetUsername(uid)
	// devo controllare se chi avanza la richiesta non è stato bloccato dall'utente richiesto
	auth := r.Header.Get("Authorization")
	int_auth, _ := strconv.Atoi(auth)

	if flag, _ := rt.db.BanExists(int_auth, uid); flag {
		response := Response{}
		sendJSONResponse(w, response, 206)
		return
	}

	if flag, _ := rt.db.BanExists(uid, int_auth); flag {
		response := Response{}
		sendJSONResponse(w, response, 400)
		return
	}

	if err != nil {
		response := Response{ErrorMessage: "Username non trovato o non valido"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	// voglio following, follower, username, numero di foto caricate
	following, err := rt.db.GetFollowing(uid)

	if err != nil {
		response := Response{ErrorMessage: "Errore server, impossibile richiedere lista following"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	followers, err := rt.db.GetFollowers(uid)

	if err != nil {
		response := Response{ErrorMessage: "Errore server, impossibile richiedere lista followers"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	var posts []database.Photo

	// posts deve essere già pronto e assemblato con tutti i dati tranne il base64
	posts, _ = rt.db.GetPosts(uid)

	profile := Profile{
		Username:  username,
		Id:        uid,
		Following: following,
		Followers: followers,
		Posts:     posts,
	}

	// Invia la risposta JSON
	sendJSONResponse(w, profile, http.StatusOK)
}

func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// controllare che sia loggato

	auth := r.Header.Get("Authorization")
	if auth == "" {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	u, _ := strconv.Atoi(auth)

	query := r.URL.Query().Get("username")

	if query == "" {
		response := Response{ErrorMessage: "Parametro di query mancante"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	users, err := rt.db.SearchUsers(u, query)
	if err != nil {
		response := Response{ErrorMessage: "Errore durante la ricerca nel database"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		response := Response{ErrorMessage: "Nessun utente trovato"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, users, http.StatusOK)
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	uid := ps.ByName("uid")

	// controllare che sia autorizzato
	auth := r.Header.Get("Authorization")
	if uid != auth || auth == "" {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	u, _ := strconv.Atoi(auth)
	photos, err := rt.db.GetStream(u)

	if err != nil {
		response := Response{ErrorMessage: "Errore server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, photos, http.StatusOK)
}

func (rt *_router) getId(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := r.URL.Query().Get("username")
	ret, err := rt.db.GetId(username)
	if err != nil {
		response := Response{ErrorMessage: "ID NOT FOUND"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	response := Response{SuccessMessage: fmt.Sprintf("%v", ret)}
	sendJSONResponse(w, response, http.StatusOK)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	// Imposta l'header Content-Type sulla risposta come application/json
	w.Header().Set("Content-Type", "application/json")

	// Imposta il codice di stato HTTP sulla risposta
	w.WriteHeader(statusCode)

	// Serializza la struttura dati in formato JSON e scrivi sulla risposta
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Se si verifica un errore durante la serializzazione JSON, invia una risposta di errore
		http.Error(w, "Errore durante la serializzazione della risposta JSON", http.StatusInternalServerError)
	}
}
