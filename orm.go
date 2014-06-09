package main

import (
	"code.google.com/p/appengine-go/appengine"
	"code.google.com/p/appengine-go/appengine/datastore"
	"net/http"
)

type SenderTypeOrm struct{}

func (s *SenderTypeOrm) FindAll(request *http.Request) ([]*SenderType, error) {
	context := appengine.NewContext(request)
	query := datastore.NewQuery("SenderType")
	var senderTypes []*SenderType
	keys, err := query.GetAll(context, &senderTypes)
	if err != nil {
		return nil, err
	}
	if senderTypes == nil {
		senderTypes = make([]*SenderType, 0)
	}
	for i, senderType := range senderTypes {
		encodedKey := keys[i].Encode()
		senderType.Key = encodedKey
	}
	return senderTypes, nil
}

func (s *SenderTypeOrm) Save(request *http.Request, senderType *SenderType) (*SenderType, error) {
	context := appengine.NewContext(request)
	incompleteKey := datastore.NewIncompleteKey(context, "SenderType", nil)
	key, err := datastore.Put(context, incompleteKey, senderType)
	if err != nil {
		return nil, err
	}
	senderType.Key = key.Encode()
	return senderType, nil
}

func (s *SenderTypeOrm) Find(request *http.Request, key string) (*SenderType, error) {
	context := appengine.NewContext(request)
	decodedKey, err := datastore.DecodeKey(key)
	if err != nil {
		return nil, err
	}
	senderType := new(SenderType)
	err = datastore.Get(context, decodedKey, senderType)
	if err != nil {
		return nil, err
	}
	senderType.Key = key
	return senderType, nil
}

func (s *SenderTypeOrm) Delete(request *http.Request, key string) (*SenderType, error) {
	context := appengine.NewContext(request)
	decodedKey, err := datastore.DecodeKey(key)
	if err != nil {
		return nil, err
	}
	senderType := new(SenderType)
	err = datastore.Get(context, decodedKey, senderType)
	if err != nil {
		return nil, err
	}
	err = datastore.Delete(context, decodedKey)
	if err != nil {
		return nil, err
	}
	senderType.Key = ""
	return senderType, nil
}

type SenderClientOrm struct{}

func (s *SenderClientOrm) FindAll(request *http.Request) ([]*SenderClient, error) {
	context := appengine.NewContext(request)
	query := datastore.NewQuery("SenderClient")
	var senderClients []*SenderClient
	keys, err := query.GetAll(context, &senderClients)
	if err != nil {
		return nil, err
	}
	if senderClients == nil {
		senderClients = make([]*SenderClient, 0)
	}
	for i, senderClient := range senderClients {
		encodedKey := keys[i].Encode()
		senderClient.Key = encodedKey
	}
	return senderClients, nil
}

func (s *SenderClientOrm) FindAllForSenderType(request *http.Request, senderTypeKey string) ([]*SenderClient, error) {
	context := appengine.NewContext(request)
	query := datastore.NewQuery("SenderClient")
	query.Filter("SenderTypeKey =", senderTypeKey)
	var senderClients []*SenderClient
	keys, err := query.GetAll(context, &senderClients)
	if err != nil {
		return nil, err
	}
	if senderClients == nil {
		senderClients = make([]*SenderClient, 0)
	}
	for i, senderClient := range senderClients {
		encodedKey := keys[i].Encode()
		senderClient.Key = encodedKey
	}
	return senderClients, nil
}

func (s *SenderClientOrm) Save(request *http.Request, senderClient *SenderClient) (*SenderClient, error) {
	context := appengine.NewContext(request)
	incompleteKey := datastore.NewIncompleteKey(context, "SenderClient", nil)
	key, err := datastore.Put(context, incompleteKey, senderClient)
	if err != nil {
		return nil, err
	}
	senderClient.Key = key.Encode()
	return senderClient, nil
}

func (s *SenderClientOrm) Find(request *http.Request, key string) (*SenderClient, error) {
	context := appengine.NewContext(request)
	decodedKey, err := datastore.DecodeKey(key)
	if err != nil {
		return nil, err
	}
	senderClient := new(SenderClient)
	err = datastore.Get(context, decodedKey, senderClient)
	if err != nil {
		return nil, err
	}
	senderClient.Key = key
	return senderClient, nil
}

func (s *SenderClientOrm) Delete(request *http.Request, key string) (*SenderClient, error) {
	context := appengine.NewContext(request)
	decodedKey, err := datastore.DecodeKey(key)
	if err != nil {
		return nil, err
	}
	senderClient := new(SenderClient)
	err = datastore.Get(context, decodedKey, senderClient)
	if err != nil {
		return nil, err
	}
	err = datastore.Delete(context, decodedKey)
	if err != nil {
		return nil, err
	}
	senderClient.Key = ""
	return senderClient, nil
}
