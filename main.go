package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Println("Ola mundo")
	r := mux.NewRouter()

	r.HandleFunc("/clientes/{id}/transacoes", TransacaoController).Methods(http.MethodPost)
	r.HandleFunc("/clientes/{id}/extrato", ExtratosController).Methods(http.MethodGet)

	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", r)

}
