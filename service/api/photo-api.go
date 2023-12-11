package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Leggi il corpo della richiesta

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response := Response{ErrorMessage: "Errore durante la lettura del corpo della richiesta"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// Decodifica il corpo JSON
	var requestBody Photo
	if err := json.Unmarshal(body, &requestBody); err != nil {
		response := Response{ErrorMessage: "Errore durante la decodifica del corpo JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// Ottieni i dati base64 dell'immagine
	photoData := requestBody.Content
	username := requestBody.Owner
	uid, err := rt.db.GetId(username)
	if err != nil {
		response := Response{ErrorMessage: "Errore interno in fase di accesso all'id associato all'utente"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	// Decodifica i dati base64
	imageData, err := base64.StdEncoding.DecodeString(photoData)
	if err != nil {
		response := Response{ErrorMessage: "Errore durante la decodifica del base64"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// Salva l'immagine su disco (in questo caso, come file PNG)
	photoId, err := rt.db.CreatePhoto(uid)
	if photoId < 0 {
		response := Response{ErrorMessage: "photoId < 0 from database has been returned"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}
	if err != nil {
		response := Response{ErrorMessage: "Errore durante la creazione della foto nel database"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	path := "photos/" + strconv.Itoa(photoId) + ".png"

	err = saveImageLocally(imageData, path)

	if err != nil {
		response := Response{ErrorMessage: "Errore durante il salvataggio in locale della foto"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// Rispondi con un messaggio di successo
	response := Response{SuccessMessage: "Foto salvata con successo e aggiunta al database"}
	sendJSONResponse(w, response, http.StatusOK)
}

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Decodifica il corpo JSON della richiesta in una struttura
	var requestBody User
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Gestisci errori di decodifica JSON
		response := Response{ErrorMessage: "Errore durante la decodifica del corpo JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	uid, err := rt.db.GetId(requestBody.Content)
	if err != nil {
		response := Response{ErrorMessage: "Username non valido"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	photoId, _ := strconv.Atoi(ps.ByName("photo_id"))
	valid, _ := rt.db.PhotoIdExists(photoId)
	if !valid {
		response := Response{ErrorMessage: "PhotoId non valido"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	err = rt.db.LikePhoto(uid, photoId)

	if errors.Is(err, database.ErrAlreadyLiked) {
		response := Response{ErrorMessage: "Hai già messo like alla foto"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	if err != nil {
		response := Response{ErrorMessage: "Errore interno del server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	// Creare una risposta JSON di successo con il messaggio di log
	response := Response{SuccessMessage: fmt.Sprintf("Like aggiunto con successo alla foto con id: %d", photoId)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusOK)
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid, _ := strconv.Atoi(ps.ByName("uid"))
	photoId, _ := strconv.Atoi(ps.ByName("photo_id"))

	err := rt.db.UnlikePhoto(uid, photoId)

	if errors.Is(err, sql.ErrNoRows) {
		response := Response{ErrorMessage: "Il tuo like non era presente per questa foto"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	if err != nil {
		response := Response{ErrorMessage: "Errore server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}
	// Creare una risposta JSON di successo con il messaggio di log
	response := Response{SuccessMessage: fmt.Sprintf("Photo successfully removed from liked list: %d", photoId)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusOK)
}

func saveImageLocally(imageData []byte, filePath string) error {
	// Apri il file in modalità di scrittura
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Scrivi i dati dell'immagine nel file
	_, err = file.Write(imageData)
	if err != nil {
		return err
	}

	fmt.Printf("Immagine salvata correttamente: %s\n", filePath)
	return nil
}
