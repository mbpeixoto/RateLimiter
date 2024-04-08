package handlers

import (
	"encoding/json"
	"net/http"
)


func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/json")
	jsonPuro := `{"Mensagem": "Limite não excedido"}`
	err := json.NewEncoder(w).Encode(jsonPuro)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}