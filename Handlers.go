package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(config)
}

func GetControllers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(config.ControllerList)
}

func GetControllersById(w http.ResponseWriter, r *http.Request) {
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

func GetLEDs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(config.LEDs)
}

func GetLEDByIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ledIndex, _ := strconv.Atoi(vars["ledindex"])

	if ledIndex >= 10460 {
		fmt.Fprintf(w, "LED not found")
	} else {
		json.NewEncoder(w).Encode(config.LEDs[ledIndex])
	}
}

func SetLEDs(w http.ResponseWriter, r *http.Request) {
	var LEDs [][]float32
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &LEDs); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	} else {
		if len(LEDs) == 10460 {
			config.LEDs = LEDs

			SendLed(config.LEDs)
		}
	}

	if err := json.NewEncoder(w).Encode(LEDs); err != nil {
		panic(err)
	}
}

func SetLEDByIndex(w http.ResponseWriter, r *http.Request) {
	var LED LEDWithId
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &LED); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	} else {
		if len(LED.LED) == 3 {
			config.LEDs[LED.ID] = LED.LED

			SendLed(config.LEDs)
		}
	}

	if err := json.NewEncoder(w).Encode(LED); err != nil {
		panic(err)
	}
}
