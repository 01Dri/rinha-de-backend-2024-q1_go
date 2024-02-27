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

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	cliente, errC := GetClientById(id)

	if errC != nil {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}

	respostaDTO, err := SaveTransaction(id, transcaoDTO, cliente)

	if err != nil {
		fmt.Println("caiu aq")
		if strings.Contains(err.Error(), "Valor da transação excede o limite") {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if strings.Contains(err.Error(), "Tipo inválido") {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if strings.Contains(err.Error(), "Descrição deve apenas conter entre 1 a 10 caracteres") {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respostaDTO)
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
	res, err := GetExtratoByClienteId(id)

	if err != nil {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
