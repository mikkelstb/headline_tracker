package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)


type General struct {
	Db_type string
	Db_name string
	Username string
	Password string
}


func Read(filename string) *General {

	var cfg General

	fileStream, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}
	defer fileStream.Close()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error: " + err.Error())
	}

	//fmt.Println(string(data))

	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Error: " + err.Error())
	}

	return &cfg
}