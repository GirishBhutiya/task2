package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	gochi "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// Entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler")
	routes().ServeHTTP(w, r)
}
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	log.Println("WriteJSON")
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
	log.Println("ReadJSON")
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

func routes() http.Handler {
	log.Println("routes")
	mux := chi.NewRouter()

	//specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(gochi.Heartbeat("/ping"))

	mux.Post("/", handleResult)
	mux.Get("/", handleResult)

	return mux
}
func handleResult(w http.ResponseWriter, r *http.Request) {
	log.Println("handleResult")
	var result jsonResponse
	result.Error = false
	result.Message = "success"
	result.Data = map[string]string{
		"message": "Hello World",
	}
	WriteJSON(w, http.StatusOK, result)
}
