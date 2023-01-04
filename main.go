package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/borerer/nlib-app-kv/database"
	nlib "github.com/borerer/nlib-go"
	"github.com/borerer/nlib-go/har"
)

var (
	mongoClient *database.MongoClient
)

func getQuery(req *nlib.FunctionIn, key string) string {
	for _, query := range req.QueryString {
		if query.Name == key {
			return query.Value
		}
	}
	return ""
}

func getKey(req *nlib.FunctionIn) (*nlib.FunctionOut, error) {
	key := getQuery(req, "key")
	val, err := mongoClient.GetKey(key)
	if errors.Is(err, database.ErrNoDocuments) {
		return har.NewResponse(http.StatusNotFound, "not found", ""), nil
	} else if err != nil {
		return har.Error(err), nil
	}
	return har.Text(val), nil
}

func setKeyGET(req *nlib.FunctionIn) (*nlib.FunctionOut, error) {
	key := getQuery(req, "key")
	value := getQuery(req, "value")
	err := mongoClient.SetKey(key, value)
	if err != nil {
		return har.Error(err), nil
	}
	return har.Text("ok"), nil
}

func setKeyPOST(req *nlib.FunctionIn) (*nlib.FunctionOut, error) {
	parseKeyValue := func(req *nlib.FunctionIn) (string, string) {
		if req.PostData != nil && req.PostData.Text != nil {
			var j map[string]interface{}
			err := json.Unmarshal([]byte(*req.PostData.Text), &j)
			if err == nil {
				key := j["key"].(string)
				switch value := j["value"].(type) {
				case string:
					return key, value
				default:
					buf, err := json.Marshal(value)
					if err == nil {
						return key, string(buf)
					}
				}
			}
		}
		return "", ""
	}

	key, value := parseKeyValue(req)
	err := mongoClient.SetKey(key, value)
	if err != nil {
		return har.Error(err), nil
	}
	return har.Text("ok"), nil
}

func setKey(req *nlib.FunctionIn) (*nlib.FunctionOut, error) {
	if req.Method == "GET" {
		return setKeyGET(req)
	} else if req.Method == "POST" || req.Method == "PUT" {
		return setKeyPOST(req)
	}
	return har.NewResponse(http.StatusMethodNotAllowed, "method not allowed", ""), nil
}

func main() {
	mongoClient = database.NewMongoClient(&database.MongoConfig{
		URI:      os.Getenv("NLIB_MONGO_URI"),
		Database: os.Getenv("NLIB_MONGO_DATABASE"),
	})
	nlib.Must(mongoClient.Start())
	nlib.SetEndpoint(os.Getenv("NLIB_SERVER"))
	nlib.SetAppID("kv")
	nlib.Must(nlib.Connect())
	nlib.RegisterFunction("get", getKey)
	nlib.RegisterFunction("set", setKey)
	nlib.Wait()
}
