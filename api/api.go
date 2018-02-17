package api

import (
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"net/http"
	"log"
	//"gopkg.in/mgo.v2/bson"
	"encoding/json"
	//"strings"
	//"strconv"
	//"reflect"
	//"fmt"
	//"gopkg.in/mgo.v2/bson"
	"strings"
	//"strings"
	//"bytes"
	"ipmserver/config"
	//"io"
	//"unicode"
)

var (
	mgoSession     *mgo.Session
	databaseName = config.MongoDatabase
	mongoHost = config.MongoHost
)

type TSearch struct {
	Skip int
	Limit int
	Query string
}

func getSession () *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(mongoHost)
		if err != nil {
			panic(err) // no, not really
		}
	}
	return mgoSession.Clone()
}

func withCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(databaseName).C(collection)
	return s(c)
}

func searchCollection (collection string, q interface{}, skip int, limit int) (searchResults []interface{}, searchErr string) {
	query := func(c *mgo.Collection) error {
		fn := c.Find(q).Skip(skip).Limit(limit).All(&searchResults)

		if limit < 0 {
			fn = c.Find(q).Skip(skip).All(&searchResults)
		}
		return fn
	}
	search := func() error {
		return withCollection(collection, query)
	}
	err := search()
	if err != nil {
		searchErr = "Database Error: " + err.Error()
	}
	return
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func get(w http.ResponseWriter, collection string, query interface{}, skip int, limit int) {
	if !Exist(collection) {
		respondWithError(w, http.StatusNotFound, "bad collection:"+collection)
		return
	}
	searchResults, searchErrors := searchCollection(collection, query, skip, limit)
	if searchErrors != "" {
		respondWithError(w, http.StatusInternalServerError, searchErrors)
	}
	respondWithJSON(w,http.StatusOK, searchResults )
}

func GetQueryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	squery := vars["squery"]
	//log.Printf("debug: get: %s, %s\n", collection, squery)
	// create a json from the query as string
	var search TSearch
	var err error
	err = json.Unmarshal([]byte(squery), &search)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if search.Limit == 0 {
		search.Limit = 100
	}

	// transform bson.M for mgo
	var query interface{}
	search.Query = strings.Replace(search.Query, "'", "\"",-1)
	err = json.Unmarshal([]byte(search.Query), &query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	get(w, collection, query, search.Skip, search.Limit)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	log.Printf("debug: get: %s\n", collection)
	get(w, collection,nil,0,100)
}
