package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

const (
	dbName     = "rinha"
	dbPort     = 5432
	dbUser     = "dridev"
	dbPassword = "130722"
)

var dbConfig = DbConfig{
	Name:     dbName,
	Port:     dbPort,
	User:     dbUser,
	Password: dbPassword,
}

func transacaoController(w http.ResponseWriter, r *http.Request) {
	var transcaoDTO TransacaoDTO

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&transcaoDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := startConnection(dbConfig)
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	respostaDTO, err := transacao(id, transcaoDTO, db)
	if err != nil {
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), "Cliente não encontrado") {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Erro na transação", http.StatusUnprocessableEntity)
		}
		return
	}
	saveTransaction(id, db, transcaoDTO)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respostaDTO)
	defer db.Close()
}

func extratosController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	db, err := startConnection(dbConfig)
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}

	res, err := getExtratoByClienteId(id, db)

	if err != nil {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	defer db.Close()

}
