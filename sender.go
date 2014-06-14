package main

import (
	"code.google.com/p/appengine-go/appengine"
	"code.google.com/p/appengine-go/appengine/datastore"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type SenderType struct {
	Key           string          `json:"key", datastore:"-"`
	Name          string          `json:"name"`
	Scope         string          `json:"scope"`
	AuthURL       string          `json:"auth_url"`
	TokenURL      string          `json:"token_url"`
	SenderClients []*SenderClient `json:"sender_clients,omitempty", datastore:"-"`
}

type SenderTypeController struct{}

func (s *SenderTypeController) Options(response http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(response)
	encoder.Encode(SenderType{})
}

func (s *SenderTypeController) Index(response http.ResponseWriter, request *http.Request) {
	senderTypes, err := app.SenderTypeOrm.FindAll(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, senderType := range senderTypes {
		senderClients, err := app.SenderClientOrm.FindAllForSenderType(request, senderType.Key)
		if err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
			return
		}
		senderType.SenderClients = senderClients
	}
	encoder := json.NewEncoder(response)
	encoder.Encode(senderTypes)
}

func (s *SenderTypeController) Create(response http.ResponseWriter, request *http.Request) {
	senderType := new(SenderType)
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&senderType)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	senderType, err = app.SenderTypeOrm.Save(request, senderType)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(response)
	encoder.Encode(senderType)
}

func (s *SenderTypeController) Show(response http.ResponseWriter, request *http.Request) {
	encodedKey := mux.Vars(request)["key"]
	senderType, err := app.SenderTypeOrm.Find(request, encodedKey)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(response)
	encoder.Encode(senderType)
}

func (s *SenderTypeController) Update(response http.ResponseWriter, request *http.Request) {
	encodedKey := mux.Vars(request)["key"]
	senderType, err := app.SenderTypeOrm.Find(request, encodedKey)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&senderType)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	senderType, err = app.SenderTypeOrm.Save(request, senderType)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(response)
	encoder.Encode(senderType)
}

func (s *SenderTypeController) Delete(response http.ResponseWriter, request *http.Request) {
	encodedKey := mux.Vars(request)["key"]
	senderType, err := app.SenderTypeOrm.Delete(request, encodedKey)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(response)
	encoder.Encode(senderType)
}

type SenderClient struct {
	Key           string `json:"key", datastore:"-"`
	ClientId      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	AuthCode      string `json:"auth_code"`
	SenderTypeKey string `json:"sender_type_key"`
}

type SenderClientController struct{}

func (s *SenderClientController) Options(response http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(response)
	encoder.Encode(SenderClient{})
}

func (s *SenderClientController) Index(response http.ResponseWriter, request *http.Request) {
	senderTypeKey, keyPresent := mux.Vars(request)["type_key"]
	var senderClients []*SenderClient
	var err error
	if keyPresent {
		senderClients, err = app.SenderClientOrm.FindAllForSenderType(request, senderTypeKey)
	} else {
		senderClients, err = app.SenderClientOrm.FindAll(request)
	}
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(response)
	encoder.Encode(senderClients)
}

func (s *SenderClientController) Create(response http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	senderClient := new(SenderClient)
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&senderClient)
	if err != nil {
		context.Infof("Decoding failed for request %v", request)
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	senderTypeKey, keyPresent := mux.Vars(request)["type_key"]
	if keyPresent {
		senderClient.SenderTypeKey = senderTypeKey
	}
	incompleteKey := datastore.NewIncompleteKey(context, "SenderClient", nil)
	key, err := datastore.Put(context, incompleteKey, senderClient)
	if err != nil {
		context.Infof("Insert failed for key '%v' and entity '%+v'", incompleteKey, senderClient)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	senderClient.Key = key.Encode()
	encoder := json.NewEncoder(response)
	encoder.Encode(senderClient)
}

func (s *SenderClientController) Show(response http.ResponseWriter, request *http.Request) {}

func (s *SenderClientController) Update(response http.ResponseWriter, request *http.Request) {}

func (s *SenderClientController) Delete(response http.ResponseWriter, request *http.Request) {}
