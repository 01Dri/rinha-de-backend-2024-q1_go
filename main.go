package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/clientes/{id}/transacoes", transacaoController).Methods(http.MethodPost)
	r.HandleFunc("/clientes/{id}/extrato", extratosController).Methods(http.MethodGet)

	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", r)

}
