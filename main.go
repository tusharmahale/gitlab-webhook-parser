package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"quillbot.com/gitlab-webhook-parser/src/gitlab"
)

const (
	path = "/webhooks"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc(path, gitlab.HandleWebhook).Methods(http.MethodPost)
	router.HandleFunc("/enable/{pId}/{branchName}", gitlab.EnableMerge).Methods(http.MethodPost)
	router.HandleFunc("/disable/{pId}/{branchName}", gitlab.DisableMerge).Methods(http.MethodPost)
	router.HandleFunc("/healthz", healthCheck).Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(invalidURL)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Error starting server ", err)
	}
}

func invalidURL(w http.ResponseWriter, r *http.Request) {
	returnPl := map[string]string{}
	w.Header().Set("Content-type", "application/json")
	returnPl["error"] = "invalid URL"
	w.WriteHeader(http.StatusNotFound)
	returnBytes, _ := json.Marshal(returnPl)
	w.Write(returnBytes)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	returnPl := map[string]string{}
	w.Header().Set("Content-type", "application/json")
	returnPl["status"] = "success"
	w.WriteHeader(http.StatusOK)
	returnBytes, _ := json.Marshal(returnPl)
	w.Write(returnBytes)
}
