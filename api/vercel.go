package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

// Entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	//app.ServeHTTP(w, r)
	var payload jsonResponse

	payload.Error = true
	payload.Message = "Darsh Bhutiya"
	WriteJSON(w, http.StatusOK, payload)
}
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {

	maxBytes := 1024 * 1024 // 1 mb

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)
	if err != nil {
		log.Println("Girish")
		log.Println(err)
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		log.Println("Giris1")
		log.Println(err)
		return errors.New("body must have single JSON value")
	}
	return nil
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
