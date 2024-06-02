package main

import (
	"encoding/json"
	"go-cmd-book/todo"
	"log"
	"net/http"
	"sync"
	"time"
)

func newMux(todoFile string) http.Handler {
	mu := sync.Mutex{}
	t := todoRouter(todoFile, &mu)

	m := http.NewServeMux()
	m.HandleFunc("/", rootHandler)
	m.Handle("/todo", http.StripPrefix("/todo", t))
	m.Handle("/todo/", http.StripPrefix("/todo/", t))
	return m
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		replyError(w, r, http.StatusNotFound, "")
		return
	}
	content := "Hello, from the the api"
	replyTextContent(w, r, http.StatusOK, content)
}

func replyTextContent(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(content))
}

func replyError(w http.ResponseWriter, r *http.Request, status int, message string) {
	log.Printf("%s %s: Error %d %s", r.URL, r.Method, status, message)
	http.Error(w, http.StatusText(status), status)
}

func replyJSONContent(w http.ResponseWriter, r *http.Request, status int, resp *todoResponse) {
	body, err := json.Marshal(resp)
	if err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

type todoResponse struct {
	Results todo.List `json:"results"`
}

func (r *todoResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		Results      todo.List `json:"results"`
		Date         int64     `json:"date"`
		TotalResults int       `json:"total_results"`
	}{
		Results:      r.Results,
		Date:         time.Now().Unix(),
		TotalResults: len(r.Results),
	}
	return json.Marshal(resp)
}
