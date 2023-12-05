package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"

	//"log"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

type User struct {
	Content string `json:"username"`
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

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Decodifica il corpo JSON della richiesta in una struttura
	var requestBody NewUsername
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Gestisci errori di decodifica JSON
		response := Response{ErrorMessage: "Errore durante la decodifica del corpo JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// Estrai il nuovo username dalla struttura
	newUsername := requestBody.Content
	if newUsername == "" {
		response := Response{ErrorMessage: "Error occured while decoding JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	uid, _ := strconv.Atoi(ps.ByName("uid"))
	// Chiamare il pacchetto del database per impostare il nuovo username
	err := rt.db.SetMyNickname(newUsername, uid)
	if err == database.ErrNameAlreadyTaken {
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

func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decodifica il corpo JSON della richiesta in una struttura
	var requestBody User
	bytes, _ := io.ReadAll(r.Body)
	//println(string(bytes))
	if err := json.Unmarshal(bytes, &requestBody); err != nil {
		// if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Gestisci errori di decodifica JSON
		response := Response{ErrorMessage: "Error occured while decoding JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	username := requestBody.Content
	if username == "" {
		response := Response{ErrorMessage: "Error occured while decoding JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	//println(username)
	err := rt.db.CreateUser(username)
	id, _ := rt.db.GetId(username)
	//println(username, id)
	//log.Printf("%v", err)
	if err == database.ErrUserAlreadyInTheDb {
		// Gestisci l'errore e invia una risposta JSON con il messaggio di errore
		response := Response{SuccessMessage: fmt.Sprintf("Successful Login. ID: %v", id)}
		sendJSONResponse(w, response, 200)
		return
	}
	if err != nil {
		response := Response{ErrorMessage: "Errore"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	// Creare una risposta JSON di successo con il messaggio di log

	response := Response{SuccessMessage: fmt.Sprintf("Utente creato con successo. ID: %v", id)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusCreated)
}

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Decodifica il corpo JSON della richiesta in una struttura
	var requestBody Identifier
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Gestisci errori di decodifica JSON
		response := Response{ErrorMessage: "Errore durante la decodifica del corpo JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	uid, _ := strconv.Atoi(ps.ByName("uid"))
	t := requestBody.Content
	followedUid, _ := strconv.Atoi(t)
	err := rt.db.FollowUser(uid, followedUid)
	if err == database.ErrCannotFollowHimself {
		response := Response{ErrorMessage: "User id e following user id coincidono"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	if err == database.ErrAlreadyFollowing {
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
	sendJSONResponse(w, response, http.StatusOK)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	uid, _ := strconv.Atoi(ps.ByName("uid"))
	followedUid, _ := strconv.Atoi(ps.ByName("following_uid"))

	err := rt.db.UnfollowUser(uid, followedUid)

	if err == sql.ErrNoRows {
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
	var requestBody Identifier
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Gestisci errori di decodifica JSON
		response := Response{ErrorMessage: "Errore durante la decodifica del corpo JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	uidBanned, _ := strconv.Atoi(requestBody.Content)

	uid, _ := strconv.Atoi(ps.ByName("uid"))
	// Chiamo il database
	err := rt.db.BanUser(uid, uidBanned)

	if err == database.ErrAlreadyBanned {
		response := Response{ErrorMessage: "Utente già bannato"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	if err == database.ErrCannotBanHimself {
		// Gestisci l'errore e invia una risposta JSON con il messaggio di errore
		response := Response{ErrorMessage: "L'utente non può bannarsi da solo"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	if err != nil {
		response := Response{ErrorMessage: "Errore interno del server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	// Creare una risposta JSON di successo con il messaggio di log
	response := Response{SuccessMessage: fmt.Sprintf("Utente bannato con successo: %d", uidBanned)}

	// Invia la risposta JSON
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
