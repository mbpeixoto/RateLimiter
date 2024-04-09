package handlers

import (
	"fmt"
	"net/http"
)

func HomeServer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Ainda não aingiu o limite de requisições.")

}
