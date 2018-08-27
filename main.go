package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Controller struct {
	ID          int `json:"id"`
	Num_Leds    int `json:"num_leds"`
	Start_Index int `json:"start_index"`
}

type DomeConfig struct {
	ControllerList []Controller `json:"Controllers"`
	LEDs           [][]float32  `json:"led_list"`
}

func UDP_Updates() {
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr, _ := net.ResolveUDPAddr("udp", ":7778")
	fmt.Println("listening on :7778")

	/* Now listen at selected port */
	ServerConn, _ := net.ListenUDP("udp", ServerAddr)
	defer ServerConn.Close()

	buf := make([]byte, 65535)

	for {
		n, _, err := ServerConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("error: ", err)
		}

		if n == 31380 {
			for i := 0; i < 10460; i++ {
				start := (i * 3)
				end := (i * 3) + 3
				config.LEDs[i][0] = float32(buf[start:end][0])
				config.LEDs[i][1] = float32(buf[start:end][1])
				config.LEDs[i][2] = float32(buf[start:end][2])
			}
		}
	}
}

var config DomeConfig

func main() {
	jsonFile, err := os.Open("dome.config")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &config)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/Controllers", Controllers)
	router.HandleFunc("/Controllers/{controllerid}", ControllersById)
	router.HandleFunc("/LED", LEDs)
	router.HandleFunc("/LED/{ledindex}", LEDByIndex)

	go UDP_Updates()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(config)
}

func Controllers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(config.ControllerList)
}

func ControllersById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	controllerId, _ := strconv.Atoi(vars["controllerid"])

	var controller Controller
	found := 0

	for i := 0; i < len(config.ControllerList); i++ {
		if config.ControllerList[i].ID == controllerId {
			controller = config.ControllerList[i]
			found = 1
			break
		}
	}

	if found == 1 {
		json.NewEncoder(w).Encode(controller)
	} else {
		fmt.Fprintf(w, "Controller not found")
	}
}

func LEDs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(config.LEDs)
}

func LEDByIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ledIndex, _ := strconv.Atoi(vars["ledindex"])

	if ledIndex >= 10460 {
		fmt.Fprintf(w, "LED not found")
	} else {
		json.NewEncoder(w).Encode(config.LEDs[ledIndex])
	}
}
