package router

import (
	"encoding/json"
	"net/http"
)

type HomeRequest struct {
	Message string
}

func Home(w http.ResponseWriter, r *http.Request) {

	msg := HomeRequest{
		Message: "Welcome to RestFul API",
	}
	json.NewEncoder(w).Encode(msg)
}
