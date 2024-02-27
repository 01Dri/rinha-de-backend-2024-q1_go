package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	if err := CreateConnectionPool("postgres://admin:admin@db:5432/rinha?sslmode=disable"); err != nil {
		fmt.Println("Erro ao criar o pool de conexões:", err)
	}

	defer db.Close()
	r.HandleFunc("/clientes/{id}/transacoes", TransacaoController).Methods(http.MethodPost)
	r.HandleFunc("/clientes/{id}/extrato", ExtratosController).Methods(http.MethodGet)

	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", r)

}
