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
	routes().ServeHTTP(w, r)
}
func writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
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
func errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload jsonResponse

	payload.Error = true
	payload.Message = err.Error()

	return writeJSON(w, statusCode, payload)
}
func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
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
	mux.Post("/result", generateResult)

	return mux
}
func handleResult(w http.ResponseWriter, r *http.Request) {
	var result jsonResponse
	result.Error = false
	result.Message = "success"
	result.Data = map[string]string{
		"message": "Hello World",
	}
	writeJSON(w, http.StatusOK, result)
}
func generateResult(w http.ResponseWriter, r *http.Request) {
	var numbers []int32
	err := readJSON(w, r, &numbers)
	if err != nil {
		if err == io.EOF {
			errorJSON(w, errors.New("please provide proper input json format like [1,2,3]"))
			return
		}
		errorJSON(w, err)
		return
	}
	if len(numbers) == 0 {
		errorJSON(w, errors.New("numbers must not be empty"))
		return
	}
	var result int32
	for _, n := range numbers {
		result += n
	}
	writeJSON(w, http.StatusOK, result)
}
