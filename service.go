package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
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

type CarteiraRespostaDTO struct {
	Saldo       int
	Limite      int
	DataExtrato time.Time
}

type ExtratoRespostaDTO struct {
	Saldo       int
	Tipo        string
	Descricao   string
	RealizadaEm time.Time
}

type ExtratoFinalRespostaDTO struct {
	Saldo             CarteiraRespostaDTO
	UltimasTransacoes []ExtratoRespostaDTO
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

func saveTrasactionService(id int, transacao TransacaoDTO, conn *sql.DB) {
	saveTransaction(id, conn, transacao)
}
