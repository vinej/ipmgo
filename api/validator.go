package api

import (
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type Users struct {
	ID			bson.ObjectId `bson:"_id,omitempty"`
	Email  		string
	Password 	string
	Admin 		bool
	RegisterOn string `json:"register_on"`
}

type validator func (bson.M) (string, bool)

type Company struct {
	ID			bson.ObjectId `bson:"_id,omitempty"`
	Name 		string
}

func validateCompany(b bson.M) (string, bool) {
	if b["name"] != nil {

	}
	return "ok", true
}

var (
	allValidators map[string]validator
	allCollections = "companies:  audits: activities: orders: plans: resources: roles: users:"
)

func init() {
	allValidators = make(map[string]validator)
	allValidators["companies"] = validateCompany
}

/*
func Validate(ot string, b bson.M) (string, bool) {
	var v = allValidators[ot]
	return v(b)
}
*/

func Exist(collection string) bool {
	return strings.Contains(allCollections, collection+":")
}
