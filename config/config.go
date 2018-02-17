package config
import "os"
import "log"

var MongoHost = "localhost"
var MongoDatabase = "restdb"

const keyMongoHost = "MONGOHOST"
const keyMongoDatabase = "MONGODB"

func getEnv(env string, key string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	} else {
		log.Printf("debug: Environment variable %s is not set, using default value", key)
		return env
	}
}

func init() {
	MongoHost = getEnv(MongoHost, keyMongoHost)
	MongoDatabase = getEnv(MongoDatabase, keyMongoDatabase)
}
