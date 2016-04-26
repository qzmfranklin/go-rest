package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	//"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Record struct {
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timestamp"`
}

var _db map[string][]Record = nil

func PutRecord(w http.ResponseWriter, r *http.Request) {
	if _db == nil {
		_db = make(map[string][]Record)
	}

	const maxBytes = 1023
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBytes))
	if err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		fmt.Fprintf(w, "Request exceeds the %s bytes.\n", maxBytes)
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
	const maxRecordNum = 10
	if records, ok := _db[name]; ok {
		if len(records) >= maxRecordNum {
			_db[name] = append(records[1:], record)
		} else {
			_db[name] = append(records, record)
		}
	} else {
		_db[name] = []Record{record}
	}
	w.WriteHeader(http.StatusCreated)
}

func GetAllRecords(w http.ResponseWriter, r *http.Request) {
	if _db == nil {
		_db = make(map[string][]Record)
	}
	if err := json.NewEncoder(w).Encode(_db); err != nil {
		panic(err)
	}
}

func GetRecord(w http.ResponseWriter, r *http.Request) {
	if _db == nil {
		_db = make(map[string][]Record)
	}
	vars := mux.Vars(r)
	name := vars["name"]
	if records, ok := _db[name]; ok {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(records); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(http.StatusGone)
		fmt.Fprintf(w, "No record found for '%s'.\n", name)
	}
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "not implemented yet\n")
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	if _db != nil {
		vars := mux.Vars(r)
		name := vars["name"]
		delete(_db, name)
	}
}

func DeleteAllRecords(w http.ResponseWriter, r *http.Request) {
	_db = nil
}
