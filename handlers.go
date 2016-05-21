package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Record struct {
	// The status is defined in bringup/collectory.py.
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timestamp"`
}

var dbMutex sync.Mutex = sync.Mutex{}
var db map[string][]Record = make(map[string][]Record)

func PutRecord(w http.ResponseWriter, r *http.Request) {
	var record Record
	record.Status = -1
	record.TimeStamp = time.Now().UTC()

	vars := mux.Vars(r)
	name := vars["name"]
	if _, ok := db[name]; ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Cannot create record for '%s'. Records already exist.\n", name)
	} else {
		dbMutex.Lock()
		db[name] = []Record{record}
		dbMutex.Unlock()
		w.WriteHeader(http.StatusCreated)
	}
}

func PostRecord(w http.ResponseWriter, r *http.Request) {
	const maxBytes = 256
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBytes))
	if err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		fmt.Fprintf(w, "Request exceeds the maximal size of %s bytes.\n", maxBytes)
		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "The server failed to close the input stream.\n", maxBytes)
		return
	}

	var record Record
	if err := json.Unmarshal(body, &record); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Cannot convert the following string into a json object:\n%s\n",
			string(body[:]))
		return
	}
	record.TimeStamp = time.Now().UTC()

	vars := mux.Vars(r)
	name := vars["name"]
	const maxRecordNum = 20
	dbMutex.Lock()
	if records, ok := db[name]; ok {
		if len(records) >= maxRecordNum {
			db[name] = append(records[1:], record)
		} else {
			db[name] = append(records, record)
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Cannot append record for '%s'. No records exist.\n", name)
	}
	dbMutex.Unlock()
}

func GetAllRecords(w http.ResponseWriter, r *http.Request) {
	dbMutex.Lock()
	if err := json.NewEncoder(w).Encode(db); err != nil {
		panic(err)
	}
	dbMutex.Unlock()
}

func GetRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	dbMutex.Lock()
	if records, ok := db[name]; ok {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(records); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(http.StatusGone)
		fmt.Fprintf(w, "No record found for '%s'.\n", name)
	}
	dbMutex.Unlock()
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "not implemented yet")
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	dbMutex.Lock()
	delete(db, name)
	dbMutex.Unlock()
}

func DeleteAllRecords(w http.ResponseWriter, r *http.Request) {
	dbMutex.Lock()
	db = nil
	db = make(map[string][]Record)
	dbMutex.Unlock()
}
