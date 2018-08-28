package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var config DomeConfig

func main() {
	jsonFile, err := os.Open("dome.config")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &config)

	go UDP_Updates()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
