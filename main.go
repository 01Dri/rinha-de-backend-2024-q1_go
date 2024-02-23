package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/clientes/{id}/transacoes", transacaoController).Methods(http.MethodPost)

	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", r)

}
