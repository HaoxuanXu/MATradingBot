package util

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func ReadYAMLFile(path string) map[interface{}]interface{} {
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(yfile, &data)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range data {
		log.Printf("%s -> %s\n", k, v)
	}

	return data
}
