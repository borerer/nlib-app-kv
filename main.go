package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/borerer/nlib-app-kv/database"
	nlibgo "github.com/borerer/nlib-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func restore(v interface{}) interface{} {
	switch t := v.(type) {
	case string:
		return t
	case bool:
		return t
	case float64:
		return t
	case primitive.A:
		return t
	case primitive.D:
		return t.Map()
	default:
		println(reflect.TypeOf(v).String())
		buf, _ := json.Marshal(v)
		println(string(buf))
		return t
	}
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
	return restore(val)
}

func setKey(in map[string]interface{}) interface{} {
	key, err := mustString(in, "key")
	if err != nil {
		return err.Error()
	}
	value, ok := in["value"]
	if !ok {
		return errors.New("missing value")
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
