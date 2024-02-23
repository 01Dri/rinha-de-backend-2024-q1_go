package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type Cliente struct {
	Id           int
	Limite       int
	SaldoInicial int
}

type TransacaoDTO struct {
	Valor     int
	Tipo      string
	Descricao string
}

type TransacaoRespostaDTO struct {
	Limite int
	Saldo  int
}

func transacao(id int, transacao TransacaoDTO, conn *sql.DB) (TransacaoRespostaDTO, error) {
	cliente, err := getClientById(id, conn)
	fmt.Println(cliente)
	if err != nil {
		return TransacaoRespostaDTO{}, errors.New("Cliente não encontrado")
	}

	if transacao.Valor > cliente.Limite {
		return TransacaoRespostaDTO{}, errors.New("Valor da transação é maior do que o limite")
	}

	cliente.SaldoInicial -= transacao.Valor

	var respostaDTO TransacaoRespostaDTO
	respostaDTO.Limite = cliente.Limite
	respostaDTO.Saldo = cliente.SaldoInicial

	return respostaDTO, nil
}
