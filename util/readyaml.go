package util

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func ReadYAMLFile(path string) map[string]interface{} {
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]interface{})
	err = yaml.Unmarshal(yfile, &data)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range data {
		log.Printf("%s -> %v\n", k, v)
	}

	return data
}
