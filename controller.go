package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func TransacaoController(w http.ResponseWriter, r *http.Request) {
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

	db, err := StartConnection("postgres://admin:admin@db:5432/rinha?sslmode=disable")
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

	_, err = SaveTransaction(id, db, transcaoDTO)

	if err != nil {
		if strings.Contains(err.Error(), "Cliente não encontrado") {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "Valor da transação excede o limite") {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if strings.Contains(err.Error(), "Tipo inválido") {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}

		if strings.Contains(err.Error(), "Descrição deve apenas conter entre 1 a 10 caracteres") {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		return
	}

	respostaDTO, err := Transacao(id, transcaoDTO, db)
	if err != nil {
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), "Cliente não encontrado") {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		} else {
			http.Error(w, "O valor da transacao é maior do que o limite do cliente", http.StatusUnprocessableEntity)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respostaDTO)
	defer db.Close()
}

func ExtratosController(w http.ResponseWriter, r *http.Request) {
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

	db, err := StartConnection("postgres://admin:admin@db:5432/rinha?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}

	res, err := GetExtratoByClienteId(id, db)

	if err != nil {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	defer db.Close()

}
