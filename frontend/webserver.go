package frontend

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/menxqk/hexarc/core"
)

const (
	publicHtml = "public_html"
	indexPage  = "index.html"
)

type webserverFrontEnd struct {
	store *core.KeyValueStore
}

func (ws *webserverFrontEnd) Start(store *core.KeyValueStore) error {
	ws.store = store

	webserverPort := os.Getenv("WEBSERVER_PORT")

	r := mux.NewRouter()
	r.Use(ws.loggingMiddleware)

	fs := http.FileServer(http.Dir(publicHtml))

	r.Handle("/", fs).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fs))
	r.HandleFunc("/api/{key}", ws.getHandler).Methods("GET")
	r.HandleFunc("/api/{key}", ws.putHandler).Methods("PUT")
	r.HandleFunc("/api/{key}", ws.deleteHandler).Methods("DELETE")
	r.HandleFunc("/all", ws.getAllHandler).Methods("GET")

	r.HandleFunc("/", ws.notAllowedHandler)
	r.HandleFunc("/api", ws.notAllowedHandler)
	r.HandleFunc("/api/{key}", ws.notAllowedHandler)

	return http.ListenAndServe(":"+webserverPort, r)
}

func (ws *webserverFrontEnd) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (ws *webserverFrontEnd) notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

type jsonResponse struct {
	Ok    bool        `json:"ok"`
	Count int         `json:"count"`
	Data  interface{} `json:"data,omitempty"`
}

func (ws *webserverFrontEnd) getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	value, err := ws.store.Get(key)
	if err == core.ErrorNoSuchKey {
		http.Error(w, core.ErrorNoSuchKey.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := jsonResponse{
		Ok:    true,
		Count: 1,
		Data:  value,
	}
	json, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ws *webserverFrontEnd) putHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = ws.store.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ws *webserverFrontEnd) deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := ws.store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ws *webserverFrontEnd) getAllHandler(w http.ResponseWriter, r *http.Request) {
	all := ws.store.GetAll()
	resp := jsonResponse{Ok: true, Count: len(all), Data: all}

	json, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
