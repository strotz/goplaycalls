package hello

import (
	"encoding/json"
	"net/http"
)

// sample of the payload
type helloPayload struct {
	Name string `json:"name"`
}

func Handler(response http.ResponseWriter, req *http.Request) {
	marshaled, err := json.Marshal(helloPayload{Name: "Double Belomor"})
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = response.Write(marshaled)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
}
