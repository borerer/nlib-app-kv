package main

import (
	"fmt"
	"os"

	"github.com/borerer/nlib-app-kv/database"
	nlibgo "github.com/borerer/nlib-go"
)

var (
	mongoClient *database.MongoClient
)

func mustString(in map[string]interface{}, key string) (string, error) {
	raw, ok := in[key]
	if !ok {
		return "", fmt.Errorf("missing %s", key)
	}
	str, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("invalid type %s", key)
	}
	return str, nil
}

func getKey(in map[string]interface{}) interface{} {
	key, err := mustString(in, "key")
	if err != nil {
		return err.Error()
	}
	val, err := mongoClient.GetKey(key)
	if err != nil {
		return err.Error()
	}
	return val
}

func setKey(in map[string]interface{}) interface{} {
	key, err := mustString(in, "key")
	if err != nil {
		return err.Error()
	}
	value, err := mustString(in, "value")
	if err != nil {
		return err.Error()
	}
	err = mongoClient.SetKey(key, value)
	if err != nil {
		return err.Error()
	}
	return "ok"
}

func wait() {
	ch := make(chan bool)
	<-ch
}

func main() {
	mongoClient = database.NewMongoClient(&database.MongoConfig{
		URI:      os.Getenv("NLIB_MONGO_URI"),
		Database: os.Getenv("NLIB_MONGO_DATABASE"),
	})
	if err := mongoClient.Start(); err != nil {
		println(err.Error())
		return
	}
	nlib := nlibgo.NewClient(os.Getenv("NLIB_SERVER"), "kv")
	nlib.RegisterFunction("get", getKey)
	nlib.RegisterFunction("set", setKey)
	wait()
}
