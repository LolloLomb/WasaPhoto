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
	"path/filepath"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

type PhotoUpload struct {
	Content string `json:"content"`
	Owner   string `json:"username_owner"`
}

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Leggi il corpo della richiesta

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response := Response{ErrorMessage: "Errore durante la lettura del corpo della richiesta"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// Decodifica il corpo JSON
	var requestBody PhotoUpload

	if err := json.Unmarshal(body, &requestBody); err != nil {
		response := Response{ErrorMessage: "Errore durante la decodifica del corpo JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// Ottieni i dati base64 dell'immagine
	photoData := requestBody.Content
	owner := requestBody.Owner

	// ottengo l'id di chi vuole pubblicare la foto
	uid, _ := rt.db.GetId(owner)

	// ottengo l'auth nell'header
	auth := r.Header.Get("Authorization")

	// converto l'auth in intero per ottenere l'username e confrontarlo con chi vuole pubblicare la foto
	int_auth, _ := strconv.Atoi(auth)
	owner_from_auth, err := rt.db.GetUsername(int_auth)

	// se non sono uguali oppure l'auth è vuoto, errore
	if owner != owner_from_auth || owner_from_auth == "" {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	if err != nil {
		response := Response{ErrorMessage: "Utente non trovato"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// Decodifica i dati base64
	imageData, err := base64.StdEncoding.DecodeString(photoData)
	if err != nil {
		response := Response{ErrorMessage: "Errore durante la decodifica del base64"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// salva l'immagine nel database
	photoId, err := rt.db.CreatePhoto(uid)

	/* if photoId < 0 {
		response := Response{ErrorMessage: "photoId < 0 from database has been returned"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}*/

	if err != nil {
		response := Response{ErrorMessage: "Errore durante la creazione della foto nel database"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// Verifica se la cartella "tmp" esiste
	_, err = os.Stat("../../tmp/")

	if os.IsNotExist(err) {
		// Se la cartella non esiste, creala
		_ = os.Mkdir("../../tmp/", 0755) // 0755 è il permesso di default per la cartella
	}

	path := "../../tmp/" + strconv.Itoa(photoId) + ".png"
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

	likedUid, err := rt.db.GetId(requestBody.Content)
	// vedere se l'utente è autorizzato e poi 2 se il suo username esiste
	auth := r.Header.Get("Authorization")
	if auth != requestBody.Content || auth == "" {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	// 2
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

	// ora devo vedere se sono il proprietario
	flag, _ := rt.db.IsPhotoOwner(likedUid, photoId)
	if flag {
		response := Response{ErrorMessage: "Non puoi mettere like ad una tua foto"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	err = rt.db.LikePhoto(likedUid, photoId)

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
	sendJSONResponse(w, response, http.StatusCreated)
}

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Decodifica il corpo JSON della richiesta in una struttura
	var requestBody database.Comment
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Gestisci errori di decodifica JSON
		response := Response{ErrorMessage: "Errore durante la decodifica del corpo JSON"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	// verificare che sei l'autore del commento che stai inviando
	auth := r.Header.Get("Authorization")
	if requestBody.Owner != auth || auth == "" {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	uid, err := rt.db.GetId(requestBody.Owner)
	if err != nil {
		response := Response{ErrorMessage: "Username non valido"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	text := requestBody.Content

	photoId, _ := strconv.Atoi(ps.ByName("photo_id"))
	valid, _ := rt.db.PhotoIdExists(photoId)
	if !valid {
		response := Response{ErrorMessage: "PhotoId non valido"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	err = rt.db.CommentPhoto(uid, photoId, text)

	if err != nil {
		response := Response{ErrorMessage: "Errore interno del server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	// Creare una risposta JSON di successo con il messaggio di log
	response := Response{SuccessMessage: fmt.Sprintf("Commento aggiunto con successo alla foto con id: %d", photoId)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusCreated)
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid, _ := strconv.Atoi(ps.ByName("uid"))
	photoId, _ := strconv.Atoi(ps.ByName("photo_id"))

	auth := r.Header.Get("Authorization")
	// devo vedere se coincide con l'uid
	authUid, _ := rt.db.GetId(auth)
	if authUid != uid || auth == "" {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

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

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	photoId, _ := strconv.Atoi(ps.ByName("photo_id"))
	commentId, _ := strconv.Atoi(ps.ByName("comment_id"))

	// l'username esiste?
	auth := r.Header.Get("Authorization")
	authUid, err := rt.db.GetId(auth)
	if authUid == 0 || auth == "" || err != nil {
		response := Response{ErrorMessage: "Username non valido nell'header"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
	// sei il proprietario?
	flag, _ := rt.db.IsPhotoOwner(authUid, photoId)
	if !flag {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	err = rt.db.UncommentPhoto(photoId, commentId)

	if errors.Is(err, sql.ErrNoRows) {
		response := Response{ErrorMessage: "Il tuo commento non era presente per questa foto"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	if err != nil {
		response := Response{ErrorMessage: "Errore server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}
	// Creare una risposta JSON di successo con il messaggio di log
	response := Response{SuccessMessage: fmt.Sprintf("Comment successfully removed from photo: %d", photoId)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusOK)
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	photoId, _ := strconv.Atoi(ps.ByName("photo_id"))

	valid, err := rt.db.PhotoIdExists(photoId)

	if errors.Is(err, database.ErrPhotoDontExists) || !valid {
		response := Response{ErrorMessage: "PhotoId non valido"}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	if err != nil {
		response := Response{ErrorMessage: "Errore server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	auth := r.Header.Get("Authorization")
	authId, _ := rt.db.GetId(auth)

	// devo vedere se authId e "uid" dell'owner nel db coincidono

	if flag, _ := rt.db.IsPhotoOwner(authId, photoId); !flag {
		response := Response{ErrorMessage: ErrForbidden.Error()}
		sendJSONResponse(w, response, http.StatusForbidden)
		return
	}

	err = rt.db.DeletePhoto(photoId)

	if err != nil {
		response := Response{ErrorMessage: "Errore server"}
		sendJSONResponse(w, response, http.StatusInternalServerError)
		return
	}

	response := Response{SuccessMessage: fmt.Sprintf("Photo successfully removed: %d", photoId)}

	// Invia la risposta JSON
	sendJSONResponse(w, response, http.StatusOK)
}

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Costruisci il percorso del file
	photoID := ps.ByName("photo_id")
	photoPath := filepath.Join("../../tmp/", photoID+".png")

	// Verifica se il file esiste
	if _, err := os.Stat(photoPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Leggi il file
	file, err := os.ReadFile(photoPath)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Codifica il contenuto del file in base64
	base64Encoding := base64.StdEncoding.EncodeToString(file)

	// Crea la risposta
	response := PhotoUpload{
		Content: base64Encoding,
		Owner:   "_",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
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

	// fmt.Printf("Immagine salvata correttamente: %s\n", filePath)
	return nil
}
