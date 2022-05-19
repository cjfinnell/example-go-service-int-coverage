package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	paramKey   = "key"
	paramValue = "value"

	routeKey      = "/:" + paramKey
	routeKeyValue = "/:" + paramKey + "/:" + paramValue
)

func (s *server) handleGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName(paramKey)
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	value, err := s.redis.Get(key)
	if err == ErrNotFound {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get value for key", http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(value))
	if err != nil {
		log.Printf("failed to write response: %s", err)
	}
}

func (s *server) handleSet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName(paramKey)
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	value := ps.ByName(paramValue)
	if value == "" {
		http.Error(w, "missing value", http.StatusBadRequest)

		return
	}

	err := s.redis.Set(key, value)
	if err != nil {
		http.Error(w, "failed to set value for key", http.StatusInternalServerError)
		return
	}

	log.Printf("key '%s' populated with value '%s'", key, value)

	w.WriteHeader(http.StatusCreated)
}

func (s *server) handleDel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName(paramKey)
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	err := s.redis.Del(key)
	if err != nil {
		http.Error(w, "failed to delete key", http.StatusInternalServerError)
		return
	}

	log.Printf("key '%s' deleted", key)

	w.WriteHeader(http.StatusOK)
}
