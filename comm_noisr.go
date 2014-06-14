package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var router *mux.Router
var app *App

func init() {
	app = &App{&SenderTypeOrm{}, &SenderClientOrm{}}
	router = mux.NewRouter()
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "Hello world")
	}).Methods("GET")

	senderTypesController := &SenderTypeController{}
	router.HandleFunc("/sender_types", senderTypesController.Options).Methods("OPTIONS")
	router.HandleFunc("/sender_types", senderTypesController.Index).Methods("GET")
	router.HandleFunc("/sender_types", senderTypesController.Create).Methods("POST")
	router.HandleFunc("/sender_types/{key}", senderTypesController.Show).Methods("GET")
	router.HandleFunc("/sender_types/{key}", senderTypesController.Update).Methods("PUT")
	router.HandleFunc("/sender_types/{key}", senderTypesController.Delete).Methods("DELETE")

	SenderClientController := &SenderClientController{}
	router.HandleFunc("/sender_clients", SenderClientController.Options).Methods("OPTIONS")
	router.HandleFunc("/sender_clients", SenderClientController.Index).Methods("GET")
	router.HandleFunc("/sender_clients", SenderClientController.Create).Methods("POST")
	router.HandleFunc("/sender_types/{type_key}/sender_clients", SenderClientController.Index).Methods("GET")
	router.HandleFunc("/sender_types/{type_key}/sender_clients", SenderClientController.Create).Methods("POST")
	http.Handle("/", router)
}
