package openapi

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/ghodss/yaml"
)

var openapiYaml []byte
var openapiJson []byte

type MIME_TYPE_SIMPLE int
type MiddlewareFn func(http.ResponseWriter, *http.Request)

const (
	ANY MIME_TYPE_SIMPLE = 0 + iota
	JSON
	YAML
)

func Openapi(force MIME_TYPE_SIMPLE) MiddlewareFn {
	if openapiYaml == nil {
		var err error
		openapiYaml, err = ioutil.ReadFile("./openapi.yaml")
		if err != nil {
			log.Fatal("Failed to load up openapi.yaml from ./openapi.yaml")
		}

		openapiJson, err = yaml.YAMLToJSON(openapiYaml)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("accept")
		requestingYaml, _ := regexp.Match("[^a-z0-9]yaml$", []byte(contentType))

		returnYaml := (force == YAML || (force == ANY && requestingYaml))
		// return json if force == "json"

		if returnYaml {
			w.Header().Set("Content-Type", "application/yaml")
			w.WriteHeader(200)
			w.Write(openapiYaml)
		} else {
			// JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(openapiJson)
		}

	}
}
