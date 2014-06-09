package main

import (
	"code.google.com/p/appengine-go/appengine"
	"code.google.com/p/appengine-go/appengine/datastore"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
)

var indexPath *regexp.Regexp = regexp.MustCompile("sender_types")

type SenderType struct {
	Key           string         `json:"key", datastore:"-"`
	Name          string         `json:"name"`
	SenderConfigs []SenderConfig `json:"sender_configs,omitempty"`
}

type SenderTypeController struct{}

func (s *SenderTypeController) Index(response http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	query := datastore.NewQuery("SenderType")
	var senderTypes []*SenderType
	keys, err := query.GetAll(context, &senderTypes)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	if senderTypes == nil {
		senderTypes = make([]*SenderType, 0)
	}
	for i, senderType := range senderTypes {
		encodedKey := keys[i].Encode()
		senderType.Key = encodedKey
	}
	encoder := json.NewEncoder(response)
	encoder.Encode(senderTypes)
}

func (s *SenderTypeController) Create(response http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	senderType := new(SenderType)
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&senderType)
	if err != nil {
		context.Infof("Decoding failed for request %v", request)
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	incompleteKey := datastore.NewIncompleteKey(context, "SenderType", nil)
	key, err := datastore.Put(context, incompleteKey, senderType)
	if err != nil {
		context.Infof("Insert failed for key '%v' and entity '%+v'", incompleteKey, senderType)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	senderType.Key = key.Encode()
	encoder := json.NewEncoder(response)
	encoder.Encode(senderType)
}

func (s *SenderTypeController) Show(response http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	encodedKey := mux.Vars(request)["key"]
	key, err := datastore.DecodeKey(encodedKey)
	if err != nil {
		context.Infof("Key decoding failed for key %s", encodedKey)
		http.NotFound(response, request)
		return
	}
	senderType := new(SenderType)
	err = datastore.Get(context, key, senderType)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusCreated)
	senderType.Key = key.Encode()
	encoder := json.NewEncoder(response)
	encoder.Encode(senderType)
}

func (s *SenderTypeController) Update(response http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	encodedKey := mux.Vars(request)["key"]
	key, err := datastore.DecodeKey(encodedKey)
	if err != nil {
		context.Infof("Key decoding failed for key %s", encodedKey)
		http.NotFound(response, request)
		return
	}
	senderType := new(SenderType)
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&senderType)
	if err != nil {
		context.Infof("Decoding failed for request %v", request)
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	key, err = datastore.Put(context, key, senderType)
	if err != nil {
		context.Infof("Update failed for key '%v' and entity '%+v'", key, senderType)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	senderType.Key = key.Encode()
	encoder := json.NewEncoder(response)
	encoder.Encode(senderType)
}

func (s *SenderTypeController) Delete(response http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	encodedKey := mux.Vars(request)["key"]
	key, err := datastore.DecodeKey(encodedKey)
	if err != nil {
		context.Infof("Key decoding failed for key %s", encodedKey)
		http.NotFound(response, request)
		return
	}
	err = datastore.Delete(context, key)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
}

type SenderConfig struct {
	Key           string `json:"key", datastore:"-"`
	ClientId      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	AuthCode      string `json:"auth_code"`
	SenderTypeKey string `json:"sender_type_key"`
}

type SenderConfigsController struct{}

func (s *SenderConfigsController) Index(response http.ResponseWriter, request *http.Request) {
	context := appengine.NewContext(request)
	senderTypeKey, keyPresent := mux.Vars(request)["type_key"]
	query := datastore.NewQuery("SenderConfig")
	if keyPresent {
		query.Filter("SenderTypeKey =", senderTypeKey)
	}
	var senderConfigs []*SenderConfig
	keys, err := query.GetAll(context, &senderConfigs)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	if senderConfigs == nil {
		senderConfigs = make([]*SenderConfig, 0)
	}
	for i, senderConfig := range senderConfigs {
		encodedKey := keys[i].Encode()
		senderConfig.Key = encodedKey
		senderConfig.SenderTypeKey = senderTypeKey
	}
	encoder := json.NewEncoder(response)
	encoder.Encode(senderConfigs)
}

func (s *SenderConfigsController) Create(response http.ResponseWriter, request *http.Request) {}

func (s *SenderConfigsController) Show(response http.ResponseWriter, request *http.Request) {}

func (s *SenderConfigsController) Update(response http.ResponseWriter, request *http.Request) {}

func (s *SenderConfigsController) Delete(response http.ResponseWriter, request *http.Request) {}
