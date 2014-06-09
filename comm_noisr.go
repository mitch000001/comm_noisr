package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var router *mux.Router

func init() {
	router = mux.NewRouter()
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "Hello world")
	}).Methods("GET")

	senderTypesController := &SenderTypeController{}
	router.HandleFunc("/sender_types", senderTypesController.Index).Methods("GET")
	router.HandleFunc("/sender_types", senderTypesController.Create).Methods("POST")
	router.HandleFunc("/sender_types/{key}", senderTypesController.Show).Methods("GET")
	router.HandleFunc("/sender_types/{key}", senderTypesController.Update).Methods("PUT")
	router.HandleFunc("/sender_types/{key}", senderTypesController.Delete).Methods("DELETE")

	senderConfigsController := &SenderConfigsController{}
	router.HandleFunc("/sender_configs", senderConfigsController.Index).Methods("GET")
	router.HandleFunc("/sender_types/{type_key}/sender_configs", senderConfigsController.Index).Methods("GET")
	http.Handle("/", router)
}
